package operatorsRepo

import (
	"api/db"
	"api/models"
)

func Get(username string) (operator models.Operator, err error) {
	return operator, db.ORM.Where("username = ?", username).First(&operator).Error
}

func GetMultiple() (operators []models.Operator, err error) {
	return operators, db.ORM.Find(&operators).Error
}

func Insert(operator *models.Operator) (err error) {
	return db.ORM.Create(&operator).Error
}
