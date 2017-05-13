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
	var longestWeight int = 0

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
			weight := fmt.Sprintf("%d", route.Weight)
			if len(weight) > longestWeight {
				longestWeight = len(weight)
			}
		}
	}

	LENGTH_IP_HEADER := 2
	LENGTH_PORT_HEADER := 4
	LENGTH_PATH_HEADER := 4
	LENGTH_VERSION_HEADER := 7


	// Calculate the spaces after the header names
	numSpacesDividerIp := 0
	if longestIp < LENGTH_IP_HEADER {
		numSpacesDividerIp = 0
	} else {
		numSpacesDividerIp = longestIp - LENGTH_IP_HEADER
	}

	numSpacesDividerPort := 0
	if longestPort < LENGTH_PORT_HEADER {
		numSpacesDividerPort = 0
	} else {
		numSpacesDividerPort = longestPort - LENGTH_PORT_HEADER
	}

	numSpacesDividerPath := 0
	if longestPath < LENGTH_PATH_HEADER {
		numSpacesDividerPath =  0
	} else {
		numSpacesDividerPath = longestPath - LENGTH_PATH_HEADER
	}

	numSpacesDividerVersion := 0
	if longestVersion < LENGTH_VERSION_HEADER {
		numSpacesDividerVersion = 0
	} else {
		numSpacesDividerVersion = longestVersion - LENGTH_VERSION_HEADER
	}
	fmt.Println("numSpacesDividerVersion: ", numSpacesDividerVersion)

	// Create the Header line
	spacesHeaderIp := strings.Repeat(" ", numSpacesDividerIp)
	spacesHeaderPort := strings.Repeat(" ", numSpacesDividerPort)
	spacesHeaderPath := strings.Repeat(" ", numSpacesDividerPath)
	spacesHeaderVersion := strings.Repeat(" ", numSpacesDividerVersion)
	headerLine := fmt.Sprintf(
		"Ip%s | Port%s | Path%s | Version%s | Weight\n",
		spacesHeaderIp,
		spacesHeaderPort,
		spacesHeaderPath,
		spacesHeaderVersion,
	)

	// First column has no space before content
	SPACE_FIRST := 1
	// The rest of the columns have a space before and after content
	SPACE_AFTER := 2


	// Calculate the length of each column
	lengthIpSection := numSpacesDividerIp + LENGTH_IP_HEADER + SPACE_FIRST
	lengthPortSection := numSpacesDividerPort + LENGTH_PORT_HEADER + SPACE_AFTER
	lengthPathSection := numSpacesDividerPath + LENGTH_PATH_HEADER + SPACE_AFTER
	lengthVersionSection := numSpacesDividerVersion + LENGTH_VERSION_HEADER + SPACE_AFTER
	dividerLine := strings.Repeat("-", lengthIpSection)
	dividerLine += "+"
	dividerLine += strings.Repeat("-", lengthPortSection)
	dividerLine += "+"
	dividerLine += strings.Repeat("-", lengthPathSection)
	dividerLine += "+"
	dividerLine += strings.Repeat("-", lengthVersionSection)
	dividerLine += "+"
	weightLength := 0
	if longestWeight + 2 < 8 {
		weightLength = 8
	}
	dividerLine += strings.Repeat("-", weightLength)

	dashLine := strings.Repeat("-", len(dividerLine))

	fmt.Println(dashLine)
	fmt.Print(headerLine)
	fmt.Println(dividerLine)


	for key := range routingTable.RouteMap {
		routeManager := routingTable.RouteMap[key]
		routes := routeManager.Routes
		for _, route := range routes {
			spacesIp :=  lengthIpSection - (len(route.Ip) + SPACE_FIRST)
			spacesPort := lengthPortSection - (len(route.Port) + SPACE_AFTER)
			spacesPath := lengthPathSection - (len(route.Path) + SPACE_AFTER)
			spacesVersion := lengthVersionSection - (len(route.Version) + SPACE_AFTER)
			fmt.Fprintf(
				os.Stdout,
				"%s%s | %s%s | %s%s | %s%s | %s\n",
				route.Ip,
				strings.Repeat(" ", spacesIp),
				route.Port,
				strings.Repeat(" ", spacesPort),
				route.Path,
				strings.Repeat(" ", spacesPath),
				route.Version,
				strings.Repeat(" ", spacesVersion),
				fmt.Sprintf("%d", route.Weight),
			);
		}
		fmt.Println(dashLine)
	}
}

func routingTableParseError(field string) {
	fmt.Fprintf(os.Stderr, "Field: [%s] has no value\n", field)
	os.Exit(1)
}

func lineToRoute(line string) Route {
	split := strings.Split(line, "|")

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
		if len(line) == 0 {
			continue
		}
		routes = append(routes, lineToRoute(line))
	}

	return createRoutingTable(routes)
}
