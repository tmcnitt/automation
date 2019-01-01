package library

import (
	"log"

	"github.com/Troy-M/automation/drivers"
	"github.com/Troy-M/go-openzwave"
)

type library map[openzwave.ProductId]drivers.DeviceRegistry

//Load registers all plugins with handlers
func load() *library {
	l := make(library)
	//TODO: load external modules
	//TODO: could this use a more vague sense that wouldn't require per device registration?
	l[openzwave.ProductId{"0x001d", "0x0334"}] = BinarySceneSwitchFactory
	l[openzwave.ProductId{"0x001d", "0x0001"}] = BinarySceneSwitchFactory
	l[openzwave.ProductId{"0x008b", "0x5442"}] = ThermostatFactory
	l[openzwave.ProductId{"0x0063", "0x3038"}] = DimmerSceneSwitchFactory
	return &l
}

//getFactory checks the library and returns the drivers.DeviceFactory for a certain ProductId
func (l *library) getFactory(id *openzwave.ProductId) drivers.DeviceRegistry {
	entry := (*l)[*id]
	if entry != nil {
		return entry
	}

	//TODO: remove debug information
	log.Println(*id)
	log.Println("Device not found")
	return nil
}

//Start loads the library and returns a callback to load device plugins
func Start(driver *drivers.Driver) openzwave.DeviceFactory {
	library := load()

	//function is called by driver when a new device needs a factory
	registry := (func(api openzwave.API, node openzwave.Node) openzwave.Device {
		factory := library.getFactory(node.GetProductId())
		if factory != nil {
			return factory((driver), node)
		}
		return unknownDevice()
	})

	return registry
}

//Default device
//TODO: better loging here
//TODO: maybe basic command class stuff?
type emptyDevice struct {
}

func unknownDevice() openzwave.Device {
	return &emptyDevice{}
}

func (*emptyDevice) NodeAdded() {
}

func (*emptyDevice) NodeChanged() {
}

func (*emptyDevice) NodeRemoved() {
}

func (*emptyDevice) ValueChanged(value openzwave.Value) {
}
