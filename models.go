package paytr

type Payment struct {
	MerchantID     string `schema:"merchant_id"`
	MerchantKey    string `schema:"merchant_key"`
	MerchantSalt   string `schema:"merchant_salt"`
	UserIP         string `schema:"user_ip"`
	MerchantOid    string `schema:"merchant_oid"`
	Email          string `schema:"email"`
	TotalAmount    int    `schema:"payment_amount"`
	Currency       string `schema:"currency"`
	Basket         string `schema:"user_basket"`
	NoInstallment  int    `schema:"no_installment"`
	MaxInstallment int    `schema:"max_installment"`
	PaytrToken     string `schema:"paytr_token"`
	UserName       string `schema:"user_name"`
	UserAddress    string `schema:"user_address"`
	UserPhone      string `schema:"user_phone"`
	OkUrl          string `schema:"merchant_ok_url"`
	FailUrl        string `schema:"merchant_fail_url"`
	TestMode       string `schema:"test_mode"`
	DebugOn        int    `schema:"debug_on"`
	Timeout        int    `schema:"timeout_limit"`
	Lang           string `schema:"lang"`
}

//NOTE: Iframe token response

type PaytrResponse struct {
	Status int
	Token  string
	Reason string
}

//NOTE: Paytr Callback

type CallbackRequest struct {
	InstallmentCount  int	`schema:"installment_count"`
	MerchantID        string `schema:"merchant_id"`
	MerchantOid       string `schema:"merchant_oid"`
	Status            string `schema:"status"`
	TotalAmount       int    `schema:"total_amount"`
	Hash              string `schema:"hash"`
	FailReasonCode    int    `schema:"failed_reason_code"`
	FailReasonMessage string `schema:"failed_reason_msg"`
	TestMode          string `schema:"test_mode"`
	PaymentType       string `schema:"payment_type"`
	Currency          string `schema:"currency"`
	PaymentAmount     int    `schema:"payment_amount"`
}
