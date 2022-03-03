package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type ExchangeResponse struct {
	Rates map[string]float64 `json:"rates"`
}

//Exchange для конвертации валюты в указанную
//используем сервис exchangeratesapi.io
//бесплатный тарифный план позволяет конвертировать только в базовую валюту EUR
//поэтому если валюта отличается от EUR мы делаем второй запрос с нужной нам валютой
//
//
func Exchange(currency string) (float64, error) {
	if len(currency) != 3 {
		return 0, fmt.Errorf("Неверно указана валюта")
	}
	urlBase := fmt.Sprintf("http://api.exchangeratesapi.io/v1/latest?access_key=98fbf5fe43eb337e647d2fa87b3ee2c2&symbols=%s", "RUB")
	respBase, err1 := http.Get(urlBase)
	defer respBase.Body.Close()

	if err1 != nil {
		return 0, err1
	}
	bodyBase, _ := ioutil.ReadAll(respBase.Body)
	exchBase := ExchangeResponse{}
	err2 := json.Unmarshal(bodyBase, &exchBase)
	if err2 != nil {
		return 0, err2
	}
	exBase := exchBase.Rates["RUB"]
	if currency == "EUR" {
		return 1 / exBase, nil
	}

	urlCurrency := fmt.Sprintf("http://api.exchangeratesapi.io/v1/latest?access_key=98fbf5fe43eb337e647d2fa87b3ee2c2&symbols=%s", currency)
	respCurrency, err3 := http.Get(urlCurrency)
	defer respCurrency.Body.Close()

	if err3 != nil {
		return 0, err3
	}

	bodyCurrency, _ := ioutil.ReadAll(respCurrency.Body)
	exchCurrency := ExchangeResponse{}
	err4 := json.Unmarshal(bodyCurrency, &exchCurrency)
	if err4 != nil {
		return 0, err4
	}
	exCurrency := exchCurrency.Rates[currency]
	currencyNow := exCurrency / exBase
	return currencyNow, nil
}
