package cmds

import (
	"encoding/json"
	"fmt"
	"path"
	"path/filepath"

	"github.com/hexops/valast"
	"github.com/spf13/cobra"
	"sigs.k8s.io/yaml"

	"github.com/zeet-dev/jsonnet-filer/internal/cabinet"
	"github.com/zeet-dev/jsonnet-filer/internal/clioptions/iostreams"
	"github.com/zeet-dev/jsonnet-filer/internal/jsonnet"
	"github.com/zeet-dev/jsonnet-filer/pkg/api/v1alpha1"
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

			resultString, err := jsonnet.EvaluateFile("./cabinet.jsonnet")
			if err != nil {
				return err
			}

			var resultMap map[string]interface{}
			err = json.Unmarshal([]byte(resultString), &resultMap)
			if err != nil {
				return err
			}

			files := cabinet.Find[v1alpha1.File](resultMap)
			fmt.Fprintln(s.ErrOut(), valast.String(files))

			for _, f := range files {
				// TODO this is just here for ease of reading
				// We only want to write the f.Content in the format given in f.EncodingStrategy
				contentBytes, err := yaml.Marshal(f)
				if err != nil {
					return err
				}

				fpath := path.Join(".", f.Metadata.Name)
				fpath, err = filepath.Abs(fpath)
				if err != nil {
					return err
				}

				fmt.Fprintf(s.ErrOut(), "writing file: %s\n", fpath)
				fmt.Fprintln(s.ErrOut(), "-----------")
				fmt.Fprintln(s.ErrOut(), string(contentBytes))
				fmt.Fprintln(s.ErrOut(), "-----------")
				//err = os.WriteFile(fpath, contentBytes, 0644)
				//if err != nil {
				//	return err
				//}
			}

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
