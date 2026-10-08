package cmd

import (
	"fmt"
	"os"

	"github.com/asorichetti/HomestylePlate/tools/interact-test/internal/browser"
	"github.com/asorichetti/HomestylePlate/tools/interact-test/internal/journeys"
	"github.com/spf13/cobra"
)

var (
	baseURL      string
	env          string
	headless     bool
	screensDir   string
	uxGates      string
	viewport     string
	colorScheme  string
	listJourneys bool
	grepFilter   string
	suiteFilter  string
)

// rootCmd is the root command
var rootCmd = &cobra.Command{
	Use:   "interact-test",
	Short: "Homestyle Plate UI interaction tester",
	Long: `A browser-automation CLI that captures screenshots, DOM, console logs,
and network failures to prove UI changes actually work.

Each command writes a debug bundle with:
- Screenshot (PNG)
- Page HTML
- Console/Network logs
- UX Gate results`,
}

// Execute runs the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&baseURL, "base-url", "", "Base URL (overrides --env)")
	rootCmd.PersistentFlags().StringVar(&env, "env", "local", "Environment (local, preview)")
	rootCmd.PersistentFlags().BoolVar(&headless, "headless", browser.ShouldRunHeadless(), "Run in headless mode")
	rootCmd.PersistentFlags().StringVar(&screensDir, "screenshot-dir", "./screenshots", "Directory for screenshots and bundles")
	rootCmd.PersistentFlags().StringVar(&uxGates, "ux-gates", "warn", "UX gates mode: off, warn, fail")
	rootCmd.PersistentFlags().StringVar(&viewport, "viewport", "desktop", "Viewport: desktop, tablet, mobile")
	rootCmd.PersistentFlags().StringVar(&colorScheme, "color-scheme", "", "Color scheme: light, dark")
}

// getBaseURL resolves the base URL from env or flag
func getBaseURL() string {
	if baseURL != "" {
		return baseURL
	}
	
	// Map environments to URLs
	envURLs := map[string]string{
		"local":  "http://localhost:8080",
		"preview": "https://homestyleplate-preview.example.com",
	}
	
	if url, ok := envURLs[env]; ok {
		return url
	}
	
	return "http://localhost:8080"
}

// setupBrowser creates a browser instance
func setupBrowser() (*browser.Browser, error) {
	b, err := browser.New(headless)
	if err != nil {
		return nil, fmt.Errorf("failed to create browser: %w", err)
	}
	return b, nil
}

// listJourneysCmd prints all registered journeys
var listJourneysCmd = &cobra.Command{
	Use:   "list-journeys",
	Short: "List all registered journeys",
	Run: func(cmd *cobra.Command, args []string) {
		all := journeys.All()
		fmt.Printf("Registered journeys (%d):\n", len(all))
		for _, j := range all {
			fmt.Printf("  %-20s [%-10s] %s\n", j.Name, j.Suite, j.Desc)
		}
	},
}

func init() {
	rootCmd.AddCommand(listJourneysCmd)
}
