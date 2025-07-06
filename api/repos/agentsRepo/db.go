package agentsRepo

import (
	"api/db"
	"api/geoip"
	"api/models"
	"time"
)

func Get(agentID string) (agent models.Agent, err error) {
	return agent, db.ORM.Where("id = ?", agentID).First(&agent).Error
}

func GetMultiple() (agents []models.Agent, err error) {
	return agents, db.ORM.Find(&agents).Error
}

func GetMultipleByCountryCode(ccs []string) (agents []models.Agent, err error) {
	return agents, db.ORM.Where("country_code IN ?", ccs).Find(&agents).Error
}

func GetMultipleUnknownOrigins() (agents []models.Agent, err error) {
	return agents, db.ORM.Where("country = ? OR country = ?", "", "unknown").Find(&agents).Error
}

func GetMultipleUnknownOriginsCount() (count int64, err error) {
	return count, db.ORM.Where("country = ? OR country = ?", "", "unknown").Model(&models.Agent{}).Count(&count).Error
}

func GetUniqueCountryCodes() (countryCodes []string, err error) {
	return countryCodes, db.ORM.Model(&models.Agent{}).
		Select("DISTINCT country_code").
		Where("country_code != ? AND country_code != ?", "", "unknown").
		Pluck("country_code", &countryCodes).Error
}

func Insert(agent *models.Agent) (err error) {
	return db.ORM.Create(&agent).Error
}

func UpdateLastSeen(agentID string) (err error) {
	agent, err := Get(agentID)
	if err != nil {
		return err
	}
	agent.LastSeen = time.Now().UnixNano()
	return db.ORM.Save(agent).Error
}

func UpdateIP(agentID, ip string) (err error) {
	agent, err := Get(agentID)
	if err != nil {
		return err
	}
	agent.IP = ip
	if len(ip) != 0 {
		country, code, err := geoip.IpToCountry(ip)
		if err == nil {
			agent.Country = country
			agent.CountryCode = code
		}
	}
	return db.ORM.Save(agent).Error
}
