package routing

type Router struct {
	RoutingTable RoutingTable
}

func FindPath(path string, routingTable RoutingTable) *[]Route {
	routingManager := routingTable.RouteMap[path]
	if routingManager != nil {
		return &routingManager.Routes
	}
	return nil
}
