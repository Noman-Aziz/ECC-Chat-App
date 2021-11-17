package main

import (
	"fmt"
	"math"
	"math/rand"
)

func main() {

	var EC EllipticCurve
	var keys Keys
	var text Text

	EC.p = 257

	EC.a = 0
	EC.b = 7

	//y^2 = x^3 + ax + b
	for {
		EC.g.x = rand.Intn(EC.p)
		EC.g.y = rand.Intn(EC.p)

		var LHS int = (int(math.Pow(float64(EC.g.y), 2)))
		LHS = mod(LHS, EC.p)
		var RHS int = (int(math.Pow(float64(EC.g.x), 3)) + (EC.a * EC.g.x) + EC.b)
		RHS = mod(RHS, EC.p)
		if LHS == RHS {
			break
		}
	}

	// 1 < PrivKey < P
	keys.PrivKey = rand.Intn(EC.p)
	for keys.PrivKey == 1 || keys.PrivKey == EC.p {
		keys.PrivKey = rand.Intn(EC.p)
	}

	//Random Positive Integer for Encryption
	var randomKey int = rand.Intn(EC.p)

	//Plain Text
	text.PlainText = ECPoint{5, 5}

	//Public Key
	keys.PubKey = mul(keys.PrivKey, EC.g, EC)

	//Encryption
	//C = { kG, M + kPub }
	text.CipherTextX = mul(randomKey, EC.g, EC)                  //kG
	text.CipherTextY = mul(randomKey, keys.PubKey, EC)           //kPub
	text.CipherTextY = add(text.CipherTextY, text.PlainText, EC) //M + kPub

	//Decryption
	//M = C2 - (PrivateKey * C1)
	var temp ECPoint = mul(keys.PrivKey, text.CipherTextX, EC) //PrivKey * C1
	temp.y = temp.y * -1                                       //Curve is symmetric about x-axis
	text.DecryptedText = add(text.CipherTextY, temp, EC)

	fmt.Println("MESSAGE :", text.PlainText)

	fmt.Println("\nPRIVATE KEY :", keys.PrivKey)
	fmt.Println("PUBLIC KEY :", keys.PubKey)

	fmt.Println("\nENCRYPTED TEXT: {", text.CipherTextX, ",", text.CipherTextY, "}")
	fmt.Println("DECRYPTED TEXT:", text.DecryptedText)
}
