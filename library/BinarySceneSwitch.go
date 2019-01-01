package library

import (
	"github.com/Troy-M/automation/drivers"

	"github.com/Troy-M/go-openzwave"
	"github.com/Troy-M/go-openzwave/CC"
)

type binarySwitch struct {
	drivers.BaseDevice
}

func BinarySceneSwitchFactory(driver *drivers.Driver, node openzwave.Node) openzwave.Device {
	device := &binarySwitch{}

	//for every CC
	//the device must send a map describing the layout and type for each
	cc := make(map[openzwave.ValueID]interface{})
	cc[openzwave.ValueID{CC.SWITCH_BINARY, 1, 0}] = false

	device.Init(driver, node, cc)

	return device
}
