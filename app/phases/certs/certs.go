package certs

import (
	"crypto"
	"crypto/x509"
	"fmt"
	"github.com/pkg/errors"
	"github.com/x893675/certadmin/app/util/pkiutil"
	"k8s.io/klog/v2"
	"sync"
)

var (
	// certPeriodValidation is used to store if period validation was done for a certificate
	certPeriodValidationMutex sync.Mutex
	certPeriodValidation      = map[string]struct{}{}
)

// CreateCACertAndKeyFiles generates and writes out a given certificate authority.
// The certSpec should be one of the variables from this package.
func CreateCACertAndKeyFiles(certSpec *CmdAdminCert, certificatesDir string) error {
	if certSpec.CAName != "" {
		return errors.Errorf("this function should only be used for CAs, but cert %s has CA %s", certSpec.Name, certSpec.CAName)
	}
	klog.V(1).Infof("creating a new certificate authority for %s", certSpec.Name)

	caCert, caKey, err := pkiutil.NewCertificateAuthority(&certSpec.Config)
	if err != nil {
		return err
	}

	return writeCertificateAuthorityFilesIfNotExist(
		certificatesDir,
		certSpec.BaseName,
		caCert,
		caKey,
	)
}

// NewCSR will generate a new CSR and accompanying key
func NewCSR(certSpec *CmdAdminCert) (*x509.CertificateRequest, crypto.Signer, error) {
	//certConfig, err := certSpec.GetConfig(cfg)
	//if err != nil {
	//	return nil, nil, errors.Wrap(err, "failed to retrieve cert configuration")
	//}
	return pkiutil.NewCSRAndKey(&certSpec.Config)
}

// CreateCSR creates a certificate signing request
func CreateCSR(certSpec *CmdAdminCert, path string) error {
	csr, key, err := NewCSR(certSpec)
	if err != nil {
		return err
	}
	return writeCSRFilesIfNotExist(path, certSpec.BaseName, csr, key)
}

// CreateCertAndKeyFilesWithCA loads the given certificate authority from disk, then generates and writes out the given certificate and key.
// The certSpec and caCertSpec should both be one of the variables from this package.
func CreateCertAndKeyFilesWithCA(certSpec *CmdAdminCert, caCertSpec *CmdAdminCert, certificatesDir string) error {
	if certSpec.CAName != caCertSpec.Name {
		return errors.Errorf("expected CAname for %s to be %q, but was %s", certSpec.Name, certSpec.CAName, caCertSpec.Name)
	}

	caCert, caKey, err := LoadCertificateAuthority(certificatesDir, caCertSpec.BaseName)
	if err != nil {
		return errors.Wrapf(err, "couldn't load CA certificate %s", caCertSpec.Name)
	}

	return certSpec.CreateFromCA(certificatesDir, caCert, caKey)
}

// LoadCertificateAuthority tries to load a CA in the given directory with the given name.
func LoadCertificateAuthority(pkiDir string, baseName string) (*x509.Certificate, crypto.Signer, error) {
	// Checks if certificate authority exists in the PKI directory
	if !pkiutil.CertOrKeyExist(pkiDir, baseName) {
		return nil, nil, errors.Errorf("couldn't load %s certificate authority from %s", baseName, pkiDir)
	}

	// Try to load certificate authority .crt and .key from the PKI directory
	caCert, caKey, err := pkiutil.TryLoadCertAndKeyFromDisk(pkiDir, baseName)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "failure loading %s certificate authority", baseName)
	}
	// Validate period
	CheckCertificatePeriodValidity(baseName, caCert)

	// Make sure the loaded CA cert actually is a CA
	if !caCert.IsCA {
		return nil, nil, errors.Errorf("%s certificate is not a certificate authority", baseName)
	}

	return caCert, caKey, nil
}

