# 🍽️ Homestyle Plate - Meal Randomizer

A self-contained desktop application for easy meal planning and recipe randomization.

![Homestyle Plate](Frontend/images/HomestylePlateLogo.png)

## Features

- **🎰 Meal Slot Machine** - Spin to randomly select protein, vegetable, and starch combinations
- **📋 Randomized Menu Planner** - Generate custom meal plans from HelloFresh and Paprika3 recipes
- **🍰 Desserts Showcase** - Browse delicious desserts with ingredients and bake times
- **🔒 Lock Favorites** - Lock meals you want and spin for new options

## How to Use

### Pre-built Releases (End Users)
Download the latest release for your platform:

| Platform | File | Size |
|----------|------|------|
| **macOS** (11.0+) | `Homestyle Plate_vX.X.X_macos.dmg` | ~3-5 MB |
| **Linux** (x86_64) | `Homestyle Plate_vX.X.X_linux-amd64.tar.gz` | ~3-5 MB |
| **Linux** (ARM64) | `Homestyle Plate_vX.X.X_linux-arm64.tar.gz` | ~3-5 MB |
| **Windows** (x64) | `Homestyle Plate_vX.X.X_windows-amd64.zip` | ~3-5 MB |

**macOS**: Drag the `.app` from the DMG to Applications, then double-click.
**Linux**: Extract and run `./launch.sh` (or `./HomestylePlate` directly).
**Windows**: Extract and run `launch.bat` (or `HomestylePlate.exe`).

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
- Download the latest release from the GitHub Releases tab
- Or build locally: `cd App && bash release.sh all`

## Building & Releasing

### Local Development
Build a single platform:
```bash
cd App
bash release.sh macos      # Build macOS DMG
bash release.sh linux-amd64 # Build Linux x86_64 tarball
bash release.sh windows-amd64 # Build Windows zip
```

Build all platforms:
```bash
bash release.sh all
# Outputs to: releases/ directory
```

### Automated Releases (GitHub)
Push a version tag to trigger automated builds:
```bash
# Update version
echo "1.0.1" > App/VERSION
git add App/VERSION
git commit -m "Bump version to 1.0.1"
git tag v1.0.1
git push origin main --tags
```

This triggers the CI/CD pipeline which:
- Builds for macOS, Linux (x86_64), Linux (ARM64), and Windows
- Creates a GitHub Release with all artifacts
- Auto-generates release notes

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

### Quick Start (from source)
```bash
cd App
go run server.go
```

### Build All Platforms (Release)
```bash
cd App
bash release.sh all
# Outputs to: releases/
# - Homestyle Plate_vX.X.X_macos.dmg
# - Homestyle Plate_vX.X.X_linux-amd64.tar.gz
# - Homestyle Plate_vX.X.X_windows-amd64.zip
```

### Build Single Platform
```bash
bash release.sh macos
bash release.sh linux-amd64
bash release.sh linux-arm64
bash release.sh windows-amd64
```

## License

Developed by Alex Sorichetti
