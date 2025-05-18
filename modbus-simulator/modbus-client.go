package main

import (
	"encoding/binary"
	"fmt"
	"reflect"

	"github.com/kubeedge/mappers-go/pkg/driver/modbus"
	"k8s.io/klog/v2"
)

func main() {
	config := modbus.ModbusTCP{
		SlaveID: 1,
		DeviceIP: "192.168.25.239",
		TCPPort: "502",
	}
	client, err := modbus.NewClient(config)
	if err != nil {
		klog.Errorf("Failed to create Modbus TCP client: %v", err)
	}
	// 尝试建立连接
	if err := client.Client.Connect(); err!= nil {
		klog.Errorf("Failed to connect to Modbus TCP device: %v", err)
	}
	// 读取保持寄存器
	data, err := client.Get("HoldingRegister", 0, 10)
	if err!= nil {
		klog.Errorf("Failed to read holding registers: %v", err)
	}
	n:=len(data)/2
	value := make([]uint16, n)
	fmt.Println(n,value)
	for i:=0;i<n;i+=2{
		value[i]=binary.BigEndian.Uint16(data[i:i+2])
	}
	klog.Infof("Read holding registers: origin:%v convert:%d", data,value)
	// 读取线圈寄存器
	data, err = client.Get("CoilRegister", 0, 1)
	if err!= nil {
		klog.Errorf("Failed to read coil registers: %v", err)
	}
	klog.Infof("Read coil registers: %v", data)

	// 写入保持寄存器
	data,err = client.Set("HoldingRegister", 0, 500)
	if err!= nil {
		klog.Errorf("Failed to write holding registers: %v", err)
	}
	klog.Infof("Write holding registers: %v", binary.BigEndian.Uint16(data))
	// 写入线圈寄存器
	data,err = client.Set("CoilRegister", 0, 0)	
	if err!= nil {
		klog.Errorf("Failed to write coil registers: %v", err)
	}
	klog.Infof("Write coil registers: %v", data)
	var cfg interface{}
	cfg=modbus.ModbusTCP{
		SlaveID: 1,
	}
	klog.Infoln(reflect.TypeOf(cfg),reflect.ValueOf(cfg))


}