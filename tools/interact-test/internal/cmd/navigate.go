package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	playwright "github.com/mxschmitt/playwright-go"
	"github.com/asorichetti/HomestylePlate/tools/interact-test/internal/browser"
	"github.com/asorichetti/HomestylePlate/tools/interact-test/internal/uxgates"
	"github.com/spf13/cobra"
)

var navigateURL string

// navigateCmd opens a URL and captures debug info
var navigateCmd = &cobra.Command{
	Use:   "navigate [url]",
	Short: "Open a URL and capture debug bundle",
	Run: func(cmd *cobra.Command, args []string) {
		url := navigateURL
		if len(args) > 0 {
			url = args[0]
		}
		
		if url == "" {
			fmt.Fprintf(os.Stderr, "Error: URL required\n")
			os.Exit(1)
		}
		
		// Setup
		base := getBaseURL()
		if url == "/" || url == "" {
			url = "/"
		}
		fullURL := base + url
		if url != "/" {
			fullURL = base + url
		}
		
		fmt.Fprintf(os.Stderr, "🌐 Navigating to: %s\n", fullURL)
		fmt.Fprintf(os.Stderr, "🖥️  Headless: %v | Viewport: %s | Color: %s\n", headless, viewport, colorScheme)
		
		// Create browser
		b, err := setupBrowser()
		if err != nil {
			fmt.Fprintf(os.Stderr, "❌ %v\n", err)
			os.Exit(1)
		}
		defer b.Close()
		
		// Create new page
		page, err := b.NewPage(viewport, colorScheme)
		if err != nil {
			fmt.Fprintf(os.Stderr, "❌ %v\n", err)
			os.Exit(1)
		}
		
		// Create run bundle
		run, err := browser.NewRun("navigate", screensDir, viewport, colorScheme)
		if err != nil {
			fmt.Fprintf(os.Stderr, "❌ %v\n", err)
			os.Exit(1)
		}
		defer run.Close()
		
		// Attach listeners
		eventLog := run.EventLog()
		eventLog.Attach(page)
		
		// Navigate
		fmt.Fprintf(os.Stderr, "⏳ Loading page...\n")
		_, err = page.Goto(fullURL, playwright.PageGotoOptions{})
		if err != nil {
			fmt.Fprintf(os.Stderr, "⚠️  Navigation warning: %v\n", err)
		}
		
		// Wait for page to settle
		page.WaitForTimeout(3000)
		
		// Wait 1 second for fonts and images
		page.WaitForTimeout(1000)
		
		// Get HTML
		html, _ := page.Content()
		run.WriteHTML(html)
		
		// Take screenshot
		screenshotData, _ := page.Screenshot(playwright.PageScreenshotOptions{
			FullPage: playwright.Bool(true),
		})
		run.WriteScreenshot(screenshotData)
		
		// Run gates
		mode := uxgates.ParseMode(uxGates)
		recorder := uxgates.NewRecorder(mode)
		
		// Pagestate gate
		pagestateFindings := uxgates.PagestateGate(page, []string{"main", "nav"})
		if len(pagestateFindings) > 0 {
			recorder.Add(uxgates.Result{
				Gate:   "pagestate",
				Target: url,
				Pass:   false,
				Findings: pagestateFindings,
			})
		} else {
			recorder.Add(uxgates.Result{
				Gate:   "pagestate",
				Target: url,
				Pass:   true,
			})
		}
		
		// Render gate
		renderFindings := uxgates.RenderGate(page)
		recorder.Add(uxgates.Result{
			Gate:   "render",
			Target: url,
			Pass:   len(renderFindings) == 0,
			Findings: renderFindings,
		})
		
		// A11y gate
		a11yFindings := uxgates.A11yGate(page)
		recorder.Add(uxgates.Result{
			Gate:   "a11y",
			Target: url,
			Pass:   len(a11yFindings) == 0,
			Findings: a11yFindings,
		})
		
		// Console/Network gates
		consoleFindings := uxgates.ConsoleGate(eventLog, []uxgates.ConsoleEntryAllowlist{})
		recorder.Add(uxgates.Result{
			Gate:   "console",
			Target: url,
			Pass:   len(consoleFindings) == 0,
			Findings: consoleFindings,
		})
		
		networkFindings := uxgates.NetworkGate(eventLog, []string{})
		recorder.Add(uxgates.Result{
			Gate:   "network",
			Target: url,
			Pass:   len(networkFindings) == 0,
			Findings: networkFindings,
		})
		
		// Write gate results
		gatesPath, _ := recorder.WriteJSON(run.RunDir())
		
		// Output JSON result
		result := map[string]interface{}{
			"cmd":      "navigate",
			"url":      fullURL,
			"viewport": viewport,
			"color_scheme": colorScheme,
			"headless": headless,
			"ux_gates": uxGates,
			"run_dir":  run.RunDir(),
			"screenshot": run.RunDir() + "/screenshot.png",
			"html":     run.RunDir() + "/page.html",
			"console":  run.RunDir() + "/console.json",
			"gates":    gatesPath,
			"passed":   !recorder.Failed(),
			"results": recorder.Results(),
		}
		
		// Print JSON to stdout
		data, _ := json.MarshalIndent(result, "", "  ")
		fmt.Println(string(data))
		
		// Print summary to stderr
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "📁 Bundle: %s\n", run.RunDir())
		fmt.Fprintf(os.Stderr, "📸 Screenshot: screenshot.png\n")
		fmt.Fprintf(os.Stderr, "📄 HTML: page.html\n")
		fmt.Fprintf(os.Stderr, "📊 Console/Network: console.json\n")
		fmt.Fprintf(os.Stderr, "🚦 UX Gates: %s\n", uxGates)
		fmt.Fprintf(os.Stderr, "✅ Passed: %v\n", !recorder.Failed())
		
		if recorder.Failed() {
			fmt.Fprintf(os.Stderr, "\n❌ Some gates failed!\n")
			os.Exit(1)
		}
	},
}

func init() {
	navigateCmd.Flags().StringVarP(&navigateURL, "url", "u", "/", "URL to navigate to")
	rootCmd.AddCommand(navigateCmd)
}
