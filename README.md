# tls-client-auth-library-golang

TLS クライアント認証を要求する Web サーバに接続するための Go 言語用ライブラリです。


## 動作環境

* Go 1.17 以上


## 使用方法

* go getで本ライブラリを取得します。

```sh
go get -u github.com/latonaio/tls-client-auth-library-golang
```

* 本ライブラリを使用するアプリケーション側に以下の import 文を追加します。

```go
import (
	tls_auth "github.com/latonaio/tls-client-auth-library-golang"
)
```

* test_auth.goの以下の部分で、鍵ペアファイルの読み込みと、サーバ証明書検証用 CA 証明書ファイルの読み込み（必要な場合）を行います。

```go

// PKCS #12 鍵ペアファイルの読み込み
keyPair, err := os.ReadFile("./keypair.p12")
if err != nil {
	// エラー処理
}

// サーバ証明書検証用 CA 証明書ファイルの読み込み (必要な場合)
rootCACert, err := os.ReadFile("./root.x509.pem")
if err != nil {
	// エラー処理
}
```

* 本ライブラリにより、HTTPSクライアントが作成されます。

```
// 本ライブラリによる HTTPS クライアントの作成
client, err := tls_auth.NewHTTPSClient(
	// PKCS #12 鍵ペアファイルの内容 ([]byte)
	keyPair,
	// PKCS #12 ファイルのパスワード (string)
	// パスワードはここに直接書かず、コンフィグファイルや環境変数などから読み込むことを推奨します
	os.Getenv("KEY_PAIR_PASSWORD"),
	// サーバ証明書検証用 CA 証明書ファイルの内容 ([]byte)
	// Go デフォルトの証明書ストアを使用する場合、nil を渡してください
	rootCACert,
)
if err != nil {
	// エラー処理
}
```

* 通常のHTTP(S)リクエストの場合と同様に、リクエストの送信を行なってください。  
なお、クライアントの設定のカスタマイズが必要な場合、ここで行なってください。

```
// 例: リクエストの送信
resp, err := client.Get("https://example.com/")
if err != nil {
	// エラー処理
}
```