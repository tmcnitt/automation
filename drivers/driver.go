package drivers

import (
	"encoding/json"
	"github.com/satori/go.uuid"
	"github.com/Troy-M/go-openzwave"
	"io/ioutil"
	"os"
	"strconv"
	"time"
)

const configFile = "config.json"

//Config holds data from configFile. Contains information on scenes
type Config struct {
	Scenes map[string]Scene `json:"scenes"`

	loaded bool
}

//Driver is the main struct to handle zwave automation
type Driver struct {
	zwaveAPI openzwave.API
	Devices  map[uint8]Device

	//Channel to update WS connections
	Refresh chan bool

	//if there is any WS connections to update
	Listen bool

	//used internally in setName command
	//TODO: use this for controller commands from underlying zwave api
	networkID uint32
}

//Device is the main interface all devices implement
type Device interface {
	GetInfo() DeviceInfo
	SetStatus(interface{}, uint8)
	GetDevice() *BaseDevice
}

//DeviceInfo is returned by GetInfo() and returns all info about device state
type DeviceInfo struct {
	Name string `json:"name,omitempty"`
	ID   uint8  `json:"ID"`

	//TODO: this should have some concrete format to it
	Sigs   map[string]string     `json:"sigs,omitempty"`
	Status map[uint8]interface{} `json:"status,omitempty"`
}

//BaseDevice contains all default fields every device has
type BaseDevice struct {
	Driver *Driver
	Node   openzwave.Node
	Info   DeviceInfo
	CC     []openzwave.ValueID
}

//Init sets general params based on given node information
//This is called from the device factory
//It gives us layout, which labels each command class (key) to what type of value it has (value)
//This gives us a format for what type of value is expected for each command class which is used later
func (device *BaseDevice) Init(driver *Driver, node openzwave.Node, layout map[openzwave.ValueID]interface{}) {
	//We get a layout that has every command class and the type of value
	//Loop through and just get the command class
	CC := []openzwave.ValueID{}
	for cc := range layout {
		CC = append(CC, cc)
	}

	//TODO: This should be done once? It probably doesn't hurt but isn't needed
	driver.networkID = node.GetHomeId()

	device.Driver = driver
	device.Node = node
	device.CC = CC

	productID := node.GetProductId()
	nodeID := node.GetId()

	productDescription := node.GetProductDescription()

	sigs := make(map[string]string)
	sigs["zwave:manufacturerId"] = productID.ManufacturerId
	sigs["zwave:nodeID"] = strconv.Itoa(int(nodeID))
	sigs["zwave:productID"] = productID.ProductId
	sigs["zwave:manufacturerName"] = productDescription.ManufacturerName
	sigs["zwave:productName"] = productDescription.ProductName
	device.Info.Sigs = sigs

	device.Info.Name = device.Node.GetNodeName()
	device.Info.ID = nodeID

	//Now use the layout to init all the status values
	device.Info.Status = make(map[uint8]interface{})
	for cc, value := range layout {
		device.Info.Status[cc.CommandClassId] = value
	}

	for _, cc := range device.CC {
		value := device.Node.GetValueWithId(cc)

		updateValue(value, device, cc)

		//TODO: how should we decide what to poll? Maybe from factory?
		device.Node.GetValueWithId(cc).SetPollingState(true)
	}

}

func updateValue(value openzwave.Value, device *BaseDevice, cc openzwave.ValueID) {
	//update value based on what type already exists
	//TODO: the other types
	switch device.Info.Status[cc.CommandClassId].(type) {
	case bool:
		val, _ := value.GetBool()
		device.Info.Status[cc.CommandClassId] = val
	case int:
		val, _ := value.GetInt()
		device.Info.Status[cc.CommandClassId] = val
	case float64:
		val, _ := value.GetFloat()
		device.Info.Status[cc.CommandClassId] = val
	case string:
		val, _ := value.GetString()
		device.Info.Status[cc.CommandClassId] = val
	case uint8:
		val, _ := value.GetUint8()
		device.Info.Status[cc.CommandClassId] = val
	}
}

//SetStatus takes a command class and value and sets the value
func (device *BaseDevice) SetStatus(status interface{}, id uint8) {
	//Using the the command class ID, find the actual value in the status
	CC := openzwave.ValueID{}
	for _, cc := range device.CC {
		if cc.CommandClassId == id {
			CC = cc
		}
	}

	//get the current value so we know the type
	value := device.Node.GetValueWithId(CC)

	//TODO: the other types
	switch device.Info.Status[CC.CommandClassId].(type) {

	case string:
		value.SetString(status.(string))
	case bool:
		s, _ := strconv.ParseBool(status.(string))
		value.SetBool(s)
	case float64:
		s, _ := strconv.ParseFloat(status.(string), 64)
		value.SetFloat(s)
	case int:
		s, _ := strconv.Atoi(status.(string))
		value.SetInt(s)
	case uint8:
		s, _ := strconv.ParseUint(status.(string), 0, 8)
		value.SetUint8(uint8(s))
	}

}

