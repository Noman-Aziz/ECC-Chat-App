package ecc

import (
	"math/big"
	"strconv"
)

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

	return Mod(int(bigX.Int64()), n)
}

func isIdentity(P ECPoint) bool {
	if P.X == 0 && P.Y == 0 {
		return true
	}
	return false
}

func neg(P ECPoint, p int) ECPoint {
	return ECPoint{P.X, p - P.Y}
}

func Mod(n int, p int) int {
	if n < 0 {
		for n < 0 {
			n += p
		}
		return n
	} else {
		return n % p
	}
}
