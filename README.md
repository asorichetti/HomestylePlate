# 🍽️ Homestyle Plate - Meal Randomizer

A self-contained desktop application for easy meal planning and recipe randomization.

![Homestyle Plate](Frontend/images/HomestylePlateLogo.png)

## Features

- **🎰 Meal Slot Machine** - Spin to randomly select protein, vegetable, and starch combinations
- **📋 Randomized Menu Planner** - Generate custom meal plans from HelloFresh and Paprika3 recipes
- **🍰 Desserts Showcase** - Browse delicious desserts with ingredients and bake times
- **🔒 Lock Favorites** - Lock meals you want and spin for new options

## How to Use

### First Time Setup (Developer)
1. Clone the repository
2. Go to the `App/` directory
3. Run the build script: `./build.sh`
4. This creates a macOS app bundle in `App/build/`

### Running as a Desktop App
1. Build the app: `cd App && ./build.sh`
2. Open `App/build/Homestyle Plate.app`
3. The app opens automatically in your browser at `http://localhost:8080`

### For End Users
- Distribute the `.dmg` file found in `App/build/`
- Users just drag the app to Applications and double-click

## Architecture

### Self-Contained Desktop App
- **Single Go binary** - No Docker, no external dependencies
- **SQLite database** - All data stored locally
- **Auto-seeding** - Recipes loaded from bundled JSON on first run
- **One-click launch** - Double-click to start

### Tech Stack
| Component | Technology |
|-----------|------------|
| Frontend | HTML, CSS, JavaScript |
| Backend | Go (single binary) |
| Database | SQLite (embedded) |
| Deployment | macOS .app bundle / Linux binary |

### Project Structure
```
MealRandomizer/
├── App/                          # Desktop application (NEW)
│   ├── server.go                 # Main Go server (frontend + API)
│   ├── build.sh                  # Build script
│   ├── launch.sh                 # Launcher script
│   ├── README.md                 # App-specific docs
│   └── go.mod
├── Frontend/                     # Web interface
│   ├── LandingPage/
│   ├── MealSlotMachine/
│   ├── DessertsPage/
│   ├── MealsGalleryPage/
│   ├── MealsListed/
│   ├── GlobalScript/
│   ├── data/
│   └── images/
├── Backend/                      # Legacy backend code
│   ├── GoCode/
│   ├── PythonCode/
│   └── initdb/
├── nginx/
├── SQL_Files/
└── docker-compose.yml            # Legacy Docker setup
```

## API Endpoints

The app runs on `http://localhost:8080`

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/` | GET | Main application (spa routing) |
| `/meals?type=HF_Meal:4,P3_Meal:4` | GET | Get random meals |
| `/api/desserts` | GET | Get all desserts |
| `/api/admin/meals?action=all` | GET | List all meals (with IDs) |
| `/api/admin/meals?action=add` | POST | Add a new meal |
| `/api/admin/meals?action=delete` | POST | Delete a meal |

### Meals API Example
```bash
curl "http://localhost:8080/meals?type=HF_Meal:2,P3_Meal:2"
```

### Admin API Example
```bash
# Add a meal
curl -X POST "http://localhost:8080/api/admin/meals?action=add" \
  -H "Content-Type: application/json" \
  -d '{"table":"P3_Meal","name":"My Meal","rating":5,"recipe_link":"url","photo_link":"url"}'

# Delete a meal
curl -X POST "http://localhost:8080/api/admin/meals?action=delete" \
  -H "Content-Type: application/json" \
  -d '{"table":"P3_Meal","id":1}'
```

## Data

The application ships with pre-loaded recipes:
- **92+ HelloFresh meals** with ratings (1-7)
- **57+ Paprika3 meals** with ratings
- **39 desserts** with ingredients and bake times

Ratings affect the probability of selection (higher rating = more likely to appear).

## Building

### macOS
```bash
cd App
./build.sh
# Creates: build/Homestyle Plate.app
# Create DMG: cd build && hdiutil create -volname "Homestyle Plate" -srcfolder "Homestyle Plate.app" -ov -format UDZO "Homestyle Plate.dmg"
```

### Linux
```bash
cd App
./build.sh
# Creates: build/HomestylePlate (binary)
```

### From Source
```bash
cd App
go run server.go
```

## License

Developed by Alex Sorichetti
