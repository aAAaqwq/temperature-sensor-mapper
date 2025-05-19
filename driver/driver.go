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
		ModbusClient:   nil,
	}
	return client, nil
}

// InitDevice 初始化设备
func (c *CustomizedClient) InitDevice() error {
	// TODO: add init operation
	// you can use c.ProtocolConfig
	//初始化modbus客户端
	var config interface{}
	klog.Infoln("Modbus CommunicateMode:", c.ProtocolConfig)
	switch c.ProtocolConfig.ConfigData.CommunicateMode {
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
		klog.Errorf("Invalid CommunicateMode: %s", c.ConfigData.CommunicateMode)
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
	c.ModbusClient = modbusClient
	klog.Infoln("InitDevice success")
	return nil
}

// GetDeviceData 获取设备数据并转换类型输出
func (c *CustomizedClient) GetDeviceData(visitor *VisitorConfig) (interface{}, error) {
	c.deviceMutex.Lock()
	defer c.deviceMutex.Unlock()

	klog.Infof("开始从设备读取数据,Register: %s, Offset: %d, Limit: %d", visitor.Register, visitor.Offset, visitor.Limit)
	data, err := c.ModbusClient.Get(visitor.Register, visitor.Offset, visitor.Limit)
	if err != nil {
		klog.Errorf("从设备读取数据失败: %v", err)
		return nil, err
	}
	var value uint16
	var warn string
	switch visitor.Register {
	case CoilRegister, InputRegister, DiscreteInputRegister:
		//返回的Coil寄存器数据占用1个字节
		value = uint16(data[0])
		// 工作状态警告处理
		if value == 0 {
			warn = "温度传感器已停止工作"
		}
	case HoldingRegister:
		//返回的Holding寄存器数据占用2个字节
		value = binary.BigEndian.Uint16(data)
		t := ZoomIn(value, visitor.Scale)
		//温度警告处理
		if visitor.Min != nil && t < *visitor.Min {
			warn = fmt.Sprintf("当前温度小于最小值: %.1f < %.1f", t, *visitor.Min)
		}
		if visitor.Max != nil && t > *visitor.Max {
			warn = fmt.Sprintf("当前温度超过最大值: %.1f > %.1f", t, *visitor.Max)
		}
	}
	// 根据数据类型进行转换输出
	v, err := DataNormalize(value, visitor)
	var out string
	if warn == "" {
		out = fmt.Sprintf("从设备读取到数据,Type:%T Value:%v", v, v)
	} else {
		out = fmt.Sprintf("从设备读取到数据,Type:%T Value:%v,Warning:%s", v, v, warn)
	}
	klog.Infof(out)
	return out, err
}

// 根据scale缩小
func ZoomIn(data uint16, scale float64) float64 {
	return cast.ToFloat64(data) * scale
}

// 根据scale放大
func ZoomOut(data float64, scale float64) uint16 {
	if scale == 0 {
		return 0
	}
	return cast.ToUint16(data / scale)
}

// DataNormalize 规范化数据类型转换
func DataNormalize(data interface{}, visitor *VisitorConfig) (interface{}, error) {

	switch visitor.DataType {
	case INT:
		return cast.ToInt(data), nil
	case STRING:
		return cast.ToString(data), nil
	case FLOAT:
		v := cast.ToFloat32(data)
		// 缩放
		if visitor.Scale != 0 {
			v = v * cast.ToFloat32(visitor.Scale)
		}
		return v, nil
	case DOUBLE:
		v := cast.ToFloat64(data)
		// 缩放
		if visitor.Scale != 0 {
			v = v * visitor.Scale
		}
		return v, nil
	case BYTES:
		return data, nil
	case BOOLEN:
		return cast.ToBool(data), nil
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
	klog.Infof("开始写入数据到设备，数据: %v,Visitor: %v", data, visitor)
	if data == "" || data == nil {
		klog.Infof("数据为空，不写入设备")
		return nil
	}
	//数据类型转换
	value := cast.ToUint16(ZoomOut(cast.ToFloat64(data), visitor.Scale))
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
