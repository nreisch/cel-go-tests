package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"time"

	"github.com/golang/glog"
	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/checker/decls"
)

func basic_policy() {
    /*
       - Configure the compiler/evaluation environment
       - Declare custom variables and functions that are used in the expression (aka policy rule)
       - Dyn represents a dynamic type. This kind only exists at type-check time?
       - Variable declarations only needed for type checking to ensure correctness/semantics of expression.
           Type checking helps improve safety/performance for evaluation though difficult to do for dynamic json
           where type inferencing is minimal.
       - Function declarations required for expression evaluation.
       - Variables are essentially like how we treat aliases we load them into the compiler?
    */
    mapAB := cel.MapType(cel.TypeParamType("A"), cel.TypeParamType("B"))
    env, err := cel.NewEnv(
        cel.Declarations(decls.NewIdent("properties", decls.NewMapType(decls.String, decls.Dyn), nil)),
        // Declare custom functions and implementation. Function useful for overloading, can also equivalently create a new Declarations and then pass FunctionBinding.
        cel.Function("contains",
            // For non-member funcs can use Overload
            cel.MemberOverload(
                // Id
                "map_contains_key_value",
                // Args
                []*cel.Type{mapAB, cel.TypeParamType("A"), cel.TypeParamType("B")},
                // Result type
                cel.BoolType,
                // Implementation using cel.FunctionBinding()
                cel.FunctionBinding(mapContainsKeyValue),
            ),
        ),
    )
    if err != nil {
        glog.Exitf("env error: %v", err)
    }

    /*
       Parse the expression (policy rule) to AST and check the expression for correctness
    */
    //exp1 := "properties.mode == 'standard1'";
    exp2 := "properties.contains('mode','standard1')"
    ast, iss := env.Parse(exp2)
    if iss.Err() != nil {
        glog.Exit(iss.Err())
    }
    /*
       Can't type check arbitrary json expression unless declaring the type for the expression.
       ERROR: <input>:1:1: undeclared reference to 'properties' (in container '')
    */
    checkedAst, iss := env.Check(ast)
    // Report semantic errors, if present.
    if iss.Err() != nil {
        glog.Exit(iss.Err())
    }
    // Check what the expression result type / output is evaluating
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

    input := make(map[string]any)
    jsonStr := `
        {
            "location" : "eastus",
            "properties": {
                "mode" : "standard"
            }
        }`
    json.Unmarshal([]byte(jsonStr), &input)
    fmt.Println("Input: ", input)

    ctx, cancel := context.WithCancel(context.Background())
    completed := make(chan bool)
    defer func() { completed <- true }()
    go func() {
        // After 50ms cancel the request unless the func has been quit when the program succeeds...
        time.Sleep(50 * time.Millisecond)
        for {
            select {
            case <-completed:
                fmt.Println("Completed")
                return
            default:
                cancel()
                fmt.Println("Request canceled")
                os.Exit(1)
            }
        }
    }()

    // Engine provides an eval timeout which can be useful
    out, det, err := program.ContextEval(ctx, input)

    report(out, det, err)
}

func compile_and_evaluate(env *cel.Env, expression string, jsonInputStr string, policyCount int) {
    checkedAsts, _ := compile(env, expression, policyCount)
    evaluate(env, checkedAsts, jsonInputStr)
}

func compile(env *cel.Env, expression string, policyCount int) ([]*cel.Ast, error) {
    checkedAsts := []*cel.Ast{}
    for i := 0; i < policyCount; i++ {
        ast, iss := env.Parse(expression)
        if iss.Err() != nil {
            glog.Exit(iss.Err())
        }
        checkedAst, iss := env.Check(ast)
        if iss.Err() != nil {
            glog.Exit(iss.Err())
        }
        if !reflect.DeepEqual(checkedAst.OutputType(), cel.BoolType) {
            glog.Exitf("Got %v, wanted %v output type", checkedAst.OutputType(), cel.BoolType)
        }
        checkedAsts = append(checkedAsts, checkedAst)
    }
    return checkedAsts, nil
}

func evaluate(env *cel.Env, checkedAsts []*cel.Ast, jsonInputStr string) {
    for i := 0; i < len(checkedAsts); i++ {
        ast := checkedAsts[i]
        program, err := env.Program(ast)
        if err != nil {
            glog.Exitf("program error: %v", err)
        }
        input := make(map[string]any)
        json.Unmarshal([]byte(jsonInputStr), &input)
        program.Eval(input)
    }
}

func oom() {
    env, _ := cel.NewEnv(cel.Declarations(
		decls.NewVar("x", decls.NewListType(decls.Int)),
	))
    // Potential for oom? If this increases exponentially...
    exp := "['foo', 'bar'].map(x, [x+x,x+x]).map(x, [x+x,x+x]).map(x, [x+x,x+x]).map(x, [x+x,x+x])"
	program, _ := env.Compile(exp)
	prg, _ := env.Program(program)
    out, det, e := prg.Eval(map[string]interface{}{})
    report(out, det, e)
}

func main() {
    basic_policy()
}
