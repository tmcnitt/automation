package main

import (
	"github.com/Troy-M/automation/drivers"
	"github.com/Troy-M/automation/library"
	"github.com/Troy-M/automation/server"
)

func main() {

	driver := drivers.Driver{}

	//Load the library
	//It returns a function for handeling the registration of devices
	register := library.Start(&driver)

	//Start http server
	//It gets passed the driver so it can listen for updates
	go server.Start(&driver)

	//Start the Driver
	driver.Start(register)
}
