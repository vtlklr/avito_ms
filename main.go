package main

import (
	"fmt"
	"github.com/gorilla/mux"

	"net/http"
	"os"
)

func main() {

	router := mux.NewRouter()

	router.HandleFunc("/api/account/new", CreateAccount).Methods("GET")
	router.HandleFunc("/api/account/balance", GetBalance).Methods("GET")
	router.HandleFunc("/api/account/credit", CreditMoney).Methods("GET")
	router.HandleFunc("/api/account/transacts", GetTransacts).Methods("GET")
	router.HandleFunc("/api/account/transfer", TransferMoney).Methods("GET")
	router.HandleFunc("/api/account/debit", DebitMoney).Methods("GET")

	//router.Use(app.JwtAuthentication) //attach JWT auth middleware

	//router.NotFoundHandler = app.NotFoundHandler

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000" //localhost
	}

	fmt.Println(port)

	err := http.ListenAndServe(":"+port, router) //Launch the app, visit localhost:8000/api
	if err != nil {
		fmt.Print(err)
	}
}
