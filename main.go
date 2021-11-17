package main

import (
	"fmt"
	"math"
	"math/big"
	"math/rand"
	"strconv"
)

type ECPoint struct {
	x int
	y int
}

type EllipticCurve struct {
	p    int     //Integer Prime Field
	a, b int     //Curve Parameters i.e (4a^3 + 27.b^2) mod p != 0 mod p
	g    ECPoint //Base Point
}

type Keys struct {
	PubKey  ECPoint
	PrivKey int
}

type Text struct {
	PlainText     ECPoint
	CipherText    ECPoint
	DecryptedText ECPoint
}

func modInv(a int, n int) int {
	bigP := big.NewInt(int64(a))
	bigBase := big.NewInt(int64(n))

	// Compute inverse of bigP modulo bigBase
	bigGcd := big.NewInt(0)
	bigX := big.NewInt(0)
	bigGcd.GCD(bigX, nil, bigP, bigBase)

	// x*bigP+y*bigBase=1
	// => x*bigP = 1 modulo bigBase

	if int(bigGcd.Int64()) != 1 {
		str := "GCD of " + strconv.Itoa(a) + " wrt " + strconv.Itoa(n) + " != 1"
		panic(str)
	}

	return mod(int(bigX.Int64()), n)
}

func isIdentity(P ECPoint) bool {
	if P.x == 0 && P.y == 0 {
		return true
	}
	return false
}

func neg(P ECPoint, p int) ECPoint {
	return ECPoint{P.x, p - P.y}
}

func mod(n int, p int) int {
	if n < 0 {
		for n < 0 {
			n += p
		}
		return n
	} else {
		return n % p
	}
}

//Point Doubling and Addition
func mul(k int, P ECPoint, EC EllipticCurve) ECPoint {
	// 	let bits = bit_representation(s) # the vector of bits (from MSB to LSB) representing s
	// 	let res = O # point at infinity
	// 	for bit in bits:
	// 		res = res + res # double
	// 		if bit == 1:
	// 			res = res + P # add
	// 		i = i - 1
	// 	return res

	var temp ECPoint = P

	kBinary := fmt.Sprintf("%b", k)
	kBinary = kBinary[2:]

	for i := 1; i < len(kBinary); i++ {
		currentBit := kBinary[i : i+1]

		//always apply doubling
		temp = add(temp, temp, EC)

		if currentBit == "1" {
			//add base point
			temp = add(temp, P, EC)
		}

	}

	return temp
}

func add(P ECPoint, Q ECPoint, EC EllipticCurve) ECPoint {
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
	if P.x != Q.x && P.y != Q.y {
		nom = mod((Q.y - P.y), EC.p)
		denom = modInv(mod(Q.x-P.x, EC.p), EC.p)
	} else {
		if P == neg(Q, EC.p) {
			return ECPoint{0, 0}
		}
		if P.y == 0 {
			return ECPoint{0, 0}
		}

		//Case 2 (P == Q) => (3 * x1^2) + a / (2 * y1) mod p
		nom = mod(((3 * int(math.Pow(float64(P.x), 2))) + EC.a), EC.p)
		denom = modInv((2 * P.y), EC.p)
	}

	lambda = nom * denom
	lambda = mod(lambda, EC.p)

	//x3 = lambda^2 - x1 - x2 mod p
	R.x = int(math.Pow(float64(lambda), 2)) - P.x - Q.x
	R.x = mod(R.x, EC.p)

	//y3 = lambda(x1 - x3) - y1 mod p
	R.y = (lambda * (P.x - R.x)) - P.y
	R.y = mod(R.y, EC.p)

	return R
}

func main() {

	var EC EllipticCurve
	var keys Keys
	var text Text

	EC.p = 199
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

	//Random Key for Encryption
	var randomKey int = rand.Intn(EC.p)

	//Plain Text
	text.PlainText = ECPoint{5, 5}

	//Public Key
	keys.PubKey = mul(keys.PrivKey, EC.g, EC)

	//Encryption
	var c1 ECPoint = mul(randomKey, EC.g, EC)
	text.CipherText = mul(randomKey, keys.PubKey, EC)
	text.CipherText = add(text.CipherText, text.PlainText, EC)

	//Decryption
	//message = c2 - secretKey * c1
	var d ECPoint = mul(keys.PrivKey, c1, EC)
	d.y = d.y * -1 //curve is symmetric about x-axis. in this way, inverse point found
	text.DecryptedText = add(text.CipherText, d, EC)

	fmt.Println("MESSAGE :", text.PlainText)

	fmt.Println("\nPUBLIC KEY :", keys.PubKey)
	fmt.Println("PRIVATE KEY :", keys.PrivKey)

	fmt.Println("\nENCRYPTED TEXT:", text.CipherText)
	fmt.Println("DECRYPTED TEXT:", text.DecryptedText)
}
