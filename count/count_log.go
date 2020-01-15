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
	sensorID   string    //传感器ID
	errorCount int       // 连续错误次数三次后 errorTag -> true
	errorTag   bool      // 错误标识符(禁止重试)
	retryTime  time.Time // 重试时间
	sync.Mutex
}

/**
 * 绝对禁用
 */
func AddErrorOperationBan(sensorID string) int {
	// 判断某个Key是否存在
	if v, ok := sensorLog[sensorID]; ok {
		v.errorTag = true
		v.errorCount ++
		return v.errorCount
	} else {
		s := SensorLog{}
		s.sensorID = sensorID
		sensorLog[sensorID] = &s
		s.errorTag = true
		s.errorCount = 1
		return s.errorCount
	}
}

/**
 * 返回重试恢复时间点
 */
func GetRetryTime(sensorID string) time.Time {
	if v, ok := sensorLog[sensorID]; ok {
		return v.retryTime
	}
	return time.Now()
}

/**
 * 重试时间
 */
func (sl *SensorLog) getRetryTime() {
	switch sl.errorCount {
	case 1:
		sl.retryTime = time.Now().Add(ERROR_DELAY_LEVEL1)
	case 2:
		sl.retryTime = time.Now().Add(ERROR_DELAY_LEVEL2)
	case 3:
		sl.retryTime = time.Now().Add(ERROR_DELAY_LEVEL3)
	}
}

var sensorLog = make(map[string]*SensorLog)

const (
	ERROR_DELAY_LEVEL1 = time.Minute
	ERROR_DELAY_LEVEL2 = time.Minute * 2
	ERROR_DELAY_LEVEL3 = time.Minute * 5
)

/*
 * 添加错误日志, 一旦通过了三次错误则不再允许出现第四次错误, 直到错误被用户处理
 * @param sensorID 传感器ID
 * @param int 错误次数
 */
func AddErrorOperation(sensorID string) int {
	// 判断某个Key是否存在
	if v, ok := sensorLog[sensorID]; ok {
		v.ForbidRequest()
		return v.errorCount
	} else {
		s := SensorLog{}
		s.sensorID = sensorID
		sensorLog[sensorID] = &s
		// 延迟
		s.ForbidRequest()
		return s.errorCount
	}
}

/**
 * 统计传感器错误次数
 * @param 传感器ID
 * @return 返回错误次数
 */
func GetErrorCount(sensorID string) int {
	if v, ok := sensorLog[sensorID]; ok {
		return v.errorCount
	}
	return 0
}

/**
 * 清除单个传感器错误信息
 * @param 传感器ID
 *
 */
func ClsErrorCount(sensorID string) {
	delete(sensorLog, sensorID)
}

const CLEAR_ALL_EXCEPTION = "all"
const CLEAR_ONE_EXCEPTION = "one"
/**
 * GC回收吧
 */
func ClsAll() {
	sensorLog = make(map[string]*SensorLog)
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
		sl.getRetryTime()
		time.AfterFunc(delayLevel, func() {
			sensorLog[sl.sensorID].errorTag = false
		})
	}
}

/**
 * 通过判断errorTag, 是否允许传感器进行通信
 * @param sensorID 传感器ID
 * @return false 允许进行/true 禁止进行查询
 */
func IsForbidden(sensorID string) bool {
	if v, ok := sensorLog[sensorID]; ok {
		if v.errorTag {
			return true
		}
	}
	return false
}

/**
 * 间隔时间处理
 */
func HandleSlot() {

}
