package storage

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/ModstDev/VVT/internal/model"
	"github.com/ModstDev/VVT/internal/schedule"
)

type StopData struct {
	Name string `json:"name"`
	Time string `json:"time"`
}

type TripData struct {
	ID          string     `json:"id"`
	RouteNumber string     `json:"route_number"`
	Direction   string     `json:"direction"`
	Stops       []StopData `json:"stops"`
}

type ScheduleData struct {
	Trips []TripData `json:"trips"`
}

func Load(path string) (*schedule.Schedule, error) {
	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return schedule.New(), nil
		}
		return nil, fmt.Errorf("open schedule file: %w", err)
	}
	defer file.Close()

	var data ScheduleData

	if err := json.NewDecoder(file).Decode(&data); err != nil {
		return nil, fmt.Errorf("decode schedule file: %w", err)
	}

	timetable := schedule.New()

	for _, tripData := range data.Trips {
		trip := model.NewTrip(
			tripData.ID,
			tripData.RouteNumber,
			tripData.Direction,
		)

		for _, stopData := range tripData.Stops {
			stop := model.Stop{
				Name: stopData.Name,
				Time: stopData.Time,
			}

			if err := trip.AddStop(stop); err != nil {
				return nil, fmt.Errorf("add stop to trip %q: %w", trip.ID, err)
			}
		}

		if err := timetable.AddTrip(trip); err != nil {
			return nil, fmt.Errorf("add trip %q to schedule: %w", trip.ID, err)
		}
	}

	return timetable, nil
}

func Save(path string, timetable *schedule.Schedule) error {
	if timetable == nil {
		return fmt.Errorf("schedule cannot be nil")
	}

	data := ScheduleData{
		Trips: make([]TripData, 0, timetable.Trips.Len()),
	}

	for element := timetable.Trips.Front(); element != nil; element = element.Next() {
		trip := element.Value.(*model.Trip)

		tripData := TripData{
			ID:          trip.ID,
			RouteNumber: trip.RouteNumber,
			Direction:   trip.Direction,
			Stops:       make([]StopData, 0, trip.Stops.Len()),
		}

		for stopElement := trip.Stops.Front(); stopElement != nil; stopElement = stopElement.Next() {
			stop := stopElement.Value.(model.Stop)

			stopData := StopData{
				Name: stop.Name,
				Time: stop.Time,
			}
			tripData.Stops = append(tripData.Stops, stopData)
		}

		data.Trips = append(data.Trips, tripData)
	}

	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create schedule file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(data); err != nil {
		file.Close()
		return fmt.Errorf("encode schedule JSON: %w", err)
	}

	if err := file.Close(); err != nil {
		return fmt.Errorf("close schedule file: %w", err)
	}

	return nil
}
