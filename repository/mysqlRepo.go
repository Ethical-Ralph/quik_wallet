package repository

import "gorm.io/gorm"

type MySqlRepository struct {
	conn *gorm.DB
}

func NewMySQLRepository(conn *gorm.DB) *MySqlRepository {
	return &MySqlRepository{conn}
}

func (mrepo *MySqlRepository) CreateWallet() {

}
