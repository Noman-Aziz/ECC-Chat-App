package ecc

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
)

func Initialization(p int) (EllipticCurve, Keys) {

	var EC EllipticCurve
	var keys Keys

	EC.P = p

	EC.A = 0
	EC.B = 7

	//y^2 = x^3 + ax + b
	for {
		EC.G.X = rand.Intn(EC.P)
		EC.G.Y = rand.Intn(EC.P)

		var LHS int = (int(math.Pow(float64(EC.G.Y), 2)))
		LHS = Mod(LHS, EC.P)
		var RHS int = (int(math.Pow(float64(EC.G.X), 3)) + (EC.A * EC.G.X) + EC.B)
		RHS = Mod(RHS, EC.P)
		if LHS == RHS {
			break
		}
	}

	// 1 < PrivKey < P
	keys.PrivKey = rand.Intn(EC.P)
	for keys.PrivKey == 1 || keys.PrivKey == EC.P {
		keys.PrivKey = rand.Intn(EC.P)
	}

	//Public Key
	keys.PubKey = Mul(keys.PrivKey, EC.G, EC)

	return EC, keys
}

//Point Doubling and Addition
func Mul(k int, P ECPoint, EC EllipticCurve) ECPoint {
	var temp ECPoint = P

	kBinary := fmt.Sprintf("%b", k)
	//kBinary = kBinary[2:]

	for i := 0; i < len(kBinary); i++ {
		currentBit := kBinary[i]

		//always apply doubling
		temp = Add(temp, temp, EC)

		if currentBit == '1' {
			//add base point
			temp = Add(temp, P, EC)
		}

	}

	return temp
}

func Add(P ECPoint, Q ECPoint, EC EllipticCurve) ECPoint {
	var R ECPoint
	var lambda int
	var nom, denom int

	//Special Cases
	if isIdentity(P) && isIdentity(Q) {
		return ECPoint{0, 0}
	}
	if isIdentity(P) {
		return Q
	}
	if isIdentity(Q) {
		return P
	}

	//Case 1 (P != Q) => (y2 - y1) / (x2 - x1) mod p
	if P.X != Q.X && P.Y != Q.Y {
		nom = Mod((Q.Y - P.Y), EC.P)
		denom = modInv(Mod(Q.X-P.X, EC.P), EC.P)
	} else {
		if P == neg(Q, EC.P) {
			return ECPoint{0, 0}
		}
		if P.Y == 0 {
			return ECPoint{0, 0}
		}

		//Case 2 (P == Q) => (3 * x1^2) + a / (2 * y1) mod p
		nom = Mod(((3 * int(math.Pow(float64(P.X), 2))) + EC.A), EC.P)
		denom = modInv((2 * P.Y), EC.P)
	}

	lambda = nom * denom
	lambda = Mod(lambda, EC.P)

	//x3 = lambda^2 - x1 - x2 mod p
	R.X = int(math.Pow(float64(lambda), 2)) - P.X - Q.X
	R.X = Mod(R.X, EC.P)

	//y3 = lambda(x1 - x3) - y1 mod p
	R.Y = (lambda * (P.X - R.X)) - P.Y
	R.Y = Mod(R.Y, EC.P)

	return R
}

func Encrypt(M ECPoint, EC EllipticCurve, keys Keys) CipherText {

	var C CipherText

	//Random Positive Integer for Encryption
	var randomKey int = rand.Intn(EC.P)

	//Encryption
	//C = { kG, M + kPub }
	C.X = Mul(randomKey, EC.G, EC)        //kG
	C.Y = Mul(randomKey, keys.PubKey, EC) //kPub
	C.Y = Add(C.Y, M, EC)                 //M + kPub

	return C
}

func Decrypt(C CipherText, EC EllipticCurve, keys Keys) ECPoint {
	//Decryption
	//M = C2 - (PrivateKey * C1)
	var temp ECPoint = Mul(keys.PrivKey, C.X, EC) //PrivKey * C1
	temp.Y = temp.Y * -1                          //Curve is symmetric about x-axis

	return Add(C.Y, temp, EC)
}

func Encoding(str string) []ECPoint {

	var arr []ECPoint
	var temp ECPoint

	for i := 0; i < len(str); i++ {
		ascii, _ := strconv.Atoi(fmt.Sprintf("%d", str[i]))

		temp.X = ascii - 1
		temp.Y = (ascii * 2) - 1

		arr = append(arr, temp)
	}

	return arr
}

func Decoding(arr []ECPoint) string {

	var str string = ""

	for i := 0; i < len(arr); i++ {
		temp := arr[i]

		char := string(temp.X + 1)
		str = str + char
	}

	return string(str)
}
