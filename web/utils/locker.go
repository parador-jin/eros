package utils

import (
	"Eros/datasource"
	"fmt"
)

func getLuckyLockKey(uid int) string {
	return fmt.Sprintf("lucky_lock_%d", uid)
}

func LockLucky(uid int) bool {
	key := getLuckyLockKey(uid)
	cacheObj := datasource.InstanceCache()
	rs, _ := cacheObj.Do("SET", key, 1, "EX", 3, "NX")
	if rs == "OK" {
		return true
	}
	return false
}

func UnLockLucky(uid int) bool {
	key := getLuckyLockKey(uid)
	cacheObj := datasource.InstanceCache()
	rs, _ := cacheObj.Do("DEL", key)
	if rs == "OK" {
		return true
	}
	return true
}
