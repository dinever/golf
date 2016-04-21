package golf

import "golang.org/x/net/context"

var _ context.Context = &Context{}

type ctxKey int

const (
	routeCtxKey ctxKey = iota
)

// A Context is the default routing context set on the root node of a
// request context to track URL parameters and an optional routing path.
type Context struct {
	context.Context

	parent context.Context

	// URL parameter key and values
	Params Parameter

	// Routing path override used by subrouters
	RoutePath string
}

// NewContext creates a new context instance
func NewContext() *Context {
	return NewContextWithParent(context.Background())
}

// NewContextWithParent creates a new context with a parent context specified
func NewContextWithParent(c context.Context) *Context {
	return &Context{
		parent: c,
	}
}

func (p *Context) reset() {
	p.parent = nil
}

func ToContext(c context.Context) *Context {
	if ctx, ok := c.(*Context); ok {
		return ctx
	}
	return NewContextWithParent(c)
}

func Param(c context.Context, key string) (string, error) {
	return ToContext(c).Params.Get(key)
}
