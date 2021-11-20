package ecc

import (
	"fmt"
	"math/big"
	"strconv"
)

func Initialization() (EllipticCurve, Keys) {

	var EC EllipticCurve
	var keys Keys
	var ok bool

	EC.P = new(big.Int)
	EC.P, ok = EC.P.SetString("115792089237316195423570985008687907852837564279074904382605163141518161494337", 10)
	if !ok {
		panic("SetString: error")
	}

	EC.A = big.NewInt(0)
	EC.B = big.NewInt(7)

	//y^2 = x^3 + ax + b
	// for {
	// 	EC.G.X = CreateRandomPrime(p - 1)
	// 	EC.G.Y = CreateRandomPrime(p - 1)

	// 	var LHS *big.Int = EC.G.Y.Exp(EC.G.Y, big.NewInt(2), EC.P)
	// 	LHS = Mod(LHS, EC.P)

	// 	var RHS *big.Int = EC.G.Y.Exp(EC.G.X, big.NewInt(3), EC.P)
	// 	RHS = RHS.Add(RHS, EC.G.X.Mul(EC.G.X, EC.A))
	// 	RHS = RHS.Add(RHS, EC.B)
	// 	RHS = Mod(RHS, EC.P)

	// 	if LHS.Cmp(RHS) == 0 {
	// 		break
	// 	}
	// }

	EC.G.X = new(big.Int)
	EC.G.Y = new(big.Int)

	EC.G.X, ok = EC.G.X.SetString("55066263022277343669578718895168534326250603453777594175500187360389116729240", 10)
	if !ok {
		panic("SetString: error")
	}

	EC.G.Y, ok = EC.G.Y.SetString("32670510020758816978083085130507043184471273380659243275938904335757337482424", 10)
	if !ok {
		panic("SetString: error")
	}

	// 1 < PrivKey < P
	keys.PrivKey = CreateRandomInt(EC.P)
	for keys.PrivKey.Cmp(big.NewInt(1)) == 0 || keys.PrivKey.Cmp(EC.P) == 0 {
		keys.PrivKey = CreateRandomInt(EC.P)
	}

	//Public Key
	keys.PubKey = Mul(keys.PrivKey, EC.G, EC)

	//Public Key not Generated Successfully
	for keys.PubKey.X.Cmp(big.NewInt(0)) == 0 && keys.PubKey.Y.Cmp(big.NewInt(0)) == 0 {

		keys.PrivKey = CreateRandomInt(EC.P)
		for keys.PrivKey.Cmp(big.NewInt(1)) == 0 || keys.PrivKey.Cmp(EC.P) == 0 {
			keys.PrivKey = CreateRandomInt(EC.P)
		}

		keys.PubKey = Mul(keys.PrivKey, EC.G, EC)

	}

	return EC, keys
}

//Point Doubling and Addition
func Mul(k *big.Int, P ECPoint, EC EllipticCurve) ECPoint {
	var temp ECPoint = P

	kBinary := fmt.Sprintf("%b", k)

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
	var lambda *big.Int
	var nom, denom *big.Int

	//Special Cases
	if isIdentity(P) && isIdentity(Q) {
		return ECPoint{big.NewInt(0), big.NewInt(0)}
	}
	if isIdentity(P) {
		return Q
	}
	if isIdentity(Q) {
		return P
	}

	//Case 1 (P != Q) => (y2 - y1) / (x2 - x1) mod p
	if P.X != Q.X && P.Y != Q.Y {

		nom = Mod(Mod(Q.Y.Sub(Q.Y, P.Y), EC.P), EC.P)

		denom = Mod(Q.X.Sub(Q.X, P.X), EC.P)
		denom = modInv(denom, EC.P)
	} else {

		if P == neg(Q, EC.P) {
			return ECPoint{big.NewInt(0), big.NewInt(0)}
		}
		if P.Y.Cmp(big.NewInt(0)) == 0 {
			return ECPoint{big.NewInt(0), big.NewInt(0)}
		}

		//Case 2 (P == Q) => (3 * x1^2) + a / (2 * y1) mod p
		nom = P.X.Exp(P.X, big.NewInt(2), EC.P)
		nom = nom.Mul(nom, big.NewInt(3))
		nom = nom.Add(nom, EC.A)
		nom = Mod(nom, EC.P)

		denom = modInv(Mod(P.Y.Mul(P.Y, big.NewInt(2)), EC.P), EC.P)
	}

	lambda = nom.Mul(nom, denom)
	lambda = Mod(lambda, EC.P)

	//x3 = lambda^2 - x1 - x2 mod p
	R.X = lambda.Exp(lambda, big.NewInt(2), EC.P)
	R.X = R.X.Sub(R.X, P.X)
	R.X = R.X.Sub(R.X, Q.X)
	R.X = Mod(R.X, EC.P)

	//y3 = lambda(x1 - x3) - y1 mod p
	R.Y = P.X.Sub(P.X, R.X)
	R.Y = R.Y.Mul(R.Y, lambda)
	R.Y = R.Y.Sub(R.Y, P.Y)
	R.Y = Mod(R.Y, EC.P)

	return R
}

func Encrypt(M ECPoint, EC EllipticCurve, publicKey ECPoint) CipherText {

	var C CipherText

	//Random Positive Integer for Encryption
	var randomKey *big.Int = new(big.Int)
	var ok bool

	randomKey, ok = randomKey.SetString("28695618543805844332113829720373285210420739438570883203839696518176414791234", 10)
	if !ok {
		panic("SetString: error")
	}

	//Encryption
	//C = { kG, M + kPub }
	C.X = Mul(randomKey, EC.G, EC)      //kG
	C.Y = Mul(randomKey, publicKey, EC) //kPub
	C.Y = Add(C.Y, M, EC)               //M + kPub

	return C
}

func Decrypt(C CipherText, EC EllipticCurve, keys Keys) ECPoint {
	//Decryption
	//M = C2 - (PrivateKey * C1)
	var temp ECPoint = Mul(keys.PrivKey, C.X, EC) //PrivKey * C1
	temp.Y = temp.Y.Mul(temp.Y, big.NewInt(-1))   //Curve is symmetric about x-axis

	return Add(C.Y, temp, EC)
}

func Encoding(str string) []ECPoint {

	var arr []ECPoint
	var temp ECPoint

	for i := 0; i < len(str); i++ {
		ascii, _ := strconv.Atoi(fmt.Sprintf("%d", str[i]))

		temp.X = big.NewInt(int64(ascii))
		temp.X = temp.X.Sub(temp.X, big.NewInt(1))

		temp.Y = big.NewInt(int64(ascii))
		temp.Y = temp.Y.Mul(temp.Y, big.NewInt(2))
		temp.Y = temp.Y.Sub(temp.Y, big.NewInt(1))

		arr = append(arr, temp)
	}

	return arr
}

func Decoding(arr []ECPoint) string {

	var str string = ""

	for i := 0; i < len(arr); i++ {
		temp := arr[i]

		char := string(temp.X.Add(temp.X, big.NewInt(1)).Int64())
		str = str + char
	}

	return string(str)
}
