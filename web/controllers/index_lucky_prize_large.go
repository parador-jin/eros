package controllers

import (
	"Eros/comm"
	"Eros/models"
)

func (c *IndexController) prizeLarge(
	ip string, loginUser *models.ObjLoginuser, userInfo *models.LtUser, blackIPInfo *models.LtBlackip) {
	nowTime := comm.NowUnix()
	blackTime := 30 * 86400
	// 更新用户的黑名单信息
	if userInfo == nil || userInfo.Id < 0 {
		userInfo = &models.LtUser{
			Id:         loginUser.Uid,
			Username:   loginUser.Username,
			Blacktime:  nowTime + blackTime,
			SysCreated: nowTime,
			SysUpdated: 0,
			SysIp:      ip,
		}
		c.ServiceUser.Create(userInfo)
	} else {
		userInfo.Id = loginUser.Uid
		userInfo.Blacktime = nowTime + blackTime
		userInfo.SysUpdated = nowTime
		c.ServiceUser.Update(userInfo, nil)
	}
	// 更新IP的黑名单信息
	if blackIPInfo == nil || blackIPInfo.Id < 0 {
		blackIPInfo = &models.LtBlackip{
			Ip:         ip,
			Blacktime:  nowTime + blackTime,
			SysCreated: nowTime,
		}
		c.ServiceBlackip.Create(blackIPInfo)
	} else {
		blackIPInfo.Blacktime = nowTime + blackTime
		blackIPInfo.SysUpdated = nowTime
		c.ServiceBlackip.Update(blackIPInfo, nil)
	}
}
