package driver

import (
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
	config:=&modbus.ModbusTCP{
		SlaveID:  byte(c.ConfigData.SlaveID),
		DeviceIP: c.ConfigData.CommunicationMode.IP,
		TCPPort:  c.ConfigData.CommunicationMode.Port,
	}
	klog.Info("Start InitDevice with config:",config)
	modbusClient,err:=modbus.NewClient(config)
	if err!=nil{
		klog.Errorf("Failed to create Modbus client: %v", err)
		return err
	}
	// 尝试建立连接
	if err:=modbusClient.Client.Connect();err!=nil{
		klog.Errorf("Failed to connect to Modbus device: %v", err)
		c.StopDevice() // 关闭客户端释放资源
		return err
	}
	c.ModbusClient=modbusClient
	klog.Info("InitDevice success")
	return nil
}

func (c *CustomizedClient) GetDeviceData(visitor *VisitorConfig) (interface{}, error) {
	// TODO: add the code to get device's data
	// you can use c.ProtocolConfig and visitor
	//获取register数据
	c.deviceMutex.Lock()
	defer c.deviceMutex.Unlock()
	return c.ModbusClient.Get(visitor.Register,0,uint16(visitor.Limit))
}

func (c *CustomizedClient) DeviceDataWrite(visitor *VisitorConfig, deviceMethodName string, propertyName string, data interface{}) error {
	// TODO: add the code to write device's data
	// you can use c.ProtocolConfig and visitor to write data to device
	return nil
}

func (c *CustomizedClient) SetDeviceData(data interface{}, visitor *VisitorConfig) error {
	// TODO: set device's data
	// you can use c.ProtocolConfig and visitor
	return nil
}

func (c *CustomizedClient) StopDevice() error {
	// TODO: stop device
	// you can use c.ProtocolConfig
	err:=c.ModbusClient.Client.Close()
	if err!=nil {
		klog.Errorf("Failed to close Modbus client: %v", err)
		return err
	}
	return nil
}

func (c *CustomizedClient) GetDeviceStates() (string, error) {
	// TODO: GetDeviceStates
	if err:=c.ModbusClient.Client.Connect();err!=nil{
		return common.DeviceStatusDisCONN,err
	}
	return common.DeviceStatusOK, nil
}
