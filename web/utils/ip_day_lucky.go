package utils

import (
	"Eros/comm"
	"Eros/datasource"
	"fmt"
	"log"
	"math"
	"time"
)

const ipFrameSize = 2

func init() {
	resetGroupIPList()
}

func resetGroupIPList() {
	log.Println("ip_day_lucky.resetGroupIPList start")
	cacheObj := datasource.InstanceCache()
	for i := 0; i < ipFrameSize; i++ {
		key := fmt.Sprintf("day_ips_%d", i)
		cacheObj.Do("DEL", key)
	}
	log.Println("ip_day_lucky.resetGroupIPList stop")
	// IP当天的统计数，零点的时候归零，设置定时器
	duration := comm.NextDayDuration()
	time.AfterFunc(duration, resetGroupIPList)

}

func IncrIPLuckyNum(strIP string) int64 {
	ip := comm.Ip4toInt(strIP)
	i := ip % ipFrameSize
	key := fmt.Sprintf("day_ips_%d", i)
	cacheObj := datasource.InstanceCache()
	rs, err := cacheObj.Do("HINCRBY", key, ip, 1)
	if err != nil {
		log.Println("ip_day_lucky redis HINCRBY error is ", err)
		return math.MaxInt32
	}
	return rs.(int64)
}
