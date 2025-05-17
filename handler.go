package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// Handle GET /problems - Returns all available problems
func handleProblems(w http.ResponseWriter, r *http.Request) {
	db := Db{}
	ConnectDb(&db)

	var problems []Problem
	err := db.GetDb().Preload("TestCases").Find(&problems).Error
	if err != nil {
		http.Error(w, "Failed to fetch problems", http.StatusInternalServerError)
		return
	}

	// Return only necessary problem data (hide test cases expected output)
	var response []map[string]interface{}
	for _, p := range problems {
		prob := map[string]interface{}{
			"id":          p.ID,
			"title":       p.Title,
			"description": p.Description,
			"test_cases":  len(p.TestCases),
		}
		response = append(response, prob)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Handle GET /problems/{id} - Returns a specific problem
func handleProblem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	problemID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid problem ID", http.StatusBadRequest)
		return
	}

	db := Db{}
	ConnectDb(&db)

	var problem Problem
	if err := db.GetDb().First(&problem, problemID).Error; err != nil {
		http.Error(w, "Problem not found", http.StatusNotFound)
		return
	}

	// Get test cases sample inputs only (not expected outputs)
	var testCases []TestCase
	if err := db.GetDb().Where("problem_id = ?", problemID).Find(&testCases).Error; err != nil {
		http.Error(w, "Failed to fetch test cases", http.StatusInternalServerError)
		return
	}

	// Create response with only visible data
	response := map[string]interface{}{
		"id":          problem.ID,
		"title":       problem.Title,
		"description": problem.Description,
		"samples":     []map[string]string{},
	}

	// Add sample inputs (first 2 test cases)
	samples := []map[string]string{}
	for i, tc := range testCases {
		if i < 2 { // Only show first 2 test cases as samples
			samples = append(samples, map[string]string{
				"input":  tc.Input,
				"output": tc.ExpectedOutput,
			})
		}
	}
	response["samples"] = samples

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Handle POST /submit/{id} - Handles code submission
func handleCodeSubmission(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	problemID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid problem ID", http.StatusBadRequest)
		return
	}

	var submission struct {
		Code string `json:"code"`
	}

	if err := json.NewDecoder(r.Body).Decode(&submission); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if submission.Code == "" {
		http.Error(w, "Code cannot be empty", http.StatusBadRequest)
		return
	}

	db := Db{}
	ConnectDb(&db)

	var problem Problem
	if err := db.GetDb().Preload("TestCases").First(&problem, problemID).Error; err != nil {
		http.Error(w, "Problem not found", http.StatusNotFound)
		return
	}

	fmt.Printf("Judging submission for problem #%d\n", problemID)
	result, err := JudgeCode(submission.Code, problem.TestCases)
	if err != nil {
		http.Error(w, "Error while judging: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Prepare detailed response with additional info
	response := map[string]interface{}{
		"success":      result.Passed == result.Total,
		"passed":       result.Passed,
		"total":        result.Total,
		"failed_cases": result.FailedCases,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}