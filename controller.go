//controller.go описаны методы работы с api
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
//api /api/account/balance
//user_id: id аккаунта
//currency: валюта, для отображения баланса. если не указано то в базовой валюте - рублях
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
	if currency != "" {
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

//CreditMoney зачисляет деньги на счет
//api /api/account/credit
//user_id: id аккаунта для зачисления денег
//sum: сумма
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
	Respond(w, data)
}

//DebitMoney списывает деньги со счета
//api /api/account/debit
//user_id: id аккаунта с которого списать деньги
//sum: сумма
//target: цель списания (корзина покупателя, товар...)
var DebitMoney = func(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(r.URL.Query().Get("user_id"))
	if err != nil {
		resp := Message(false, "error")
		resp["err"] = "не корректно введен user_id"
		Respond(w, resp)
		return
	}
	err1, balance := ReturnBalance(uint(id))
	//проверка что аккаунт существует
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

//GetTransact возвращает все транзакции по аккаунту
//api /api/account/transacts
//user_id: id аккаунта для которого запрашиваются транзакции
//
//
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

//TransferMoney переводит деньги с баланса одного пользователя другому
//api /api/account/transfer
//user_id: id аккаунта с которого списать деньги
//user_id_to: id аккаунта куда зачислить деньги
//sum: сумма
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
