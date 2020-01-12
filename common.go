package sensor

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"os"
	"regexp"
	"strconv"
	"sync"
)

/**
 * bigEndian transfer
 */
func ToBigEndian(num uint16) []byte {
	int16buf := new(bytes.Buffer)
	if err := binary.Write(int16buf, binary.BigEndian, num); err != nil {
		panic("error type: num")
	}
	return int16buf.Bytes()
}

/**
 * littleEndian transfer
 */
func ToLittleEndian(num uint16) []byte {
	int16buf := new(bytes.Buffer)
	if err := binary.Write(int16buf, binary.LittleEndian, num); err != nil {
		panic("error type: num")
	}
	return int16buf.Bytes()
}

/**
 * bytes to uint(bigEndian)
 */
func BytesToIntU(b []byte) (uint16, error) {
	bytesBuffer := bytes.NewBuffer(b)
	var tmp uint16
	err := binary.Read(bytesBuffer, binary.BigEndian, &tmp)
	return tmp, err
}

/**
 * generate a default request command
 */
func ComposeBody(DeviceAddr, FuncCode, Data []byte) []byte {
	ret := append(DeviceAddr, FuncCode...)
	ret = append(ret, Data...)
	ret = append(ret, CreateCRC(ret)...)
	return ret
}

/**
 * separate config command parameters
 * @params src is a measure respond data
 * @return DeviceAddr and FuncCode
 * @return RegisterAddr
 * @return ConfigData
 * @return error it will also use crc-16 to validate whether it got a true data, if not returns error(wrong data)
 */
func SplitConfig(src []byte) ([]byte, []byte, []byte, error) {
	if ValidateCRC(src[:6], src[6:8]) {
		return src[:2], src[2:4], src[4:6], nil
	}
	return nil, nil, nil, errors.New("got error data")
}

/**
 * CRC验证与拆分过程
 * 验证respond的类型, crc校验以及是否为可被类型识别的指令返回体
 * @param src 接收到的数据
 * @return DeviceMeta 验证过程中拆分到的meta数据
 * @return []byte 拆分到的数据体
 * @return CRC验证失败时的错误返回
 */
func SplitAndValidate(src []byte) (DeviceMeta, []byte, error) {
	base := len(src) - 2
	if ValidateCRC(src[:base], src[base:]) {
		var meta DeviceMeta
		meta.Addr = src[0]
		meta.FuncCode = src[1]
		if src[1] > 0x80 {
			return meta, src[2:base], nil
		} else if src[1] == 0x03 {
			return meta, src[3:base], nil
		} else if src[1] == 0x06 {
			return meta, src[2:base], nil
		}
	}
	return DeviceMeta{}, nil, errors.New("unreachable validate")
}

const HexValue = 256

func ByteToFloat(v []byte) float64 {
	return float64(v[0])*HexValue + float64(v[1])
}

func FourByteToFloat(v []byte) ([]float64, error) {
	count := len(v)
	if count%4 != 0 || count == 0 {
		return nil, errors.New("get error type")
	}
	var iter = 0
	var ret []float64
	for {
		i := DecimalFloat(ByteToFloat(v[0+iter:2+iter]) * math.Pow(0.1, ByteToFloat(v[2+iter:4+iter])))
		ret = append(ret, i)
		if iter+4 == count {
			break
		} else {
			iter += 4
		}
	}
	return ret, nil
}

func DecimalFloat(value float64) float64 {
	value, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", value), 64)
	return value
}

func TwoByteToFloatX1000(v []byte) (float64, error) {
	count := len(v)
	if count > 2 || count == 0 {
		return 0, errors.New("get error type")
	}
	var ret float64
	ret = ByteToFloat(v) / 1000
	return ret, nil
}

func T() {

}

// ========================json-part-start===========================

// 传感器可能出现的状态
const (
	STATUS_NORMAL = iota // 正常运行的
	STATUS_DETACH        // 异常断开的
	STATUS_CLOSED        // 人为关闭的 备注: 该状态不写入CONFIG文件中, 因此该状态仅持续到下位机重启时刻
)

// 传感器参数 包含自定义任务和状态等信息
type LocalSensorInformation struct {
	// ==========OPTIONS============
	TaskHandler func(body TaskSensorBody, wg *sync.WaitGroup) // 自定义传感器任务
	Status      int                                           // 传感器状态
	// ==========CONFIGS============
	Addr     byte   `json:"addr"`     // 传感器设备地址
	Type     byte   `json:"type"`     // 传感器类型
	Attach   string `json:"attach"`   // 传感器附着的透传设备
	Interval int64  `json:"interval"` // 最大间隔时间(秒)
	SensorID string `json:"sensorID"` // 传感器ID
}

// 下位机参数
type LocalDeviceDetail struct {
	Name                   string                   `json:"name"`                   // 收集器名称
	LocalSensorInformation []*LocalSensorInformation `json:"localSensorInformation"` // 传感器集合
}

// 加载测试
func GetConfigTest() *LocalDeviceDetail {
	config := LoadConfig("cnf/conf.json")
	return config
}

// 本地传感器信息()
var localDeviceDetail *LocalDeviceDetail = nil

// 加载参数
func GetLocalDevicesInstance() *LocalDeviceDetail {
	if localDeviceDetail == nil {
		localDeviceDetail = GetConfigTest()
		return localDeviceDetail
	} else {
		return localDeviceDetail
	}
}

const configFileSizeLimit = 10 << 20

// Config加载
func LoadConfig(path string) *LocalDeviceDetail {
	var config LocalDeviceDetail
	configFile, err := os.Open(path)
	if err != nil {
		emit("Failed to open config file '%s': %s\n", path, err)
		return &config
	}

	fi, _ := configFile.Stat()
	if size := fi.Size(); size > (configFileSizeLimit) {
		emit("config file (%q) size exceeds reasonable limit (%d) - aborting", path, size)
		return &config
	}

	if fi.Size() == 0 {
		emit("config file (%q) is empty, skipping", path)
		return &config
	}

	buffer := make([]byte, fi.Size())
	_, err = configFile.Read(buffer)
	buffer, err = StripComments(buffer)
	if err != nil {
		emit("Failed to strip comments from json: %s\n", err)
		return &config
	}

	buffer = []byte(os.ExpandEnv(string(buffer)))

	err = json.Unmarshal(buffer, &config)
	if err != nil {
		emit("Failed unmarshalling json: %s\n", err)
		return &config
	}
	return &config
}

// 注释清除
func StripComments(data []byte) ([]byte, error) {
	data = bytes.Replace(data, []byte("\r"), []byte(""), 0)
	lines := bytes.Split(data, []byte("\n"))
	filtered := make([][]byte, 0)

	for _, line := range lines {
		match, err := regexp.Match(`^\s*#`, line)
		if err != nil {
			return nil, err
		}
		if !match {
			filtered = append(filtered, line)
		}
	}
	return bytes.Join(filtered, []byte("\n")), nil
}

func emit(msg string, args ...interface{}) {
	log.Printf(msg, args...)
}

// test
func ResultConfig(test []map[string]interface{}) (port_password []map[string]interface{}) {
	return
}

// ========================json-part-end===========================
