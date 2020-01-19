go-bootpay
===========

golang 부트페이 API

#### 설치

	go get github.com/on-future/go-bootpay

#### 임포트

```go
import "github.com/on-future/go-bootpay"
```

#### 사용법

```go
bp := go_bootpay.Bootpay {
    ApplicationID: "<어플리케이션 아이디>",
    PrivateKey:    "<개인 키>",
}

at, _ := bp.AccessToken()
verify, _ := bp.Verify(*at, "<receipt id>", <price>)
cancel, _ := bp.Cancel(*at, "<receipt id>", "<name>", "<reason>")
```

#### LICENSE
MIT license.
