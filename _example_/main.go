package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/schema"
	"github.com/root27/go-paytr"
)

var (
	PORT = "6969"
)

func main() {

	r := http.NewServeMux()

	r.HandleFunc("/payment", handlePayment)

	r.HandleFunc("/paymentCallback", handleCallback)

	log.Println("Test server is running port ", PORT)

	log.Fatal(http.ListenAndServe(":"+PORT, r))

}

func handlePayment(w http.ResponseWriter, r *http.Request) {
	var req Request

	var basketData [][]any

	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {

		log.Println("Error parsing body: ", err)

		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	for _, data := range req.Data {

		basketData = append(basketData, []any{
			data.Name,
			data.Price,
			data.Amount,
		})

	}

	p := paytr.Payment{

		MerchantID:    "272972",
		MerchantKey:   "fKbPJLPc2UtyUNet",
		MerchantSalt:  "gTp58ZLH7Wjhxrhi",
		UserIP:        "78.163.140.60",
		MerchantOid:   "test123",
		Email:         "test@test.com",
		TotalAmount:   10 * 100,
		Currency:      "TL",
		NoInstallment: 1,
		UserName:      "test",
		UserAddress:   "test address",
		UserPhone:     "123345567",
		OkUrl:         "https://apps.uniqgene.com/checkout",
		FailUrl:       "https://apps.uniqgene.com/",
		TestMode:      "1",
		DebugOn:       0,
		Timeout:       30,
		Lang:          "tr",
	}

	p.BasketConfig(basketData)

	p.GenerateToken(p.MerchantKey, p.MerchantSalt)

	token, err := p.GetIframe()

	if err != nil {

		log.Printf("Error fetching iframe: %s\n", err)

		return
	}
	//TODO: Return iframe token to client
	log.Println(token)
}

//NOTE: Paytr Callback API

func handleCallback(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()

	if err != nil {

		log.Printf("Error parsing request form:%s\n", err.Error())

		http.Error(w, "Bad request", http.StatusBadRequest)

		return

	}

	decoder := schema.NewDecoder()

	var payment paytr.CallbackRequest

	err = decoder.Decode(&payment, r.PostForm)

	if err != nil {

		log.Printf("Error decoding postform:%s\n", err.Error())

		http.Error(w, "Error decoding postform", http.StatusInternalServerError)

		return
	}

	valid := payment.IsValid(merchantKey, merchantSalt)

	if !valid {

		return

	}

	if payment.Status != "success" {

		log.Printf("Error payment process\n")

		w.Write([]byte("OK"))

		return
	}

	//NOTE:Payment Successfull

	w.Write([]byte("OK"))

}

type Request struct {
	Data []Cart `json:"data"`
}

type Cart struct {
	Name   string `json:"name"`
	Price  int    `json:"price"`
	Amount int    `json:"amount"`
}
