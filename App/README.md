# Homestyle Plate - Meal Randomizer

A self-contained desktop application for meal planning and recipe randomization.

## What This Is

Homestyle Plate is a meal randomizer designed to make dinner planning simple and fun. It features:

- **Meal Slot Machine** - Spin to get random protein, vegetable, and starch combinations
- **Randomized Menu Planner** - Generate HelloFresh and Paprika3 meal recommendations
- **Desserts Showcase** - Browse a collection of dessert recipes with ingredients and bake times

## How to Run

### macOS
1. Double-click `Homestyle Plate.app`
2. The app will open automatically in your browser at `http://localhost:8080`
3. That's it! No installation required.

### Linux
```bash
chmod +x launch.sh
./launch.sh
```

### Direct Execution
```bash
./HomestylePlate
```

Then open `http://localhost:8080` in your browser.

## How It Works

- **Single Binary** - Everything runs from one executable
- **SQLite Database** - Meals and desserts are stored locally in a SQLite database
- **No Docker Required** - Just run the binary
- **No External Dependencies** - Once built, it works offline (except for recipe images and links)
- **First Run** - Automatically seeds the database with recipes from the included JSON data

## For Developers

### Building

```bash
cd App
chmod +x build.sh
./build.sh
```

This creates:
- `build/Homestyle Plate.app/` on macOS
- `build/HomestylePlate` on Linux

### Running from Source

```bash
cd App
go run server.go
```

The server will start on `http://localhost:8080`

### Architecture

```
App/
├── server.go          # Main application (frontend + API)
├── go.mod             # Dependencies
├── build.sh           # Build script
├── launch.sh          # Launcher script
└── README.md          # This file
```

The server:
1. Seeds a SQLite database from JSON files on first run
2. Serves all static frontend files (HTML/CSS/JS)
3. Provides REST API endpoints for meals and desserts
4. Runs on port 8080 by default (configurable via `PORT` env var)

## Data

The application ships with pre-loaded recipes:
- 92+ HelloFresh meals with ratings
- 57+ Paprika3 meals with ratings  
- 39 desserts with ingredients and bake times

## License

Developed by Alex Sorichetti
