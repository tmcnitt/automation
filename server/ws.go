package server

import (
	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
	"github.com/Troy-M/automation/drivers"
	"log"
	"net/http"
)

type hub struct {
	conn    map[*connection]bool
	env     *Env
	started bool
}

var h hub

type connection struct {
	//The main conn
	ws *websocket.Conn
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

type updateMessage struct {
	Devices []drivers.DeviceInfo     `json:"devices"`
	Scenes  map[string]drivers.Scene `json:"scenes"`
}

func (h *hub) watchDriver() {
	for {
		select {
		//This is called when a value is updated from the driver
		//Propogate this change to all connections in hub
		case _ = <-h.env.driver.Refresh:
			info := updateMessage{}
			for _, device := range h.env.driver.Devices {
				info.Devices = append(info.Devices, device.GetInfo())
			}

			info.Scenes = drivers.GetConfig().Scenes

			for c := range h.conn {
				err := c.ws.WriteJSON(info)
				if err != nil {
					log.Println(err)
				}
			}

		}
	}
}

//serveWs handles websocket requests from the peer.
func (env *Env) serveWs(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	if !h.started {
		h = hub{env: env, conn: make(map[*connection]bool)}
		h.started = true

		//Tell the driver to start listeners
		h.env.driver.Listen = true

		go h.watchDriver()
	}

	if req.Method != "GET" {
		http.Error(res, "Method not allowed", 405)
		return
	}
	ws, err := upgrader.Upgrade(res, req, nil)
	if err != nil {
		log.Println(err)
		return
	}

	c := &connection{ws: ws}
	h.conn[c] = true

	defer func() {
		c.ws.Close()
	}()
	select {}
}
