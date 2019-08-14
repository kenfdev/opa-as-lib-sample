package main

import (
	"encoding/json"
	"github.com/open-policy-agent/opa/rego"
	"os"
)

import (
	"context"
	"fmt"
)

func main() {
	ctx := context.Background()

	// Unmarshal the input JSON
	var input map[string]interface{}
	inputStr := os.Args[1]
	err := json.Unmarshal([]byte(inputStr), &input)

	if err != nil {
		// Handle error.
		fmt.Println(err.Error())
		os.Exit(1)
	}

	// Create a simple query
	r := rego.New(
		rego.Query("data.example.allow"),
		rego.Module("example.rego",
			`package example
default allow = false
allow { input.x == 1 }`,
		),
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
