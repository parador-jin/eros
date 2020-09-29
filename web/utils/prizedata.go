package utils

import (
	"Eros/comm"
	"Eros/conf"
	"Eros/datasource"
	"Eros/models"
	"Eros/services"
	"encoding/json"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"log"
	"time"
)

func PrizeGift(id, leftNum int) bool {
	ok := false
	ok = prizeServGift(id)
	if ok {
		giftService := services.NewGiftService()
		rows, err := giftService.DecrLeftNum(id, leftNum)
		if rows < 1 || err != nil {
			log.Println("prizeData.PrizeGift giftService.DecrLeftNum error is ", err, ", rows is ", rows)
			return false
		}
	}
	return ok
}

func PrizeCodeDiff(id int, codeService services.CodeService) string {
	return prizeServCodeDiff(id, codeService)
}

func GetGiftPoolNum(id int) int {
	key := "gift_pool"
	cacheObj := datasource.InstanceCache()
	rs, err := cacheObj.Do("HGET", key, id)
	if err != nil {
		log.Println("prizedata.GetGiftPoolNum error, error is ", err)
		return 0
	}
	num := comm.GetInt64(rs, 0)
	return int(num)
}

func prizeServGift(id int) bool {
	key := "gift_pool"
	cacheObj := datasource.InstanceCache()
	rs, err := cacheObj.Do("HINCRBY", key, id, -1)
	if err != nil {
		log.Println("prizedata.prizeServGift HINCRBY error, error is ", err)
		return false
	}
	num := comm.GetInt64(rs, -1)
	if num >= 0 {
		return true
	}
	return false
}

func ImportCacheCodes(id int, code string) bool {
	key := fmt.Sprintf("gift_code_%d", id)
	cacheObj := datasource.InstanceCache()
	_, err := cacheObj.Do("SADD", key, code)
	if err != nil {
		log.Println("prizedata.ImportCacheCodes SADD error is ", err)
		return false
	}
	return true
}

func ReCacheCodes(id int, codeService services.CodeService) (sucNum, errNum int) {
	list := codeService.Search(id)
	if list == nil || len(list) <= 0 {
		return 0, 0
	}
	key := fmt.Sprintf("gift_code_%d", id)
	cacheObj := datasource.InstanceCache()
	for _, data := range list {
		if data.SysStatus == 0 {
			code := data.Code
			_, err := cacheObj.Do("SADD", fmt.Sprintf("tmp_%s", key), code)
			if err != nil {
				log.Println("prizedata.ReCacheCodes SADD error is ", err)
				errNum++
				continue
			}
			sucNum++
		}
	}
	_, err := cacheObj.Do("RENAME", fmt.Sprintf("tmp_%s", key), key)
	if err != nil {
		log.Println("prizedata.ReCacheCodes RENAME error is ", err)
	}
	return sucNum, errNum
}

func GetCacheCodeNum(id int, codeService services.CodeService) (num, cacheNum int) {
	// 从数据库中获取数据
	list := codeService.Search(id)
	if len(list) > 0 {
		for _, data := range list {
			if data.SysStatus == 0 {
				num++
			}
		}
	}
	// 从缓存中获取数据
	key := fmt.Sprintf("gift_code_%d", id)
	cacheObj := datasource.InstanceCache()
	rs, err := cacheObj.Do("SCARD", key)
	if err != nil {
		log.Println("prizedata.GetCacheCodeNum SCARD error, error is ", err)
	} else {
		cacheNum = int(comm.GetInt64(rs, 0))
	}
	return num, cacheNum
}

func prizeServCodeDiff(id int, codeService services.CodeService) string {
	key := fmt.Sprintf("gift_code_%d", id)
	cacheObj := datasource.InstanceCache()
	rs, err := cacheObj.Do("SPOP", key)
	if err != nil {
		log.Println("prizedata.prizeServCodeDiff SPOP error, error is ", err)
		return ""
	}
	code := comm.GetString(rs, "")
	if code == "" {
		log.Println("prizedata.prizeServCodeDiff rs is ", rs)
	}
	codeService.UpdateByCode(&models.LtCode{Code: code, SysStatus: 2, SysUpdated: comm.NowUnix()}, nil)
	return code
}

