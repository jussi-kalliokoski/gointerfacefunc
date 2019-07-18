package interfacefunc_test

import (
	"bytes"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/jussi-kalliokoski/gointerfacefunc/internal/interfacefunc"
)

func TestGenerate(t *testing.T) {
	const testDataDir = "testdata"
	testCases, err := ioutil.ReadDir(testDataDir)
	mustNotError(t, err)
	for _, testCase := range testCases {
		name := testCase.Name()
		t.Run(name, func(t *testing.T) {
			dirName := path.Join(testDataDir, name)
			fileSet := token.NewFileSet()
			pkgs, err := parser.ParseDir(fileSet, dirName, nil, 0)
			mustNotError(t, err)
			interfaceNames := make([]string, 0)
			for _, pkg := range pkgs {
				for _, file := range pkg.Files {
					for _, obj := range file.Scope.Objects {
						typeSpec, isTypeSpec := obj.Decl.(*ast.TypeSpec)
						if !isTypeSpec {
							continue
						}
						_, isInterfaceType := typeSpec.Type.(*ast.InterfaceType)
						if !isInterfaceType {
							continue
						}
						interfaceNames = append(interfaceNames, typeSpec.Name.Name)
					}
				}
			}
			for _, interfaceName := range interfaceNames {
				t.Run(interfaceName, func(t *testing.T) {
					decls, err := interfacefunc.Generate(dirName, interfaceName, pkgs)
					mustNotError(t, err)
					buf := &bytes.Buffer{}
					for i, decl := range decls {
						printer.Fprint(buf, token.NewFileSet(), decl)
						if i < len(decls)-1 {
							buf.WriteString("\n\n")
						}
					}
					buf.WriteString("\n")
					received := buf.String()
					expectedFilename := path.Join(dirName, strings.ToLower(interfaceName)+"_func.go.txt")
					expectedData, err := ioutil.ReadFile(expectedFilename)
					if os.IsNotExist(err) {
						mustNotError(t, ioutil.WriteFile(expectedFilename, []byte(received), 0644))
						t.Logf("\n%s", received)
						t.Fatal("recorded new snapshot - rerun test to make it pass without recording to make it pass")
					}
					mustNotError(t, err)
					expected := string(expectedData)
					if expected != received {
						t.Fatalf("expected \n%s\n\nreceived\n\n%s", expected, received)
					}
				})
			}
		})
	}

	t.Run("error cases", func(t *testing.T) {
		tests := []struct {
			name          string
			interfaceName string
			pkgs          map[string]*ast.Package
		}{
			{
				"name not in scope",
				"Missing",
				map[string]*ast.Package{},
			},
			{
				"name is not a type",
				"NotAType",
				map[string]*ast.Package{
					"foo": {
						Files: map[string]*ast.File{
							"foo.go": {
								Scope: &ast.Scope{
									Objects: map[string]*ast.Object{
										"NotAType": &ast.Object{Decl: &ast.FuncDecl{}},
									},
								},
							},
						},
					},
				},
			},
			{
				"name is not an interface",
				"NotAnInterfaceType",
				map[string]*ast.Package{
					"foo": {
						Files: map[string]*ast.File{
							"foo.go": {
								Scope: &ast.Scope{
									Objects: map[string]*ast.Object{
										"NotAnInterfaceType": &ast.Object{
											Decl: &ast.TypeSpec{
												Name: &ast.Ident{Name: "NotAType"},
												Type: &ast.Ident{Name: "123"},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			{
				"interface has multiple methods",
				"MultipleMethods",
				map[string]*ast.Package{
					"foo": {
						Files: map[string]*ast.File{
							"foo.go": {
								Scope: &ast.Scope{
									Objects: map[string]*ast.Object{
										"MultipleMethods": &ast.Object{
											Decl: &ast.TypeSpec{
												Name: &ast.Ident{Name: "NotAType"},
												Type: &ast.InterfaceType{
													Methods: &ast.FieldList{
														List: []*ast.Field{
															{},
															{},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			{
				"interface has multiple method names",
				"MultipleMethodNames",
				map[string]*ast.Package{
					"foo": {
						Files: map[string]*ast.File{
							"foo.go": {
								Scope: &ast.Scope{
									Objects: map[string]*ast.Object{
										"MultipleMethodNames": &ast.Object{
											Decl: &ast.TypeSpec{
												Name: &ast.Ident{Name: "NotAType"},
												Type: &ast.InterfaceType{
													Methods: &ast.FieldList{
														List: []*ast.Field{
															{
																Names: []*ast.Ident{
																	{Name: "A"},
																	{Name: "B"},
																},
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				_, err := interfacefunc.Generate("", tt.interfaceName, tt.pkgs)
				if err == nil {
					t.Fatal("expected an error")
				}
			})
		}
	})
}

func mustNotError(tb testing.TB, err error) {
	tb.Helper()
	if err != nil {
		tb.Fatal(err)
	}
}
