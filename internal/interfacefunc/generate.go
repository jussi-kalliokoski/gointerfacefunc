package interfacefunc

import (
	"errors"
	"go/ast"
	"go/token"
	"strconv"
	"strings"
	"unicode"
)

// Generate a function that implements the specified interface.
//
// Returns the AST to be printed.
func Generate(sourceDirectory, interfaceName string, pkgs map[string]*ast.Package) ([]ast.Decl, error) {
	for _, pkg := range pkgs {
		for _, file := range pkg.Files {
			if obj, ok := file.Scope.Objects[interfaceName]; ok {
				typeSpec, isTypeSpec := obj.Decl.(*ast.TypeSpec)
				if !isTypeSpec {
					return nil, errors.New(sourceDirectory + "." + interfaceName + " is not a type")
				}
				interfaceType, isInterfaceType := typeSpec.Type.(*ast.InterfaceType)
				if !isInterfaceType {
					return nil, errors.New(sourceDirectory + "." + interfaceName + " is not an interface")
				}
				if len(interfaceType.Methods.List) != 1 || len(interfaceType.Methods.List[0].Names) != 1 {
					return nil, errors.New(sourceDirectory + "." + interfaceName + " has more or less than 1 method")
				}

				interfaceFunc := generateInterfaceFunc(typeSpec, interfaceType)
				implementation := generateImplementation(typeSpec, interfaceType)
				return []ast.Decl{interfaceFunc, implementation}, nil
			}
		}
	}
	return nil, errors.New("did not find " + sourceDirectory + "." + interfaceName)
}

func generateInterfaceFunc(typeSpec *ast.TypeSpec, interfaceType *ast.InterfaceType) *ast.GenDecl {
	return &ast.GenDecl{
		Tok: token.TYPE,
		Specs: []ast.Spec{
			&ast.TypeSpec{
				Name: &ast.Ident{
					Name: interfaceFuncName(typeSpec),
				},
				Type: &ast.FuncType{
					Params:  generateInterfaceFuncParams(interfaceType),
					Results: generateResults(interfaceType),
				},
			},
		},
	}
}

func generateInterfaceFuncParams(interfaceType *ast.InterfaceType) *ast.FieldList {
	mt := interfaceType.Methods.List[0].Type.(*ast.FuncType)
	return mt.Params
}

func generateImplementation(typeSpec *ast.TypeSpec, interfaceType *ast.InterfaceType) *ast.FuncDecl {
	return &ast.FuncDecl{
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{
						{Name: "fn"},
					},
					Type: &ast.Ident{Name: interfaceFuncName(typeSpec)},
				},
			},
		},
		Name: interfaceType.Methods.List[0].Names[0],
		Type: &ast.FuncType{
			Params:  generateImplementationParams(interfaceType),
			Results: generateResults(interfaceType),
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				generateImplementationBody(interfaceType),
			},
		},
	}
}

func generateImplementationBody(interfaceType *ast.InterfaceType) ast.Stmt {
	mt := interfaceType.Methods.List[0].Type.(*ast.FuncType)
	callExpr := &ast.CallExpr{
		Fun:  &ast.Ident{Name: "fn"},
		Args: generateImplementationBodyCallArgs(interfaceType),
	}
	if mt.Results != nil {
		return &ast.ReturnStmt{Results: []ast.Expr{callExpr}}
	}
	return &ast.ExprStmt{X: callExpr}
}

func generateImplementationBodyCallArgs(interfaceType *ast.InterfaceType) []ast.Expr {
	args := make([]ast.Expr, 0)
	for _, f := range generateImplementationParams(interfaceType).List {
		for _, ident := range f.Names {
			args = append(args, ident)
		}
	}
	return args
}

func generateImplementationParams(interfaceType *ast.InterfaceType) *ast.FieldList {
	mt := interfaceType.Methods.List[0].Type.(*ast.FuncType)
	list := make([]*ast.Field, 0, len(mt.Params.List))
	for _, f := range mt.Params.List {
		if f.Names != nil {
			list = append(list, f)
			continue
		}
		name := inferParamName(f.Type)
		list = append(list, &ast.Field{
			Names: []*ast.Ident{{Name: name}},
			Type:  f.Type,
		})
	}
	fieldList := &ast.FieldList{List: list}
	deduplicateNames(fieldList)
	return fieldList
}

func generateResults(interfaceType *ast.InterfaceType) *ast.FieldList {
	mt := interfaceType.Methods.List[0].Type.(*ast.FuncType)
	return mt.Results
}

func interfaceFuncName(typeSpec *ast.TypeSpec) string {
	return typeSpec.Name.Name + "Func"
}

func inferParamName(typ ast.Expr) (name string) {
	switch t := typ.(type) {
	case *ast.ArrayType:
		return inferParamName(t.Elt) + "s"
	case *ast.MapType:
		return inferParamName(t.Value) + "s"
	case *ast.StarExpr:
		return inferParamName(t.X)
	case *ast.Ident:
		return unexported(t.Name)
	case *ast.SelectorExpr:
		lhs := t.X.(*ast.Ident)
		rhs := t.Sel
		if lhs.Name == "context" && rhs.Name == "Context" {
			return "ctx"
		}
		return unexported(rhs.Name)
	default:
		// here we pretty much just give up, the user can figure out these names themselves if they are bothersome
		return "v"
	}
}

func unexported(name string) string {
	prev := 0
	for pos, char := range name {
		if unicode.IsLower(char) {
			if prev == 0 {
				return strings.ToLower(name[:pos]) + name[pos:]
			}
			return strings.ToLower(name[:prev]) + name[prev:]
		}
		prev = pos
	}
	return strings.ToLower(name)
}

func deduplicateNames(fieldList *ast.FieldList) {
	for {
		total := 0
		uniques := map[string][2]int{}
		occurrences := map[string]int{}
		for i, field := range fieldList.List {
			for j, ident := range field.Names {
				total++
				o := occurrences[ident.Name] + 1
				occurrences[ident.Name] = o
				if o == 1 {
					uniques[ident.Name] = [2]int{i, j}
				} else if o == 2 {
					pos := uniques[ident.Name]
					delete(uniques, ident.Name)
					fieldList.List[pos[0]].Names[pos[1]].Name += "1"
					ident.Name += strconv.Itoa(o)
				} else {
					ident.Name += strconv.Itoa(o)
				}
			}
		}
		if len(uniques) == total {
			return
		}
	}
}
