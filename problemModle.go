package main

import "gorm.io/gorm"

type Problem struct {
	gorm.Model
	Title       string     `json:"title"`
	Description string     `json:"description"`
	HaderFile   string     `json:"hader_file"`
	FuncBody    string     `json:"func_body"`
	MainFunc    string     `json:"main_func"`
	TestCases   []TestCase `json:"test_cases" gorm:"foreignKey:ProblemID"`
}

type TestCase struct {
	gorm.Model
	Input          string `json:"input"`
	ExpectedOutput string `json:"expected_output"`
	ProblemID      uint   `json:"problem_id"` // Foreign key to Problem
}
