package server

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/Troy-M/automation/drivers"
	"log"
	"net/http"
	"strconv"
	"strings"

	"time"
)

type Env struct {
	driver *drivers.Driver
}

func Start(driver *drivers.Driver) {
	router := httprouter.New()

	env := Env{driver: driver}

	router.GET("/set/:id/:cc/:status", env.setStatus)
	router.POST("/name", env.name)
	router.GET("/devices", env.getDevices)

	router.POST("/scenes/:id/schedule", env.scheduleScene)
	router.POST("/scenes", env.newScene)
	router.DELETE("/scenes/:id", env.deleteScene)
	router.GET("/scenes", env.getScenes)
	router.GET("/scenes/run/*params", env.runScene)

	router.GET("/add", env.add)

	router.GET("/", env.serveWs)

	log.Fatalln(http.ListenAndServe(":80", router))
}

//CORS sets headers to allow the frontend to be served on a different port
func CORS(res http.ResponseWriter) {
	res.Header().Add("Content-Type", "application/json")
	res.Header().Add("Access-Control-Allow-Origin", "*")
	res.Header().Add("Access-Control-Allow-Headers", "Authorization, Content-Type")
	res.Header().Add("Access-Control-Allow-Content-Type", "*")
}

func (env *Env) add(res http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	CORS(res)

	result := env.driver.AddNode()

	//TODO: better output here
	if result {
		res.WriteHeader(200)
	} else {
		res.WriteHeader(408)
	}
}

func (env *Env) setStatus(res http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	CORS(res)

	//Hack to pad number. ie 8 to 008
	i, _ := strconv.Atoi(ps.ByName("id"))
	device := env.driver.Devices[uint8(i)]

	cc, _ := strconv.Atoi(ps.ByName("cc"))
	if device != nil {
		device.SetStatus(ps.ByName("status"), uint8(cc))
	}
}

func (env *Env) getDevices(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	CORS(res)

	var devices []drivers.DeviceInfo
	for _, device := range env.driver.Devices {
		devices = append(devices, device.GetInfo())
	}

	send, _ := json.Marshal(devices)
	fmt.Fprintf(res, "%s", send)
}

type NameRequest struct {
	ID   uint8  `json:"ID"`
	Name string `json:"name"`
}

func (env *Env) name(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	CORS(res)

	request := NameRequest{}
	json.NewDecoder(req.Body).Decode(&request)

	device := env.driver.Devices[request.ID]
	device.GetDevice().Node.SetNodeName(request.Name)
}

func (env *Env) newScene(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	CORS(res)

	scene := drivers.Scene{}
	json.NewDecoder(req.Body).Decode(&scene)

	//scene ID is generated in save()
	scene.ID = ""
	id := scene.Save()

	scene = drivers.Scene{ID: id}
	send, _ := json.Marshal(scene)
	fmt.Fprintf(res, "%s", send)
}

func (env *Env) runScene(res http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	CORS(res)

	params := strings.Split(ps.ByName("params"), "/")
	id := params[1]

	env.driver.RunScene(drivers.GetConfig().Scenes[id])
}

type schedule struct {
	Time string `json:"time"`
}

func (env *Env) scheduleScene(res http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	CORS(res)

	id := ps.ByName("id")

	schedule := schedule{}
	json.NewDecoder(req.Body).Decode(&schedule)

	_, err := time.Parse(time.Kitchen, schedule.Time)
	if err != nil {
		res.WriteHeader(400)
		return
	}

	config := drivers.GetConfig()

	//make sure this is a valid scene to schedule
	scene := config.Scenes[id]
	if scene.ID == "" {
		res.WriteHeader(400)
		return
	}

	//This route is used for toggle schedule
	//So if it's already in the list, we take it that the user wants to unschedule it
	found := false
	for i, x := range scene.Schedules {
		if x == schedule.Time {
			found = true
			scene.Schedules = append(scene.Schedules[:i], scene.Schedules[i+1:]...)
		}
	}

	if !found {
		scene.Schedules = append(scene.Schedules, schedule.Time)
	}

	scene.Save()
}

func (env *Env) deleteScene(res http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	CORS(res)

	id := ps.ByName("id")
	drivers.DeleteScene(id)
}

func (env *Env) getScenes(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	CORS(res)

	scenes := drivers.GetConfig().Scenes

	send, _ := json.Marshal(scenes)
	fmt.Fprintf(res, "%s", send)
}
