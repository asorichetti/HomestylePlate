package uxgates

import (
	playwright "github.com/mxschmitt/playwright-go"
)

// A11yGate checks for structural accessibility issues
func A11yGate(page playwright.Page) []Finding {
	var findings []Finding
	
	// Single evaluation for all a11y checks
	result, err := page.Evaluate(`() => {
		const findings = [];
		
		// 1. Exactly one h1
		const h1s = document.querySelectorAll('h1');
		if (h1s.length === 0) {
			findings.push({ severity: 'error', message: 'No h1 element found' });
		} else if (h1s.length > 1) {
			findings.push({ severity: 'error', message: 'Multiple h1 elements found (' + h1s.length + ')' });
		}
		
		// 2. Heading level skips
		const headings = document.querySelectorAll('h1, h2, h3, h4, h5, h6');
		let prevLevel = 0;
		headings.forEach(h => {
			const level = parseInt(h.tagName.substring(1));
			if (prevLevel > 0 && level > prevLevel + 1) {
				findings.push({
					severity: 'warning',
					message: 'Heading level skip: h' + prevLevel + ' to h' + level,
					detail: h.textContent.substring(0, 50)
				});
			}
			prevLevel = level;
		});
		
		// 3. Images with alt
		const imgs = document.querySelectorAll('img');
		imgs.forEach(img => {
			if (!img.hasAttribute('alt') || img.alt.trim() === '') {
				// Empty alt is valid for decorative images
				if (!img.hasAttribute('decorative')) {
					findings.push({
						severity: 'warning',
						message: 'Image missing alt text',
						detail: img.src ? img.src.substring(0, 100) : 'no-src'
					});
				}
			}
		});
		
		// 4. Buttons/links with accessible names
		const buttons = document.querySelectorAll('button, [role="button"]');
		buttons.forEach(btn => {
			const name = btn.textContent.trim() || btn.getAttribute('aria-label') || btn.getAttribute('title');
			if (!name) {
				findings.push({ severity: 'error', message: 'Button has no accessible name' });
			}
		});
		
		const links = document.querySelectorAll('a');
		links.forEach(link => {
			const name = link.textContent.trim() || link.getAttribute('aria-label') || link.getAttribute('title');
			if (!name) {
				findings.push({ severity: 'warning', message: 'Link has no accessible name' });
			}
		});
		
		// 5. HTML lang attribute
		const html = document.documentElement;
		if (!html.getAttribute('lang')) {
			findings.push({ severity: 'error', message: 'html element missing lang attribute' });
		}
		
		// 6. Exactly one main
		const mains = document.querySelectorAll('main');
		if (mains.length === 0) {
			findings.push({ severity: 'warning', message: 'No main element found' });
		} else if (mains.length > 1) {
			findings.push({ severity: 'warning', message: 'Multiple main elements found (' + mains.length + ')' });
		}
		
		return findings;
	}`)
	
	if err != nil {
		findings = append(findings, Finding{
			Severity: "critical",
			Message:  "Failed to evaluate a11y checks",
			Evidence: map[string]interface{}{"error": err.Error()},
		})
		return findings
	}
	
	// Parse findings
	if results, ok := result.([]interface{}); ok {
		for _, r := range results {
			if f, ok := r.(map[string]interface{}); ok {
				severity, _ := f["severity"].(string)
				message, _ := f["message"].(string)
				
				findings = append(findings, Finding{
					Severity: severity,
					Message:  message,
					Evidence: f,
				})
			}
		}
	}
	
	return findings
}
