package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"sort"
	"time"

	"github.com/golang/glog"
	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/checker/decls"
	"github.com/google/cel-go/common/types/ref"
)

func basic_policy_test() {
    /*
        Configure the compiler/evaluation environment
        Declare custom variables and functions that are used in the expression (aka policy rule)
        Dyn represents a dynamic type. This kind only exists at type-check time?
        Variable declarations only needed for type checking to ensure correctness/semantics of expression
    */
    var propertiesType = decls.NewMapType(decls.String, decls.Dyn)
    env, err := cel.NewEnv(
        cel.Declarations(decls.NewIdent("properties", propertiesType, nil)),
    )
    if err != nil {
        glog.Exitf("env error: %v", err)
    }

    /*
        Parse the expression (policy rule) to AST and check the expression for correctness
    */
    exp := "properties.mode == 'standard'";
    ast, iss := env.Parse(exp)
    if iss.Err() != nil {
        glog.Exit(iss.Err())
    }
    checkedAst, iss := env.Check(ast)
    // Report semantic errors, if present.
    if iss.Err() != nil {
        // Can't type check arbitrary json expression unless declaring the type for the expression.
        // ERROR: <input>:1:1: undeclared reference to 'properties' (in container '')
        glog.Exit(iss.Err())
    }
    // Check what the expression output is evaluating based on its operands
    if !reflect.DeepEqual(checkedAst.OutputType(), cel.BoolType) {
        glog.Exitf("Got %v, wanted %v output type", checkedAst.OutputType(), cel.BoolType)
    }

    /*
        Evaluate the program and its expression compiled to AST against the input
    */
    program, err := env.Program(ast)
    if err != nil {
        glog.Exitf("program error: %v", err)
    }

    input := make(map[string]any);
    jsonStr := `
        {
            "location" : "eastus",
            "properties": {
                "mode" : "standard"
            }
        }`
    json.Unmarshal([]byte(jsonStr), &input)
    fmt.Println(input)
    ctx, cancel := context.WithCancel(context.Background())
    go func() {
        // After 500ms cancel the request...
        time.Sleep(500 * time.Millisecond)
        cancel()
        fmt.Println("Cancel the request, should not report")
        os.Exit(1);
    }()

    // Engine provides an eval timeout which can be useful
    out, det, err := program.ContextEval(ctx, input)
    report(out, det, err)
}


// Taken from here - https://github.com/google/cel-go/blob/master/codelab/codelab.go#L196
func report(result ref.Val, details *cel.EvalDetails, err error) {
    fmt.Println("------ result ------")
    if err != nil {
        fmt.Printf("error: %s\n", err)
    } else {
        fmt.Printf("value: %v (%T)\n", result, result)
    }
    if details != nil {
        fmt.Printf("\n------ eval states ------\n")
        state := details.State()
        stateIDs := state.IDs()
        ids := make([]int, len(stateIDs), len(stateIDs))
        for i, id := range stateIDs {
            ids[i] = int(id)
        }
        sort.Ints(ids)
        for _, id := range ids {
            v, found := state.Value(int64(id))
            if !found {
                continue
            }
            fmt.Printf("%d: %v (%T)\n", id, v, v)
        }
    }
}

func main() {
    basic_policy_test()
}