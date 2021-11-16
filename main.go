package main

import (
	"fmt"
	"math"
	"math/big"
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

	return int(bigX.Int64())
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

func mul(k int, P ECPoint, EC EllipticCurve) ECPoint {
	var R ECPoint
	var temp []ECPoint

	for i := 1; i <= k-1; i++ {
		//1P = P
		if i == 1 {
			temp = append(temp, P)
		} else {
			//nP = (nP-1) + P => 2P = P + P => 3P = 2P + P
			temp = append(temp, (add(temp[i-2], P, EC)))
		}
	}

	for i := 0; i < len(temp); i++ {
		if i == 0 {
			R = temp[0]
		} else {
			R = add(R, temp[i], EC)
		}
	}

	return R
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
	if P.x != Q.x {
		nom = Q.y - P.y
		denom = modInv(mod(Q.x-P.x, EC.p), EC.p)
	} else {
		if P == neg(Q, EC.p) {
			return ECPoint{0, 0}
		}
		if P.y == 0 {
			return ECPoint{0, 0}
		}

		//Case 2 (P == Q) => (3 * x1^2) + a / (2 * y1) mod p
		nom = (3 * int(math.Pow(float64(P.x), 2))) + EC.a
		denom = modInv((2 * P.y), EC.p)
	}

	lambda = nom * denom
	lambda = mod(lambda, EC.p)

	//x3 = lambda^2 - x1 - x2 mod p
	R.x = (int(math.Pow(float64(lambda), 2))) - P.x - Q.x
	R.x = mod(R.x, EC.p)

	//y3 = lambda(x1 - x3) - y1 mod p
	R.y = (lambda * (P.x - R.x)) - P.y
	R.y = mod(R.y, EC.p)

	return R
}

func main() {
	//EQ
	// y^2 = x^3 + ax + b

	var EC EllipticCurve
	var keys Keys

	// EC.p = 256
	// EC.a = 9
	// EC.b = 17

	//y^2 = x^3 + ax + b
	// for {
	// 	EC.g.x = rand.Intn(256)
	// 	EC.g.y = rand.Intn(256)

	// 	var LHS int = (int(math.Pow(float64(EC.g.y), 2)))
	// 	LHS = mod(LHS, EC.p)
	// 	var RHS int = (int(math.Pow(float64(EC.g.x), 3)) + (EC.a * EC.g.x) + EC.b)
	// 	RHS = mod(RHS, EC.p)
	// 	if LHS == RHS {
	// 		break
	// 	}
	// }

	// 1 < PrivKey < P
	// keys.PrivKey = rand.Intn(EC.p)
	// for keys.PrivKey == 1 || keys.PrivKey == EC.p {
	// 	keys.PrivKey = rand.Intn(EC.p)
	// }

	EC.p = 751
	EC.a = -1
	EC.b = 188
	keys.PrivKey = 386

	// Public Key
	keys.PubKey = mul(keys.PrivKey, EC.g, EC)

	var enc ECPoint = add(ECPoint{562, 201}, keys.PubKey, EC)
	var dec ECPoint = mul(keys.PrivKey, enc, EC)

	fmt.Println("MESSAGE : {562, 201}")
	fmt.Println("PUBLIC KEY :", keys.PubKey)
	fmt.Println("ENCRYPTION :", enc)
	fmt.Println("DECRYPT :", dec)
}
