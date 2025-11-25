package router

type Router interface {
	Group(prefix string, group func(group RouteGroup)) RouteGroup
}

type RouteGroup interface {
	Router
}