func prizeLocalCodeDiff(id int, codeService services.CodeService) string {
	lockUid := 0 - id - 1000000000
	LockLucky(lockUid)
	defer UnLockLucky(lockUid)

	codeId := 0
	codeInfo := codeService.NextUsingCode(id, codeId)
	if codeInfo != nil && codeInfo.Id > 0 {
		codeInfo.SysStatus = 2
		codeInfo.SysUpdated = comm.NowUnix()
		codeService.Update(codeInfo, nil)
	} else {
		log.Println("prizedata.prizeCodeDiff num codeInfo, gift_id is ", id)
		return ""
	}
	return codeInfo.Code
}

func ResetGiftPrizeData(giftInfo *models.LtGift, giftService services.GiftService) {
	if giftInfo == nil || giftInfo.Id < 1 {
		return
	}
	id := giftInfo.Id
	nowTime := comm.NowUnix()
	if giftInfo.SysStatus == 1 || giftInfo.TimeBegin >= nowTime || giftInfo.TimeEnd <= nowTime ||
		giftInfo.LeftNum <= 0 || giftInfo.PrizeNum <= 0 {
		if giftInfo.PrizeData != "" {
			// 清空旧的发奖数据
			clearGiftPrizeData(giftInfo, giftService)
		}
		return
	}
	// 没有设置发江周期
	dayNum := giftInfo.PrizeTime
	if dayNum < 0 {
		setGiftPool(id, giftInfo.LeftNum)
		return
	}
	// 重置发奖计划数据
	setGiftPool(id, 0)
	// 实际的奖品计划分布运算
	prizeNum := giftInfo.PrizeNum
	avgNum := prizeNum / dayNum
	// 每天可以分配的奖品数
	dayPrizeNum := make(map[int]int)
	if avgNum >= 1 {
		for day := 0; day < dayNum; day++ {
			dayPrizeNum[day] = avgNum
		}
	}
	// 剩下的随机分配到任意哪天
	prizeNum -= dayNum * avgNum
	for prizeNum > 0 {
		prizeNum--
		day := comm.Random(dayNum)
		_, ok := dayPrizeNum[day]
		if !ok {
			dayPrizeNum[day] = 1
		} else {
			dayPrizeNum[day] += 1
		}
	}
	// 每天的map，每小时的map，60分钟的数组，奖品数
	prizeData := make(map[int]map[int][60]int)
	for day, num := range dayPrizeNum {
		// 计算出这一天的发奖计划
		dayPrizeData := getGiftDataOneDay(num)
		prizeData[day] = dayPrizeData
	}
	// 将周期内的每天，每小时，每分钟的数据 prizeData 格式化([时间:数量])
	dataList := formatGiftPrizeData(nowTime, dayNum, prizeData)
	str, err := json.Marshal(dataList)
	if err != nil {
		log.Println("prizedata.ResetGiftPrizeData json error, error info is ", err)
	} else {
		info := &models.LtGift{
			Id:         giftInfo.Id,
			LeftNum:    giftInfo.PrizeNum,
			PrizeData:  string(str),
			PrizeBegin: nowTime,
			PrizeEnd:   nowTime + dayNum*86400,
			SysUpdated: nowTime,
		}
		err := giftService.Update(info, nil)
		if err != nil {
			log.Println("prizedata.ResetGiftPrizeData giftService.Update error, error info is ", err)
		}
	}
}

func clearGiftPrizeData(giftInfo *models.LtGift, giftService services.GiftService) {
	info := &models.LtGift{Id: giftInfo.Id, PrizeData: ""}
	err := giftService.Update(info, []string{"prize_data"})
	if err != nil {
		log.Println("prizedata.clearGiftPrizeData giftService.Update is ", info, ", error info is ", err)
	}
	setGiftPool(giftInfo.Id, 0)
}

func setGiftPool(id int, num int) {
	// 设置奖品池的库存数量
	key := "gift_pool"
	cacheObj := datasource.InstanceCache()
	_, err := cacheObj.Do("HSET", key, id, num)
	if err != nil {
		log.Println("prizedata.setGiftPool error info is ", err)
	}

}

