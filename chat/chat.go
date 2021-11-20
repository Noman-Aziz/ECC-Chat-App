package chat

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/Noman-Aziz/ECC-Chat-App/ecc"
)

type Config struct {
	Name            string `json:"name"`
	BroadcastName   bool   `json:"broadcast_name"`
	IsLocal         bool   `json:"is_local"`
	IsHost          bool   `json:"is_host"`
	Port            uint16 `json:"port"`
	PortDescription string `json:"port_description"`
}

type chatapp struct {
	EC         *ecc.EllipticCurve `json:"ecc_config"`
	ECCKeyPair *ecc.Keys          `json:"ecc_key_pair"`
	AppConfig  *Config            `json:"app_config"`
	Other      *partner           `json:"other"`
}

type DataHeader struct {
	Name      string      `json:"name"`
	PublicKey ecc.ECPoint `json:"public_key"`
}

func CreateChatApp(config *Config) *chatapp {
	// var port = config.Port
	// var description = config.PortDescription

	var EC ecc.EllipticCurve
	var keys ecc.Keys

	EC, keys = ecc.Initialization()

	fmt.Println("\nMY PRIVATE KEY :", keys.PrivKey)
	fmt.Println("MY PUBLIC KEY :", keys.PubKey)

	var app = &chatapp{
		EC:         &EC,
		ECCKeyPair: &keys,
		AppConfig:  config,
		Other:      nil,
	}

	termChan := make(chan os.Signal)
	signal.Notify(termChan, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		<-termChan
		log.Print("[Info]: Ctrl-C shutting down...\n")
		os.Exit(0)
	}()

	return app
}

func (c *chatapp) runHost() {
	//Open TCP Socket
	var listener, err = net.Listen("tcp4", fmt.Sprintf("%s:%d", "127.0.0.1", c.AppConfig.Port))
	if err != nil {
		log.Panicf("[Error: %s]: Error opening TCP socket\n", err.Error())
	}
	defer listener.Close()
	//Accept Client Connection
	conn, err := listener.Accept()
	if err != nil {
		log.Panicf("[Error: %s]: Error accepting client connection\n", err.Error())
	}
	defer conn.Close()

	var connWriter = json.NewEncoder(conn)
	var connReader = json.NewDecoder(conn)

	//Send host data to client
	var name = "Host"
	if c.AppConfig.BroadcastName {
		name = c.AppConfig.Name
	}
	err = connWriter.Encode(DataHeader{name, *&c.ECCKeyPair.PubKey})
	if err != nil {
		log.Panicf("[Error: %s]: Error sending data to client\n", err.Error())
	}

	//Read client data from client
	var clientData DataHeader
	err = connReader.Decode(&clientData)
	if err != nil {
		log.Panicf("[Error: %s]: Error receiving data from client\n", err.Error())
	}
	c.Other = CreatePartner(clientData.Name, clientData.PublicKey)

	fmt.Println("\nOther Peer Public Key:", clientData.PublicKey)
	fmt.Println()

	//Custom Input Reader
	cin := bufio.NewScanner(os.Stdin)

	//Chatting can begin
	//Taking Text Input
	for {
		fmt.Printf("Enter Message to Send: ")
		cin.Scan()
		message := cin.Text()

		Send(message, c, connWriter)

		fmt.Printf("Decrypted Message: %v\n\n", Recv(connReader, c))

	}

}

func (c *chatapp) runClient() {
	var conn, err = net.Dial("tcp4", fmt.Sprintf("%s:%d", "127.0.0.1", c.AppConfig.Port))
	if err != nil {
		log.Panicf("[Error: %s]: Error connecting to host\n", err.Error())
	}
	defer conn.Close()

	var connWriter = json.NewEncoder(conn)
	var connReader = json.NewDecoder(conn)

	//Recieve host data from host
	var hostData DataHeader
	err = connReader.Decode(&hostData)
	if err != nil {
		log.Panicf("[Error: %s]: Error recieving data from host\n", err.Error())
	}
	c.Other = CreatePartner(hostData.Name, hostData.PublicKey)

	fmt.Println("\nOther Peer Public Key:", hostData.PublicKey)

	//Send client data to host
	var name = "Client"
	if c.AppConfig.BroadcastName {
		name = c.AppConfig.Name
	}
	err = connWriter.Encode(DataHeader{name, *&c.ECCKeyPair.PubKey})
	if err != nil {
		log.Panicf("[Error: %s]: Error sending data to host\n", err.Error())
	}

	//Custom Input Reader
	cin := bufio.NewScanner(os.Stdin)

	//Chatting can begin
	for {
		fmt.Printf("Decrypted Message: %v\n\n", Recv(connReader, c))

		fmt.Printf("Enter Message to Send: ")
		cin.Scan()
		message := cin.Text()

		Send(message, c, connWriter)
	}

}

func (c *chatapp) Run() {
	if c.AppConfig.IsHost {
		c.runHost()
	} else {
		c.runClient()
	}
}
