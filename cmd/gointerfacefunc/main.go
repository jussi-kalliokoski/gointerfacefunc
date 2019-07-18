package main

import (
	"fmt"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"path"
	"strings"

	"github.com/jussi-kalliokoski/gointerfacefunc/internal/interfacefunc"
)

var helpText = strings.TrimSpace(fmt.Sprintf(`
Usage:
%s %s InterfaceName

Flags:
  --help Shows this message.
`,
	os.Args[0],
	path.Join("file", "path", "to", "package"),
))

func main() {
	if len(os.Args) != 3 || (len(os.Args) > 1 && os.Args[1] == "--help") {
		fmt.Println(helpText)
		os.Exit(1)
	}
	sourceDirectory := os.Args[1]
	interfaceName := os.Args[2]
	fileSet := token.NewFileSet()
	pkgs, err := parser.ParseDir(fileSet, sourceDirectory, nil, 0)
	if err != nil {
		panic(err)
	}
	decls, err := interfacefunc.Generate(sourceDirectory, interfaceName, pkgs)
	if err != nil {
		panic(err)
	}
	for i, decl := range decls {
		printer.Fprint(os.Stdout, token.NewFileSet(), decl)
		fmt.Println("")
		if i < len(decls)-1 {
			fmt.Println("")
		}
	}
}
