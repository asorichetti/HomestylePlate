package uxgates

import (
	"strings"

	playwright "github.com/mxschmitt/playwright-go"
)

// RenderGate checks that the page is styled, not just present
func RenderGate(page playwright.Page) []Finding {
	var findings []Finding
	
	// Single evaluation for all render checks
	result, err := page.Evaluate(`() => {
		const checks = {
			hasStylesheets: false,
			hasCSSRules: false,
			bodyMarginZero: false,
			bodyBgTransparent: false,
			defaultSerif: false,
			imagesLoaded: [],
			fontsLoaded: [],
			fontsErrored: [],
			overflow: false
		};
		
		// 1. Check stylesheets
		checks.hasStylesheets = document.styleSheets.length > 0;
		
		// Check CSS rules (try/catch for cross-origin)
		try {
			let totalRules = 0;
			for (let i = 0; i < document.styleSheets.length; i++) {
				totalRules += document.styleSheets[i].cssRules.length;
			}
			checks.hasCSSRules = totalRules > 0;
		} catch (e) {
			// Cross-origin, can't check rules
			checks.hasCSSRules = checks.hasStylesheets;
		}
		
		// 2. Body margin (reset should be 0)
		const bodyStyle = getComputedStyle(document.body);
		checks.bodyMarginZero = bodyStyle.marginTop === '0px';
		checks.bodyBgTransparent = bodyStyle.backgroundColor === 'rgba(0, 0, 0, 0)';
		
		// 3. Default serif font detection
		const fontFamily = bodyStyle.fontFamily;
		checks.defaultSerif = /^(Times|serif|-webkit-standard)/i.test(fontFamily);
		
		// 4. Images
		const images = document.querySelectorAll('img');
		images.forEach((img, idx) => {
			checks.imagesLoaded.push({
				complete: img.complete,
				naturalWidth: img.naturalWidth,
				naturalHeight: img.naturalHeight,
				src: img.src ? img.src.substring(0, 100) : 'no-src'
			});
		});
		
		// 5. Fonts
		if (document.fonts) {
			document.fonts.forEach((font) => {
				checks.fontsLoaded.push(font.family);
				if (font.status === 'error') {
					checks.fontsErrored.push(font.family);
				}
			});
		}
		
		// 6. Overflow check
		checks.overflow = document.documentElement.scrollWidth > (document.documentElement.clientWidth + 1);
		
		return checks;
	}`)
	
	if err != nil {
		findings = append(findings, Finding{
			Severity: "critical",
			Message:  "Failed to evaluate render checks",
			Evidence: map[string]interface{}{"error": err.Error()},
		})
		return findings
	}
	
	// Parse results safely
	checksBytes, ok := result.(map[string]interface{})
	if !ok {
		findings = append(findings, Finding{
			Severity: "critical",
			Message:  "Failed to parse render checks result",
		})
		return findings
	}
	
	// Helper to safely get bool
	getBool := func(key string, def bool) bool {
		if v, ok := checksBytes[key]; ok {
			if b, ok := v.(bool); ok {
				return b
			}
		}
		return def
	}
	
	// Check stylesheets
	if !getBool("hasStylesheets", false) {
		findings = append(findings, Finding{
			Severity: "critical",
			Message:  "No stylesheets loaded",
			Evidence: map[string]interface{}{"stylesheets": 0},
		})
	}
	
	if !getBool("bodyMarginZero", false) {
		findings = append(findings, Finding{
			Severity: "critical",
			Message:  "CSS reset did not apply (body margin is not 0)",
			Evidence: map[string]interface{}{"body_margin_top": "non-zero"},
		})
	}
	
	if getBool("defaultSerif", false) {
		findings = append(findings, Finding{
			Severity: "critical",
			Message:  "Page rendering in default serif font (CSS not applied)",
			Evidence: map[string]interface{}{"font_family": "serif fallback"},
		})
	}
	
	// Check images
	if images, ok := checksBytes["imagesLoaded"].([]interface{}); ok {
		for _, imgData := range images {
			if imgMap, ok := imgData.(map[string]interface{}); ok {
				naturalWidth := 0
				if w, ok := imgMap["naturalWidth"]; ok {
					switch v := w.(type) {
					case float64:
						naturalWidth = int(v)
					case int:
						naturalWidth = v
					}
				}
				complete := false
				if c, ok := imgMap["complete"]; ok {
					complete = c.(bool)
				}
				src := ""
				if s, ok := imgMap["src"]; ok {
					src = s.(string)
				}
				
				// Skip external images and empty-src images
				if src == "no-src" || (strings.Contains(src, "http") && !strings.Contains(src, "localhost")) {
					continue
				}
				
				if complete && naturalWidth == 0 {
					findings = append(findings, Finding{
						Severity: "warning",
						Message:  "Image failed to load (404 or broken)",
						Evidence: map[string]interface{}{"src": src, "natural_width": 0},
					})
				}
			}
		}
	}
	
	// Check overflow
	if getBool("overflow", false) {
		findings = append(findings, Finding{
			Severity: "warning",
			Message:  "Page has horizontal overflow",
			Evidence: map[string]interface{}{"scrollWidth": "exceeds clientWidth"},
		})
	}
	
	return findings
}
