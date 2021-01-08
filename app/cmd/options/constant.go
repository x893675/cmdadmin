package options


const (
	// CertificatesDir flag sets the path where to save and read the certificates.
	CertificatesDir = "cert-dir"

	// CfgPath flag sets the path to kubeadm config file.
	CfgPath = "config"


	// CSROnly flag instructs kubeadm to create CSRs instead of automatically creating or renewing certs
	CSROnly = "csr-only"

	// CSRDir flag sets the location for CSRs and flags to be output
	CSRDir = "csr-dir"
)