package uxgates

import (
	"fmt"
	"regexp"

	"github.com/asorichetti/HomestylePlate/tools/interact-test/internal/browser"
)

// ConsoleEntryAllowlist defines an allowlist entry for console errors
type ConsoleEntryAllowlist struct {
	Pattern string `json:"pattern"`
	Reason  string `json:"reason"`
}

// ConsoleGate checks for console errors not in the allowlist
func ConsoleGate(eventLog *browser.EventLog, allowlist []ConsoleEntryAllowlist) []Finding {
	var findings []Finding
	
	for _, entry := range eventLog.Console {
		if entry.Type != "error" {
			continue
		}
		
		// Check if it's in the allowlist
		allowed := false
		for _, allowEntry := range allowlist {
			if regexp.MustCompile(allowEntry.Pattern).MatchString(entry.Text) {
				allowed = true
				break
			}
		}
		if allowed {
			continue
		}
		
		findings = append(findings, Finding{
			Severity: "error",
			Message:  fmt.Sprintf("Console error: %s", entry.Text),
			Evidence: map[string]interface{}{
				"type": entry.Type,
				"text": entry.Text,
				"line": entry.Line,
				"column": entry.Column,
			},
		})
	}
	
	// Check page errors
	for _, errMsg := range eventLog.PageErrors {
		// Check allowlist
		allowed := false
		for _, entry := range allowlist {
			if regexp.MustCompile(entry.Pattern).MatchString(errMsg) {
				allowed = true
				break
			}
		}
		if allowed {
			continue
		}
		
		findings = append(findings, Finding{
			Severity: "error",
			Message:  fmt.Sprintf("Page error: %s", errMsg),
			Evidence: map[string]interface{}{"error": errMsg},
		})
	}
	
	return findings
}
