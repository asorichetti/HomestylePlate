package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	playwright "github.com/mxschmitt/playwright-go"
	"github.com/asorichetti/HomestylePlate/tools/interact-test/internal/browser"
	"github.com/asorichetti/HomestylePlate/tools/interact-test/internal/journeys"
	"github.com/asorichetti/HomestylePlate/tools/interact-test/internal/uxgates"
	"github.com/spf13/cobra"
)

var journeyName string
var journeySuite string
var journeyList bool
var journeyGrep string

// journeyCmd runs test journeys
var journeyCmd = &cobra.Command{
	Use:   "journey [name]",
	Short: "Run UI test journeys",
	Run: func(cmd *cobra.Command, args []string) {
		// List mode
		if journeyList {
			all := journeys.All()
			fmt.Printf("Registered journeys (%d):\n", len(all))
			for _, j := range all {
				fmt.Printf("  %-25s [%-10s] %s\n", j.Name, j.Suite, j.Desc)
			}
			return
		}
		
		// Filter journeys
		var js []journeys.Journey
		if journeyName != "" {
			for _, j := range journeys.All() {
				if j.Name == journeyName {
					js = append(js, j)
					break
				}
			}
			if len(js) == 0 {
				fmt.Fprintf(os.Stderr, "❌ Journey not found: %s\n", journeyName)
				fmt.Fprintf(os.Stderr, "Available: ")
				for _, j := range journeys.All() {
					fmt.Fprintf(os.Stderr, "%s ", j.Name)
				}
				fmt.Fprintf(os.Stderr, "\n")
				os.Exit(1)
			}
		} else if journeySuite != "" {
			js = journeys.BySuite(journeySuite)
		} else if journeyGrep != "" {
			js = journeys.ByGrep(journeyGrep)
		} else {
			js = journeys.All()
		}
		
		if len(js) == 0 {
			fmt.Fprintf(os.Stderr, "❌ No journeys match the filter\n")
			os.Exit(1)
		}
		
		fmt.Fprintf(os.Stderr, "🚀 Running %d journey(s)...\n", len(js))
		
		base := getBaseURL()
		allPassed := true
		
		for _, j := range js {
			fmt.Fprintf(os.Stderr, "\n📍 Journey: %s\n", j.Name)
			fmt.Fprintf(os.Stderr, "   %s\n", j.Desc)
			
			// Fresh browser for each journey
			b, err := setupBrowser()
			if err != nil {
				fmt.Fprintf(os.Stderr, "   ❌ Browser error: %v\n", err)
				allPassed = false
				continue
			}
			
			page, err := b.NewPage(viewport, colorScheme)
			if err != nil {
				fmt.Fprintf(os.Stderr, "   ❌ Page error: %v\n", err)
				b.Close()
				allPassed = false
				continue
			}
			
			run, err := browser.NewRun(j.Name, screensDir, viewport, colorScheme)
			if err != nil {
				fmt.Fprintf(os.Stderr, "   ❌ Run error: %v\n", err)
				b.Close()
				allPassed = false
				continue
			}
			
			// Attach listeners
			eventLog := run.EventLog()
			eventLog.Attach(page)
			
			// Run the journey
			journeyPassed := false
			func() {
				defer func() {
					if r := recover(); r != nil {
						fmt.Fprintf(os.Stderr, "   ❌ Panic: %v\n", r)
					}
				}()
				
				// Create context
				ctx := &journeys.Context{
					BaseURL: base,
					Page:    page,
					Gates:   nil,
					Logf: func(format string, args ...interface{}) {
						fmt.Fprintf(os.Stderr, "   "+format+"\n", args...)
					},
				}
				
				// Run journey
				err := j.Run(ctx)
				if err != nil {
					fmt.Fprintf(os.Stderr, "   ❌ %v\n", err)
				} else {
					journeyPassed = true
				}
			}()
			
			// Capture results if journey passed assertions
			if journeyPassed {
				// Wait for fonts and images
				page.WaitForTimeout(500)
				
				// Get HTML
				html, _ := page.Content()
				run.WriteHTML(html)
				
				// Screenshot
				screenshotData, _ := page.Screenshot(playwright.PageScreenshotOptions{
					FullPage: playwright.Bool(true),
				})
				run.WriteScreenshot(screenshotData)
				
				// Run gates
				mode := uxgates.ParseMode(uxGates)
				recorder := uxgates.NewRecorder(mode)
				
				// Pagestate
				pagestateFindings := uxgates.PagestateGate(page, []string{"main", "nav"})
				recorder.Add(uxgates.Result{Gate: "pagestate", Target: j.Name, Pass: len(pagestateFindings) == 0, Findings: pagestateFindings})
				
				// Render
				renderFindings := uxgates.RenderGate(page)
				recorder.Add(uxgates.Result{Gate: "render", Target: j.Name, Pass: len(renderFindings) == 0, Findings: renderFindings})
				
				// A11y
				a11yFindings := uxgates.A11yGate(page)
				recorder.Add(uxgates.Result{Gate: "a11y", Target: j.Name, Pass: len(a11yFindings) == 0, Findings: a11yFindings})
				
				// Console
				consoleFindings := uxgates.ConsoleGate(eventLog, []uxgates.ConsoleEntryAllowlist{})
				recorder.Add(uxgates.Result{Gate: "console", Target: j.Name, Pass: len(consoleFindings) == 0, Findings: consoleFindings})
				
				// Network
				networkFindings := uxgates.NetworkGate(eventLog, []string{})
				recorder.Add(uxgates.Result{Gate: "network", Target: j.Name, Pass: len(networkFindings) == 0, Findings: networkFindings})
				
				// Write gates
				gatesPath, _ := recorder.WriteJSON(run.RunDir())
				
				// Output JSON
				result := map[string]interface{}{
					"journey":    j.Name,
					"suite":      j.Suite,
					"desc":       j.Desc,
					"url":        base,
					"viewport":   viewport,
					"passed":     !recorder.Failed(),
					"run_dir":    run.RunDir(),
					"screenshot": run.RunDir() + "/screenshot.png",
					"gates":      gatesPath,
				}
				
				data, _ := json.MarshalIndent(result, "", "  ")
				fmt.Println(string(data))
				
				fmt.Fprintf(os.Stderr, fmt.Sprintf("   ✅ %s passed\n", j.Name))
				fmt.Fprintf(os.Stderr, fmt.Sprintf("   📁 Bundle: %s\n", run.RunDir()))
			} else {
				allPassed = false
			}
			
			// Clean up
			run.Close()
			b.Close()
		}
		
		fmt.Fprintf(os.Stderr, "\n")
		if allPassed {
			fmt.Fprintf(os.Stderr, "🎉 All journeys passed!\n")
		} else {
			fmt.Fprintf(os.Stderr, "❌ Some journeys failed!\n")
			os.Exit(1)
		}
	},
}

func init() {
	journeyCmd.Flags().StringVar(&journeyName, "name", "", "Run specific journey by name")
	journeyCmd.Flags().StringVar(&journeySuite, "suite", "", "Run all journeys in suite")
	journeyCmd.Flags().BoolVar(&journeyList, "list", false, "List all journeys")
	journeyCmd.Flags().StringVar(&journeyGrep, "grep", "", "Filter journeys by name/description")
	rootCmd.AddCommand(journeyCmd)
}
