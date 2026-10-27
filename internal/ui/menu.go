package ui

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/ModstDev/VVT/internal/model"
	"github.com/ModstDev/VVT/internal/schedule"
	"github.com/ModstDev/VVT/internal/storage"
)

func Run(timetable *schedule.Schedule, filePath string) error {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Println("\n========== VIEŠOJO TRANSPORTO TVARKARAŠTIS ==========")
		fmt.Println("1. Peržiūrėti visus reisus")
		fmt.Println("2. Peržiūrėti pasirinkto reiso informaciją")
		fmt.Println("3. Ištrinti reisą")
		fmt.Println("4. Išsaugoti ir išeiti")
		fmt.Print("Pasirinkimas: ")

		if !scanner.Scan() {
			if err := scanner.Err(); err != nil {
				return err
			}

			return storage.Save(filePath, timetable)
		}

		switch strings.TrimSpace(scanner.Text()) {
		case "1":
			timetable.PrintTrips()
		case "2":
			viewTriP(scanner, timetable)
		case "3":
			deleteTrip(scanner, timetable)
		case "4":
			if err := storage.Save(filePath, timetable); err != nil {
				return err
			}
			fmt.Println("Tvarkaraštis išsaugotas. Programa baigta.")
			return nil
		default:
			fmt.Println("Neteisingas pasirinkimas. Bandykite dar kartą.")
		}
	}
}

func viewTriP(scanner *bufio.Scanner, timetable *schedule.Schedule) {
	fmt.Print("Įveskite reiso ID: ")

	if !scanner.Scan() {
		return
	}

	id := strings.TrimSpace(scanner.Text())

	trip, found := timetable.FindTrip(id)
	if !found {
		fmt.Printf("Reisas su ID %s nerastas.\n", id)
		return
	}

	fmt.Printf("\nReiso ID: %s\n", trip.ID)
	fmt.Printf("Maršruto numeris: %s\n", trip.RouteNumber)
	fmt.Printf("Kryptis: %s\n", trip.Direction)
	fmt.Println("Stotelės:")

	for element := trip.Stops.Front(); element != nil; element = element.Next() {
		stop := element.Value.(model.Stop)
		fmt.Printf("  %s — %s\n", stop.Name, stop.Time)
	}
}

func deleteTrip(scanner *bufio.Scanner, timetable *schedule.Schedule) {
	fmt.Print("Įveskite norimo ištrinti reiso ID: ")

	if !scanner.Scan() {
		return
	}

	id := strings.TrimSpace(scanner.Text())

	if timetable.RemoveTrip(id) {
		fmt.Println("Reisas sėkmingai ištrintas.")
	} else {
		fmt.Println("Reisas nerastas.")
	}
}
