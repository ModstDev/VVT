package schedule

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/ModstDev/VVT/internal/model"
)

type Arrival struct {
	TripID       string
	RouteNumber  string
	Direction    string
	ArrivalTime  string
	MinutesUntil int
}

func (s *Schedule) UpcomingArrivals(
	stopName string,
	now time.Time,
	limit int,
) ([]Arrival, error) {
	stopName = strings.TrimSpace(stopName)
	if stopName == "" {
		return nil, fmt.Errorf("stop name cannot be empty")
	}

	if limit < 0 {
		return nil, fmt.Errorf("limit cannot be negative")
	}

	const minutesPerDay = 24 * 60
	currentMinute := now.Hour()*60 + now.Minute()

	arrivals := make([]Arrival, 0)

	for tripElement := s.Trips.Front(); tripElement != nil; tripElement = tripElement.Next() {
		trip := tripElement.Value.(*model.Trip)

		for stopElement := trip.Stops.Front(); stopElement != nil; stopElement = stopElement.Next() {
			stop := stopElement.Value.(model.Stop)

			if !strings.EqualFold(strings.TrimSpace(stop.Name), stopName) {
				continue
			}

			parsedTime, err := time.Parse("15:04", stop.Time)
			if err != nil {
				return nil, fmt.Errorf(
					"invalid time %q in trip %q: %w",
					stop.Time,
					trip.ID,
					err,
				)
			}

			arrivalMinute := parsedTime.Hour()*60 + parsedTime.Minute()
			minutesUntil := (arrivalMinute - currentMinute + minutesPerDay) % minutesPerDay

			arrivals = append(arrivals, Arrival{
				TripID:       trip.ID,
				RouteNumber:  trip.RouteNumber,
				Direction:    trip.Direction,
				ArrivalTime:  stop.Time,
				MinutesUntil: minutesUntil,
			})

			break
		}
	}

	sort.Slice(arrivals, func(i, j int) bool {
		return arrivals[i].MinutesUntil < arrivals[j].MinutesUntil
	})

	if limit > 0 && len(arrivals) > limit {
		arrivals = arrivals[:limit]
	}

	return arrivals, nil
}
