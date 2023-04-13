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
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"github.com/google/cel-go/common/types/traits"
)

func basic_policy_test() {
    /*
        - Configure the compiler/evaluation environment
        - Declare custom variables and functions that are used in the expression (aka policy rule)
        - Dyn represents a dynamic type. This kind only exists at type-check time?
        - Variable declarations only needed for type checking to ensure correctness/semantics of expression.
            Type checking helps improve safety/performance for evaluation though difficult to do for dynamic json
            where type inferencing is minimal.
        - Function declarations required for expression evaluation.
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
    exp2 := "properties.contains('mode','standard1')";
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

    input := make(map[string]any);
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
    defer func() {completed <- true}();
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
                os.Exit(1);
            }
        }
    }()

    // Engine provides an eval timeout which can be useful
    out, det, err := program.ContextEval(ctx, input)

    report(out, det, err)
}


// Taken from here -> https://github.com/google/cel-go/blob/master/codelab/codelab.go#L196
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

// Taken from here -> https://github.com/google/cel-go/blob/master/codelab/solution/codelab.go#L466
func mapContainsKeyValue(args ...ref.Val) ref.Val {
    // The declaration of the function ensures that only arguments which match
    // the mapContainsKey signature will be provided to the function.
    m := args[0].(traits.Mapper)
    fmt.Println(m)
  
    // CEL has many interfaces for dealing with different type abstractions.
    // The traits.Mapper interface unifies field presence testing on proto
    // messages and maps.
    key := args[1]
    v, found := m.Find(key)
  
    // If not found and the value was non-nil, the value is an error per the
    // `Find` contract. Propagate it accordingly. Such an error might occur with
    // a map whose key-type is listed as 'dyn'.
    if !found {
      if v != nil {
        return types.ValOrErr(v, "unsupported key type")
      }
      // Return CEL False if the key was not found.
      return types.False
    }
    // Otherwise whether the value at the key equals the value provided.
    return v.Equal(args[2])
  }

func main() {
    basic_policy_test()
}
