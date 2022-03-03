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
	//resp := Message(true, "success")
	//resp["transact_id"] = transact.ID
	tx.Commit()
	return nil, transact.ID
}

func GetTransact(id uint) *Transact {

	transact := &Transact{}
	err := GetDB().Table("transact").Where("account_id = ?", id).First(transact).Error
	if err != nil {
		return nil
	}
	return transact
}

func GetTransactsFor(account uint) []*Transact {

	transacts := make([]*Transact, 0)
	err := GetDB().Table("transacts").Where("account_id = ?", account).Find(&transacts).Error
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return transacts
}
