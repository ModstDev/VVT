package model

import (
	"container/list"
	"fmt"
)

type Trip struct {
	ID          string     `json:"id"`
	RouteNumber string     `json:"route_number"`
	Direction   string     `json:"direction"`
	Stops       *list.List `json:"-"`
}

func NewTrip(id, routeNumber, direction string) *Trip {
	return &Trip{
		ID:          id,
		RouteNumber: routeNumber,
		Direction:   direction,
		Stops:       list.New(),
	}
}

func (t *Trip) AddStop(stop Stop) error {
	if t.Stops == nil {
		t.Stops = list.New()
	}

	if stop.Name == "" {
		return fmt.Errorf("stop name cannot be empty")
	}

	t.Stops.PushBack(stop)
	return nil
}

func (t *Trip) PrintStops() {
	for element := t.Stops.Front(); element != nil; element = element.Next() {
		stop := element.Value.(Stop)
		fmt.Printf("%s - %s\n", stop.Name, stop.Time)
	}
}
