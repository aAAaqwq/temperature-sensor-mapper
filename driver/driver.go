package driver

import (
	"encoding/binary"
	"fmt"
	"sync"

	"github.com/kubeedge/mapper-framework/pkg/common"
	"github.com/kubeedge/mappers-go/pkg/driver/modbus"
	"github.com/spf13/cast"
	"k8s.io/klog/v2"
)

func NewClient(protocol ProtocolConfig) (*CustomizedClient, error) {
	client := &CustomizedClient{
		ProtocolConfig: protocol,
		deviceMutex:    sync.Mutex{},
		// TODO initialize the variables you added
		ModbusClient: nil,
	}
	return client, nil
}

func (c *CustomizedClient) InitDevice() error {
	// TODO: add init operation
	// you can use c.ProtocolConfig
	//初始化modbus客户端
	var config interface{}
	switch c.ProtocolConfig.ConfigData.Mode {
	case "TCP":
		config = modbus.ModbusTCP{
			SlaveID:  byte(c.ConfigData.SlaveID),
			DeviceIP: c.ConfigData.IP,
			TCPPort:  c.ConfigData.Port,
		}
	case "RTU":
		config = modbus.ModbusRTU{
			SlaveID:      byte(c.ConfigData.SlaveID),
			SerialName:   c.ConfigData.SerialName,
			BaudRate:     c.ConfigData.BaudRate,
			DataBits:     c.ConfigData.DataBits,
			StopBits:     c.ConfigData.StopBits,
			Parity:       c.ConfigData.Parity,
			RS485Enabled: c.ConfigData.RS485Enabled,
		}
	default:
		klog.Errorf("Invalid CommunicateMode: %s", c.ProtocolConfig.ConfigData.Mode)
	}
	klog.Infoln("Start InitDevice with config:", config)
	klog.Infoln("ConfigType:", fmt.Sprintf("%T", config))
	modbusClient, err := modbus.NewClient(config)
	if err != nil {
		klog.Errorf("Failed to create Modbus client: %v", err)
		return err
	}
	// 尝试建立连接
	if err := modbusClient.Client.Connect(); err != nil {
		klog.Errorf("Failed to connect to Modbus device: %v", err)
		c.StopDevice() // 关闭客户端释放资源
		return err
	}
	klog.Infoln("InitDevice success")
	return nil
}

// GetDeviceData 获取设备数据并转换类型输出
func (c *CustomizedClient) GetDeviceData(visitor *VisitorConfig) (interface{}, error) {
	c.deviceMutex.Lock()
	defer c.deviceMutex.Unlock()

	klog.Infof("开始从设备读取数据，Register: %s, Offset: %d, Limit: %d", visitor.Register, visitor.Offset, visitor.Limit)
	data, err := c.ModbusClient.Get(visitor.Register, visitor.Offset, visitor.Limit)
	//返回的data是实际寄存器数据，每个寄存器占用2个字节
	if err != nil {
		klog.Errorf("从设备读取数据失败: %v", err)
		return nil, err
	}
	// 根据每2个字节1个uint16数据,依次转换为uint16
	value := binary.BigEndian.Uint16(data)

	// 根据数据类型进行转换
	switch visitor.DataType {
	case "uint8", "uint16", "uint32", "uint64", "uint", "int8", "int", "int16", "int32", "int64":
		return cast.ToInt(value), nil
	case "string":
		return cast.ToString(value), nil
	case "float32", "float64":
		v := cast.ToFloat64(value)
		// 缩放
		if visitor.Scale != 0 {
			v = v * visitor.Scale
		}
		return v, nil
	default:
		if len(data) > 0 {
			klog.Infof("成功从设备读取数据: %v", value)
			return value, nil
		}
	}
	return nil, fmt.Errorf("无效的数据或数据类型: %v", visitor.DataType)
}

// DeviceDataWrite 外部调用DeviceMethod写入数据到设备，DeviceTwins管理层
func (c *CustomizedClient) DeviceDataWrite(visitor *VisitorConfig, deviceMethodName string, propertyName string, data interface{}) error {
	c.deviceMutex.Lock()
	defer c.deviceMutex.Unlock()

	klog.Infof("开始写入数据到设备，方法名: %s, 属性名: %s, 数据: %v, 寄存器: %s", deviceMethodName, propertyName, data, visitor.Register)

	return c.SetDeviceData(data, visitor)
}

// SetDeviceData 内部写入数据到设备,Device协议层
func (c *CustomizedClient) SetDeviceData(data interface{}, visitor *VisitorConfig) error {
	// TODO: set device's data
	// you can use c.ProtocolConfig and visitor
	klog.Infof("开始写入数据到设备，数据: %v,类型:%s 寄存器: %s", data, visitor.DataType, visitor.Register)
	//数据类型转换
	value := cast.ToUint16(data)
	// 写入到设备的寄存器
	res, err := c.ModbusClient.Set(visitor.Register, visitor.Offset, value)
	if err != nil {
		klog.Errorf("写入数据到设备失败: %v", err)
		return err
	}
	klog.Infof("成功写入数据到设备: %v", binary.BigEndian.Uint16((res)))
	return nil
}

// StopDevice 停止设备连接
func (c *CustomizedClient) StopDevice() error {
	// TODO: stop device
	// you can use c.ProtocolConfig
	err := c.ModbusClient.Client.Close()
	if err != nil {
		klog.Errorf("Failed to close Modbus client: %v", err)
		return err
	}
	return nil
}

// GetDeviceStates 获取设备状态
func (c *CustomizedClient) GetDeviceStates() (string, error) {
	// TODO: GetDeviceStates
	klog.Infoln("开始检查设备状态")
	if err := c.ModbusClient.Client.Connect(); err != nil {
		klog.Errorf("设备连接失败: %v", err)
		return common.DeviceStatusDisCONN, err
	}
	klog.Infoln("设备状态正常")
	return common.DeviceStatusOK, nil
}
