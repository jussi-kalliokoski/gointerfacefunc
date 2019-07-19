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

## Editor/IDE integrations

### vim

The easiest way to integrate is to install the command line utility and then add a keybind for it in your `.vimrc`, for example:

```vim

autocmd FileType go nnoremap goiffn "*yiw:let @* = system('gointerfacefunc "'.expand('%:h').'" '.@*)<CR>
```

This will add a keybinding that gets triggered when you type `goiffn` in normal mode, using the interface name under cursor and the directory of the current buffer as inputs. The resulting type definition will be copied to your `*` (clipboard) register, so you can then paste it using `"*p` or just `p` if you have `set clipboard=unnamed`.

## License

MIT License. See [LICENSE](LICENSE) for more details.
