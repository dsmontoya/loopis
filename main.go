package main

import (
	"fmt"

	"github.com/kylelemons/gousb/usb"
)

func main() {
	var controllers []*MIDIController
	c := usb.NewContext()
	defer c.Close()
	c.Debug(0)
	devs, err := c.ListDevices(func(desc *usb.Descriptor) bool {
		return true
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(devs)
	for _, d := range devs {
		desc := d.Descriptor
		fmt.Println(desc)
		fmt.Printf("Vendor %d\n", desc.Vendor)
		fmt.Printf("Product %d\n", desc.Product)
		fmt.Println("Spec", desc.Spec)
		ctrl := GetMIDIControllerFromDevice(d)
		if ctrl != nil {
			controllers = append(controllers, ctrl)
			continue
		}
		d.Close()
	}
	controller := controllers[0]
	defer controller.Close()
	fmt.Println("configInfo", controller.Configs)
	fmt.Println("configs length", len(controller.Configs))
	desc1, err := controller.GetStringDescriptor(1)
	desc2, err := controller.GetStringDescriptor(2)
	if err != nil {
		panic(err)
	}
	fmt.Println("getStringDescriptor 1", desc1)
	fmt.Println("getStringDescriptor 2", desc2)
	fmt.Println("interfaces", controller.Configs[0].Interfaces)
	//endpoint := controller.OpenEndpoint(controller.Descriptor.Configs, iface, setup, epoint)
}

var knownMIDIControllers = map[ID][]ID{
	6048: []ID{
		13640,
	},
}

type ID uint16
type IDs []ID

type Endpoint struct {
	usb.Endpoint
}

type MIDIController struct {
	*usb.Device
}

func GetMIDIControllerFromDevice(device *usb.Device) *MIDIController {
	var vendor ID
	var product ID
	desc := device.Descriptor
	vendor = ID(desc.Vendor)
	product = ID(desc.Product)
	if getListOfProducts(vendor).Include(product) {
		return &MIDIController{device}
	}
	return nil
}

func getListOfProducts(vendor ID) IDs {
	return knownMIDIControllers[vendor]
}

func (ids IDs) Include(product ID) bool {
	for i := 0; i < len(ids); i++ {
		if ids[i] == product {
			return true
		}
	}
	return false
}

func (c *MIDIController) NewEnpoint() *Endpoint {
	//c.OpenEndpoint(c.Descriptor.Configs[0].Config, iface, setup, epoint)
	return nil
}
