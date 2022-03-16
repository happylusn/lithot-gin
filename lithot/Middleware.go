package lithot

import (
	"github.com/gin-gonic/gin"
	"sync"
)

type Middleware interface {
	OnRequest(*gin.Context) error
	OnResponse(result interface{}) (interface{}, error)
}

var middlewareHandler *MiddlewareHandler
var middleware_once sync.Once

type MiddlewareHandler struct {
	middlewares []Middleware
}

func getMiddlewareHandler() *MiddlewareHandler {
	middleware_once.Do(func() {
		middlewareHandler = &MiddlewareHandler{}
	})
	return middlewareHandler
}

func NewMiddlewareHandler() *MiddlewareHandler {
	return &MiddlewareHandler{}
}

func (this *MiddlewareHandler) AddMiddleware(m ...Middleware) {
	if m != nil && len(m) > 0 {
		this.middlewares = append(this.middlewares, m...)
	}
}

func (this *MiddlewareHandler) before(ctx *gin.Context) {
	for _, m := range this.middlewares {
		err := m.OnRequest(ctx)
		if err != nil {
			//Throw(err.Error(), 400, ctx)
			panic(err)
		}
	}
}

func (this *MiddlewareHandler) after(ctx *gin.Context, ret interface{}) interface{} {
	var result = ret
	for _, m := range this.middlewares {
		r, err := m.OnResponse(result)
		if err != nil {
			//Throw(err.Error(), 400, ctx)
			panic(err)
		}
		result = r
	}
	return result
}

func (this *MiddlewareHandler) handlerMiddleware(responder Responder, ctx *gin.Context) interface{} {
	this.before(ctx)
	var ret interface{}
	innerNode := getInnerRouter().getRoute(ctx.Request.Method, ctx.Request.URL.Path)
	var innerMiddlewareHandler *MiddlewareHandler
	if innerNode.fullPath != "" && innerNode.handlers != nil { //create inner MiddlewareHandler for route-level middlerware.  hook like
		if fs, ok := innerNode.handlers.([]Middleware); ok {
			innerMiddlewareHandler = NewMiddlewareHandler()
			innerMiddlewareHandler.AddMiddleware(fs...)
		}
	}
	// exec route-level middleware
	if innerMiddlewareHandler != nil {
		innerMiddlewareHandler.before(ctx)
	}
	if s1, ok := responder.(StringResponder); ok {
		ret = s1(ctx)
	}
	if s2, ok := responder.(JsonResponder); ok {
		ret = s2(ctx)
	}
	if s3, ok := responder.(SqlResponder); ok {
		ret = s3(ctx)
	}
	if s4, ok := responder.(SqlQueryResponder); ok {
		ret = s4(ctx)
	}
	if s5, ok := responder.(VoidResponder); ok {
		s5(ctx)
		ret = struct{}{}
	}
	// exec route-level middleware
	if innerMiddlewareHandler != nil {
		ret = innerMiddlewareHandler.after(ctx, ret)
	}
	return this.after(ctx, ret)
}
