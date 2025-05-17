package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type JudgeResult struct {
	Passed      int      `json:"passed"`
	Total       int      `json:"total"`
	FailedCases []int    `json:"failed_cases"`
}

func JudgeCode(code string, testCases []TestCase) (JudgeResult, error) {
	// Create temp directory for this submission
	tempDir, err := os.MkdirTemp("", "submission_*")
	if err != nil {
		return JudgeResult{}, fmt.Errorf("failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir) // Clean up when done

	// Define file paths
	codeFile := filepath.Join(tempDir, "submission.cpp")
	binFile := "submission.out"
	inputFile := filepath.Join(tempDir, "input.txt")
	outputFile := filepath.Join(tempDir, "output.txt")
	expectedFile := filepath.Join(tempDir, "expected.txt")

	// Write code to file
	err = os.WriteFile(codeFile, []byte(code), 0644)
	if err != nil {
		return JudgeResult{}, fmt.Errorf("failed to write code: %v", err)
	}

	// Write a main function wrapper for C++ if it doesn't have one
	if !strings.Contains(code, "main(") {
		return JudgeResult{}, fmt.Errorf("code must contain a main() function")
	}

	// Compile inside Docker
	fmt.Println("Compiling code...")
	cmd := exec.Command("docker", "run", "--rm",
		"-v", fmt.Sprintf("%s:/code", tempDir),
		"-w", "/code",
		"gcc:latest",
		"g++", "-o", binFile, "submission.cpp", "-std=c++17")

	compileOutput, err := cmd.CombinedOutput()
	if err != nil {
		return JudgeResult{}, fmt.Errorf("compilation failed: %s", string(compileOutput))
	}

	passed := 0
	var failedCases []int

	fmt.Printf("Running %d test cases...\n", len(testCases))
	for i, tc := range testCases {
		// Write input to file
		err = os.WriteFile(inputFile, []byte(tc.Input), 0644)
		if err != nil {
			return JudgeResult{}, fmt.Errorf("failed to write input file: %v", err)
		}

		// Write expected output to file
		err = os.WriteFile(expectedFile, []byte(tc.ExpectedOutput), 0644)
		if err != nil {
			return JudgeResult{}, fmt.Errorf("failed to write expected output file: %v", err)
		}

		fmt.Printf("Running test case %d...\n", i+1)
		// Run binary in Docker with input
		cmd := exec.Command("docker", "run", "--rm",
			"-v", fmt.Sprintf("%s:/code", tempDir),
			"-w", "/code",
			"--memory=128m", // Set memory limit
			"--cpus=0.5",    // Set CPU limit
			"gcc:latest",
			"timeout", "2", "sh", "-c", fmt.Sprintf("./%s < input.txt > output.txt", binFile))

		runErr := cmd.Run()
		if runErr != nil {
			fmt.Printf("Test case %d execution error: %v\n", i+1, runErr)
			failedCases = append(failedCases, i+1)
			continue
		}

		// Read output
		output, err := os.ReadFile(outputFile)
		if err != nil {
			fmt.Printf("Failed to read output file for test case %d: %v\n", i+1, err)
			failedCases = append(failedCases, i+1)
			continue
		}

		// Read expected output
		expectedOutput, err := os.ReadFile(expectedFile)
		if err != nil {
			fmt.Printf("Failed to read expected output file for test case %d: %v\n", i+1, err)
			failedCases = append(failedCases, i+1)
			continue
		}

		// Compare output (trimming whitespace)
		actual := strings.TrimSpace(string(output))
		expected := strings.TrimSpace(string(expectedOutput))
		
		if actual == expected {
			passed++
			fmt.Printf("Test case %d: PASSED\n", i+1)
		} else {
			failedCases = append(failedCases, i+1)
			fmt.Printf("Test case %d: FAILED\n", i+1)
			fmt.Printf("  Expected: '%s'\n", expected)
			fmt.Printf("  Actual  : '%s'\n", actual)
		}
	}

	return JudgeResult{
		Passed:      passed,
		Total:       len(testCases),
		FailedCases: failedCases,
	}, nil
}