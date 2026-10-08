package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// smokeCmd runs a quick smoke test
var smokeCmd = &cobra.Command{
	Use:   "smoke",
	Short: "Run smoke test (core journeys with gates=FAIL)",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Fprintf(os.Stderr, "🔥 Running smoke test...\n")
		fmt.Fprintf(os.Stderr, "   Suite: core\n")
		fmt.Fprintf(os.Stderr, "   UX Gates: fail\n")
		fmt.Fprintf(os.Stderr, "   Headless: true\n")
		
		// Override flags for smoke test
		os.Args = []string{os.Args[0], "journey", "--suite", "core", "--ux-gates", "fail", "--headless", "true"}
		rootCmd.SetArgs(os.Args[1:])
		rootCmd.Execute()
	},
}

func init() {
	rootCmd.AddCommand(smokeCmd)
}
