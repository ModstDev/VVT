package schedule

import (
	"container/list"
	"fmt"

	"github.com/ModstDev/VVT/internal/model"
)

type Schedule struct {
	Trips *list.List
}

func New() *Schedule {
	return &Schedule{
		Trips: list.New(),
	}
}

func (s *Schedule) AddTrip(trip *model.Trip) error {
	if trip == nil {
		return fmt.Errorf("trip cannot be nil")

	}
	if trip.ID == "" {
		return fmt.Errorf("trip ID cannot be empty")
	}

	for element := s.Trips.Front(); element != nil; element = element.Next() {
		existingTrip, ok := element.Value.(*model.Trip)
		if !ok {
			return fmt.Errorf("invalid trip in schedule")
		}
		if existingTrip.ID == trip.ID {
			return fmt.Errorf("trip with ID %s already exists in schedule", trip.ID)
		}
	}

	s.Trips.PushBack(trip)

	return nil
}

func (s *Schedule) FindTrip(id string) (*model.Trip, bool) {
	for element := s.Trips.Front(); element != nil; element = element.Next() {
		trip := element.Value.(*model.Trip)
		if trip.ID == id {
			return trip, true
		}
	}
	return nil, false
}

func (s *Schedule) RemoveTrip(id string) bool {
	for element := s.Trips.Front(); element != nil; element = element.Next() {
		trip := element.Value.(*model.Trip)
		if trip.ID == id {
			s.Trips.Remove(element)
			return true
		}
	}
	return false
}

func (s *Schedule) PrintTrips() {
	if s.Trips.Len() == 0 {
		fmt.Println("No trips in the schedule.")
		return
	}

	for element := s.Trips.Front(); element != nil; element = element.Next() {
		trip := element.Value.(*model.Trip)
		fmt.Printf(
			"ID: %s | Route: %s | Direction: %s\n",
			trip.ID,
			trip.RouteNumber,
			trip.Direction,
		)
	}
}
