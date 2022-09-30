package pki

import "crypto/x509"

func cloneCert(cert *x509.Certificate) *x509.Certificate {
	// These are the properties x509.CreateCertificate says it pays attention to
	return &x509.Certificate{
		AuthorityKeyId:              cert.AuthorityKeyId,
		BasicConstraintsValid:       cert.BasicConstraintsValid,
		CRLDistributionPoints:       cert.CRLDistributionPoints,
		DNSNames:                    cert.DNSNames,
		EmailAddresses:              cert.EmailAddresses,
		ExcludedDNSDomains:          cert.ExcludedDNSDomains,
		ExcludedEmailAddresses:      cert.ExcludedEmailAddresses,
		ExcludedIPRanges:            cert.ExcludedIPRanges,
		ExcludedURIDomains:          cert.ExcludedURIDomains,
		ExtKeyUsage:                 cert.ExtKeyUsage,
		ExtraExtensions:             cert.ExtraExtensions,
		IPAddresses:                 cert.IPAddresses,
		IsCA:                        cert.IsCA,
		IssuingCertificateURL:       cert.IssuingCertificateURL,
		KeyUsage:                    cert.KeyUsage,
		MaxPathLen:                  cert.MaxPathLen,
		MaxPathLenZero:              cert.MaxPathLenZero,
		NotAfter:                    cert.NotAfter,
		NotBefore:                   cert.NotBefore,
		OCSPServer:                  cert.OCSPServer,
		PermittedDNSDomains:         cert.PermittedDNSDomains,
		PermittedDNSDomainsCritical: cert.PermittedDNSDomainsCritical,
		PermittedEmailAddresses:     cert.PermittedEmailAddresses,
		PermittedIPRanges:           cert.PermittedIPRanges,
		PermittedURIDomains:         cert.PermittedURIDomains,
		PolicyIdentifiers:           cert.PolicyIdentifiers,
		SerialNumber:                cert.SerialNumber,
		SignatureAlgorithm:          cert.SignatureAlgorithm,
		Subject:                     cert.Subject,
		SubjectKeyId:                cert.SubjectKeyId,
		URIs:                        cert.URIs,
		UnknownExtKeyUsage:          cert.UnknownExtKeyUsage,
	}
}
