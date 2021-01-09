package app

import (
	"flag"
	"os"

	"github.com/spf13/pflag"
	cmd2 "github.com/x893675/cmdadmin/app/cmd"
	cliflag "k8s.io/component-base/cli/flag"
	"k8s.io/klog/v2"
)

func Run() error {
	klog.InitFlags(nil)
	pflag.CommandLine.SetNormalizeFunc(cliflag.WordSepNormalizeFunc)
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)

	pflag.Set("logtostderr", "true")
	// We do not want these flags to show up in --help
	// These MarkHidden calls must be after the lines above
	pflag.CommandLine.MarkHidden("version")
	pflag.CommandLine.MarkHidden("log-flush-frequency")
	pflag.CommandLine.MarkHidden("alsologtostderr")
	pflag.CommandLine.MarkHidden("log-backtrace-at")
	pflag.CommandLine.MarkHidden("log-dir")
	pflag.CommandLine.MarkHidden("logtostderr")
	pflag.CommandLine.MarkHidden("stderrthreshold")
	pflag.CommandLine.MarkHidden("vmodule")

	cmd := cmd2.NewCertAdminCommand(os.Stdin, os.Stdout, os.Stderr)
	return cmd.Execute()
}
