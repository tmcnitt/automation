package library

import (
	"github.com/Troy-M/automation/drivers"

	"github.com/Troy-M/go-openzwave"
	"github.com/Troy-M/go-openzwave/CC"
)

type dimmerSwitch struct {
	drivers.BaseDevice
}

func DimmerSceneSwitchFactory(driver *drivers.Driver, node openzwave.Node) openzwave.Device {
	device := &dimmerSwitch{}

	//for every CC
	//the device must send a map describing the layout and type for each
	cc := make(map[openzwave.ValueID]interface{})
	cc[openzwave.ValueID{CC.SWITCH_MULTILEVEL, 1, 0}] = uint8(0)

	device.Init(driver, node, cc)

	return device
}
