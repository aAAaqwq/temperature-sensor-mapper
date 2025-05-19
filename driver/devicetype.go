package driver

import (
	"sync"

	"github.com/kubeedge/mapper-framework/pkg/common"
	"github.com/kubeedge/mappers-go/pkg/driver/modbus"
)
const (
	CoilRegister          = "CoilRegister"
	InputRegister         = "InputRegister"
	HoldingRegister       = "HoldingRegister"
	DiscreteInputRegister = "DiscreteInputRegister"
	INT                   = "int"
	FLOAT                 = "float"
	DOUBLE                = "double"
	STRING                = "string"
	BOOLEN                = "boolen"
	BYTES                 = "bytes"

)

// CustomizedDev is the customized device configuration and client information.
type CustomizedDev struct {
	Instance         common.DeviceInstance
	CustomizedClient *CustomizedClient
}

type CustomizedClient struct {
	// TODO add some variables to help you better implement device drivers
	deviceMutex sync.Mutex
	ProtocolConfig
	ModbusClient *modbus.ModbusClient // Modbus客户端
}

type ProtocolConfig struct {
	ProtocolName string `json:"protocolName"`
	ConfigData      `json:"configData"`
}

type ConfigData struct {
	// TODO: add your protocol config data
	SlaveID int `json:"slaveID,omitempty"` // 从站ID
	CommunicateMode string `json:"communicateMode"`   // 通信模式:enum:TCP/RTU
	// TCP mode
	IP string `json:"ip,omitempty"`
	Port string `json:"port,omitempty"`
	Timeout int `json:"timeout,omitempty"`
	// RTU mode
	SerialName string `json:"serialName,omitempty"` // 串口名称
	BaudRate int `json:"baudRate,omitempty"`// 波特率
	DataBits int `json:"dataBits,omitempty"`// 数据位
	StopBits int `json:"stopBits,omitempty"`// 停止位
	Parity string `json:"parity,omitempty"`// 校验位:enum:None/Even/Odd
	RS485Enabled bool `json:"rs485Enabled,omitempty"`// 是否启用RS485
}

type VisitorConfig struct {
	ProtocolName      string `json:"protocolName"`
	VisitorConfigData `json:"configData"`
}

type VisitorConfigData struct {
	// TODO: add your visitor config data
	DataType string `json:"dataType"` // 数据类型:enum:Int/Float/String
	Register string `json:"register"` // 寄存器类型:enum:CoilRegister/HoldingRegister
	Offset uint16 `json:"offset"`        // 寄存器偏移量
	Scale float64 `json:"scale"`      // 数据缩放比例
	IsSwap bool `json:"isSwap,omitempty"`       // 是否交换字节
	Limit uint16 `json:"limit"`          // 读取数量
	IsRegisterSwap bool `json:"isRegisterSwap,omitempty"` // 是否交换寄存器
	Max *float64 `json:"max,omitempty"`            // 最大值
	Min *float64 `json:"min,omitempty"`            // 最小值
}
