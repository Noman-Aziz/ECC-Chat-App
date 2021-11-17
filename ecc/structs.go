package ecc

type ECPoint struct {
	X int
	Y int
}

type EllipticCurve struct {
	P    int     //Integer Prime Field
	A, B int     //Curve Parameters i.e (4a^3 + 27.b^2) mod p != 0 mod p
	G    ECPoint //Base Point
}

type Keys struct {
	PubKey        ECPoint
	PrivKey       int
	RandomInteger int //Used in Encryption
}

type CipherText struct {
	X ECPoint
	Y ECPoint
}