func getGiftDataOneDay(num int) map[int][60]int {
	rs := make(map[int][60]int)
	// 计算24小时各自的奖品数
	hourData := [24]int{}
	if num > 100 {
		var hourNum int
		for _, h := range conf.PrizeDataRandomDayTime {
			hourData[h]++
		}
		for h := 0; h < 24; h++ {
			d := hourData[h]
			n := num * d / 100
			hourData[h] = n
			hourNum += n
		}
		num -= hourNum
	}
	for num > 0 {
		num--
		hourIndex := comm.Random(100)
		h := conf.PrizeDataRandomDayTime[hourIndex]
		hourData[h]++
	}
	// 将每小时内的奖品数量分配到60分钟
	for h, hNum := range hourData {
		if hNum <= 0 {
			continue
		}
		minuteData := [60]int{}
		if hNum >= 60 {
			avgMinute := hNum / 60
			for i := 0; i < 60; i++ {
				minuteData[i] = avgMinute
			}
			hNum -= avgMinute * 60
		}
		for hNum > 0 {
			hNum--
			m := comm.Random(60)
			minuteData[m]++
		}
		rs[h] = minuteData
	}
	return rs
}

func formatGiftPrizeData(nowTime, dayNum int, prizeData map[int]map[int][60]int) [][2]int {
	rs := make([][2]int, 0)
	nowHour := time.Now().Hour()
	// 处理日期的数据
	for dn := 0; dn < dayNum; dn++ {
		dayData, ok := prizeData[dn]
		if !ok {
			continue
		}
		dayTime := nowTime + dn*86400
		// 处理小时的数据
		for hn := 0; hn < 24; hn++ {
			hourData, ok := dayData[(hn+nowHour)%24]
			if !ok {
				continue
			}
			hourTime := dayTime + hn*3600
			// 处理分钟的数据
			for mn := 0; mn < 60; mn++ {
				num := hourData[mn]
				if num <= 0 {
					continue
				}
				minuteTime := hourTime + mn*60
				rs = append(rs, [2]int{minuteTime, num})
			}
		}
	}
	return rs
}

func DistributionGiftPool() int {
	var totalNum int
	now := comm.NowUnix()
	giftService := services.NewGiftService()
	list := giftService.GetAll(false)
	if list != nil && len(list) > 0 {
		for _, gift := range list {
			if gift.SysStatus != 0 {
				continue
			}
			if gift.PrizeNum < 1 {
				continue
			}
			if gift.TimeBegin > now || gift.TimeEnd < now {
				continue
			}
			if len(gift.PrizeData) <= 7 {
				continue
			}
			var cronData [][2]int
			err := json.Unmarshal([]byte(gift.PrizeData), &cronData)
			if err != nil {
				log.Println("prizeData.DistributionGiftPool Unmarshal error, error info is", err)
			} else {
				var giftNum, index int
				for i, data := range cronData {
					ct := data[0]
					num := data[1]
					if ct <= now {
						giftNum += num
						index = i + 1
					} else {
						break
					}
				}
				// 更新奖品池
				if giftNum > 0 {
					incrGiftPool(gift.Id, giftNum)
					totalNum += giftNum
				}
				// 更新奖品的发奖计划
				if index > 0 {
					if index >= len(cronData) {
						cronData = make([][2]int, 0)
					} else {
						cronData = cronData[index:]
					}
				}
				str, err := json.Marshal(cronData)
				if err != nil {
					log.Println("prizeData.DistributionGiftPool Marshal is ", cronData, ", error info is", err)
				}
				columns := []string{"prize_data"}
				err = giftService.Update(&models.LtGift{Id: gift.Id, PrizeData: string(str)}, columns)
				if err != nil {
					log.Println("prizeData.DistributionGiftPool Update error, error info is", err)
				}
			}
		}
		if totalNum > 0 {
			giftService.GetAll(true)
		}
	}
	return totalNum
}

func incrGiftPool(id, num int) int {
	key := "gift_pool"
	cacheObj := datasource.InstanceCache()
	rsNum, err := redis.Int(cacheObj.Do("HINCRBY", key, id, num))
	if err != nil {
		log.Println("prizeData.incrGiftPool error is", err)
		return 0
	}
	if rsNum < num {
		// 递增少于预期值，补偿一次
		num2 := num - rsNum
		rsNum, err = redis.Int(cacheObj.Do("HINCRBY", key, id, num2))
		if err != nil {
			log.Println("prizeData.incrGiftPool2 error is ", err)
			return 0
		}
	}
	return rsNum
}
