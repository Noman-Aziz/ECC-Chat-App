package chat

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/Noman-Aziz/ECC-Chat-App/ecc"
)

func Send(message string, c *chatapp, connWriter *json.Encoder) {
	encrypted := c.Other.EncryptMessage(message, *c.EC)

	err := connWriter.Encode(encrypted)
	if err != nil {
		log.Printf("[Error: %s]: Error sending message to client\n", err.Error())
	}
}

func Recv(connReader *json.Decoder, c *chatapp) string {
	var cipher []ecc.CipherText
	err := connReader.Decode(&cipher)
	if err != nil {
		log.Printf("[Error: %s]: Error receiving message from host\n", err.Error())
	}

	fmt.Println("\nEncrypted Message:", cipher)

	decrypted := DecryptMessage(cipher, *c.EC, *c.ECCKeyPair)

	return string(decrypted)
}
