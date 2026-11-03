# Public Transport Schedule

## Requirements

- Go installed.

## How to run

Open the project directory in a terminal and run:

```powershell
go run ./cmd
```

## Build executable

To compile the application into a Windows executable:

```powershell
go build -o vvt.exe ./cmd
```

This creates `vvt.exe` in the project's root directory.

To run the executable:

```powershell
.\vvt.exe
```

The timetable is loaded automatically from the JSON file at startup and saved when the application exits.
