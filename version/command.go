package version

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
	"github.com/tliron/puccini/format"
)

func NewCommand(name string) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: fmt.Sprintf("Show the version of %s", name),
		Long:  fmt.Sprintf(`Shows the version of %s.`, name),
		Run: func(cmd *cobra.Command, args []string) {
			// Why not use the version from runtime/debug.ReadBuildInfo? See:
			// https://github.com/golang/go/issues/29228
			if GitVersion != "" {
				fmt.Fprintf(format.Stdout, "version=%s\n", GitVersion)
			}
			if GitRevision != "" {
				fmt.Fprintf(format.Stdout, "revision=%s\n", GitRevision)
			}
			fmt.Fprintf(format.Stdout, "arch=%s\n", runtime.GOARCH)
			fmt.Fprintf(format.Stdout, "os=%s\n", runtime.GOOS)
			fmt.Fprintf(format.Stdout, "compiler=%s\n", runtime.Compiler)
		},
	}
}
