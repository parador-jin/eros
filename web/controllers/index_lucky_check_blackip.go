package controllers

import (
	"Eros/models"
	"time"
)

func (c *IndexController) checkBlackIP(ip string) (bool, *models.LtBlackip) {
	info := c.ServiceBlackip.GetByIP(ip)
	if info == nil || info.Ip == "" {
		return true, nil
	}
	if info.Blacktime > int(time.Now().Unix()) {
		// IP黑名单存在，并且还在黑名单有效期内
		return false, info
	}
	return true, info
}
