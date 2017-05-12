package routing

import "sync"

type Scheduler struct {
	Index          int
	Counter        int
	NumberRoutes   int
	Mutex          *sync.Mutex
}

func CreateScheduler(numberRoutes int) Scheduler {
	return Scheduler{0, 0, numberRoutes, new(sync.Mutex)}
}

func NextRoute(scheduler *Scheduler, routes []Route) Route {
	scheduler.Mutex.Lock()
	defer scheduler.Mutex.Unlock()

	currentRoute := routes[scheduler.Index]

	if scheduler.Counter < currentRoute.Weight {
		scheduler.Counter += 1
		return currentRoute
	} else {
		scheduler.Counter = 1
		scheduler.Index = (scheduler.Index + 1) % scheduler.NumberRoutes
		return routes[scheduler.Index]
	}
}





