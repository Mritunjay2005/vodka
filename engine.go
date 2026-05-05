package vodka

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// httprouter wrapper
type Engine struct {
	router      *httprouter.Router
	middlewares []HandlerFunc
}

// creates a new router
func New() *Engine {
	return &Engine{
		router:      httprouter.New(),
		middlewares: make([]HandlerFunc, 0),
	}
}

func (e *Engine) Use(middleware ...HandlerFunc) {
	e.middlewares = append(e.middlewares, middleware...)
}

// Runs the http server
func (e *Engine) Run(addr string) error {
	if addr == "" {
		addr = ":8080"
	}

	log.Printf(Green+"Pouring Vodka on %s\n"+Reset, addr)

	// Using net/http
	return http.ListenAndServe(addr, e.router)
}

func makeHandlers(e *Engine, handler HandlerFunc) []HandlerFunc {
	handlers := make([]HandlerFunc, 0, len(e.middlewares)+1)
	handlers = append(handlers, e.middlewares...)
	handlers = append(handlers, handler)

	return handlers
}

func (e *Engine) addRoute(method string, path string, handler HandlerFunc) {
	e.router.Handle(method, path, func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		c := &Context{
			Writer:   w,
			Request:  r,
			Params:   params,
			handlers: makeHandlers(e, handler),
			index:    -1,
		}

		c.Next()
	})
}

func (e *Engine) GET(path string, handler HandlerFunc) {
	e.addRoute(http.MethodGet, path, handler)
}

func (e *Engine) POST(path string, handler HandlerFunc) {
	e.addRoute(http.MethodPost, path, handler)
}

func (e *Engine) PUT(path string, handler HandlerFunc) {
	e.addRoute(http.MethodPut, path, handler)
}

func (e *Engine) DELETE(path string, handler HandlerFunc) {
	e.addRoute(http.MethodDelete, path, handler)
}
