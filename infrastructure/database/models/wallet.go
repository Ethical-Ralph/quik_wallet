package models

import "gorm.io/gorm"

type Wallet struct {
	Id               uint    `gorm:"primary_key"`
	WalletIdentifier string  `json:"walletIdentifier" gorm:"unique;not null;type:varchar(100);default:null"`
	Balance          float64 `json:"balance" gorm:"not null;default:0"`
	gorm.Model
}
