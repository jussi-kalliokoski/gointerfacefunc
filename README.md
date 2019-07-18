# gointerfacefunc

`gointerfacefunc` is a small command line utility for generating function types that implement an interface.
This is useful especially if you're writing tests and need to mock an interface, you can just define inline functions in the tests instead of writing complicated and tedious mocks that cover all cases.

For example if you have the following code:
```go
// awesome/awesome.go

package awesome

import (
    "context"
    "io"
)

type Awesomer interface {
    Awesome(context.Context, io.Reader) error
}
```

And run

    gointerfacefunc awesome Awesomer

You will get the following output:

```go
type AwesomerFunc func(context.Context, io.Reader) error

func (fn AwesomerFunc) Awesome(ctx context.Context, reader io.Reader) error {
    return fn(ctx, reader)
}
```

## Installation

```
go get -u github.com/jussi-kalliokoski/gointerfacefunc/...
```

## License

MIT License. See [LICENSE](LICENSE) for more details.
