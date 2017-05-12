package routing

import (
	"os"
	"fmt"
	"io/ioutil"
	"strings"
	"runtime"
	"path"
	"strconv"
	"reverseproxy/utils"
)

/*  Routing Table Format
 *  .-----------------------------------.
 *  |IP | Port | Path | Version | Weight|
 *  '-----------------------------------'
 *
 *  For optional fields Version and Weight if one is not specified
 *  use the $ character
 */

type Route struct {
	Ip 			string
	Port 		string
	Path 		string
	Version string
	Weight  int
}

type RouteManager struct {
	Scheduler Scheduler
	Routes    []Route
}

type RoutingTable struct {
	RouteMap map[string]*RouteManager
}


func PrintRoutingTable(routingTable RoutingTable) {

	var longestIp int = 0
	var longestPort int = 0
	var longestPath int = 0
	var longestVersion int = 0

	for key := range routingTable.RouteMap {
		routeManager := routingTable.RouteMap[key]
		routes := routeManager.Routes
		for _, route := range routes {
			if len(route.Ip) > longestIp {
				longestIp = len(route.Ip)
			}
			if len(route.Port) > longestPort {
				longestPort = len(route.Port)
			}
			if len(route.Path) > longestPath {
				longestPath = len(route.Path)
			}
			if len(route.Version) > longestVersion {
				longestVersion = len(route.Version)
			}
		}
	}

	spacesIp := strings.Repeat(" ", longestIp - 2)
	spacesPort := strings.Repeat(" ", longestPort - 4)
	spacesPath := strings.Repeat(" ", longestPath - 4)
	spacesVersion := strings.Repeat(" ", longestVersion - 7)

	headerLine := fmt.Sprintf(
		"Ip%s | Port%s | Path%s | Version%s | Weight\n",
		spacesIp,
		spacesPort,
		spacesPath,
		spacesVersion,
	)
	dividerLine := strings.Repeat("-", len(headerLine))
	fmt.Print(headerLine)
	fmt.Println(dividerLine)

	for key := range routingTable.RouteMap {
		routeManager := routingTable.RouteMap[key]
		routes := routeManager.Routes
		for _, route := range routes {
			fmt.Fprintf(
				os.Stdout,
				"%s | %s | %s | %s | %s\n",
				route.Ip,
				route.Port,
				route.Path,
				route.Version,
				fmt.Sprintf("%d", route.Weight),
			);
		}
	}
}

func routingTableParseError(field string) {
	fmt.Fprintf(os.Stderr, "Field: [%s] has no value\n", field)
	os.Exit(1)
}

func lineToRoute(line string) Route {
	fmt.Println("line: ", line)
	split := strings.Split(line, "|")
	fmt.Println("split: ", split)


	route := new(Route)
	if ip := strings.TrimSpace(split[0]) ; ip == "$" {
		routingTableParseError("ip")
	} else {
		route.Ip = ip
	}

	if port := strings.TrimSpace(split[1]) ; port == "$" {
		routingTableParseError("port")
	} else {
		route.Port = port
	}

	if path := strings.TrimSpace(split[2]) ; path == "$" {
		routingTableParseError("path")
	} else {
		route.Path = path
	}

	if version := strings.TrimSpace(split[3]) ; version == "$" {
		route.Version = utils.ShortUid(&route.Path)
	} else {
		route.Version = version
	}

	if weight := strings.TrimSpace(split[4]) ; weight == "$" {
		route.Weight = 1
	} else {
		route.Weight, _ = strconv.Atoi(weight)
	}

	return *route
}


func createRoutingTable(routes []Route) RoutingTable {
	routingTable := new(RoutingTable)
	routingTable.RouteMap = make(map[string]*RouteManager)
	for _, route := range routes {
		routeManager := routingTable.RouteMap[route.Path]
		if routeManager == nil {
			routingTable.RouteMap[route.Path] = new(RouteManager)
			routeManager = routingTable.RouteMap[route.Path]
		}
		routeManager.Routes = append(routeManager.Routes, route)
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
