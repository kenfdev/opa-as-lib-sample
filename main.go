package main

import "github.com/open-policy-agent/opa/rego"

import (
	"context"
	"fmt"
)

func main() {
	ctx := context.Background()

	// Create a simple query
	r := rego.New(
		rego.Query("input.x == 1"),
	)

	// Prepare for evaluation
	pq, err := r.PrepareForEval(ctx)

	if err != nil {
		// Handle error.
		panic(err)
	}

	// Raw input data that will be used in the first evaluation
	input := map[string]interface{}{"x": 2}

	// Run the evaluation
	rs, err := pq.Eval(ctx, rego.EvalInput(input))

	if err != nil {
		// Handle error.
		panic(err)
	}

	// Inspect results.
	fmt.Println("initial result:", rs[0].Expressions[0])

	// Update input
	input["x"] = 1

	// Run the evaluation with new input
	rs, err = pq.Eval(ctx, rego.EvalInput(input))

	if err != nil {
		// Handle error.
		panic(err)
	}

	// Inspect results.
	fmt.Println("updated result:", rs[0].Expressions[0])
}
