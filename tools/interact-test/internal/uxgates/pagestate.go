package uxgates

import (
	"fmt"
	"regexp"

	playwright "github.com/mxschmitt/playwright-go"
)

// PagestateGate checks that the page has content
func PagestateGate(page playwright.Page, expectedSelectors []string) []Finding {
	var findings []Finding
	
	// Check for content markers
	_, err := page.Evaluate(`() => {
		// Check for main content area
		const main = document.querySelector('main, .content, #content');
		if (main) {
			return {
				hasMain: true,
				textLength: main.textContent.trim().length,
				hasVisibleChildren: main.children.length > 0
			};
		}
		return { hasMain: false, textLength: 0 };
	}`)
	
	if err != nil {
		findings = append(findings, Finding{
			Severity: "critical",
			Message:  "Failed to evaluate page content",
			Evidence: map[string]interface{}{"error": err.Error()},
		})
		return findings
	}
	
	// Check expected selectors
	for _, sel := range expectedSelectors {
		js := "(sel) => document.querySelectorAll(sel).length"
		matchCount, err := page.Evaluate(js, sel)
		if err != nil || matchCount == nil {
			findings = append(findings, Finding{
				Severity: "warning",
				Message:  fmt.Sprintf("Could not evaluate selector: %s", sel),
			})
			continue
		}
		
		var count int
		switch v := matchCount.(type) {
		case float64:
			count = int(v)
		case int:
			count = v
		default:
			count = 0
		}
		
		if count == 0 {
			findings = append(findings, Finding{
				Severity: "warning",
				Message:  fmt.Sprintf("Expected selector not found: %s", sel),
				Evidence: map[string]interface{}{"selector": sel, "matches": 0},
			})
		}
	}
	
	return findings
}

// IsExpected404 checks if a URL contains the expected 404 marker
func IsExpected404(url string) bool {
	return regexp.MustCompile(`__expected-404`).MatchString(url)
}
