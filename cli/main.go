package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/minecraft-cn/MC_RCON/rcon"
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

	var line = bufio.NewReader(os.Stdin)
	for {
		line, _, err := line.ReadLine()
		if err != nil {
			log.Println(err)
			continue
		}

		if len(line) == 0 {
			continue
		}

		result, err := c.RunCmd(string(line))
		if err != nil {
			log.Printf("run command(%s) failed: %s", line, err)
			continue
		}
		fmt.Println(result)
	}
}
