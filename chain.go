package Golf

type Chain struct {
	middlewareHandlers []middlewareHandler
}

func NewChain(handlerArray []middlewareHandler) Chain {
	c := Chain{}
	c.middlewareHandlers = handlerArray
	return c
}

func (c Chain) Final(fn requestHandler) requestHandler {
	final := fn

	for i := len(c.middlewareHandlers) - 1; i >= 0; i-- {
		final = c.middlewareHandlers[i](final)
	}

	return final
}
