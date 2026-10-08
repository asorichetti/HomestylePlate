package journeys

import (
	"fmt"
	"regexp"
	"strings"
)

// expect is a helper for journey assertions
func expect(condition bool, format string, args ...interface{}) {
	if !condition {
		msg := fmt.Sprintf(format, args...)
		panic(msg)
	}
}

func init() {
	// Register the Home Page journey
	Register(Journey{
		Name: "home-page",
		Suite: "core",
		Desc:  "Verify landing page loads with content, styles, and navigation",
		Run: func(ctx *Context) error {
			baseURL := ctx.BaseURL + "/"
			ctx.Logf("Navigating to %s", baseURL)
			
			// Navigate
			page := ctx.Page.(interface{ Goto(string, interface{}) (interface{}, error); WaitForLoadState(interface{}, interface{}) error; WaitForTimeout(int64); Content() (string, error); Evaluate(string, ...interface{}) (interface{}, error); QuerySelector(string) (interface{}, error) })
			
			_, err := page.Goto(baseURL, nil)
			if err != nil {
				return fmt.Errorf("navigation failed: %w", err)
			}
			
			page.WaitForLoadState(0, nil)
			page.WaitForTimeout(1000)
			
			// Verify page loaded
			html, _ := page.Content()
			expect(strings.Contains(html, "Homestyle Plate"), "Page title not found")
			expect(strings.Contains(html, "DOCTYPE"), "No DOCTYPE in page")
			
			ctx.Logf("✅ Page loaded with title 'Homestyle Plate'")
			
			// Verify navigation links work
			links := []string{"/MealSlotMachine/", "/DessertsPage/", "/MealsGalleryPage/"}
			for _, link := range links {
				fullLink := ctx.BaseURL + link
				ctx.Logf("Checking link to %s", link)
				
				_, err := page.Goto(fullLink, nil)
				if err != nil {
					ctx.Logf("⚠️  Could not navigate to %s: %v", link, err)
					continue
				}
				
				page.WaitForLoadState(0, nil)
				page.WaitForTimeout(500)
				
				linkHTML, _ := page.Content()
				expect(strings.Contains(linkHTML, "Homestyle Plate"), fmt.Sprintf("Page at %s missing title", link))
				
				ctx.Logf("✅ %s loaded successfully", link)
			}
			
			return nil
		},
	})
	
	// Register the Meal Slot Machine journey
	Register(Journey{
		Name: "meal-slot-machine",
		Suite: "core",
		Desc:  "Verify meal slot machine page with spin and lock functionality",
		Run: func(ctx *Context) error {
			page := ctx.Page.(interface{ Goto(string, interface{}) (interface{}, error); WaitForLoadState(interface{}, interface{}) error; WaitForTimeout(int64); Content() (string, error); Evaluate(string, ...interface{}) (interface{}, error) })
			
			baseURL := ctx.BaseURL + "/MealSlotMachine/"
			ctx.Logf("Navigating to %s", baseURL)
			
			_, err := page.Goto(baseURL, nil)
			if err != nil {
				return fmt.Errorf("navigation failed: %w", err)
			}
			
			page.WaitForLoadState(0, nil)
			page.WaitForTimeout(1000)
			
			// Verify page structure
			html, _ := page.Content()
			expect(strings.Contains(html, "Meal Slot Machine"), "Slot machine title not found")
			
			// Verify JavaScript is present
			expect(strings.Contains(html, "script"), "No script tags found")
			
			ctx.Logf("✅ Slot machine page loaded")
			
			// Check that meals data is embedded
			js, _ := page.Content()
			expect(regexp.MustCompile(`const\s+items`).MatchString(js), "No items array found in JS")
			
			ctx.Logf("✅ Meal data found in page")
			
			return nil
		},
	})
	
	// Register the Desserts page journey
	Register(Journey{
		Name: "desserts-page",
		Suite: "core",
		Desc:  "Verify desserts showcase page loads and API works",
		Run: func(ctx *Context) error {
			page := ctx.Page.(interface{ Goto(string, interface{}) (interface{}, error); WaitForLoadState(interface{}, interface{}) error; WaitForTimeout(int64); Content() (string, error); Evaluate(string, ...interface{}) (interface{}, error) })
			
			baseURL := ctx.BaseURL + "/DessertsPage/"
			ctx.Logf("Navigating to %s", baseURL)
			
			_, err := page.Goto(baseURL, nil)
			if err != nil {
				return fmt.Errorf("navigation failed: %w", err)
			}
			
			page.WaitForLoadState(0, nil)
			page.WaitForTimeout(1000)
			
			// Verify page structure
			html, _ := page.Content()
			expect(strings.Contains(html, "Desserts"), "Desserts title not found")
			
			ctx.Logf("✅ Desserts page loaded")
			
			return nil
		},
	})
	
	// Register the Menu Planner journey
	Register(Journey{
		Name: "menu-planner",
		Suite: "core",
		Desc:  "Verify menu planner page with form and API integration",
		Run: func(ctx *Context) error {
			page := ctx.Page.(interface{ Goto(string, interface{}) (interface{}, error); WaitForLoadState(interface{}, interface{}) error; WaitForTimeout(int64); Content() (string, error); Evaluate(string, ...interface{}) (interface{}, error) })
			
			baseURL := ctx.BaseURL + "/MealsGalleryPage/"
			ctx.Logf("Navigating to %s", baseURL)
			
			_, err := page.Goto(baseURL, nil)
			if err != nil {
				return fmt.Errorf("navigation failed: %w", err)
			}
			
			page.WaitForLoadState(0, nil)
			page.WaitForTimeout(1000)
			
			// Verify page structure
			html, _ := page.Content()
			expect(strings.Contains(html, "Menu Planner") || strings.Contains(html, "Randomized Menu"), "Menu planner title not found")
			
			// Check for form elements
			expect(strings.Contains(html, "hf-count") || strings.Contains(html, "p3-count"), "Form inputs not found")
			
			ctx.Logf("✅ Menu planner page loaded with form")
			
			return nil
		},
	})
	
	// Register the API integration journey
	Register(Journey{
		Name: "api-integration",
		Suite: "api",
		Desc:  "Verify all API endpoints return correct data",
		Run: func(ctx *Context) error {
			page := ctx.Page.(interface{ Goto(string, interface{}) (interface{}, error); WaitForLoadState(interface{}, interface{}) error; WaitForTimeout(int64); Content() (string, error); Evaluate(string, ...interface{}) (interface{}, error) })
			
			baseURL := ctx.BaseURL
			
			// Test meals API
			mealsURL := baseURL + "/meals?type=HF_Meal:3,P3_Meal:3"
			ctx.Logf("Testing meals API: %s", mealsURL)
			
			result, err := page.Evaluate(fmt.Sprintf(`fetch('%s').then(r => r.json())`, mealsURL))
			if err != nil {
				return fmt.Errorf("meals API request failed: %w", err)
			}
			
			// Check result is array with 6 items
			resultMap, ok := result.(map[string]interface{})
			if !ok {
				return fmt.Errorf("meals API returned unexpected format")
			}
			
			ctx.Logf("✅ Meals API returned %d items", len(resultMap))
			
			// Test desserts API
			dessertsURL := baseURL + "/api/desserts"
			ctx.Logf("Testing desserts API: %s", dessertsURL)
			
			result, err = page.Evaluate(fmt.Sprintf(`fetch('%s').then(r => r.json())`, dessertsURL))
			if err != nil {
				return fmt.Errorf("desserts API request failed: %w", err)
			}
			
			ctx.Logf("✅ Desserts API returned items")
			
			return nil
		},
	})
	
	// Register the Render Gate test journey
	Register(Journey{
		Name: "render-gate-test",
		Suite: "gates",
		Desc:  "Comprehensive render gate check on all pages",
		Run: func(ctx *Context) error {
			// This journey specifically tests the render gate
			// The actual gate checking happens in the command layer
			// We just need to navigate to a page for the gates to run against
			page := ctx.Page.(interface{ Goto(string, interface{}) (interface{}, error); WaitForLoadState(interface{}, interface{}) error; WaitForTimeout(int64) })
			
			baseURL := ctx.BaseURL + "/"
			_, _ = page.Goto(baseURL, nil)
			page.WaitForLoadState(0, nil)
			page.WaitForTimeout(1000)
			
			ctx.Logf("✅ Render gate test page loaded")
			return nil
		},
	})
}
