package options

import "github.com/spf13/pflag"

// AddCertificateDirFlag adds the --certs-dir flag to the given flagset
func AddCertificateDirFlag(fs *pflag.FlagSet, certsDir *string) {
	fs.StringVar(certsDir, CertificatesDir, *certsDir, "The path where to save the certificates")
}

// AddCSRFlag adds the --csr-only flag to the given flagset
func AddCSRFlag(fs *pflag.FlagSet, csr *bool) {
	fs.BoolVar(csr, CSROnly, *csr, "Create CSRs instead of generating certificates")
}

// AddCSRDirFlag adds the --csr-dir flag to the given flagset
func AddCSRDirFlag(fs *pflag.FlagSet, csrDir *string) {
	fs.StringVar(csrDir, CSRDir, *csrDir, "The path to output the CSRs and private keys to")
}