// writeCertificateAuthorityFilesIfNotExist write a new certificate Authority to the given path.
// If there already is a certificate file at the given path; kubeadm tries to load it and check if the values in the
// existing and the expected certificate equals. If they do; kubeadm will just skip writing the file as it's up-to-date,
// otherwise this function returns an error.
func writeCertificateAuthorityFilesIfNotExist(pkiDir string, baseName string, caCert *x509.Certificate, caKey crypto.Signer) error {

	// If cert or key exists, we should try to load them
	if pkiutil.CertOrKeyExist(pkiDir, baseName) {

		// Try to load .crt and .key from the PKI directory
		caCert, _, err := pkiutil.TryLoadCertAndKeyFromDisk(pkiDir, baseName)
		if err != nil {
			return errors.Wrapf(err, "failure loading %s certificate", baseName)
		}
		// Validate period
		CheckCertificatePeriodValidity(baseName, caCert)

		// Check if the existing cert is a CA
		if !caCert.IsCA {
			return errors.Errorf("certificate %s is not a CA", baseName)
		}

		// kubeadm doesn't validate the existing certificate Authority more than this;
		// Basically, if we find a certificate file with the same path; and it is a CA
		// kubeadm thinks those files are equal and doesn't bother writing a new file
		fmt.Printf("[certs] Using the existing %q certificate and key\n", baseName)
	} else {
		// Write .crt and .key files to disk
		fmt.Printf("[certs] Generating %q certificate and key\n", baseName)

		if err := pkiutil.WriteCertAndKey(pkiDir, baseName, caCert, caKey); err != nil {
			return errors.Wrapf(err, "failure while saving %s certificate and key", baseName)
		}
	}
	return nil
}

// writeCertificateFilesIfNotExist write a new certificate to the given path.
// If there already is a certificate file at the given path; kubeadm tries to load it and check if the values in the
// existing and the expected certificate equals. If they do; kubeadm will just skip writing the file as it's up-to-date,
// otherwise this function returns an error.
func writeCertificateFilesIfNotExist(pkiDir string, baseName string, signingCert *x509.Certificate, cert *x509.Certificate, key crypto.Signer, cfg *pkiutil.CertConfig) error {

	// Checks if the signed certificate exists in the PKI directory
	if pkiutil.CertOrKeyExist(pkiDir, baseName) {
		// Try to load signed certificate .crt and .key from the PKI directory
		signedCert, _, err := pkiutil.TryLoadCertAndKeyFromDisk(pkiDir, baseName)
		if err != nil {
			return errors.Wrapf(err, "failure loading %s certificate", baseName)
		}
		// Validate period
		CheckCertificatePeriodValidity(baseName, signedCert)

		// Check if the existing cert is signed by the given CA
		if err := signedCert.CheckSignatureFrom(signingCert); err != nil {
			return errors.Errorf("certificate %s is not signed by corresponding CA", baseName)
		}

		// Check if the certificate has the correct attributes
		if err := validateCertificateWithConfig(signedCert, baseName, cfg); err != nil {
			return err
		}

		fmt.Printf("[certs] Using the existing %q certificate and key\n", baseName)
	} else {
		// Write .crt and .key files to disk
		fmt.Printf("[certs] Generating %q certificate and key\n", baseName)

		if err := pkiutil.WriteCertAndKey(pkiDir, baseName, cert, key); err != nil {
			return errors.Wrapf(err, "failure while saving %s certificate and key", baseName)
		}
		if pkiutil.HasServerAuth(cert) {
			fmt.Printf("[certs] %s serving cert is signed for DNS names %v and IPs %v\n", baseName, cert.DNSNames, cert.IPAddresses)
		}
	}

	return nil
}

