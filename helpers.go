package main

import (
	"fmt"
	"sort"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/checker/decls"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"github.com/google/cel-go/common/types/traits"
)

func setupGlobalEnv() (*cel.Env, error) {
    mapAB := cel.MapType(cel.TypeParamType("A"), cel.TypeParamType("B"))
    return cel.NewEnv(
        cel.Declarations(decls.NewIdent("input.properties", decls.NewMapType(decls.String, decls.Dyn), nil)),
        cel.Declarations(decls.NewIdent("input.type", decls.String, nil)),
        cel.Function("contains",
            cel.MemberOverload(
                "map_contains_key_value",
                []*cel.Type{mapAB, cel.TypeParamType("A"), cel.TypeParamType("B")},
                cel.BoolType,
                cel.FunctionBinding(mapContainsKeyValue),
            ),
        ),
    )
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
