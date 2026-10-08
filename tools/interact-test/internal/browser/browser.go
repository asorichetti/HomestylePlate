package browser

import (
	"fmt"
	"os"
	"runtime"

	playwright "github.com/mxschmitt/playwright-go"
)

// ViewportPreset defines standard viewport sizes
type ViewportPreset struct {
	Name   string
	Width  int
	Height int
}

// Viewports are the standard viewport presets
var Viewports = map[string]playwright.Size{
	"desktop": {Width: 1280, Height: 900},
	"tablet":  {Width: 820, Height: 1180},
	"mobile":  {Width: 390, Height: 844},
}

// Browser wraps Playwright browser instance
type Browser struct {
	pw      *playwright.Playwright
	browser playwright.Browser
}

// New launches a headless Playwright browser
func New(headless bool) (*Browser, error) {
	pw, err := playwright.Run()
	if err != nil {
		return nil, fmt.Errorf("failed to start playwright: %w", err)
	}

	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(headless),
		Args:     []string{"--disable-dev-shm-usage"},
	})
	if err != nil {
		pw.Stop()
		return nil, fmt.Errorf("failed to launch browser: %w", err)
	}

	return &Browser{pw: pw, browser: browser}, nil
}

// NewPage creates a new page with the specified viewport and color scheme
func (b *Browser) NewPage(vpName string, colorScheme string) (playwright.Page, error) {
	vp, ok := Viewports[vpName]
	if !ok {
		vp = Viewports["desktop"]
	}

	// Override height from env for full-page audits
	if h := os.Getenv("VIEWPORT_HEIGHT"); h != "" {
		vp.Height = 0 // will be set dynamically
	}

	page, err := b.browser.NewPage(playwright.BrowserNewPageOptions{
		Viewport:    &playwright.Size{Width: vp.Width, Height: vp.Height},
		ColorScheme: translateColorScheme(colorScheme),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create page: %w", err)
	}

	return page, nil
}

// Close shuts down the browser
func (b *Browser) Close() {
	// Log errors instead of returning them (safe for defer)
	if b.browser != nil {
		_ = b.browser.Close()
	}
	if b.pw != nil {
		_ = b.pw.Stop()
	}
}

// Screenshot captures a screenshot of the page
func Screenshot(page playwright.Page, outPath, dir string) (string, error) {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("failed to create screenshot dir: %w", err)
	}

	_, err := page.Screenshot(playwright.PageScreenshotOptions{
		Path: playwright.String(outPath),
		FullPage: playwright.Bool(true),
	})
	if err != nil {
		return "", fmt.Errorf("failed to take screenshot: %w", err)
	}

	return outPath, nil
}

// GetBrowser returns the underlying browser (for advanced use)
func (b *Browser) GetBrowser() playwright.Browser {
	return b.browser
}

// GetPlaywright returns the underlying Playwright instance
func (b *Browser) GetPlaywright() *playwright.Playwright {
	return b.pw
}

func translateColorScheme(s string) *playwright.ColorScheme {
	switch s {
	case "dark":
		return playwright.ColorSchemeDark
	case "light":
		return playwright.ColorSchemeLight
	default:
		return nil // browser default
	}
}

// ShouldRunHeadless computes headless default from display presence
func ShouldRunHeadless() bool {
	if os.Getenv("CI") != "" {
		return true
	}
	if os.Getenv("DISPLAY") != "" || os.Getenv("WAYLAND_DISPLAY") != "" {
		return false
	}
	return runtime.GOOS != "darwin" // macOS draws via Quartz
}
