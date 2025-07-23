# Weather CLI

A simple command-line weather application written in Go.  
It automatically detects your location or accepts a city name, then fetches and displays the weather forecast from [wttr.in](https://wttr.in).

## Features
- Automatic location detection via IP
- Manual city input
- Configurable output format (1–4)
- Clean, readable output with support for ANSI/weather symbols
- Unit-tested code

#### 1. Run with default settings (detect location automatically)
`go run main.go`

#### 2. Specify a city
`go run main.go -city <your_city>`

Example: `go run main.go -city London`

#### 3. Change output format
Format 1: Brief summary (default)
`go run main.go -format 1`

Format 2:
`go run main.go -format 2`

Format 3: 
`go run main.go -format 3`

Format 4: 
`go run main.go -format 4`

#### 4. Combine city and format
You can use both flags together:
`go run main.go -city London -format 4`