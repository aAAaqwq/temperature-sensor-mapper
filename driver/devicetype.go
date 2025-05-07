package driver

import (
	"sync"

	"github.com/kubeedge/mapper-framework/pkg/common"
	"github.com/kubeedge/mappers-go/pkg/driver/modbus"
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
	ConfigData   `json:"configData"`
}

type ConfigData struct {
	// TODO: add your protocol config data
	SlaveID int `json:"slaveID,omitempty"`
	CommunicationMode `json:"communication"` // 通信模式:TCP/RTU 
}
type CommunicationMode struct{
	Mode string `json:"mode"`
	IP string `json:"ip,omitempty"`
	Port string `json:"port,omitempty"`
	// Todo：add RTU mode
}

type VisitorConfig struct {
	ProtocolName      string `json:"protocolName"`
	VisitorConfigData `json:"configData"`
}

type VisitorConfigData struct {
	// TODO: add your visitor config data
	DataType string `json:"dataType"`
	Register string `json:"register"`
	Offset int `json:"offset"`
	Scale float64 `json:"scale"`
	IsSwap bool `json:"isSwap"`
	Limit int `json:"limit"`
	IsRegisterSwap bool `json:"isRegisterSwap"`
}
