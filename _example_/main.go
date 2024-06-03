package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/root27/go-paytr"
	_ "github.com/root27/test-pay/docs"
	"github.com/swaggo/http-swagger"
	"log"
	"net/http"
	"os"
)

var (
	PORT = "6969"
)

// @title API Docs
// @description Payment API Documentation
// @version 0.1
// @host http://localhost:6969
// @BasePath /
func main() {

	r := mux.NewRouter()

	r.HandleFunc("/payment", handlePayment)

	r.HandleFunc("/paymentCallback", handleCallback)

	r.PathPrefix("/swagger").Handler(httpSwagger.WrapHandler)

	log.Println("Test server is running port ", PORT)

	log.Fatal(http.ListenAndServe(":"+PORT, r))

}

// @Payment godoc
// @Summary Payment request process
// @Description Request payment to get iframe token
// @Tags Payment
// @Accept json
// @Produce json
// @Failure 400 {object} HttpError
// @Param request body Request true "Request Body" example(Request{Data: []Cart{{Name: "test product", Price: 1000, Amount: 1}}})
// @Router /payment [post]
func handlePayment(w http.ResponseWriter, r *http.Request) {
	var req Request

	var basketData [][]any

	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {

		log.Println("Error parsing body: ", err)

		resp := HttpError{
			Code:    400,
			Message: "Bad Request",
		}

		json.NewEncoder(w).Encode(resp)

		return
	}

	for _, data := range req.Data {

		basketData = append(basketData, []any{
			data.Name,
			data.Price,
			data.Amount,
		})

	}

	var email, username, useraddress, userIP, userPhone, merchantOid string
	var totalAmount int

	p := paytr.Payment{

		MerchantID:    os.Getenv("merchantId"),
		MerchantKey:   os.Getenv("merchantKey"),
		MerchantSalt:  os.Getenv("merchantSalt"),
		UserIP:        userIP,
		MerchantOid:   merchantOid,
		Email:         email,
		TotalAmount:   totalAmount,
		Currency:      "TL",
		NoInstallment: 1,
		UserName:      username,
		UserAddress:   useraddress,
		UserPhone:     userPhone,
		OkUrl:         os.Getenv("okurl"),
		FailUrl:       os.Getenv("failurl"),
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

// @Tags Paytr Callback
// @Description Paytr Callback API (No request and response needed)
// @Router /paymentCallback [post]
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

		return
	}

	valid := payment.IsValid(os.Getenv("merchantKey"), os.Getenv("merchantSalt"))

	if !valid {

		//NOTE:Payment Hash not matched. Error handling
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

type HttpError struct {
	Code    int    `json:"code" example:"400"`
	Message string `json:"message" example:"bad request"`
}

// @description Request represents the request payload containing multiple cart items
type Request struct {
	Data []Cart `json:"data"`
}

// @description Cart represents a single item in the cart
type Cart struct {
	Name   string `json:"name" example:"test product"`
	Price  int    `json:"price" example:"1000"`
	Amount int    `json:"amount" example:"1"`
}
