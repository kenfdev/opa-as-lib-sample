package main

import (
	"github.com/open-policy-agent/opa/rego"
	"os"
)

import (
	"context"
	"fmt"
)

func main() {
	ctx := context.Background()

	r := rego.New(
		rego.Query("x = data.example.authz.allow"),
		rego.Module("example.rego",
			`package example.authz

default allow = false
allow { input.subject.user == "bob" }`,
		))

	// Prepare for evaluation
	pq, err := r.PrepareForEval(ctx)

	if err != nil {
		// Handle error.
		fmt.Println(err.Error())
        os.Exit(1)
	}

	// Raw input data that will be used in the first evaluation
	input := map[string]interface{}{
		"method": "GET",
		"path":   []interface{}{"salary", "alice"},
		"subject": map[string]interface{}{
			"user":   "alice",
			"groups": []interface{}{},
		},
	}

	// Run the evaluation
	result, err := evalAndInspect(ctx, pq, input)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	// Handle result/decision.
	fmt.Println("initial result:", result)

	// Update input
	input = map[string]interface{}{
		"method": "GET",
		"path":   []interface{}{"salary", "bob"},
		"subject": map[string]interface{}{
			"user":   "bob",
			"groups": []interface{}{"sales", "marketing"},
		},
	}

	// Run the evaluation with new input
	result, err = evalAndInspect(ctx, pq, input)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	// Handle result/decision.
	fmt.Println("updated result:", result)
}

func evalAndInspect(ctx context.Context, pq rego.PreparedEvalQuery, input map[string]interface{}) (bool, error) {
	rs, err := pq.Eval(ctx, rego.EvalInput(input))

	if err != nil {
		return false, err
	} else if len(rs) == 0 {
		return false, fmt.Errorf("the result was undefined")
	} else if result, ok := rs[0].Bindings["x"].(bool); !ok {
		return false, fmt.Errorf("unexpected result type")
	} else {
		return result, nil
	}

}
