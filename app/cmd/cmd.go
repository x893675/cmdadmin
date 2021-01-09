package cmd

import (
	"io"

	"github.com/spf13/cobra"
	"github.com/x893675/cmdadmin/app/cmd/util"
)

func NewCertAdminCommand(in io.Reader, out, err io.Writer) *cobra.Command {
	var rootfsPath string

	cmds := &cobra.Command{
		Use:           "cmdadmin",
		Short:         "cmdadmin: cmdline program for some utils",
		Long:          "TODO",
		SilenceErrors: true,
		SilenceUsage:  true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if rootfsPath != "" {
				if err := util.Chroot(rootfsPath); err != nil {
					return err
				}
			}
			return nil
		},
	}

	cmds.ResetFlags()
	cmds.AddCommand(newCmdCerts())
	cmds.AddCommand(newCmdVersion(out))

	return cmds
}
