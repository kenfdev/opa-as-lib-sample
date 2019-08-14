package main

import (
	"encoding/json"
	"github.com/open-policy-agent/opa/rego"
	"io/ioutil"
	"os"
	"path/filepath"
)

import (
	"context"
	"fmt"
)

// createModules creates rego.Module slices from the file paths provided
func createModules(filePaths []string) ([]func(r *rego.Rego), error) {
	var modules []func(r *rego.Rego)

	for _, path := range filePaths {
		file, err := os.Stat(path)
		if err != nil {
			// Handle error.
			return nil, err
		}

		var dirPath string
		dirPath = filepath.Dir(path)

		out, err := ioutil.ReadFile(dirPath + "/" + file.Name())
		if err != nil {
			// Handle error.
			return nil, err
		}

		module := rego.Module(file.Name(), string(out))
		modules = append(modules, module)
	}

	return modules, nil
}

func main() {
	ctx := context.Background()

	// Unmarshal the input JSON
	var input map[string]interface{}
	inputStr := os.Args[1]
	err := json.Unmarshal([]byte(inputStr), &input)

	filePaths := os.Args[2:]

	modules, err := createModules(filePaths)
	if err != nil {
		// Handle error.
		fmt.Println(err.Error())
		os.Exit(1)
	}

	// Create a simple query
	options := append([]func(r *rego.Rego){rego.Query("data.example.allow")}, modules...)
	r := rego.New(
		options...,
	)

	// Prepare for evaluation
	pq, err := r.PrepareForEval(ctx)

	if err != nil {
		// Handle error.
		fmt.Println(err.Error())
		os.Exit(1)
	}

	// Run the evaluation
	rs, err := pq.Eval(ctx, rego.EvalInput(input))

	if err != nil {
		// Handle error.
		fmt.Println(err.Error())
		os.Exit(1)
	}

	// Inspect results.
	fmt.Println("result:", rs[0].Expressions[0])
}
