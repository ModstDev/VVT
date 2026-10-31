package ui

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/ModstDev/VVT/internal/model"
	"github.com/ModstDev/VVT/internal/schedule"
	"github.com/ModstDev/VVT/internal/storage"
)

func clearCLI() {
	var cmd *exec.Cmd

	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	} else {
		cmd = exec.Command("clear")
	}

	cmd.Stdout = os.Stdout
	_ = cmd.Run()
}

func Run(timetable *schedule.Schedule, filePath string) error {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Println("\n========== VIEŠOJO TRANSPORTO TVARKARAŠTIS ==========")
		fmt.Println("1. Peržiūrėti visus reisus")
		fmt.Println("2. Peržiūrėti pasirinkto reiso informaciją")
		fmt.Println("3. Ištrinti reisą")
		fmt.Println("4. Redaguoti reiso stotelės laiką")
		fmt.Println("5. Pakeisti maršruto stotelės pavadinimą")
		fmt.Println("6. Ištrinti maršruto stotelę")
		fmt.Println("7. Peržiūrėti stotelės tvarkaraštį")
		fmt.Println("8. Greičiausio tiesioginio reiso paieška")
		fmt.Println("9. Išsaugoti ir išeiti")
		fmt.Print("Pasirinkimas: ")

		if !scanner.Scan() {
			if err := scanner.Err(); err != nil {
				return err
			}
			return storage.Save(filePath, timetable)
		}

		switch strings.TrimSpace(scanner.Text()) {
		case "1":
			clearCLI()
			timetable.PrintTrips()
		case "2":
			clearCLI()
			viewTrip(scanner, timetable)
		case "3":
			clearCLI()
			deleteTrip(scanner, timetable)
		case "4":
			clearCLI()
			editTripStopTime(scanner, timetable)
		case "5":
			clearCLI()
			renameRouteStop(scanner, timetable)
		case "6":
			clearCLI()
			removeRouteStop(scanner, timetable)
		case "7":
			clearCLI()
			showStopSchedule(scanner, timetable)
		case "8":
			clearCLI()
			findFastestTrips(scanner, timetable)
		case "9":
			clearCLI()
			if err := storage.Save(filePath, timetable); err != nil {
				return err
			}
			fmt.Println("Tvarkaraštis išsaugotas. Programa baigta.")
			return nil
		default:
			clearCLI()
			fmt.Println("Neteisingas pasirinkimas. Bandykite dar kartą.")
		}
	}
}

