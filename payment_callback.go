package paytr

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"strconv"
)

func (c *CallbackRequest) IsValid(merchantKey, merchantSalt string) bool {

	tokenStr := c.MerchantOid + merchantSalt + c.Status + strconv.Itoa(c.TotalAmount)

	tokenHmac := hmac.New(sha256.New, []byte(merchantKey))

	tokenHmac.Write([]byte(tokenStr))

	hash := base64.StdEncoding.EncodeToString(tokenHmac.Sum(nil))

	return hash == c.Hash

}
