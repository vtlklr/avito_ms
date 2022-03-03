package main

import (
	"net/http"
	"strconv"
	"strings"
)

//CreateAccount Для создания аккаунтов
var CreateAccount = func(w http.ResponseWriter, r *http.Request) {

	account := &Account{}
	account.Balance = 0

	resp := account.Create()

	Respond(w, resp)
}

//GetBalance возвращает баланс аккаунта
var GetBalance = func(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(r.URL.Query().Get("user_id"))
	if err != nil {
		resp := Message(false, "error")
		resp["err"] = "не корректно введен id"
		Respond(w, resp)
		return
	}
	err1, data := ReturnBalance(uint(id))
	if err1 != nil {
		resp := Message(false, "error")
		resp["err"] = err1.Error()
		Respond(w, resp)
		return
	}
	currency := r.URL.Query().Get("currency")
	if currency != "" || currency != "RUB" {
		currency = strings.ToUpper(currency)
		ex, err2 := Exchange(currency)
		if err2 != nil {
			resp := Message(false, "error currency")
			resp["err"] = err2.Error()
			Respond(w, resp)
			return
		}
		resp := Message(true, "success")

		resp["balance"] = float64(data) * ex
		resp["currency"] = currency
		Respond(w, resp)
		return
	}

	resp := Message(true, "success")
	resp["balance"] = data
	Respond(w, resp)
}

var CreditMoney = func(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(r.URL.Query().Get("user_id"))
	if err != nil {
		resp := Message(false, "error")
		resp["err"] = "не корректно введен user_id"
		Respond(w, resp)
		return
	}
	err1, _ := ReturnBalance(uint(id))
	if err1 != nil {
		resp := Message(false, "error")
		resp["err"] = err1.Error()
		Respond(w, resp)
		return
	}
	sum, err2 := strconv.Atoi(r.URL.Query().Get("sum"))
	if err2 != nil {
		resp := Message(false, "error")
		resp["err"] = "не корректно введена sum"
		Respond(w, resp)
		return
	}
	data := CreditMoneyFor(uint(id), sum)

	//resp := Message(true, "success")
	//resp["data"] = data
	Respond(w, data)
}
var DebitMoney = func(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(r.URL.Query().Get("user_id"))
	if err != nil {
		resp := Message(false, "error")
		resp["err"] = "не корректно введен user_id"
		Respond(w, resp)
		return
	}
	err1, balance := ReturnBalance(uint(id))
	if err1 != nil {
		resp := Message(false, "error")
		resp["err"] = err1.Error()
		Respond(w, resp)
		return
	}
	sum, err2 := strconv.Atoi(r.URL.Query().Get("sum"))
	if err2 != nil {
		resp := Message(false, "error")
		resp["err"] = "не корректно введена sum"
		Respond(w, resp)
		return
	}
	if sum <= 0 {
		resp := Message(false, "error")
		resp["err"] = "сумма должна быть больше 0"
		Respond(w, resp)
		return
	}
	target := r.URL.Query().Get("target")
	if target == "" {
		resp := Message(false, "error")
		resp["err"] = "не указана цель покупки(корзина, товар)"
		Respond(w, resp)
		return
	}
	if balance < sum {
		resp := Message(false, "error")
		resp["err"] = "не достаточно денег на балансе"
		Respond(w, resp)

	} else {
		data := DebitMoneyFrom(uint(id), sum, target)

		Respond(w, data)
	}
}
var GetTransacts = func(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(r.URL.Query().Get("user_id"))
	if err != nil {
		resp := Message(false, "error")
		resp["err"] = "не корректно введен user_id"
		Respond(w, resp)
		return
	}
	err1, _ := ReturnBalance(uint(id))
	if err1 != nil {
		resp := Message(false, "error")
		resp["err"] = err1.Error()
		Respond(w, resp)
		return
	}
	data := GetTransactsFor(uint(id))
	resp := Message(true, "success")
	resp["data"] = data
	Respond(w, resp)
}
var TransferMoney = func(w http.ResponseWriter, r *http.Request) {

	idFrom, errId1 := strconv.Atoi(r.URL.Query().Get("user_id"))
	if errId1 != nil {
		resp := Message(false, "error")
		resp["err"] = "не корректно введен user_id"
		Respond(w, resp)
		return
	}
	idTo, errId2 := strconv.Atoi(r.URL.Query().Get("user_id_to"))
	if errId2 != nil {
		resp := Message(false, "error")
		resp["err"] = "не корректно введен user_id_to"
		Respond(w, resp)
		return
	}
	sum, errSum := strconv.Atoi(r.URL.Query().Get("sum"))
	if errSum != nil {
		resp := Message(false, "error")
		resp["err"] = "не корректно введена sum"
		Respond(w, resp)
		return
	}
	if sum <= 0 {
		resp := Message(false, "error")
		resp["err"] = "сумма должна быть больше 0"
		Respond(w, resp)
		return
	}
	//проверим что аккаунт есть и узнаем баланс
	err, balanceIdFrom := ReturnBalance(uint(idFrom))
	if err != nil {
		resp := Message(false, "error user_id")
		resp["err"] = err.Error()
		Respond(w, resp)
		return
	}
	//проверим что аккаунт на который зачисляются деньги существует
	if err, _ := ReturnBalance(uint(idTo)); err != nil {
		resp := Message(false, "error user_id_to")
		resp["err"] = err.Error()
		Respond(w, resp)
		return
	}
	//проверим что денег на балансе достаточно для списания
	if balanceIdFrom < sum {
		resp := Message(false, "error")
		resp["err"] = "не достаточно средств для списания"
		Respond(w, resp)
	} else {

		data := DebitMoneyFromTo(uint(idFrom), uint(idTo), sum)

		Respond(w, data)
	}

}
