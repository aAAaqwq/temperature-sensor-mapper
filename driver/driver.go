package driver

import (
	"encoding/binary"
	"fmt"
	"sync"

	"github.com/kubeedge/mapper-framework/pkg/common"
	"github.com/kubeedge/mappers-go/pkg/driver/modbus"
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
	config := modbus.ModbusTCP{
		SlaveID:  byte(c.ConfigData.SlaveID),
		DeviceIP: c.ConfigData.IP,
		TCPPort:  c.ConfigData.Port,
		// DeviceIP: "192.168.25.239",
		// TCPPort:  "502",
		Timeout: 10,
	}
	klog.Infoln("Start InitDevice with config:",config)
	klog.Infoln("ConfigType:",fmt.Sprintf("%T",config))
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

	// 根据数据类型进行转换
	switch visitor.DataType {
	case "int":
		if len(data) > 0 {
			value := binary.BigEndian.Uint16(data)
			if visitor.Scale != 0 {
				value = uint16(float64(value) * visitor.Scale)
			}
			klog.Infof("成功从设备读取并转换数据: %v", value)
			return value, nil
		}
	case "float":
	default:
		if len(data) > 0 {
			klog.Infof("成功从设备读取数据: %v", data)
			return data, nil
		}
	}

	return nil, fmt.Errorf("无效的数据或数据类型: %v", visitor.DataType)
}

func (c *CustomizedClient) DeviceDataWrite(visitor *VisitorConfig, deviceMethodName string, propertyName string, data interface{}) error {
	c.deviceMutex.Lock()
	defer c.deviceMutex.Unlock()

	klog.Infof("开始写入数据到设备，方法名: %s, 属性名: %s, 数据: %v, 寄存器: %s", deviceMethodName, propertyName, data, visitor.Register)

	// 根据数据类型进行转换
	var value uint16
	switch v := data.(type) {
	case int:
		if visitor.Scale != 0 {
			value = uint16(float64(v) / visitor.Scale)
		} else {
			value = uint16(v)
		}
	case float64:
		if visitor.Scale != 0 {
			value = uint16(v / visitor.Scale)
		} else {
			value = uint16(v)
		}
	case string:
		// 尝试将字符串转换为整数
		var intValue int
		_, err := fmt.Sscanf(v, "%d", &intValue)
		if err != nil {
			klog.Errorf("字符串转换为整数失败: %v", err)
			return err
		}
		value = uint16(intValue)
	default:
		klog.Errorf("不支持的数据类型: %T", data)
		return fmt.Errorf("不支持的数据类型: %T", data)
	}

	klog.Infof("准备写入数据到寄存器，值: %d", value)
	// 写入到设备的寄存器
	res, err := c.ModbusClient.Set(visitor.Register, uint16(visitor.Offset), value)
	if err != nil {
		klog.Errorf("写入数据到设备失败: %v", err)
		return err
	}
	klog.Infof("成功写入数据到设备: %v", res)

	return nil
}

func (c *CustomizedClient) SetDeviceData(data interface{}, visitor *VisitorConfig) error {
	// TODO: set device's data
	// you can use c.ProtocolConfig and visitor
	return c.DeviceDataWrite(visitor, "", "", data)
}

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
