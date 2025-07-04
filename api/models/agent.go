package models

import (
	"api/geoip"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Agent struct {
	ID          string `json:"id" gorm:"primaryKey" sql:"type:text"`
	Hostname    string `json:"hostname"`
	OS          string `json:"os"`
	IP          string `json:"ip"`
	Country     string `json:"country"`
	CountryCode string `json:"countryCode"`
	RAM         string `json:"ram"`
	OSVersion   string `json:"osVersion"`
	CPUName     string `json:"cpuName"`
	Arch        string `json:"arch"`
	LastSeen    int64  `json:"lastSeen"`
}

func (agent *Agent) BeforeCreate(tx *gorm.DB) (err error) {
	agent.ID = uuid.New().String()
	if len(agent.IP) != 0 {
		country, code, err := geoip.IpToCountry(agent.IP)
		if err == nil {
			agent.Country = country
			agent.CountryCode = code
		}
	}
	return
}
