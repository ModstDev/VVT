package schedule

import (
	"bufio"
	"fmt"
	"strings"
	"time"

	"github.com/ModstDev/VVT/internal/model"
)

func (s *Schedule) EditTripStopTime(
	tripID string,
	stopIndex int,
	newTime string,
	ignoreDuplicate bool,
) (*model.Trip, error) {
	if _, err := time.Parse("15:04", newTime); err != nil {
		return nil, fmt.Errorf("invalid time %q: use 06:06 format", newTime)
	}

	trip, found := s.FindTrip(tripID)
	if !found {
		return nil, fmt.Errorf("trip with ID %q not found", tripID)
	}

	if stopIndex < 0 || stopIndex >= trip.Stops.Len() {
		return nil, fmt.Errorf("invalid stop index")
	}

	element := trip.Stops.Front()
	for i := 0; i < stopIndex; i++ {
		element = element.Next()
	}

	stop := element.Value.(model.Stop)

	oldTime := stop.Time

	stop.Time = newTime
	element.Value = stop

	if !ignoreDuplicate {
		duplicate := s.FindDuplicateTripExclude(trip, trip.ID)

		if duplicate != nil {
			stop.Time = oldTime
			element.Value = stop

			return duplicate, nil
		}
	}
	return nil, nil
}

func (s *Schedule) RenameRouteStop(
	routeNumber string,
	oldName string,
	newName string,
	scanner *bufio.Scanner,
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

	for tripElement := s.Trips.Front(); tripElement != nil; {
		nextTrip := tripElement.Next()
		trip := tripElement.Value.(*model.Trip)
		if trip.RouteNumber != routeNumber {
			tripElement = nextTrip
			continue
		}

		tripUpdated := false

		for stopElement := trip.Stops.Front(); stopElement != nil; stopElement = stopElement.Next() {
			stop := stopElement.Value.(model.Stop)

			if stop.Name == oldName {
				stop.Name = newName
				trip.Direction = strings.ReplaceAll(trip.Direction, oldName, newName)
				stopElement.Value = stop
				updated++
				tripUpdated = true
			}
		}
		if tripUpdated {
			duplicate := s.FindDuplicateTripExclude(trip, trip.ID)
			if duplicate != nil {
				fmt.Printf(
					"\nReisas %s tapo identiškas reisui %s.\n",
					trip.ID,
					duplicate.ID,
				)

				fmt.Println("1. Pašalinti redaguotą reisą")
				fmt.Println("2. Palikti pakeistą reisą")
				fmt.Print("Pasirinkimas: ")

				if !scanner.Scan() {
					return updated, fmt.Errorf("nepavyko nuskaityti įvesties")
				}

				switch strings.TrimSpace(scanner.Text()) {
				case "1":
					s.Trips.Remove(tripElement)
					updated--
					fmt.Println("Redaguotas reisas pašalintas.")
				case "2":
					fmt.Println("Reisas paliktas.")
				default:
					fmt.Println("Neteisingas pasirinkimas. Paliekamas pakeistas reisas.")
				}
			}
		}
		tripElement = nextTrip
	}

	if updated == 0 {
		return 0, fmt.Errorf("no stops with name %q found for route %q", oldName, routeNumber)
	}
	return updated, nil
}

func (s *Schedule) RemoveRouteStop(
	routeNumber string,
	stopName string,
	scanner *bufio.Scanner,
) (int, error) {
	stopName = strings.TrimSpace(stopName)
	if stopName == "" {
		return 0, fmt.Errorf("stop name cannot be empty")
	}

	removed := 0

	for tripElement := s.Trips.Front(); tripElement != nil; {
		nextTrip := tripElement.Next()
		trip := tripElement.Value.(*model.Trip)

		if trip.RouteNumber != routeNumber {
			tripElement = nextTrip
			continue
		}

		tripRemoved := 0

		for stopElement := trip.Stops.Front(); stopElement != nil; {
			next := stopElement.Next()
			stop := stopElement.Value.(model.Stop)

			if stop.Name == stopName {
				trip.Stops.Remove(stopElement)
				removed++
				tripRemoved++
			}

			stopElement = next
		}

		if tripRemoved > 0 {
			duplicate := s.FindDuplicateTripExclude(trip, trip.ID)
			if duplicate != nil {
				fmt.Printf(
					"\nReisas %s tapo identiškas reisui %s.\n",
					trip.ID,
					duplicate.ID,
				)

				fmt.Println("1. Pašalinti redaguotą reisą")
				fmt.Println("2. Palikti pakeistą reisą")
				fmt.Print("Pasirinkimas: ")

				if !scanner.Scan() {
					return removed, fmt.Errorf("nepavyko nuskaityti įvesties")
				}

				switch strings.TrimSpace(scanner.Text()) {
				case "1":
					s.Trips.Remove(tripElement)
					removed -= tripRemoved
					fmt.Println("Redaguotas reisas pašalintas.")
				case "2":
					fmt.Println("Reisas paliktas.")
				default:
					fmt.Println("Neteisingas pasirinkimas. Paliekamas pakeistas reisas.")
				}
			}
		}
		tripElement = nextTrip
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
