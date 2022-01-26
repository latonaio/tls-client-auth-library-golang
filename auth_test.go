package tls_auth_test

import (
	"crypto/tls"
	"crypto/x509"
	"net/http"
	"os"
	"reflect"
	"testing"

	tls_auth "github.com/latonaio/tls-client-auth-library-golang"
)

func TestNewHTTPSClient(t *testing.T) {
	// PKCS #12 形式 鍵ペア (証明書 + 秘密鍵)
	p12, err := os.ReadFile("./testcerts/leaf.p12")
	if err != nil {
		t.Fatal(err)
	}

	// PKCS #8 PEM 形式 証明書
	cert, err := os.ReadFile("./testcerts/leaf.x509.pem")
	if err != nil {
		t.Fatal(err)
	}

	// PKCS #8 PEM 形式 秘密鍵
	key, err := os.ReadFile("./testcerts/leaf.key")
	if err != nil {
		t.Fatal(err)
	}

	// TLS ライブラリ用鍵ペア
	keyPair, err := tls.X509KeyPair(cert, key)
	if err != nil {
		t.Fatal(err)
	}
	keyPair.Leaf, err = x509.ParseCertificate(keyPair.Certificate[0])
	if err != nil {
		t.Fatal(err)
	}

	type args struct {
		p12File    []byte
		password   string
		rootCACert []byte
	}
	tests := []struct {
		name    string
		args    args
		want    *http.Client
		wantErr bool
	}{
		{
			name: "cert without root",
			args: args{
				p12File:    p12,
				password:   "password-for-test",
				rootCACert: nil,
			},
			want: &http.Client{
				Transport: &http.Transport{
					TLSClientConfig: &tls.Config{
						Certificates: []tls.Certificate{keyPair},
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tls_auth.NewHTTPSClient(tt.args.p12File, tt.args.password, tt.args.rootCACert)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewHTTPSClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewHTTPSClient() = %v, want %v", got, tt.want)
			}
		})
	}
}
