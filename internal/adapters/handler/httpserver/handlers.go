package httpserver

type HttpHandlers interface {
	SetupRoutes(router *Router)
}
