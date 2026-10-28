package schedule

import (
	"fmt"
	"strings"
	"time"

	"github.com/ModstDev/VVT/internal/model"
)

func (s *Schedule) EditTripStopTime(
	tripID string,
	stopIndex int,
	newTime string,
) error {
	if _, err := time.Parse("15:04", newTime); err != nil {
		return fmt.Errorf("invalid time %q: use 06:06 format", newTime)
	}

	trip, found := s.FindTrip(tripID)
	if !found {
		return fmt.Errorf("trip with ID %q not found", tripID)
	}

	if stopIndex < 0 || stopIndex >= trip.Stops.Len() {
		return fmt.Errorf("invalid stop index")
	}

	element := trip.Stops.Front()
	for i := 0; i < stopIndex; i++ {
		element = element.Next()
	}

	stop := element.Value.(model.Stop)
	stop.Time = newTime
	element.Value = stop

	return nil
}

func (s *Schedule) RenameRouteStop(
	routeNumber string,
	oldName string,
	newName string,
) (int, error) {
	oldName = strings.TrimSpace(oldName)
	newName = strings.TrimSpace(newName)

	if oldName == " " || newName == "" {
		return 0, fmt.Errorf("stop names cannot be empty")
	}

	if oldName == newName {
		return 0, fmt.Errorf("new stop name is the same as the old name")
	}

	updated := 0

	for tripElement := s.Trips.Front(); tripElement != nil; tripElement = tripElement.Next() {
		trip := tripElement.Value.(*model.Trip)
		if trip.RouteNumber != routeNumber {
			continue
		}

		for stopElement := trip.Stops.Front(); stopElement != nil; stopElement = stopElement.Next() {
			stop := stopElement.Value.(model.Stop)

			if stop.Name == oldName {
				stop.Name = newName
				stopElement.Value = stop
				updated++
			}
		}
	}

	if updated == 0 {
		return 0, fmt.Errorf("no stops with name %q found for route %q", oldName, routeNumber)
	}
	return updated, nil
}

func (s *Schedule) RemoveRouteStop(
	routeNumber string,
	stopName string,
) (int, error) {
	stopName = strings.TrimSpace(stopName)
	if stopName == "" {
		return 0, fmt.Errorf("stop name cannot be empty")
	}

	removed := 0

	for tripElement := s.Trips.Front(); tripElement != nil; tripElement = tripElement.Next() {
		trip := tripElement.Value.(*model.Trip)

		if trip.RouteNumber != routeNumber {
			continue
		}

		for stopElement := trip.Stops.Front(); stopElement != nil; {
			next := stopElement.Next()
			stop := stopElement.Value.(model.Stop)

			if stop.Name == stopName {
				trip.Stops.Remove(stopElement)
				removed++
			}

			stopElement = next
		}
	}

	if removed == 0 {
		return 0, fmt.Errorf(
			"stop %q was not found on route %q",
			stopName,
			routeNumber,
		)
	}

	return removed, nil
}
