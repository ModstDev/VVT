package schedule

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/ModstDev/VVT/internal/model"
)

type DirectTrip struct {
	TripID        string
	RouteNumber   string
	Direction     string
	FromStop      string
	ToStop        string
	DepartureTime string
	ArrivalTime   string
	Duration      int
	ArrivalOffset int
}

func (s *Schedule) FindFastestTrips(
	firstStop,
	secondStop string,
	departure time.Time,
	limit int,
) ([]DirectTrip, error) {
	firstStop = strings.TrimSpace(firstStop)
	secondStop = strings.TrimSpace(secondStop)

	if firstStop == "" || secondStop == "" {
		return nil, fmt.Errorf("stop names cannot be empty")
	}

	if strings.EqualFold(firstStop, secondStop) {
		return nil, fmt.Errorf("departure and destination stops must be different")
	}

	if limit < 0 {
		return nil, fmt.Errorf("limit cannot be negative")
	}

	const minutesPerDay = 24 * 60
	requestedMinute := departure.Hour()*60 + departure.Minute()

	results := make([]DirectTrip, 0)

	for tripElement := s.Trips.Front(); tripElement != nil; tripElement = tripElement.Next() {
		trip := tripElement.Value.(*model.Trip)

		stops := make([]model.Stop, 0, trip.Stops.Len())

		for element := trip.Stops.Front(); element != nil; element = element.Next() {
			stop := element.Value.(model.Stop)
			stops = append(stops, stop)
		}

		if len(stops) < 2 {
			continue
		}

		offsets := make([]int, len(stops))
		previousMinute := -1
		dayOffset := 0

		for i, stop := range stops {
			parsed, err := time.Parse("15:04", stop.Time)
			if err != nil {
				return nil, fmt.Errorf(
					"invalid time %q in trip %q: %w",
					stop.Time,
					trip.ID,
					err,
				)
			}

			minute := parsed.Hour()*60 + parsed.Minute()

			if i > 0 && minute < previousMinute {
				dayOffset += minutesPerDay
			}

			offsets[i] = dayOffset + minute
			previousMinute = minute
		}

		var best *DirectTrip

		for i := 0; i < len(stops); i++ {
			if !strings.EqualFold(stops[i].Name, firstStop) {
				continue
			}

			for j := i + 1; j < len(stops); j++ {
				if !strings.EqualFold(stops[j].Name, secondStop) {
					continue
				}

				departureClock := offsets[i] % minutesPerDay

				wait := (departureClock - requestedMinute + minutesPerDay) % minutesPerDay

				actualDeparture := requestedMinute + wait

				duration := offsets[j] - offsets[i]

				actualArrival := actualDeparture + duration

				candidate := DirectTrip{
					TripID:        trip.ID,
					RouteNumber:   trip.RouteNumber,
					Direction:     trip.Direction,
					FromStop:      stops[i].Name,
					ToStop:        stops[j].Name,
					DepartureTime: formatMinutes(actualDeparture),
					ArrivalTime:   formatMinutes(actualArrival),
					Duration:      duration,
					ArrivalOffset: actualArrival,
				}

				if best == nil || candidate.ArrivalOffset < best.ArrivalOffset {
					copy := candidate
					best = &copy
				}
			}
		}

		if best != nil {
			results = append(results, *best)
		}
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].ArrivalOffset < results[j].ArrivalOffset
	})

	if limit > 0 && len(results) > limit {
		results = results[:limit]
	}

	return results, nil
}

func formatMinutes(minutes int) string {
	day := minutes / (24 * 60)
	minuteOfDay := minutes % (24 * 60)

	hours := minuteOfDay / 60
	minutesPart := minuteOfDay % 60

	if day > 0 {
		return fmt.Sprintf("%02d:%02d (+%d d)",
			hours,
			minutesPart,
			day,
		)
	}
	return fmt.Sprintf("%02d:%02d", hours, minutesPart)
}
