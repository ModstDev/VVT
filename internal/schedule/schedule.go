package schedule

import (
	"container/list"
	"fmt"
	"strings"

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

func (s *Schedule) FindDuplicateTrip(candidate *model.Trip) *model.Trip {
	return s.FindDuplicateTripExclude(candidate, "")
}

func (s *Schedule) FindDuplicateTripExclude(candidate *model.Trip, tripIDExclude string) *model.Trip {
	for element := s.Trips.Front(); element != nil; element = element.Next() {
		existing := element.Value.(*model.Trip)

		if existing.ID == tripIDExclude {
			continue
		}

		if existing.ID == candidate.ID {
			return existing
		}

		if existing.RouteNumber != candidate.RouteNumber ||
			existing.Direction != candidate.Direction {
			continue
		}

		if existing.Stops.Len() != candidate.Stops.Len() {
			continue
		}

		existingStop := existing.Stops.Front()
		candidateStop := candidate.Stops.Front()

		identical := true

		for existingStop != nil && candidateStop != nil {
			a := existingStop.Value.(model.Stop)
			b := candidateStop.Value.(model.Stop)

			if !strings.EqualFold(a.Name, b.Name) || a.Time != b.Time {
				identical = false
				break
			}

			existingStop = existingStop.Next()
			candidateStop = candidateStop.Next()
		}

		if identical {
			return existing
		}
	}

	return nil
}