// writeCSRFilesIfNotExist writes a new CSR to the given path.
// If there already is a CSR file at the given path; kubeadm tries to load it and check if it's a valid certificate.
// otherwise this function returns an error.
func writeCSRFilesIfNotExist(csrDir string, baseName string, csr *x509.CertificateRequest, key crypto.Signer) error {
	if pkiutil.CSROrKeyExist(csrDir, baseName) {
		_, _, err := pkiutil.TryLoadCSRAndKeyFromDisk(csrDir, baseName)
		if err != nil {
			return errors.Wrapf(err, "%s CSR existed but it could not be loaded properly", baseName)
		}

		fmt.Printf("[certs] Using the existing %q CSR\n", baseName)
	} else {
		// Write .key and .csr files to disk
		fmt.Printf("[certs] Generating %q key and CSR\n", baseName)

		if err := pkiutil.WriteKey(csrDir, baseName, key); err != nil {
			return errors.Wrapf(err, "failure while saving %s key", baseName)
		}

		if err := pkiutil.WriteCSR(csrDir, baseName, csr); err != nil {
			return errors.Wrapf(err, "failure while saving %s CSR", baseName)
		}
	}

	return nil
}

type certKeyLocation struct {
	pkiDir     string
	caBaseName string
	baseName   string
	uxName     string
}

// UsingExternalCA determines whether the user is relying on an external CA.  We currently implicitly determine this is the case
// when the CA Cert is present but the CA Key is not.
// This allows us to, e.g., skip generating certs or not start the csr signing controller.
// In case we are using an external front-proxy CA, the function validates the certificates signed by front-proxy CA that should be provided by the user.
//func UsingExternalCA(cfg *apis.ClusterConfiguration) (bool, error) {
//
//	if err := validateCACert(certKeyLocation{cfg.CertificatesDir, constants.CACertAndKeyBaseName, "", "CA"}); err != nil {
//		return false, err
//	}
//
//	caKeyPath := filepath.Join(cfg.CertificatesDir, constants.CAKeyName)
//	if _, err := os.Stat(caKeyPath); !os.IsNotExist(err) {
//		return false, nil
//	}
//
//	if err := validateSignedCert(certKeyLocation{cfg.CertificatesDir, constants.CACertAndKeyBaseName, constants.APIServerCertAndKeyBaseName, "API server"}); err != nil {
//		return true, err
//	}
//
//	if err := validateSignedCert(certKeyLocation{cfg.CertificatesDir, constants.CACertAndKeyBaseName, constants.APIServerKubeletClientCertAndKeyBaseName, "API server kubelet client"}); err != nil {
//		return true, err
//	}
//
//	return true, nil
//}

// UsingExternalFrontProxyCA determines whether the user is relying on an external front-proxy CA.  We currently implicitly determine this is the case
// when the front proxy CA Cert is present but the front proxy CA Key is not.
// In case we are using an external front-proxy CA, the function validates the certificates signed by front-proxy CA that should be provided by the user.
//func UsingExternalFrontProxyCA(cfg *apis.ClusterConfiguration) (bool, error) {
//
//	if err := validateCACert(certKeyLocation{cfg.CertificatesDir, constants.FrontProxyCACertAndKeyBaseName, "", "front-proxy CA"}); err != nil {
//		return false, err
//	}
//
//	frontProxyCAKeyPath := filepath.Join(cfg.CertificatesDir, constants.FrontProxyCAKeyName)
//	if _, err := os.Stat(frontProxyCAKeyPath); !os.IsNotExist(err) {
//		return false, nil
//	}
//
//	if err := validateSignedCert(certKeyLocation{cfg.CertificatesDir, constants.FrontProxyCACertAndKeyBaseName, constants.FrontProxyClientCertAndKeyBaseName, "front-proxy client"}); err != nil {
//		return true, err
//	}
//
//	return true, nil
//}

// validateCACert tries to load a x509 certificate from pkiDir and validates that it is a CA
func validateCACert(l certKeyLocation) error {
	// Check CA Cert
	caCert, err := pkiutil.TryLoadCertFromDisk(l.pkiDir, l.caBaseName)
	if err != nil {
		return errors.Wrapf(err, "failure loading certificate for %s", l.uxName)
	}
	// Validate period
	CheckCertificatePeriodValidity(l.uxName, caCert)

	// Check if cert is a CA
	if !caCert.IsCA {
		return errors.Errorf("certificate %s is not a CA", l.uxName)
	}
	return nil
}

