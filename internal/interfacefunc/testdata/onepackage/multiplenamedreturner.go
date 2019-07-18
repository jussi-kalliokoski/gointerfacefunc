package onepackage

// MultipleNamedReturner ...
type MultipleNamedReturner interface {
	ReturnMultipleNamed() (err1, err2, err3 error)
}
