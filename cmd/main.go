package main

import (
	"log"

	"github.com/ModstDev/VVT/internal/storage"
	"github.com/ModstDev/VVT/internal/ui"
)

func main() {
	const filePath = "data/schedule.json"

	timetable, err := storage.Load(filePath)
	if err != nil {
		log.Fatal(err)
	}

	if err := ui.Run(timetable, filePath); err != nil {
		log.Fatal(err)
	}
}
