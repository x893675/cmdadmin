package certs

import (
	"crypto"
	"crypto/x509"
	"github.com/pkg/errors"
	"github.com/x893675/certadmin/app/util/pkiutil"
)

const (
	errInvalid = "invalid argument"
	errExist   = "file already exists"
)

type configMutatorsFunc func(*pkiutil.CertConfig) error

type CmdAdminCert struct {
	Name     string `yaml:"name"`
	LongName string `yaml:"-"`
	BaseName string `yaml:"baseName"`
	CAName   string `yaml:"caName"`
	// Some attributes will depend on the InitConfiguration, only known at runtime.
	// These functions will be run in series, passed both the InitConfiguration and a cert Config.
	configMutators []configMutatorsFunc `yaml:"-"`
	Config         pkiutil.CertConfig   `yaml:"cfg"`
}

// Certificates is a list of Certificates that Kubeadm should create.
type Certificates []*CmdAdminCert

// AsMap returns the list of certificates as a map, keyed by name.
func (c Certificates) AsMap() CertificateMap {
	certMap := make(map[string]*CmdAdminCert)
	for _, cert := range c {
		certMap[cert.Name] = cert
	}

	return certMap
}

// CertificateMap is a flat map of certificates, keyed by Name.
type CertificateMap map[string]*CmdAdminCert

// CertTree returns a one-level-deep tree, mapping a CA cert to an array of certificates that should be signed by it.
func (m CertificateMap) CertTree() (CertificateTree, error) {
	caMap := make(CertificateTree)

	for _, cert := range m {
		if cert.CAName == "" {
			if _, ok := caMap[cert]; !ok {
				caMap[cert] = []*CmdAdminCert{}
			}
		} else {
			ca, ok := m[cert.CAName]
			if !ok {
				return nil, errors.Errorf("certificate %q references unknown CA %q", cert.Name, cert.CAName)
			}
			caMap[ca] = append(caMap[ca], cert)
		}
	}

	return caMap, nil
}

// CertificateTree is represents a one-level-deep tree, mapping a CA to the certs that depend on it.
type CertificateTree map[*CmdAdminCert]Certificates

// CreateTree creates the CAs, certs signed by the CAs, and writes them all to disk.
func (t CertificateTree) CreateTree(certificatesDir string) error {
	for ca, leaves := range t {
		cfg := &ca.Config

		var caKey crypto.Signer

		caCert, err := pkiutil.TryLoadCertFromDisk(certificatesDir, ca.BaseName)
		if err == nil {
			// Validate period
			CheckCertificatePeriodValidity(ca.BaseName, caCert)

			// Cert exists already, make sure it's valid
			if !caCert.IsCA {
				return errors.Errorf("certificate %q is not a CA", ca.Name)
			}
			// Try and load a CA Key
			caKey, err = pkiutil.TryLoadKeyFromDisk(certificatesDir, ca.BaseName)
			if err != nil {
				// If there's no CA key, make sure every certificate exists.
				for _, leaf := range leaves {
					cl := certKeyLocation{
						pkiDir:   certificatesDir,
						baseName: leaf.BaseName,
						uxName:   leaf.Name,
					}
					if err := validateSignedCertWithCA(cl, caCert); err != nil {
						return errors.Wrapf(err, "could not load expected certificate %q or validate the existence of key %q for it", leaf.Name, ca.Name)
					}
				}
				continue
			}
			// CA key exists; just use that to create new certificates.
		} else {
			// CACert doesn't already exist, create a new cert and key.
			caCert, caKey, err = pkiutil.NewCertificateAuthority(cfg)
			if err != nil {
				return err
			}

			err = writeCertificateAuthorityFilesIfNotExist(
				certificatesDir,
				ca.BaseName,
				caCert,
				caKey,
			)
			if err != nil {
				return err
			}
		}

		for _, leaf := range leaves {
			if err := leaf.CreateFromCA(certificatesDir, caCert, caKey); err != nil {
				return err
			}
		}
	}
	return nil
}

// CreateFromCA makes and writes a certificate using the given CA cert and key.
func (k *CmdAdminCert) CreateFromCA(certificatesDir string, caCert *x509.Certificate, caKey crypto.Signer) error {
	cert, key, err := pkiutil.NewCertAndKey(caCert, caKey, &k.Config)
	if err != nil {
		return err
	}
	err = writeCertificateFilesIfNotExist(
		certificatesDir,
		k.BaseName,
		caCert,
		cert,
		key,
		&k.Config,
	)

	if err != nil {
		return errors.Wrapf(err, "failed to write or validate certificate %q", k.Name)
	}

	return nil
}
