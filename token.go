package paytr

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/gorilla/schema"
)

func (p *Payment) BasketConfig(cart [][]any) {

	cartbytes, _ := json.Marshal(cart)

	p.Basket = base64.StdEncoding.EncodeToString(cartbytes)

}

func (p *Payment) GenerateToken(merchantKey, merchantSalt string) string {

	hashedStr := p.MerchantID + p.UserIP + p.MerchantOid + p.Email + strconv.Itoa(p.TotalAmount) + p.Basket +
		strconv.Itoa(p.NoInstallment) + strconv.Itoa(p.MaxInstallment) + p.Currency + p.TestMode

	paytr_token := hashedStr + merchantSalt

	hmacToken := hmac.New(sha256.New, []byte(merchantKey))

	hmacToken.Write([]byte(paytr_token))

	p.PaytrToken = base64.StdEncoding.EncodeToString(hmacToken.Sum(nil))

	p.MerchantKey = merchantKey
	p.MerchantSalt = merchantSalt

	return p.PaytrToken

}

func (p *Payment) GetIframe() (PaytrResponse, error) {

	var response PaytrResponse

	formData := url.Values{}

	encoder := schema.NewEncoder()

	err := encoder.Encode(p, formData)

	log.Println(formData)

	if err != nil {

		log.Println("Error encode form data: ", err)

		return response, err

	}

	res, err := http.PostForm("https://www.paytr.com/odeme/api/get-token", formData)

	if err != nil {

		log.Println("Error requesting payment: ", err)

		return response, err

	}

	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)

	if err != nil {

		log.Println("Error reading body: ", err)

		return response, err

	}

	json.Unmarshal(resBody, &response)

	return response, nil

}
