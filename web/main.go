package main

import (
	"fmt"
	"mcRcon/rcon"
	"net/http"
)

var (
	rconIP = "127.0.0.1"
	//rconIP       = "193.112.63.225"
	rconPort = "25575"

	listenIP   = "0.0.0.0"
	listenPort = "25576"
)

func main() {

	http.Handle("/", http.FileServer(http.Dir(".")))
	http.HandleFunc("/cmd", func(writer http.ResponseWriter, request *http.Request) {

		addr := request.FormValue("addr")
		password := request.FormValue("pswd")
		cmd := request.FormValue("cmd")

		fmt.Println(addr, password, cmd)

		c := &rcon.RconClient{
			ServerAddr: addr,
			Password:   password,
		}

		err := c.Connect()
		if err != nil {
			writer.Write([]byte(err.Error()))
			return
		}

		err = c.Login()
		if err != nil {
			writer.Write([]byte(err.Error()))
			return
		} else {
			//writer.Write([]byte("login successfully"))
		}

		result, err := c.RunCmd(cmd)
		if err != nil {
			writer.Write([]byte(fmt.Sprintf("run command(%s) failed: %s", cmd, err)))
			return
		}
		writer.Write([]byte(result))
	})

	panic(http.ListenAndServe(fmt.Sprintf("%s:%s", listenIP, listenPort), nil))
}
