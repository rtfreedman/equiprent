package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"equiprent/internal/util/log"

	"github.com/gorilla/mux"
)

func init() {
	initMiddleware()
}

func writeResponse(w http.ResponseWriter, data interface{}) {
	b, err := json.Marshal(data)
	if err != nil {
		log.Logger.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(b)
}

func writeError(w http.ResponseWriter, msg string) {
	b, err := json.Marshal(struct {
		Error string `json:"error"`
	}{msg})
	if err != nil {
		log.Logger.Warn(err.Error())
		return
	}
	w.Write(b)
}

var ctx = context.Background()
var cancel context.CancelFunc

// Start the api
func Start(port int) (err error) {
	ctx, cancel = context.WithCancel(ctx)
	r := mux.NewRouter()
	linkRoutes(r.PathPrefix("/api").Subrouter())
	srv := &http.Server{
		Handler:      r,
		Addr:         fmt.Sprintf("0.0.0.0:%d", port),
		WriteTimeout: 45 * time.Second,
		ReadTimeout:  45 * time.Second,
	}
	var apiExited = make(chan bool, 1)
	go func() {
		err = srv.ListenAndServe()
		if err != nil {
			log.Logger.Panic(err.Error())
		}
		apiExited <- true
	}()
	log.Logger.Debug("API started on port ", port)
	select {
	case <-ctx.Done():
	case <-apiExited:
	}
	return
}

func Stop() {
	cancel()
}

func linkRoutes(r *mux.Router) {
	routeList := getRoutes()
	for _, routes := range routeList {
		for _, route := range routes.routes {
			for method, handler := range route.method_handlers {
				/* middleware is resolved in a 0->len(routes middleware) then 0->len(route middleware)
				* example:
				* routesMiddleware: a, b, c
				* routeMiddleware: d, e, f
				* when the http request is sent it will experience the middleware in the following order:
				* a, b, c, d, e, f
				* therefore we need to reverse the order of these middlewares as they're applied to the handler:
				* f, e, d, c, b, a
				* so that we hit a first
				 */
				handlerFunc := handler.ServeHTTP
				for idx := len(route.middleware) - 1; idx >= 0; idx-- {
					handlerFunc = route.middleware[idx].Handle(handlerFunc)
				}
				for idx := len(routes.middleware) - 1; idx >= 0; idx-- {
					handlerFunc = routes.middleware[idx].Handle(handlerFunc)
				}
				if len(route.queries) == 0 {
					r.HandleFunc(routes.prefix+route.endpoint, handlerFunc).Methods(method)
				} else {
					r.HandleFunc(routes.prefix+route.endpoint, handlerFunc).Methods(method).Queries(route.queries...)
				}
			}
		}
	}
}
