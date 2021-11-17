package ecc

import (
	"fmt"
	"math"
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
	CipherTextX   ECPoint
	CipherTextY   ECPoint
	DecryptedText ECPoint
}

//Point Doubling and Addition
func mul(k int, P ECPoint, EC EllipticCurve) ECPoint {
	var temp ECPoint = P

	kBinary := fmt.Sprintf("%b", k)
	//kBinary = kBinary[2:]

	for i := 0; i < len(kBinary); i++ {
		currentBit := kBinary[i]

		//always apply doubling
		temp = add(temp, temp, EC)

		if currentBit == '1' {
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
