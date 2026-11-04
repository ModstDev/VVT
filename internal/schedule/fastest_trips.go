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
	firstStop string,
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

		for fromElement := trip.Stops.Front(); fromElement != nil; fromElement = fromElement.Next() {
			from := fromElement.Value.(model.Stop)

			if !strings.EqualFold(from.Name, firstStop) {
				continue
			}

			fromTime, err := time.Parse("15:04", from.Time)
			if err != nil {
				return nil, fmt.Errorf("invalid time %q in trip %q", from.Time, trip.ID)
			}

			fromMinute := fromTime.Hour()*60 + fromTime.Minute()

			wait := (fromMinute - requestedMinute + minutesPerDay) % minutesPerDay
			actualDeparture := requestedMinute + wait

			previousMinute := fromMinute
			dayOffset := 0

			for toElement := fromElement.Next(); toElement != nil; toElement = toElement.Next() {
				to := toElement.Value.(model.Stop)

				toTime, err := time.Parse("15:04", to.Time)
				if err != nil {
					return nil, fmt.Errorf("invalid time %q in trip %q", to.Time, trip.ID)
				}

				toMinute := toTime.Hour()*60 + toTime.Minute()

				if toMinute < previousMinute {
					dayOffset += minutesPerDay
				}

				if strings.EqualFold(to.Name, secondStop) {
					duration := dayOffset + toMinute - fromMinute
					actualArrival := actualDeparture + duration

					candidate := DirectTrip{
						TripID:        trip.ID,
						RouteNumber:   trip.RouteNumber,
						Direction:     trip.Direction,
						FromStop:      from.Name,
						ToStop:        to.Name,
						DepartureTime: formatMinutes(actualDeparture),
						ArrivalTime:   formatMinutes(actualArrival),
						Duration:      duration,
						ArrivalOffset: actualArrival,
					}

					results = append(results, candidate)
				}

				previousMinute = toMinute
			}
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
		return fmt.Sprintf("%02d:%02d (+%d d (teks palaukti kitos dienos))",
			hours,
			minutesPart,
			day,
		)
	}
	return fmt.Sprintf("%02d:%02d", hours, minutesPart)
}
