// set-gen is an example usage of go2idl.
//
// Structs in the input directories with the below line in their comments will
// have sets generated for them.
// // +genset
//
// Any builtin type referenced anywhere in the input directories will have a
// set generated for it.
package main

import (
	"fmt"
	"os"
	// "path/filepath"

	"k8s.io/gengo/args"
	"k8s.io/gengo/examples/error-code/generators"
)

func main() {
	arguments := args.Default()
	arguments.InputDirs = []string{"inputs"}
	arguments.OutputBase = "outputs"
	// arguments.GoHeaderFilePath = filepath.Join(args.DefaultSourceTree(),
	// 	"k8s.io/gengo/boilerplate/no-boilerplate.go.txt")
	// arguments.IncludeTestFiles = false

	if err := arguments.Execute(
		generators.NameSystems(),
		generators.DefaultNameSystem(),
		generators.Packages,
	); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
