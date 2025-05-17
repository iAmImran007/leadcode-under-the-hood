package main

import (
	"fmt"
	"log"
	"net/http"
	"github.com/gorilla/mux"
)

func main() {
	// Load the env variables
	LoadEnv()

	// Initialize the database
	db := Db{}
	ConnectDb(&db)

	// Check if any problems exist
	var count int64
	db.GetDb().Model(&Problem{}).Count(&count)
	
	// Insert sample coding problems if none exist
	if count == 0 {
		// Problem 1: Sum of Two Numbers
		problem1 := Problem{
			Title:       "Sum of Two Numbers",
			Description: "Write a program that reads two integers from standard input and outputs their sum.\n\n**Input**: Two integers a and b (separated by a space)\n**Output**: The sum of a and b",
			TestCases: []TestCase{
				{Input: "1 2", ExpectedOutput: "3"},
				{Input: "10 5", ExpectedOutput: "15"},
				{Input: "-3 3", ExpectedOutput: "0"},
				{Input: "100 -50", ExpectedOutput: "50"},
			},
		}

		// Problem 2: Reverse String
		problem2 := Problem{
			Title:       "Reverse String",
			Description: "Write a program that reads a string from standard input and outputs the string in reverse order.\n\n**Input**: A string (up to 100 characters)\n**Output**: The string in reverse order",
			TestCases: []TestCase{
				{Input: "hello", ExpectedOutput: "olleh"},
				{Input: "algorithm", ExpectedOutput: "mhtirogla"},
				{Input: "a", ExpectedOutput: "a"},
				{Input: "12345", ExpectedOutput: "54321"},
			},
		}

		// Problem 3: Check Prime Number
		problem3 := Problem{
			Title:       "Check Prime Number",
			Description: "Write a program that determines if a given number is a prime number.\n\n**Input**: An integer n\n**Output**: 'Prime' if n is a prime number, 'Not Prime' otherwise",
			TestCases: []TestCase{
				{Input: "7", ExpectedOutput: "Prime"},
				{Input: "15", ExpectedOutput: "Not Prime"},
				{Input: "2", ExpectedOutput: "Prime"},
				{Input: "1", ExpectedOutput: "Not Prime"},
				{Input: "97", ExpectedOutput: "Prime"},
			},
		}

		// Create the problems
		if err := db.GetDb().Create(&problem1).Error; err != nil {
			log.Println("Failed to insert problem 1:", err)
		}
		if err := db.GetDb().Create(&problem2).Error; err != nil {
			log.Println("Failed to insert problem 2:", err)
		}
		if err := db.GetDb().Create(&problem3).Error; err != nil {
			log.Println("Failed to insert problem 3:", err)
		}

		fmt.Println("Sample coding problems inserted successfully.")
	} else {
		fmt.Println("Database already contains problems. Skipping sample insertion.")
	}

	// Create router
	r := mux.NewRouter()

	// Set up routes
	r.HandleFunc("/problems", handleProblems).Methods("GET")
	r.HandleFunc("/problems/{id}", handleProblem).Methods("GET")
	r.HandleFunc("/submit/{id}", handleCodeSubmission).Methods("POST")
	
	// Add CORS middleware
	r.Use(corsMiddleware)

	// Start server
	fmt.Println("Server is listening on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

// CORS middleware to allow frontend connections
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}
