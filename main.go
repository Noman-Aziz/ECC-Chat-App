package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/Noman-Aziz/ECC-Chat-App/ecc"
)

func main() {

	var EC ecc.EllipticCurve
	var keys ecc.Keys

	var cipher []ecc.CipherText
	var decipher []ecc.ECPoint
	var encoded []ecc.ECPoint

	//Custom Input Reader
	cin := bufio.NewScanner(os.Stdin)

	EC, keys = ecc.Initialization(1844677)

	fmt.Println("\nPRIVATE KEY :", keys.PrivKey)
	fmt.Println("PUBLIC KEY :", keys.PubKey)

	//Taking Text Input
	fmt.Printf("\nEnter Message to Send: ")
	cin.Scan()
	message := cin.Text()

	//Encoding the Text
	encoded = ecc.Encoding(message)

	fmt.Print("\nENCRYPTED TEXT: ")
	for i := 0; i < len(encoded); i++ {
		cipher = append(cipher, ecc.Encrypt(encoded[i], EC, keys))
		fmt.Print("{", cipher[i].X, ",", cipher[i].Y, "} ")
	}
	fmt.Println()

	for i := 0; i < len(cipher); i++ {
		decipher = append(decipher, ecc.Decrypt(cipher[i], EC, keys))
	}
	fmt.Println("\nDECRYPTED TEXT:", ecc.Decoding((decipher)))
}
