package ui

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"

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
		fmt.Println("7. Išsaugoti ir išeiti")
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
