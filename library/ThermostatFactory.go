package library

import (
	"github.com/Troy-M/automation/drivers"

	"gitlab.com/SurfingNinjas/go-openzwave"
	"gitlab.com/SurfingNinjas/go-openzwave/CC"
)

type thermostat struct {
	drivers.BaseDevice
}

func ThermostatFactory(driver *drivers.Driver, node openzwave.Node) openzwave.Device {
	device := &thermostat{}

	//for every CC
	//the device must send a map describing the layout and type for each
	cc := make(map[openzwave.ValueID]interface{})
	cc[openzwave.ValueID{CC.THERMOSTAT_MODE, 1, 0}] = ""
	cc[openzwave.ValueID{CC.THERMOSTAT_SETPOINT, 1, 2}] = float64(0.0)
	cc[openzwave.ValueID{CC.THERMOSTAT_OPERATING_STATE, 1, 0}] = ""
	cc[openzwave.ValueID{CC.SENSOR_MULTILEVEL, 1, 1}] = float64(0.0)

	device.Init(driver, node, cc)

	return device
}
