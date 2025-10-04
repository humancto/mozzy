package cmd

import (
	"os"

	"gopkg.in/yaml.v3"
	"github.com/spf13/cobra"
	"github.com/humancto/mozzy/internal/chain"
)

var runCmd = &cobra.Command{
	Use:   "run <flow.yaml>",
	Short: "Run a YAML workflow (steps with captures and vars)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		path := args[0]
		b, err := os.ReadFile(path)
		if err != nil { return err }
		var flow chain.Flow
		if err := yaml.Unmarshal(b, &flow); err != nil { return err }
		flow.EnvName = envName
		flow.BaseURL = baseURL
		flow.GlobalAuth = authToken
		return chain.Run(cmd.Context(), flow)
	},
}

func init() { rootCmd.AddCommand(runCmd) }
