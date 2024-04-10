package web

type MiddlewareHandler func(Handler) Handler

// wrapMiddlewares creates a new handler by wrapping middleware around a final
// handler. The middlewares' Handlers will be executed by requests in the order
// they are provided.
func wrapMiddlewares(mdw []MiddlewareHandler, handler Handler) Handler {
	// Usually we write from outter to inner.
	// Loop backwards through the middleware invoking each one. Replace the
	// handler with the new wrapped handler. Looping backwards ensures that the
	// first middleware of the slice is the first to be executed by requests.
	for i := len(mdw) - 1; i >= 0; i-- {
		mwFunc := mdw[i]
		if mwFunc != nil {
			handler = mwFunc(handler)
		}
	}

	return handler
}
