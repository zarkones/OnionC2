package channelsRepo

import (
	"api/db"
	"api/models"
)

func Get(name string) (channel models.Channel, err error) {
	return channel, db.ORM.Where("name = ?", name).First(&channel).Error
}

func GetMultiple() (channels []models.Channel, err error) {
	return channels, db.ORM.Find(&channels).Error
}

func Insert(channel *models.Channel) (err error) {
	return db.ORM.Create(&channel).Error
}

func Delete(name string) (err error) {
	return db.ORM.Where("name = ?", name).Delete(models.Channel{}).Error
}
