package onepackage

import "io"

// DuplicateNamer ...
type DuplicateNamer interface {
	DuplicateName(io.Reader, io.Reader, io.Reader)
}
