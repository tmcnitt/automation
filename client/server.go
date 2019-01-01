package main

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"net/http"
)

//TODO: refactor this garbage

func main() {
	router := httprouter.New()

	router.GET("/", func(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {
		webpage, err := ioutil.ReadFile("./build/index.html")
		if err != nil {
			http.Error(res, fmt.Sprintf("home.html file error %v", err), 500)
		}
		fmt.Fprint(res, string(webpage))

	})

	router.GET("/bundle.js", func(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {
		webpage, err := ioutil.ReadFile("./build/bundle.js")
		if err != nil {
			http.Error(res, fmt.Sprintf("home.html file error %v", err), 500)
		}
		fmt.Fprint(res, string(webpage))

	})

	router.GET("/bundle.js.map", func(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {
		webpage, err := ioutil.ReadFile("./build/bundle.js.map")
		if err != nil {
			http.Error(res, fmt.Sprintf("home.html file error %v", err), 500)
		}
		fmt.Fprint(res, string(webpage))

	})

	http.ListenAndServe("0.0.0.0:5000", router)

}

/*
  Open bash on directory
  $go run server.go
*/
