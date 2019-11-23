// error-code is an example usage
//
// Consts in the input directories with the below line in their comments will
// have sets generated for them.
// // +ErrCode
// // +ErrCode=xxx
// // +ErrCode=xxx,yyy
package main

import (
	"fmt"
	"os"

	"k8s.io/gengo/args"
	"k8s.io/gengo/examples/error-code/generators"
)

func main() {
	arguments := args.Default()
	arguments.InputDirs = []string{"inputs"}
	arguments.OutputBase = "outputs"

	if err := arguments.Execute(
		generators.NameSystems(),
		generators.DefaultNameSystem(),
		generators.Packages,
	); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
