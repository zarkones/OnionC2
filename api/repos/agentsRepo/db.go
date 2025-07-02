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
	return agents, db.ORM.Where("ip = ? OR ip = ?", "", "unknown").Find(&agents).Error
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
