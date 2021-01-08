package options

import "github.com/spf13/pflag"

// AddConfigFlag adds the --config flag to the given flagset
func AddConfigFlag(fs *pflag.FlagSet, cfgPath *string) {
	fs.StringVar(cfgPath, CfgPath, *cfgPath, "Path to a kubeadm configuration file.")
}