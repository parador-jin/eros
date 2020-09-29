package utils

import (
	"Eros/comm"
	"Eros/datasource"
	"fmt"
	"log"
	"math"
	"time"
)

const userFrameSize = 2

func init() {
	resetGroupUserList()
}

func resetGroupUserList() {
	log.Println("user_day_lucky.resetGroupUserList start")
	cacheObj := datasource.InstanceCache()
	for i := 0; i < userFrameSize; i++ {
		key := fmt.Sprintf("day_users_%d", i)
		cacheObj.Do("DEL", key)
	}
	log.Println("user_day_lucky.resetGroupUserList stop")
	// 用户抽奖当天的统计数，零点的时候归零，设置定时器
	duration := comm.NextDayDuration()
	time.AfterFunc(duration, resetGroupUserList)
}

// 初始化用户每日参与次数
func InitUserLuckyNum(uid int, num int64) {
	if num <= 1 {
		return
	}
	i := uid % userFrameSize
	key := fmt.Sprintf("day_users_%d", i)
	cacheObj := datasource.InstanceCache()
	_, err := cacheObj.Do("HSET", key, uid, num)
	if err != nil {
		log.Println("user_day_lucky redis HSET key is ", key, ", uid is ", uid, ", error is ", err)
	}
}

func IncrUserLuckyNum(uid int) int64 {
	i := uid % userFrameSize
	key := fmt.Sprintf("day_users_%d", i)
	cacheObj := datasource.InstanceCache()
	rs, err := cacheObj.Do("HINCRBY", key, uid, 1)
	if err != nil {
		log.Println("user_day_lucky redis HINCRBY key is ", key, ", uid is ", uid, ", error is ", err)
		return math.MaxInt32
	}
	num := rs.(int64)
	return num
}
