package models

import "gorm.io/gorm"

type Wallet struct {
	gorm.Model
	WalletIdentifier string  `json:"walletIdentifier" gorm:"unique;not null;type:varchar(100)"`
	Balance          float64 `json:"balance" gorm:"not null;default:0"`
}
