package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

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

	var (
		line = bufio.NewReader(os.Stdin)

		timeout = time.Minute * 1
		timer   = time.NewTimer(timeout)

		sig = make(chan os.Signal, 1)
	)

	signal.Notify(sig, os.Interrupt)
	defer timer.Stop()

	go func() {
		for {
			select {
			case <-timer.C:
				log.Printf("exit for time out")
				os.Exit(0)
			case <-sig:
				log.Printf("exit for ctrl+c")
				os.Exit(0)
			default:
			}

			time.Sleep(time.Second * 1)
		}
	}()

	for {
		line, _, err := line.ReadLine()
		if err != nil {
			log.Println(err)
			continue
		}
		timer.Reset(timeout)

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
