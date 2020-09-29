package controllers

import (
	"Eros/comm"
	"Eros/conf"
	"Eros/models"
	"Eros/web/utils"
	"fmt"
	"log"
)

// http://localhost:8080/lucky
func (c *IndexController) GetLucky() map[string]interface{} {
	result := make(map[string]interface{}, 0)
	result["code"] = 0
	result["msg"] = ""
	// 1.验证登录用户
	loginUser := comm.GetLoginUser(c.Ctx.Request())
	if loginUser == nil || loginUser.Uid < 1 {
		result["code"] = 101
		result["msg"] = "请先登录，再来抽奖"
		return result
	}
	// 2.用户抽奖分布式锁定
	ok := utils.LockLucky(loginUser.Uid)
	if ok {
		defer utils.UnLockLucky(loginUser.Uid)
	} else {
		result["code"] = 102
		result["msg"] = "正在抽奖，请稍后重试"
		return result
	}
	// 3.验证用户今日参与次数
	userDayNum := utils.IncrUserLuckyNum(loginUser.Uid)
	if userDayNum > conf.UserPrizeMax {
		result["code"] = 103
		result["msg"] = "今日的抽奖次数已用完，明天再来吧"
		return result
	}
	if ok := c.checkUserDay(loginUser.Uid, int(userDayNum)); !ok {
		result["code"] = 104
		result["msg"] = "今日的抽奖次数已用完，明天再来吧"
		return result
	}
	// 4.验证IP今日参与次数
	ip := comm.ClientIP(c.Ctx.Request())
	ipDayNum := utils.IncrIPLuckyNum(ip)
	if ipDayNum > conf.IpLimitMax {
		result["code"] = 105
		result["msg"] = "相同IP参与次数太多，明天再来吧"
		return result
	}
	// 黑名单开关
	limitBlack := false
	if ipDayNum > conf.IpPrizeMax {
		limitBlack = true
	}
	// 5.验证IP黑名单
	var blackIPInfo *models.LtBlackip
	if !limitBlack {
		ok, blackIPInfo = c.checkBlackIP(ip)
		if !ok {
			fmt.Println("黑名单中的IP", ip, limitBlack)
			limitBlack = true
		}
	}
	// 6.验证用户黑名单
	var userInfo *models.LtUser
	if !limitBlack {
		ok, userInfo = c.checkBlackUser(loginUser.Uid)
		if !ok {
			fmt.Println("黑名单中的用户", loginUser.Uid, limitBlack)
			limitBlack = true
		}
	}
	// 7.获得抽奖编码
	prizeCode := comm.Random(10000)
	// 8.匹配奖品是否中奖
	prizeGift := c.prize(prizeCode, limitBlack)
	if prizeGift == nil || prizeGift.PrizeNum < 0 || (prizeGift.PrizeNum > 0 && prizeGift.LeftNum <= 0) {
		result["code"] = 205
		result["msg"] = "很遗憾，没有中奖，请下次再试"
		return result
	}
	// 9.有限制奖品发放
	if prizeGift.PrizeNum > 0 {
		if utils.GetGiftPoolNum(prizeGift.Id) <= 0 {
			result["code"] = 206
			result["msg"] = "很遗憾，没有中奖，请下次再试"
			return result
		}
		ok = utils.PrizeGift(prizeGift.Id, 1)
		if !ok {
			result["code"] = 207
			result["msg"] = "很遗憾，没有中奖，请下次再试"
			return result
		}
	}
	// 10.不同编码优惠券的发放
	if prizeGift.Gtype == conf.GtypeCodeDiff {
		code := utils.PrizeCodeDiff(prizeGift.Id, c.ServiceCode)
		if code == "" {
			result["code"] = 208
			result["msg"] = "很遗憾，没有中奖，请下次再试"
			return result
		}
		prizeGift.Gdata = code
	}
	// 11.记录中奖记录
	finResult := models.LtResult{
		GiftId:     prizeGift.Id,
		GiftName:   prizeGift.Title,
		GiftType:   prizeGift.Gtype,
		Uid:        loginUser.Uid,
		Username:   loginUser.Username,
		PrizeCode:  prizeCode,
		GiftData:   prizeGift.Gdata,
		SysCreated: comm.NowUnix(),
		SysIp:      ip,
		SysStatus:  0,
	}
	err := c.ServiceResult.Create(&finResult)
	if err != nil {
		log.Println("index_lucky.GetLucky ServiceResult.Create ", finResult, ", error is ", err)
		result["code"] = 209
		result["msg"] = "很遗憾，没有中奖，请下次再试"
		return result
	}
	if prizeGift.Gtype == conf.GtypeGiftLarge {
		// 如果获得实物大奖，需要将用户、IP设置成黑名单一段时间
		c.prizeLarge(ip, loginUser, userInfo, blackIPInfo)
	}
	// 12.返回抽奖结果
	result["gift"] = prizeGift
	return result
}
