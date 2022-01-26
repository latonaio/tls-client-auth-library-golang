package tls_auth

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"net/http"

	"golang.org/x/crypto/pkcs12"
)

func readPKCS12File(p12 []byte, password string) (*tls.Certificate, error) {
	privateKey, x509Cert, err := pkcs12.Decode(p12, password)
	if err != nil {
		return nil, err
	}

	keyPair := &tls.Certificate{
		Certificate: [][]byte{x509Cert.Raw},
		Leaf:        x509Cert,
		PrivateKey:  privateKey,
	}

	return keyPair, nil
}

func NewHTTPSClient(p12File []byte, password string, rootCACert []byte) (*http.Client, error) {
	// PKCS #12 形式のファイルから公開鍵と秘密鍵を取り出す
	keyPair, err := readPKCS12File(p12File, password)
	if err != nil {
		return nil, err
	}

	// サーバ証明書検証用の信頼できる Root CA を読み込む
	var rootCAs *x509.CertPool
	if rootCACert != nil {
		rootCAs = x509.NewCertPool()
		if !rootCAs.AppendCertsFromPEM(rootCACert) {
			return nil, errors.New("failed to append ca certificate")
		}
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				Certificates: []tls.Certificate{*keyPair},
				RootCAs:      rootCAs,
			},
		},
	}

	return client, nil
}
