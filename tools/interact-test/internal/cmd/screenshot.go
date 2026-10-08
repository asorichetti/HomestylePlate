package cmd

import (
	"fmt"
	"os"
	"strings"

	playwright "github.com/mxschmitt/playwright-go"
	"github.com/spf13/cobra"
)

var screenshotURL string
var viewports string
var schemes string

// screenshotCmd takes screenshots at different viewports/schemes
var screenshotCmd = &cobra.Command{
	Use:   "screenshot [url]",
	Short: "Take screenshots at different viewports/color schemes",
	Run: func(cmd *cobra.Command, args []string) {
		url := "/"
		if len(args) > 0 {
			url = args[0]
		} else if screenshotURL != "" {
			url = screenshotURL
		}
		
		base := getBaseURL()
		fullURL := base + url
		
		// Parse viewports and schemes
		vps := []string{"desktop"}
		if viewports != "" {
			for _, v := range strings.Split(viewports, ",") {
				vps = append(vps, strings.TrimSpace(v))
			}
		}
		
		sch := []string{""} // default
		if schemes != "" {
			for _, s := range strings.Split(schemes, ",") {
				sch = append(sch, strings.TrimSpace(s))
			}
		}
		
		fmt.Fprintf(os.Stderr, "📸 Taking screenshots for: %s\n", fullURL)
		fmt.Fprintf(os.Stderr, "   Viewports: %v\n", vps)
		fmt.Fprintf(os.Stderr, "   Schemes: %v\n", sch)
		
		for _, vp := range vps {
			for _, sc := range sch {
				label := vp
				if sc != "" {
					label += "-" + sc
				}
				
				fmt.Fprintf(os.Stderr, "   🖼️  %s...", label)
				
				b, err := setupBrowser()
				if err != nil {
					fmt.Fprintf(os.Stderr, " ❌ %v\n", err)
					continue
				}
				
				page, err := b.NewPage(vp, sc)
				if err != nil {
					fmt.Fprintf(os.Stderr, " ❌ %v\n", err)
					b.Close()
					continue
				}
				
				_, err = page.Goto(fullURL, playwright.PageGotoOptions{
					WaitUntil: playwright.WaitUntilStateLoad,
				})
				if err != nil {
					fmt.Fprintf(os.Stderr, " ⚠️  %v\n", err)
				}
				
				page.WaitForTimeout(500)
				
				filename := fmt.Sprintf("%s%s-%s.png", screensDir, strings.TrimPrefix(url, "/"), label)
				os.MkdirAll(screensDir, 0755)
				
				_, err = page.Screenshot(playwright.PageScreenshotOptions{
					Path:     playwright.String(filename),
					FullPage: playwright.Bool(true),
				})
				
				b.Close()
				
				if err != nil {
					fmt.Fprintf(os.Stderr, " ❌ %v\n", err)
				} else {
					fmt.Fprintf(os.Stderr, " ✅ %s\n", filename)
				}
			}
		}
	},
}

func init() {
	screenshotCmd.Flags().StringVar(&screenshotURL, "url", "", "URL to screenshot")
	screenshotCmd.Flags().StringVar(&viewports, "viewports", "desktop", "Comma-separated viewports")
	screenshotCmd.Flags().StringVar(&schemes, "schemes", "", "Comma-separated color schemes")
	rootCmd.AddCommand(screenshotCmd)
}