// validateCACertAndKey tries to load a x509 certificate and private key from pkiDir,
// and validates that the cert is a CA. Failure to load the key produces a warning.
func validateCACertAndKey(l certKeyLocation) error {
	if err := validateCACert(l); err != nil {
		return err
	}

	_, err := pkiutil.TryLoadKeyFromDisk(l.pkiDir, l.caBaseName)
	if err != nil {
		klog.Warningf("assuming external key for %s: %v", l.uxName, err)
	}
	return nil
}

// validateSignedCert tries to load a x509 certificate and private key from pkiDir and validates
// that the cert is signed by a given CA
func validateSignedCert(l certKeyLocation) error {
	// Try to load CA
	caCert, err := pkiutil.TryLoadCertFromDisk(l.pkiDir, l.caBaseName)
	if err != nil {
		return errors.Wrapf(err, "failure loading certificate authority for %s", l.uxName)
	}
	// Validate period
	CheckCertificatePeriodValidity(l.uxName, caCert)

	return validateSignedCertWithCA(l, caCert)
}

// validateSignedCertWithCA tries to load a certificate and validate it with the given caCert
func validateSignedCertWithCA(l certKeyLocation, caCert *x509.Certificate) error {
	// Try to load key and signed certificate
	signedCert, _, err := pkiutil.TryLoadCertAndKeyFromDisk(l.pkiDir, l.baseName)
	if err != nil {
		return errors.Wrapf(err, "failure loading certificate for %s", l.uxName)
	}
	// Validate period
	CheckCertificatePeriodValidity(l.uxName, signedCert)

	// Check if the cert is signed by the CA
	if err := signedCert.CheckSignatureFrom(caCert); err != nil {
		return errors.Wrapf(err, "certificate %s is not signed by corresponding CA", l.uxName)
	}
	return nil
}

// validatePrivatePublicKey tries to load a private key from pkiDir
func validatePrivatePublicKey(l certKeyLocation) error {
	// Try to load key
	_, _, err := pkiutil.TryLoadPrivatePublicKeyFromDisk(l.pkiDir, l.baseName)
	if err != nil {
		return errors.Wrapf(err, "failure loading key for %s", l.uxName)
	}
	return nil
}

// validateCertificateWithConfig makes sure that a given certificate is valid at
// least for the SANs defined in the configuration.
func validateCertificateWithConfig(cert *x509.Certificate, baseName string, cfg *pkiutil.CertConfig) error {
	for _, dnsName := range cfg.AltNames.DNSNames {
		if err := cert.VerifyHostname(dnsName); err != nil {
			return errors.Wrapf(err, "certificate %s is invalid", baseName)
		}
	}
	for _, ipAddress := range cfg.AltNames.IPs {
		if err := cert.VerifyHostname(ipAddress.String()); err != nil {
			return errors.Wrapf(err, "certificate %s is invalid", baseName)
		}
	}
	return nil
}

// CheckCertificatePeriodValidity takes a certificate and prints a warning if its period
// is not valid related to the current time. It does so only if the certificate was not validated already
// by keeping track with a cache.
func CheckCertificatePeriodValidity(baseName string, cert *x509.Certificate) {
	certPeriodValidationMutex.Lock()
	if _, exists := certPeriodValidation[baseName]; exists {
		certPeriodValidationMutex.Unlock()
		return
	}
	certPeriodValidation[baseName] = struct{}{}
	certPeriodValidationMutex.Unlock()

	klog.V(5).Infof("validating certificate period for %s certificate", baseName)
	if err := pkiutil.ValidateCertPeriod(cert, 0); err != nil {
		klog.Warningf("WARNING: could not validate bounds for certificate %s: %v", baseName, err)
	}
}
