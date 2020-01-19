package go_bootpay_test

import (
	go_bootpay "github.com/on-future/go-bootpay"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Verify(t *testing.T) {
	bp := go_bootpay.Bootpay{
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
