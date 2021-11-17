package main

import (
	"fmt"

	"github.com/Noman-Aziz/ECC-Chat-App/ecc"
)

func main() {

	var EC ecc.EllipticCurve
	var keys ecc.Keys

	var cipher ecc.CipherText
	var decipher ecc.ECPoint

	EC, keys = ecc.Initialization(257)

	cipher = ecc.Encrypt(ecc.ECPoint{5, 5}, EC, keys)
	decipher = ecc.Decrypt(cipher, EC, keys)

	fmt.Println("MESSAGE : {5, 5}")

	fmt.Println("\nPRIVATE KEY :", keys.PrivKey)
	fmt.Println("PUBLIC KEY :", keys.PubKey)

	fmt.Println("\nENCRYPTED TEXT: {", cipher.X, ",", cipher.Y, "}")
	fmt.Println("DECRYPTED TEXT:", decipher)
}
