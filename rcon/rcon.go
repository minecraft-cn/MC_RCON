package rcon

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math/rand"
	"net"
)

type PacketSize int32

const (
	PacketTypeLogin               PacketSize = 3
	PacketTypeCmd                 PacketSize = 2
	PacketTypeMultiPacketResponse PacketSize = 0

	WrongPasswordRequestID PacketSize = -1
)

var (
	Debug bool

	PacketEndPad = make([]byte, 2, 2)
)

type RconClient struct {
	ServerIP   string
	ServerPort string
	ServerAddr string
	Password   string

	Connection net.Conn

	RequestPacket *RconPacket
}

type RconPacket struct {
	Length    PacketSize
	RequestID PacketSize
	Type      PacketSize
	Payload   string
}

func (p *RconPacket) fromBytes(byteData []byte) {

	var (
		buffer = bytes.NewBuffer(byteData)
	)

	binary.Read(buffer, binary.LittleEndian, &p.Length)
	binary.Read(buffer, binary.LittleEndian, &p.RequestID)
	binary.Read(buffer, binary.LittleEndian, &p.Type)

	cmdResult := make([]byte, p.Length-10)

	buffer.Read(cmdResult)
	p.Payload = string(cmdResult)
}

func (p *RconPacket) toBytes() []byte {

	var (
		data = bytes.NewBuffer([]byte{})
	)

	if p.Length == 0 {
		p.Length = PacketSize(4 + 4 + len(p.Payload) + 2)
	}

	if p.RequestID == 0 {
		p.RequestID = PacketSize(rand.Int31())
	}

	binary.Write(data, binary.LittleEndian, p.Length)
	binary.Write(data, binary.LittleEndian, p.RequestID)
	binary.Write(data, binary.LittleEndian, p.Type)
	data.Write([]byte(p.Payload))
	data.Write(PacketEndPad)

	return data.Bytes()
}

func (c *RconClient) Connect() error {
	conn, err := net.Dial("tcp", c.ServerAddr)
	if err != nil {
		return err
	}

	c.Connection = conn
	return nil
}

func (c *RconClient) write(requestType PacketSize, data string) error {
	packet := RconPacket{
		Type:    requestType,
		Payload: data,
	}
	c.RequestPacket = &packet

	_, err := c.Connection.Write(packet.toBytes())
	if err != nil {
		return err
	}

	return nil
}

func (c *RconClient) read() (*RconPacket, error) {
	var (
		respPacket = RconPacket{}
		respData   = make([]byte, 4096)
	)

	_, err := c.Connection.Read(respData)
	if err != nil {
		return nil, err
	}

	respPacket.fromBytes(respData)

	if Debug {
		fmt.Printf("%+v\n", *c.RequestPacket)
		fmt.Printf("%+v\n", respPacket)
	}

	//fmt.Printf("%s\n", respPacket.Payload)

	return &respPacket, nil
}

func (c *RconClient) Login() error {

	err := c.write(PacketTypeLogin, c.Password)
	if err != nil {
		return err
	}

	packet, err := c.read()
	if err != nil {
		return err
	}

	if packet.RequestID == WrongPasswordRequestID ||
		packet.RequestID != c.RequestPacket.RequestID ||
		packet.Type != PacketTypeCmd {
		return fmt.Errorf("wrong password: %s", c.Password)
	}

	return nil
}

func (c *RconClient) RunCmd(cmd string) (string, error) {

	err := c.write(PacketTypeCmd, cmd)
	if err != nil {
		return "", err
	}

	rp, err := c.read()
	if err != nil {
		return "", err
	}

	return rp.Payload, nil
}
