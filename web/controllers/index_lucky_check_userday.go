package controllers

import (
	"Eros/conf"
	"Eros/models"
	"Eros/web/utils"
	"fmt"
	"log"
	"strconv"
	"time"
)

func (c *IndexController) checkUserDay(uid, num int) bool {
	userdayInfo := c.ServiceUserday.GetUserToday(uid)
	if userdayInfo != nil && userdayInfo.Uid == uid {
		// 今天存在抽奖记录
		if userdayInfo.Num >= conf.UserPrizeMax {
			if num < userdayInfo.Num {
				utils.InitUserLuckyNum(uid, int64(userdayInfo.Num))
			}
			return false
		}
		userdayInfo.Num++
		if num < userdayInfo.Num {
			utils.InitUserLuckyNum(uid, int64(userdayInfo.Num))
		}
		if err103 := c.ServiceUserday.Update(userdayInfo, nil); err103 != nil {
			log.Println("index_lucky_check_userday ServiceUserDay.Update error, error info is ", err103)
		}
	} else {
		// 创建今天的用户参与记录
		y, m, d := time.Now().Date()
		strDay := fmt.Sprintf("%d%02d%02d", y, m, d)
		day, _ := strconv.Atoi(strDay)
		userdayInfo = &models.LtUserday{
			Uid:        uid,
			Day:        day,
			Num:        1,
			SysCreated: int(time.Now().Unix()),
		}
		utils.InitUserLuckyNum(uid, 1)
		if err103 := c.ServiceUserday.Create(userdayInfo); err103 != nil {
			log.Println("index_lucky_check_userday ServiceUserDay.Create error, error info is ", err103)
		}
	}
	return true
}
