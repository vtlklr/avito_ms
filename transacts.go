//transacts.go функции работы с транзакциями
package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
)

type Transact struct {
	gorm.Model
	Operation string `json:"operation"`
	Sum       int    `json:"sum"`
	Date      string `json:"date"`
	Comment   string `json:"comment"`
	AccountId uint   `json:"account_id"` //The account that this transact belongs to
}

//Create создает новую транзакцию
func (transact *Transact) Create() (error, uint) {
	tx := GetDB().Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		//response := Message(false, "error bd")
		return err, 0
	}

	if err := tx.Create(transact).Error; err != nil {
		tx.Rollback()
		//response := Message(false, "error bd")
		return err, 0
	}
	tx.Commit()
	return nil, transact.ID
}

//GetTrancacts возвращает все транзакции по указанному id аккаунта
func GetTransactsFor(account uint) []*Transact {

	transacts := make([]*Transact, 0)
	err := GetDB().Table("transacts").Where("account_id = ?", account).Find(&transacts).Error
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return transacts
}
