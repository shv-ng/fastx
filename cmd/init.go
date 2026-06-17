package cmd

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
	"github.com/shv-ng/fastx/internal/tui"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "initiate a new fastapi project",
	RunE: func(cmd *cobra.Command, args []string) error {
		p := tea.NewProgram(tui.InitialModel())
		_, err := p.Run()
		if err != nil {
			return fmt.Errorf("fail to run bubbletea: %w", err)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
