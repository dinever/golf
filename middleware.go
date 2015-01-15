package Golf

import (
	"log"
	"net/http"
	"time"
)

type requestHandler func(req *Request, res *Response)

type middlewareHandler func(next requestHandler) requestHandler

var (
	middlewareChain = []middlewareHandler{loggingHandler, recoverHandler}
	chain           = NewChain(middlewareChain)
)

func loggingHandler(next requestHandler) requestHandler {
	fn := func(req *Request, res *Response) {
		t1 := time.Now()
		next(req, res)
		t2 := time.Now()
		log.Printf("[%s] %q %v\n", req.Method, req.URL.String(), t2.Sub(t1))
	}
	return fn
}

func recoverHandler(next requestHandler) requestHandler {
	fn := func(req *Request, res *Response) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("panic: %+v", err)
				http.Error(res, http.StatusText(500), 500)
			}
		}()
		next(req, res)
	}

	return fn
}
