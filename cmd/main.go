package main

import (
	"fmt"
	"log"

	"github.com/ModstDev/VVT/internal/storage"
)

func main() {
	const filePath = "data/schedule.json"

	timetable, err := storage.Load(filePath)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("=== VIEŠOJO TRANSPORTO TVARKARAŠTIS ===")
	timetable.PrintTrips()

	if err := storage.Save(filePath, timetable); err != nil {
		log.Fatal(err)
	}

	fmt.Println("\nTvarkaraštis išsaugotas.")
}
