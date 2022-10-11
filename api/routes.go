package api

import "net/http"

type routes struct {
	routes     []route
	prefix     string
	middleware []middleware
}

type route struct {
	endpoint        string
	method_handlers map[string]http.Handler // map[method]handler i.e. "GET": get_foo()
	middleware      []middleware
	queries         []string
}

// routes definitions:
var registeredRotues = []routes{
	{ // user routes
		prefix:     "/user",
		middleware: []middleware{&logMiddleware, &authMiddleware},
		routes:     []route{},
	},
	{ // auth - login/logout routes
		prefix:     "/auth",
		middleware: []middleware{&logMiddleware},
		routes:     []route{},
	},
}

func getRoutes() []routes {
	return registeredRotues
}
