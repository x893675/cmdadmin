package cmd

import (
	"github.com/spf13/cobra"
	"github.com/x893675/certadmin/app/cmd/options"
	certsphase "github.com/x893675/certadmin/app/phases/certs"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type certsOptions struct {
	CertificatesDir string
	ConfigPath      string
}

func newCertsOptions() *certsOptions {
	return &certsOptions{
		CertificatesDir: "",
		ConfigPath:      "config.yaml",
	}
}

func newCmdCerts() *cobra.Command {
	opts := newCertsOptions()
	cmd := &cobra.Command{
		Use:     "certs",
		Short:   "Generate certs with config",
		Aliases: []string{"certificates", "c"},
		RunE: func(cmd *cobra.Command, args []string) error {
			return RunCerts(opts)
		},
		Args: cobra.NoArgs,
	}

	// TODO: add sub command to generate example config
	options.AddCertificateDirFlag(cmd.Flags(), &opts.CertificatesDir)
	options.AddConfigFlag(cmd.Flags(), &opts.ConfigPath)
	return cmd
}

func RunCerts(opts *certsOptions) error {
	data, err := parseConfig(opts.ConfigPath)
	if err != nil {
		return err
	}
	caMap, crts := filterCAAndCerts(data)
	for _, current := range caMap {
		if err = certsphase.CreateCACertAndKeyFiles(current, opts.CertificatesDir); err != nil {
			return err
		}
	}
	for _, current := range crts {
		if err = certsphase.CreateCertAndKeyFilesWithCA(current, caMap[current.CAName], opts.CertificatesDir); err != nil {
			return err
		}
	}
	return nil
}

func parseConfig(file string) ([]*certsphase.CmdAdminCert, error) {
	var certsCfg []*certsphase.CmdAdminCert
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(data, &certsCfg)
	if err != nil {
		return nil, err
	}
	return certsCfg, nil
}

func filterCAAndCerts(data []*certsphase.CmdAdminCert) (map[string]*certsphase.CmdAdminCert, []*certsphase.CmdAdminCert) {
	ca := make(map[string]*certsphase.CmdAdminCert)
	cert := make([]*certsphase.CmdAdminCert, 0)
	for _, currentCert := range data {
		if currentCert.CAName == "" {
			ca[currentCert.Name] = currentCert
		} else {
			cert = append(cert, currentCert)
		}
	}
	return ca, cert
}
