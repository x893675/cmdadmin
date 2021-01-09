package cmd

import (
	"github.com/spf13/cobra"
	"github.com/x893675/certadmin/app/cmd/util"
	"io"
)

func NewCertAdminCommand(in io.Reader, out, err io.Writer) *cobra.Command {
	var rootfsPath string

	cmds := &cobra.Command{
		Use:           "certadmin",
		Short:         "certadmin: easily generated cert for server and client",
		Long:          "todo",
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
