package cert

import (
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"io/ioutil"
	"net"
	"testing"

	"github.com/x893675/cmdadmin/pkg/util/keyutil"
)

func TestMakeCSR(t *testing.T) {
	keyFile := "testdata/dontUseThisKey.pem"
	subject := &pkix.Name{
		CommonName: "kube-worker",
	}
	dnsSANs := []string{"localhost"}
	ipSANs := []net.IP{net.ParseIP("127.0.0.1")}

	keyData, err := ioutil.ReadFile(keyFile)
	if err != nil {
		t.Fatal(err)
	}
	key, err := keyutil.ParsePrivateKeyPEM(keyData)
	if err != nil {
		t.Fatal(err)
	}
	csrPEM, err := MakeCSR(key, subject, dnsSANs, ipSANs)
	if err != nil {
		t.Error(err)
	}
	csrBlock, rest := pem.Decode(csrPEM)
	if csrBlock == nil {
		t.Fatal("Unable to decode MakeCSR result.")
	}
	if len(rest) != 0 {
		t.Error("Found more than one PEM encoded block in the result.")
	}
	if csrBlock.Type != CertificateRequestBlockType {
		t.Errorf("Found block type %q, wanted 'CERTIFICATE REQUEST'", csrBlock.Type)
	}
	csr, err := x509.ParseCertificateRequest(csrBlock.Bytes)
	if err != nil {
		t.Errorf("Found %v parsing MakeCSR result as a CertificateRequest.", err)
	}
	if csr.Subject.CommonName != subject.CommonName {
		t.Errorf("Wanted %v, got %v", subject, csr.Subject)
	}
	if len(csr.DNSNames) != 1 {
		t.Errorf("Wanted 1 DNS name in the result, got %d", len(csr.DNSNames))
	} else if csr.DNSNames[0] != dnsSANs[0] {
		t.Errorf("Wanted %v, got %v", dnsSANs[0], csr.DNSNames[0])
	}
	if len(csr.IPAddresses) != 1 {
		t.Errorf("Wanted 1 IP address in the result, got %d", len(csr.IPAddresses))
	} else if csr.IPAddresses[0].String() != ipSANs[0].String() {
		t.Errorf("Wanted %v, got %v", ipSANs[0], csr.IPAddresses[0])
	}
}
