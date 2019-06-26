package main

import (
	"flag"
	"fmt"
	"log"
	"mcRcon/rcon"
	"os"
)

var (
	addr     string
	password string
)

func flagArgs() {
	flag.StringVar(&addr, "addr", "127.0.0.1:25575", "ip:port")
	flag.StringVar(&password, "password", "rconPassword", "rcon-password in server.properties")
	flag.Parse()
}

func main() {

	flagArgs()



	c := &rcon.RconClient{
		ServerAddr: addr,
		Password:   password,
	}

	err := c.Connect()
	if err != nil {
		log.Println(err)
		return
	}

	err = c.Login()
	if err != nil {
		log.Println(err)
		return
	} else {
		log.Println("login successfully")
	}

	var cmd string
	for {
		_, err := fmt.Fscanf(os.Stdin, "%s\n", &cmd)
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}

		result, err := c.RunCmd(cmd)
		if err != nil {
			log.Printf("run command(%s) failed: %s", cmd, err)
			continue
		}
		fmt.Println(result)
	}
}
