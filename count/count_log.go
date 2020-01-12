package count

import (
	"sync"
	"time"
)

/**
 * 对单位sensor进行detach断言
 * @param sensorID 传感器ID
 * 可以实现单位sensor的增删查改
 *
 */

type SensorLog struct {
	sensorID   string //传感器ID
	errorCount int    // 连续错误次数三次后 errorTag -> true
	errorTag   bool   // 错误标识符(禁止重试)
	sync.Mutex
}

var sensorLog = make(map[string]*SensorLog)

const (
	ERROR_DELAY_LEVEL1 = time.Minute
	ERROR_DELAY_LEVEL2 = time.Minute * 2
	ERROR_DELAY_LEVEL3 = time.Minute * 3
)

/*
 * 添加错误日志, 一旦通过了三次错误则不再允许出现第四次错误, 直到错误被用户处理
 * @param sensorID 传感器ID
 *
 */
func AddErrorOperation(sensorID string) {
	// 判断某个Key是否存在
	if v, ok := sensorLog[sensorID]; ok {
		v.ForbidRequest()
	} else {
		s := SensorLog{}
		s.sensorID = sensorID
		sensorLog[sensorID] = &s
		// 延迟
		s.ForbidRequest()
	}
}

/**
 * 禁止请求传感器, 关闭errorTag
 * @param sensorID 传感器ID
 * @param delayLevel 延迟级别
 *
 */
func (sl *SensorLog) ForbidRequest() {
	if sl.errorTag {
		return
	}
	sl.errorCount++
	sl.errorTag = true
	var delayLevel time.Duration
	if sensorLog[sl.sensorID].errorCount == 1 {
		delayLevel = ERROR_DELAY_LEVEL1
	} else if sensorLog[sl.sensorID].errorCount == 2 {
		delayLevel = ERROR_DELAY_LEVEL2
	} else if sensorLog[sl.sensorID].errorCount == 3 {
		delayLevel = ERROR_DELAY_LEVEL3
	}
	if delayLevel != 0 {
		time.AfterFunc(delayLevel, func() {
			sensorLog[sl.sensorID].errorTag = false
		})
	}
}

/**
 * 通过判断errorTag, 是否允许传感器进行通信
 * @param sensorID 传感器ID
 * @return true 允许进行/false 禁止进行查询
 */
func IsForbidden(sensorID string) bool {
	if v, ok := sensorLog[sensorID]; ok {
		if v.errorTag {
			return false
		}
	}
	return true
}
