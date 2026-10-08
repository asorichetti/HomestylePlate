package uxgates

import (
	"fmt"
	"regexp"

	"github.com/asorichetti/HomestylePlate/tools/interact-test/internal/browser"
)

// NetworkGate checks for failed network requests
func NetworkGate(eventLog *browser.EventLog, expected404Patterns []string) []Finding {
	var findings []Finding
	
	// Check 400+ responses
	for _, resp := range eventLog.Network400 {
		// Skip if it's an expected 404
		if IsExpected404(resp.URL) {
			continue
		}
		
		// Check against expected patterns
		skip := false
		for _, pattern := range expected404Patterns {
			if regexp.MustCompile(pattern).MatchString(resp.URL) {
				skip = true
				break
			}
		}
		if skip {
			continue
		}
		
		findings = append(findings, Finding{
			Severity: "error",
			Message:  fmt.Sprintf("HTTP %d response: %s", resp.Status, resp.URL),
			Evidence: map[string]interface{}{
				"url": resp.URL,
				"status": resp.Status,
				"method": resp.Method,
			},
		})
	}
	
	// Check failed requests
	for _, req := range eventLog.NetworkFailed {
		if IsExpected404(req.URL) {
			continue
		}
		
		findings = append(findings, Finding{
			Severity: "error",
			Message:  fmt.Sprintf("Network request failed: %s", req.URL),
			Evidence: map[string]interface{}{
				"url": req.URL,
				"failure_text": req.FailureText,
			},
		})
	}
	
	return findings
}
