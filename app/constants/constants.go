package constants

import "time"

const (
	// CertificateValidity defines the validity for all the signed certificates generated by kubeadm
	CertificateValidity = time.Hour * 24 * 365
)