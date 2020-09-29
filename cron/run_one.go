package cron

import (
	"Eros/comm"
	"Eros/services"
	"Eros/web/utils"
	"log"
	"time"
)

func ConfigureAppOneCron() {
	go resetAllGiftPrizeData()
	go distributionAllGiftPool()
}

func resetAllGiftPrizeData() {
	giftService := services.NewGiftService()
	nowTime := comm.NowUnix()
	list := giftService.GetAll(false)
	for _, giftInfo := range list {
		if giftInfo.PrizeTime > 0 && (giftInfo.PrizeData == "" || giftInfo.PrizeEnd <= nowTime) {
			log.Println("crontab start utils.resetAllGiftPrizeData, giftInfo is ", giftInfo)
			utils.ResetGiftPrizeData(&giftInfo, giftService)
			giftService.GetAll(true)
			log.Println("crontab end utils.resetAllGiftPrizeData, giftInfo is ", giftInfo)
		}
	}
	// 每5分钟执行一次
	time.AfterFunc(5*time.Minute, resetAllGiftPrizeData)
}

func distributionAllGiftPool() {
	log.Println("crontab start utils.distributionAllGiftPool")
	num := utils.DistributionGiftPool()
	log.Println("crontab end utils.distributionAllGiftPool num is ", num)
	// 每分钟执行一次
	time.AfterFunc(time.Minute, distributionAllGiftPool)
}
