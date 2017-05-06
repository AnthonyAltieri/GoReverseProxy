package routing

import (
	"os"
	"fmt"
	"io/ioutil"
	"strings"
	"runtime"
	"path"
)

/*  Routing Table Format
 *  .--------------------.
 *  |IP Port Path Gateway|
 *  '--------------------'
 */

type Route struct {
	Ip 			string
	Port 		string
	Path 		string
	Gateway	string
}

type RoutingTable struct {
	Routes map[string]Route
}

func PrintRoutingTable(routingTable RoutingTable) {
	fmt.Println("Ip         Port  Path        Gateway")

	for key := range routingTable.Routes {
		route := routingTable.Routes[key]
		fmt.Fprintf(
			os.Stdout,
			"%s  %s  %s  %s\n",
			route.Ip,
			route.Port,
			route.Path,
			route.Gateway,
		);
	}
}

func lineToRoute(line string) Route {
	split := strings.Split(line, " ")

	route := new(Route)
	route.Ip      = split[0]
	route.Port    = split[1]
	route.Path    = split[2]
	route.Gateway = split[3]

	return *route
}

func FindPath(path string, routingTable RoutingTable) *Route {
	if route, exists := routingTable.Routes[path] ; exists {
		return &route
	} else {
		return nil
	}
}

func createRoutingTable(routes []Route) RoutingTable {
	routingTable := new(RoutingTable)
	routingTable.Routes = make(map[string]Route)
	for _, route := range routes {
		routingTable.Routes[route.Path] = route
	}
	return *routingTable
}

func ParseRoutingTable() RoutingTable {

	_, filename, _, ok:= runtime.Caller(1)
	if !ok {
		fmt.Fprint(os.Stderr, "Problem with runtime.Caller()\n")
		os.Exit(1)
	}

	data, err := ioutil.ReadFile(path.Dir(filename) + "/config/rtable")

	if err != nil {
		fmt.Fprint(os.Stderr, "Error reading rtable: %s\n", err)
		os.Exit(1)
	}

	routes := make([]Route, 0)
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		routes = append(routes, lineToRoute(line))
	}

	return createRoutingTable(routes)
}
