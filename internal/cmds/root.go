package cmds

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/zeet-dev/jsonnet-filer/internal/clioptions/iostreams"
	"github.com/zeet-dev/jsonnet-filer/internal/jsonnet"
)

func NewRootCmd(s iostreams.IOStreams) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:           "jsonnet-filer",
		Short:         "Generate configuration files using jsonnet",
		Args:          cobra.ArbitraryArgs,
		SilenceErrors: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmd.ParseFlags(args); err != nil {
				return err
			}

			//if len(args) == 0 {
			//	return runHelp(cmd)
			//}

			result, err := jsonnet.EvaluateFile("./cabinet.jsonnet")
			if err != nil {
				return err
			}

			fmt.Fprintf(s.Out(), "%+v\n", result)

			return nil
		},
	}

	rootCmd.SetHelpTemplate(`{{.Long}}
{{.Short}}

{{if or .Runnable .HasSubCommands}}{{.UsageString}}{{end}}`)

	return rootCmd
}

func runHelp(cmd *cobra.Command) error {
	return cmd.Help()
}
