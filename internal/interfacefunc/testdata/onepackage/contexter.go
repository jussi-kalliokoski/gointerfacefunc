package onepackage

import "context"

// Contexter ...
type Contexter interface {
	Context(context.Context) context.Context
}
