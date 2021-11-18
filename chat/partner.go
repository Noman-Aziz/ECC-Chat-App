package chat

import (
	"github.com/Noman-Aziz/ECC-Chat-App/ecc"
)

type partner struct {
	Name      string       `json:"name"`
	PublicKey *ecc.ECPoint `json:"public_key"`
}

func CreatePartner(name string, publicKey ecc.ECPoint) *partner {
	//TODO: Add IP and Port as arguments, decode the code from host to get IP and Port
	return &partner{
		Name:      name,
		PublicKey: &publicKey,
	}
}

func (p *partner) EncryptMessage(message string, EC ecc.EllipticCurve) []ecc.CipherText {

	var cipher []ecc.CipherText
	var encoded []ecc.ECPoint

	//Encoding the Text
	encoded = ecc.Encoding(message)

	for i := 0; i < len(encoded); i++ {
		cipher = append(cipher, ecc.Encrypt(encoded[i], EC, *p.PublicKey))
	}

	return cipher
}

func DecryptMessage(cipher []ecc.CipherText, EC ecc.EllipticCurve, keys ecc.Keys) string {

	var decipher []ecc.ECPoint

	for i := 0; i < len(cipher); i++ {
		decipher = append(decipher, ecc.Decrypt(cipher[i], EC, keys))
	}

	return ecc.Decoding((decipher))
}
