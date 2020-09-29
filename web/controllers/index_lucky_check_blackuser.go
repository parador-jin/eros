package controllers

import (
	"Eros/models"
	"time"
)

func (c *IndexController) checkBlackUser(uid int) (bool, *models.LtUser) {
	info := c.ServiceUser.Get(uid)
	if info != nil && info.Blacktime > int(time.Now().Unix()) {
		// 黑名单存在并且有效
		return false, info
	}
	return true, info
}
