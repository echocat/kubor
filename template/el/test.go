package main

import (
	"fmt"
	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/checker/decls"
	exprpb "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
	"log"
)

func main() {
	decl := cel.Declarations(
		decls.NewIdent("i", decls.String, nil),
		decls.NewIdent("you", decls.String, nil),
		decls.NewFunction("foo", decls.NewOverload("foo", []*exprpb.Type{decls.String}, decls.String)),
	)
	env, err := cel.NewEnv(decl)
	if err != nil {
		log.Fatalf("environment creation error: %s\n", err)
	}

	p, iss := env.Parse(`"Hello " + you + "! I'm " + i + "." + foo("bar")`)
	if iss != nil && iss.Err() != nil {
		panic(iss.Err())
	}
	c, iss := env.Check(p)
	if iss != nil && iss.Err() != nil {
		panic(iss.Err())
	}

	prg, err := env.Program(c)
	if err != nil {
		panic(err)
	}

	out, _, err := prg.Eval(cel.Vars(map[string]interface{}{
		"i":   "CEL",
		"you": "world",
		"foo": func(in string) string {
			return in + "bar"
		},
	}))
	if err != nil {
		panic(err)
	}

	fmt.Println(out.Value(), out.Type())
}