//NodeAdded is called when the zwave api tells us there is a new device
//We add it to the driver list and call the refresh
func (device *BaseDevice) NodeAdded() {
	device.Driver.Devices[device.Node.GetId()] = device
	if device.Driver.Listen {
		device.Driver.Refresh <- true
	}
}

//ValueChanged is called when the zwave api tells us there is a change to a value
//If there is a difference we call refresh and set the value
func (device *BaseDevice) ValueChanged(v openzwave.Value) {
	changed := false

	for _, cc := range device.CC {
		if v.Id() == cc {
			current := device.Info.Status[cc.CommandClassId]

			updateValue(v, device, cc)

			if device.Info.Status[cc.CommandClassId] != current {
				changed = true
			}
		}
	}

	if device.Driver.Listen && changed {
		device.Driver.Refresh <- true
	}
}

//GetInfo is called from the interface and we return all the state information we have
func (device *BaseDevice) GetInfo() DeviceInfo {
	return device.Info
}

//TODO: is this needed?
func (device *BaseDevice) GetDevice() *BaseDevice {
	return &BaseDevice{Driver: device.Driver, Info: device.Info, Node: device.Node}
}

//TODO: What value does this respond to?
func (device *BaseDevice) NodeChanged() {
	//TODO: how can we check for this?
	device.Info.Name = device.Node.GetNodeName()
}

//NodeRemoved is called when the device gets removed from the network
//We just remove it from the array and call refresh
func (device *BaseDevice) NodeRemoved() {
	delete(device.Driver.Devices, device.Node.GetId())

	if device.Driver.Listen {
		device.Driver.Refresh <- true
	}
}

//DeviceRegistry is called from the DeviceFactory in order to setup a new device
type DeviceRegistry func(driver *Driver, node openzwave.Node) openzwave.Device

//Start takes a deviceFactroy and starts the event loop
func (driver *Driver) Start(factory openzwave.DeviceFactory) {
	API := openzwave.
		BuildAPI("/etc/openzwave", ".", "").
		SetDeviceFactory(factory).
		SetDeviceName("/dev/ttyACM0")

	driver.zwaveAPI, _ = API.(openzwave.API)
	driver.Devices = make(map[uint8]Device)
	driver.Refresh = make(chan bool)

	driver.scheduleLoop()

	//this is a blocking function
	API.Run()

	os.Exit(0)
}

func (driver *Driver) scheduleLoop() {
	timer := time.NewTicker(time.Minute * 1)

	go func() {
		for {
			select {
			case <-timer.C:
				current := time.Now().Format(time.Kitchen)
				for _, scene := range config.Scenes {
					for _, schedule := range scene.Schedules {
						scheduleT, _ := time.Parse(time.Kitchen, schedule)
						check := scheduleT.Format(time.Kitchen)

						if check == current {
							driver.RunScene(scene)
						}
					}
				}
			}
		}
	}()
}

//AddNode is a command that calls the zwave api addnode function
func (driver *Driver) AddNode() bool {
	return driver.zwaveAPI.GetNetwork(driver.networkID).AddNode(false)
}

var config = Config{}

//GetConfig loads the config from the file from the const configFile
//If it has already been loaded, we just return the cached version
func GetConfig() Config {
	if config.loaded {
		return config
	}

	file, err := ioutil.ReadFile(configFile)
	if err != nil {
		writeConfig(Config{Scenes: make(map[string]Scene)})
		return Config{Scenes: make(map[string]Scene)}
	}

	json.Unmarshal(file, &config)

	config.loaded = true

	return config
}

//writeConfig Writes the updated config to the file system
//Also updates the local var so it doesn't need to be reloaded
func writeConfig(save Config) {
	config = save

	write, _ := json.Marshal(save)
	ioutil.WriteFile(configFile, write, os.ModePerm)
}

//RunScene steps through each action of the provided scene
func (driver *Driver) RunScene(scene Scene) {
	for _, action := range scene.Actions {
		driver.Devices[action.ID].SetStatus(action.Status, action.CC)
	}
}

//Scene is a type of stored actions
//It can also be schedulded to run at certain times
type Scene struct {
	ID string `json:"ID,omitempty"`

	//TODO: does something with name
	Name string `json:"name,omitempty"`

	//The actual actions
	Actions []SceneAction `json:"actions,omitempty"`

	//Times to run the scene in kitchen format (hh:mm PM/AM)
	Schedules []string `json:"schedules,omitempty"`
}

//SceneAction details the action to take in a scene
type SceneAction struct {
	ID     uint8       `json:"ID,omitempty"`
	Status interface{} `json:"status,omitempty"`
	CC     uint8
}

//Save adds a scene to the config
//If it is passed an ID, it just updates the scene
func (scene Scene) Save() string {
	if scene.ID == "" {
		id, _ := uuid.NewV4()
		scene.ID = id.String()
	}

	config := GetConfig()
	config.Scenes[scene.ID] = scene

	writeConfig(config)

	return scene.ID
}

//DeleteScene removes a scene from the config
//This also deletes the schedules
func DeleteScene(id string) {
	config := GetConfig()
	delete(config.Scenes, id)

	writeConfig(config)
}