func viewTrip(scanner *bufio.Scanner, timetable *schedule.Schedule) {
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
		fmt.Printf("  %s - %s\n", stop.Name, stop.Time)
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

func editTripStopTime(scanner *bufio.Scanner, timetable *schedule.Schedule) {
	fmt.Print("Reiso ID: ")
	if !scanner.Scan() {
		return
	}
	id := strings.TrimSpace(scanner.Text())

	trip, found := timetable.FindTrip(id)
	if !found {
		fmt.Println("Reisas nerastas.")
		return
	}

	fmt.Println("Reiso stotelės:")
	for element, index := trip.Stops.Front(), 1; element != nil; element, index = element.Next(), index+1 {
		stop := element.Value.(model.Stop)
		fmt.Printf("%d. %s - %s\n", index, stop.Name, stop.Time)
	}

	fmt.Print("Stotelės numeris: ")
	if !scanner.Scan() {
		return
	}

	index, err := strconv.Atoi(strings.TrimSpace(scanner.Text()))
	if err != nil || index < 1 || index > trip.Stops.Len() {
		fmt.Println("Neteisingas stotelės numeris.")
		return
	}

	fmt.Print("Naujas laikas (HH:MM): ")
	if !scanner.Scan() {
		return
	}
	newTime := strings.TrimSpace(scanner.Text())

	if err := timetable.EditTripStopTime(id, index-1, newTime); err != nil {
		fmt.Println("Klaida:", err)
		return
	}

	fmt.Println("Reiso laikas pakeistas.")
}

func renameRouteStop(scanner *bufio.Scanner, timetable *schedule.Schedule) {
	fmt.Print("Maršruto numeris: ")
	if !scanner.Scan() {
		return
	}
	route := strings.TrimSpace(scanner.Text())

	fmt.Print("Dabartinis stotelės pavadinimas: ")
	if !scanner.Scan() {
		return
	}
	oldName := strings.TrimSpace(scanner.Text())

	fmt.Print("Naujas stotelės pavadinimas: ")
	if !scanner.Scan() {
		return
	}
	newName := strings.TrimSpace(scanner.Text())

	count, err := timetable.RenameRouteStop(route, oldName, newName)
	if err != nil {
		fmt.Println("Klaida:", err)
		return
	}

	fmt.Printf("Pakeista stotelių: %d\n", count)
}

func removeRouteStop(scanner *bufio.Scanner, timetable *schedule.Schedule) {
	fmt.Print("Maršruto numeris: ")
	if !scanner.Scan() {
		return
	}
	route := strings.TrimSpace(scanner.Text())

	fmt.Print("Stotelės pavadinimas, kurį norite pašalinti: ")
	if !scanner.Scan() {
		return
	}
	stopName := strings.TrimSpace(scanner.Text())

	count, err := timetable.RemoveRouteStop(route, stopName)
	if err != nil {
		fmt.Println("Klaida:", err)
		return
	}

	fmt.Printf("Pašalinta stotelių: %d\n", count)
}

func showStopSchedule(scanner *bufio.Scanner, timetable *schedule.Schedule) {
	fmt.Print("Stotelės pavadinimas: ")
	if !scanner.Scan() {
		return
	}
	stopName := strings.TrimSpace(scanner.Text())

	now := time.Now()

	arrivals, err := timetable.UpcomingArrivals(stopName, now, 5)
	if err != nil {
		fmt.Println("Klaida: ", err)
		return
	}

	fmt.Printf("\nStotelė: %s\n", stopName)
	fmt.Printf("Dabartinis laikas: %s\n", now.Format("15:04"))
	fmt.Println("------------------------------------------------")

	if len(arrivals) == 0 {
		fmt.Println("Šioje stotelėje reisų nerasta.")
		return
	}

	fmt.Printf("%-12s %-12s %-8s %s\n",
		"Reisas", "Maršrutas", "Laikas", "Liko")

	for _, arrival := range arrivals {
		fmt.Printf("%-12s %-12s %-8s %d min.\n",
			arrival.TripID,
			arrival.RouteNumber,
			arrival.ArrivalTime,
			arrival.MinutesUntil,
		)
	}
}

func findFastestTrips(scanner *bufio.Scanner, timetable *schedule.Schedule) {
	fmt.Print("Pradinė stotelė: ")
	if !scanner.Scan() {
		return
	}
	from := strings.TrimSpace(scanner.Text())

	fmt.Print("Galutinė stotelė: ")
	if !scanner.Scan() {
		return
	}
	to := strings.TrimSpace(scanner.Text())

	fmt.Println("\nIšvykimo laikas:")
	fmt.Println("1. Dabartinis kompiuterio laikas")
	fmt.Println("2. Vartotojo nurodytas laikas")
	fmt.Print("Pasirinkimas: ")

	if !scanner.Scan() {
		return
	}

	var departure time.Time

	switch strings.TrimSpace(scanner.Text()) {
	case "1":
		departure = time.Now()

	case "2":
		fmt.Print("Įveskite išvykimo laiką (HH:MM): ")

		if !scanner.Scan() {
			return
		}

		input := strings.TrimSpace(scanner.Text())

		parsed, err := time.Parse("15:04", input)
		if err != nil {
			fmt.Println("Neteisingas laiko formatas. Naudokite HH:MM")
			return
		}

		departure = parsed

	default:
		fmt.Println("Neteisingas pasirinkimas.")
		return
	}

	results, err := timetable.FindFastestTrips(
		from,
		to,
		departure,
		3,
	)
	if err != nil {
		fmt.Println("Klaida: ", err)
		return
	}

	fmt.Printf("\nGreičiausi reisai: %s -> %s\n", from, to)
	fmt.Printf("Išvykimo paieškos laikas: %s\n", departure.Format("15:04"))
	fmt.Println("======================================================")

	for i, result := range results {
		fmt.Printf("\n%d. Maršrutas %s (reisas %s)\n",
			i+1,
			result.RouteNumber,
			result.TripID,
		)

		fmt.Printf("   Kryptis: %s\n", result.Direction)
		fmt.Println("   Kelias:")

		for tripElement := timetable.Trips.Front(); tripElement != nil; tripElement = tripElement.Next() {
			trip := tripElement.Value.(*model.Trip)

			if trip.ID == result.TripID {
				printPath := false

				for stopElement := trip.Stops.Front(); stopElement != nil; stopElement = stopElement.Next() {
					stop := stopElement.Value.(model.Stop)

					if strings.EqualFold(stop.Name, result.FromStop) {
						printPath = true
					}

					if printPath {
						fmt.Printf("      %s - %s\n", stop.Name, stop.Time)
					}

					if strings.EqualFold(stop.Name, result.ToStop) && printPath {
						break
					}
				}

				break
			}
		}
		fmt.Printf("   Išvykimas: %s\n", result.DepartureTime)
		fmt.Printf("   Atvykimas: %s\n", result.ArrivalTime)
		fmt.Printf("   Trukmė: %d min.\n", result.Duration)
	}
}
