package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"strconv"
	"time"
)

type Account struct {
	gorm.Model
	Balance int `json:"balance"`
}

//Create создает аккаунт в бд
func (account *Account) Create() map[string]interface{} {

	GetDB().Create(account)
	response := Message(true, "Аккаунт создан")
	response["id"] = account.ID

	return response
}

//CreditMoneyFor зачисляет деньги на баланс аккаунта
func CreditMoneyFor(u uint, sum int) (error, map[string]interface{}) {
	tx := GetDB().Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		response := Message(false, err.Error())
		tx.Rollback()
		return err, response
	}

	var acc Account
	if err := tx.Table("accounts").Where("id = ?", u).First(&acc).Error; err != nil {
		tx.Rollback()
		response := Message(false, "error")
		return err, response
	}
	bal := acc.Balance + sum
	if err := tx.Model(&acc).Update("balance", bal).Error; err != nil {
		tx.Rollback()
		response := Message(false, "error")
		return err, response
	}
	now := time.Now()
	trans := &Transact{Operation: "пополнение", Sum: sum, Date: now.Format("02-01-2006 15:04:05"), AccountId: u}
	err, transactId := trans.Create()
	if err != nil {
		tx.Rollback()
		response := Message(false, "error")
		return err, response
	}
	response := Message(true, "success")
	response["transact_id"] = transactId
	tx.Commit()
	return nil, response
}

//DebitMoneyFromTo переводит деньги с баланса одного аккаунта другому
func DebitMoneyFromTo(u, u2 uint, sum int) (error, map[string]interface{}) {
	tx := GetDB().Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		response := Message(false, "error")
		return err, response
	}
	var acc1, acc2 Account
	if err := tx.Table("accounts").Where("id = ?", u).First(&acc1).Error; err != nil {
		tx.Rollback()
		response := Message(false, "error user_id")
		return err, response
	}
	if acc1.Balance < sum {
		tx.Rollback()
		response := Message(false, "error")
		return fmt.Errorf("не достаточно средств"), response
	}
	bal1 := acc1.Balance - sum
	if err := tx.Model(&acc1).Update("balance", bal1).Error; err != nil {
		tx.Rollback()
		response := Message(false, "error")
		return err, response
	}
	if err := tx.Table("accounts").Where("id = ?", u2).First(&acc2).Error; err != nil {
		tx.Rollback()
		response := Message(false, "error user_id_to")
		return err, response
	}
	bal2 := acc2.Balance + sum
	if err := tx.Model(&acc2).Update("balance", bal2).Error; err != nil {
		tx.Rollback()
		response := Message(false, "error")
		return err, response
	}
	now := time.Now()
	trans1 := &Transact{Operation: "перевод(списание) для аккаунта id: " + strconv.Itoa(int(u2)), Sum: sum, Date: now.Format("02-01-2006 15:04:05"), AccountId: u}
	err1, transactId1 := trans1.Create()
	if err1 != nil {
		tx.Rollback()
		response := Message(false, "error")

		return err1, response
	}
	trans2 := &Transact{Operation: "перевод(пополнение) от аккаунта id: " + strconv.Itoa(int(u)), Sum: sum, Date: now.Format("02-01-2006 15:04:05"), AccountId: u2}
	err2, transactId2 := trans2.Create()
	if err2 != nil {
		tx.Rollback()
		response := Message(false, "error")

		return err2, response
	}
	response := Message(true, "success")
	response["transact_id_user"] = transactId1
	response["transact_id_user_to"] = transactId2
	tx.Commit()
	return nil, response
}

//DebitMoneyFrom списывает деньги с баланса аккаунта
func DebitMoneyFrom(u uint, sum int, target string) (error, map[string]interface{}) {
	tx := GetDB().Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		response := Message(false, "error")
		return err, response
	}

	var acc Account
	if err := tx.Table("accounts").Where("id = ?", u).First(&acc).Error; err != nil {
		tx.Rollback()
		response := Message(false, "error")
		return err, response
	}
	if acc.Balance < sum {
		tx.Rollback()
		response := Message(false, "error")
		return fmt.Errorf("не достаточно средств"), response
	}
	bal := acc.Balance - sum
	if err := tx.Model(&acc).Update("balance", bal).Error; err != nil {
		tx.Rollback()
		response := Message(false, "error")
		return err, response
	}
	now := time.Now()
	trans := &Transact{Operation: "списание на :" + target, Sum: sum, Date: now.Format("02-01-2006 15:04:05"), AccountId: u}
	err, transactId := trans.Create()
	if err != nil {
		tx.Rollback()
		response := Message(false, "error")
		return err, response
	}
	response := Message(true, "success")
	response["transact_id"] = transactId
	tx.Commit()
	return nil, response
}

//ReturnBalance возвращает баланс аккаунта
func ReturnBalance(u uint) (error, int) {
	tx := GetDB().Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		//response := Message(false, err.Error())
		return err, 0
	}
	acc := &Account{}
	if err := tx.Table("accounts").Where("id = ?", u).First(acc).Error; err != nil {
		tx.Rollback()

		return err, 0
	}

	return nil, acc.Balance

}
