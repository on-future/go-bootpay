package go_bootpay

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var (
	httpClient = &http.Client{Timeout: 3 * time.Second}
)

type (
	paymentData struct {
		CardName   string `json:"card_name"`
		CardNo     string `json:"card_no"`
		CardQuota  string `json:"card_quota"`
		CardAuthNo string `json:"card_auth_no"`
		ReceiptId  string `json:"receipt_id"`
		N          string `json:"n"`
		P          int    `json:"p"`
		Tid        string `json:"tid"`
		Pg         string `json:"pg"`
		Pm         string `json:"pm"`
		PgA        string `json:"pg_a"`
		PmA        string `json:"pm_a"`
		OId        string `json:"o_id"`
		PAt        string `json:"p_at"`
		S          int    `json:"s"`
		G          int    `json:"g"`
	}
	accessTokenData struct {
		response
		Data struct {
			Token      string `json:"token"`
			ServerTime int    `json:"server_time"`
			ExpiredAt  int    `json:"expired_at"`
		}
	}
	receiptData struct {
		response
		Data struct {
			ReceiptId        string      `json:"receipt_id"`
			OrderId          string      `json:"order_id"`
			Name             string      `json:"name"`
			Price            int         `json:"price"`
			TaxFree          int         `json:"tax_free"`
			RemainPrice      int         `json:"remain_price"`
			RemainTaxFree    int         `json:"remain_tax_free"`
			CancelledPrice   int         `json:"cancelled_price"`
			CancelledTaxFree int         `json:"cancelled_tax_free"`
			ReceiptUrl       string      `json:"receipt_url"`
			Unit             string      `json:"unit"`
			Pg               string      `json:"pg"`
			Method           string      `json:"method"`
			PgName           string      `json:"pg_name"`
			MethodName       string      `json:"method_name"`
			PaymentData      paymentData `json:"payment_data"`
			RequestedAt      string      `json:"requested_at"`
			PurchasedAt      string      `json:"purchase_at"`
			Status           int         `json:"status"`
			StatusEn         string      `json:"status_en"`
			StatusKo         string      `json:"status_ko"`
		}
	}
	cancelData struct {
		response
		ReceiptId          string
		RequestCancelPrice int
		RemainPrice        int
		RemainTaxFree      int
		CancelledPrice     int
		CancelledTaxFree   int
		RevokedAt          string
		Tid                string
	}
	response struct {
		Status  int    `json:"status"`
		Code    int    `json:"code"`
		Message string `json:"message"`
	}
)

type Bootpay struct {
	ApplicationID string
	PrivateKey    string
}

func (c Bootpay) AccessToken() (*string, error) {
	f := url.Values{}
	f.Set("application_id", c.ApplicationID)
	f.Set("private_key", c.PrivateKey)

	req, err := http.NewRequest(http.MethodPost, "https://api.bootpay.co.kr/request/token", strings.NewReader(f.Encode()))
	if err != nil {
		return nil, err
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusOK {
		r := accessTokenData{}
		c.bind(resp, &r)

		if r.Status == http.StatusOK {
			return &r.Data.Token, nil
		}

		return nil, errors.New(fmt.Sprintf("unpected response status %d, got %d", http.StatusOK, r.Status))
	}

	return nil, errors.New(fmt.Sprintf("unexpected http status code %d, got %d", http.StatusOK, resp.StatusCode))
}

func (c Bootpay) Verify(token string, receipt string, price int) (*receiptData, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://api.bootpay.co.kr/receipt/%s", receipt), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", token)
	r, err := c.response(req)
	if err != nil {
		return nil, err
	}

	m := receiptData{}
	c.bind(r, &m)

	if m.Data.Status == 1 {
		if m.Data.Price == price {
			return &m, nil
		}

		return nil, errors.New(fmt.Sprintf("unpected price %d, got %d", price, m.Data.Price))
	}

	return nil, errors.New(fmt.Sprintf("unexpected http status code %d, got %d", http.StatusOK, m.Status))
}

func (c Bootpay) MustVerify(token string, receipt string, price int) *receiptData {
	verify, err := c.Verify(token, receipt, price)
	if err != nil {
		panic(err)
	}
	return verify
}

func (c Bootpay) Cancel(token string, receiptID string, name string, reason string) (*cancelData, error) {
	f := url.Values{}
	f.Set("receipt_id", receiptID)
	f.Set("name", name)
	f.Set("reason", reason)

	req, err := http.NewRequest(http.MethodGet, "https://api.bootpay.co.kr/cancel", strings.NewReader(f.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", token)
	r, err := c.response(req)
	if err != nil {
		return nil, err
	}

	m := cancelData{}
	c.bind(r, &m)

	return &m, nil
}

func (c Bootpay) MustCancel(token string, receiptID string, name string, reason string) *cancelData {
	cancel, err := c.Cancel(token, receiptID, name, reason)
	if err != nil {
		panic(err)
	}
	return cancel
}

func (c Bootpay) response(r *http.Request) (*http.Response, error) {
	resp, err := httpClient.Do(r)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusOK {
		return resp, nil
	}

	return nil, errors.New(fmt.Sprintf("unexpected http status code %d, got %d", http.StatusOK, resp.StatusCode))
}

func (c Bootpay) bind(resp *http.Response, bindPtr interface{}) {
	defer resp.Body.Close()
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(bytes, bindPtr)
	if err != nil {
		panic(err)
	}
}
