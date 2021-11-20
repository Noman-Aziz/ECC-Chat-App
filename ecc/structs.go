package ecc

import "math/big"

type ECPoint struct {
	X *big.Int
	Y *big.Int
}

type EllipticCurve struct {
	P    *big.Int //Integer Prime Field
	A, B *big.Int //Curve Parameters i.e (4a^3 + 27.b^2) mod p != 0 mod p
	G    ECPoint  //Base Point
}

type Keys struct {
	PubKey        ECPoint
	PrivKey       *big.Int
	RandomInteger *big.Int //Used in Encryption
}

type CipherText struct {
	X ECPoint
	Y ECPoint
}
