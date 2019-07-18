package onepackage

// Foo ...
type Foo int

// XFoo ...
type XFoo int

// XXFoo ...
type XXFoo int

// X ...
type X int

// Uppercaser ...
type Uppercaser interface {
	Uppercase(X, Foo, XFoo, XXFoo)
}
