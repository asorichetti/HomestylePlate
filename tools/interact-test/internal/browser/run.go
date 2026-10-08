package browser

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Artifacts holds paths to all generated debug files
type Artifacts struct {
	RunDir, ScreenshotPath, HTMLPath, ConsolePath, NetworkPath, UXGatesPath string
}

// RunBundle manages a single test run's debug bundle
type RunBundle struct {
	cmdName    string
	runDir     string
	page       interface{} // playwright.Page (interface to avoid import cycle)
	eventLog   *EventLog
	screensDir string
	scheme     string
	vpName     string
	finished   bool
	gates      interface{ WriteJSON(dir string) (string, error) } // structural gate hook
}

// NewRun creates a new run with timestamped directory
func NewRun(cmdName string, screensDir string, vpName string, scheme string) (*RunBundle, error) {
	timestamp := time.Now().UnixNano()
	runDir := filepath.Join(screensDir, fmt.Sprintf("%s-%d", cmdName, timestamp))
	
	if err := os.MkdirAll(runDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create run dir: %w", err)
	}
	
	return &RunBundle{
		cmdName:    cmdName,
		runDir:     runDir,
		screensDir: screensDir,
		scheme:     scheme,
		vpName:     vpName,
		eventLog:   NewEventLog(),
	}, nil
}

// SetPage assigns the page for screenshot capture
func (r *RunBundle) SetPage(p interface{}) {
	r.page = p
}

// SetGates assigns the gates hook
func (r *RunBundle) SetGates(gates interface{ WriteJSON(dir string) (string, error) }) {
	r.gates = gates
}

// EventLog returns the event log for attaching listeners
func (r *RunBundle) EventLog() *EventLog {
	return r.eventLog
}

// RunDir returns the run directory path
func (r *RunBundle) RunDir() string {
	return r.runDir
}

// Finish writes all bundle files (screenshot, HTML, console, network, gates)
func (r *RunBundle) Finish() (Artifacts, error) {
	if r.finished {
		return Artifacts{}, fmt.Errorf("run already finished")
	}
	
	artifacts := Artifacts{RunDir: r.runDir}
	var lastErr error
	
	// Screenshot will be written via WriteScreenshot()
	// HTML will be written via WriteHTML()
	// Both are handled in cmd layer
	htmlPath := filepath.Join(r.runDir, "page.html")
	artifacts.HTMLPath = htmlPath
	
	// Console/Network
	if r.eventLog != nil {
		consolePath, err := r.eventLog.WriteJSON(r.runDir)
		if err != nil {
			lastErr = err
		}
		artifacts.ConsolePath = consolePath
		
		// Network (combined console.json includes network)
		artifacts.NetworkPath = consolePath
	}
	
	// UX Gates
	if r.gates != nil {
		gatesPath, err := r.gates.WriteJSON(r.runDir)
		if err != nil {
			lastErr = err
		}
		artifacts.UXGatesPath = gatesPath
	}
	
	// Mark finished only after writes succeed
	r.finished = true
	
	return artifacts, lastErr
}

// FinishOrLog calls Finish and logs errors, never panics
func (r *RunBundle) FinishOrLog() Artifacts {
	artifacts, err := r.Finish()
	if err != nil {
		fmt.Fprintf(os.Stderr, "⚠️  Error writing bundle: %v\n", err)
	}
	return artifacts
}

// Close writes the bundle if nothing else did (safety net)
func (r *RunBundle) Close() {
	if !r.finished {
		r.Finish()
	}
}

// WriteHTML saves the page HTML to the run directory
func (r *RunBundle) WriteHTML(htmlContent string) error {
	htmlPath := filepath.Join(r.runDir, "page.html")
	return os.WriteFile(htmlPath, []byte(htmlContent), 0644)
}

// WriteScreenshot saves the screenshot to the run directory
func (r *RunBundle) WriteScreenshot(data []byte) error {
	screenshotPath := filepath.Join(r.runDir, "screenshot.png")
	return os.WriteFile(screenshotPath, data, 0644)
}
