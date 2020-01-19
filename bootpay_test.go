package go_bootpay_test

import (
	"github.com/stretchr/testify/assert"
	gobootpay "go-bootpay"
	"testing"
)

func Test_Verify(t *testing.T) {
	bp := gobootpay.Bootpay{
		ApplicationID: "",
		PrivateKey:    "",
	}

	at, _ := bp.AccessToken()
	assert.NotNil(t, at)

	verify, _ := bp.Verify(*at, "", 1000)
	assert.NotNil(t, verify)

	cancel, _ := bp.Cancel(*at, "", "", "")
	assert.NotNil(t, cancel)
}
