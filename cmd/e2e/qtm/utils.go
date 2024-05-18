package main

import (
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"math/big"
	"os"
)

func PemToJwk(pemBytes []byte) string {
	keyStruct := processBlock(pemBytes)
	return structToJWK(keyStruct)
}

func findPemBlock(pemBytes []byte) (*pem.Block, []byte) {
	// Only get the first PEM block
	pemBlock, rest := pem.Decode(pemBytes)

	if pemBlock == nil {
		throwParseError("invalid PEM file format.")
	}

	if x509.IsEncryptedPEMBlock(pemBlock) {
		throwParseError("the given PEM file is encrypted. Please decrypt first.")
	}

	return pemBlock, rest
}

func processBlock(pemBytes []byte) interface{} {
	pemBlock, rest := findPemBlock(pemBytes)

	var keyStruct interface{}

	switch pemBlock.Type {
	case "PUBLIC KEY":
		keyStruct = processPublicKey(pemBlock.Bytes)
	case "RSA PRIVATE KEY":
		key, err := x509.ParsePKCS1PrivateKey(pemBlock.Bytes)
		stopOnParseError(err)

		keyStruct = processRSAPrivate(key)
	case "EC PARAMETERS":
		// The EC PARAMETERS section appears to not be needed…
		ecKey, _ := findPemBlock(rest)

		if ecKey.Type != "EC PRIVATE KEY" {
			throwParseError("unsupported EC PEM format.")
		}

		key, err := x509.ParseECPrivateKey(ecKey.Bytes)
		stopOnParseError(err)

		keyStruct = processECPrivate(key)
	default:
		throwParseError("unsupported PEM type.")
	}

	return keyStruct
}

func processPublicKey(bytes []byte) interface{} {
	key, err := x509.ParsePKIXPublicKey(bytes)
	stopOnParseError(err)

	var keyStruct interface{}

	switch key := key.(type) {
	case *rsa.PublicKey:
		keyStruct = processRSAPublic(key)
	case *ecdsa.PublicKey:
		keyStruct = processECPublic(key)
	default:
		throwParseError("Unknown key type.")
	}

	return keyStruct
}

func throwParseError(message string) {
	fmt.Fprintf(os.Stderr, "Could not parse key: %v\n", message)
	os.Exit(1)
}

func stopOnParseError(err error) {
	if err != nil {
		throwParseError(err.Error())
	}
}

type ECPublic struct {
	*JWK
	Curve string `json:"crv"`
	X     string `json:"x"`
	Y     string `json:"y"`
}

type ECPrivate struct {
	*ECPublic
	D string `json:"d"`
}

func processECPublic(key *ecdsa.PublicKey) *ECPublic {
	return &ECPublic{
		JWK:   &JWK{KeyType: "EC"},
		Curve: key.Params().Name,
		X:     b64EncodeBigInt(key.X),
		Y:     b64EncodeBigInt(key.Y),
	}
}

func processECPrivate(key *ecdsa.PrivateKey) *ECPrivate {
	return &ECPrivate{
		ECPublic: processECPublic(&key.PublicKey),
		D:        b64EncodeBigInt(key.D),
	}
}

func b64EncodeInt64(i64 int64) string {
	return b64EncodeBigInt(big.NewInt(i64))
}

func b64EncodeBigInt(bigInt *big.Int) string {
	return base64.RawURLEncoding.EncodeToString(bigInt.Bytes())
}

type JWK struct {
	KeyType string `json:"kty"`
}

func JwkToPem(jwkBytes []byte) string {
	throwParseError("JWK parsing is not implemented yet :(")
	return ""
}

func structToJWK(keyStruct interface{}) string {
	jwk, err := json.Marshal(keyStruct)

	if err != nil {
		fmt.Fprintf(os.Stderr, "error: could not render JSON.", err)
		os.Exit(4)
	}

	return string(jwk)
}

type RSAPublic struct {
	*JWK
	Modulus        string `json:"n"`
	PublicExponent string `json:"e"`
}

type RSAPrivate struct {
	*RSAPublic
	PrivateExponent string `json:"d"`
	Prime1          string `json:"p"`
	Prime2          string `json:"q"`
	Exponent1       string `json:"dp"`
	Exponent2       string `json:"dq"`
	Coefficient     string `json:"qi"`
}

func processRSAPublic(key *rsa.PublicKey) *RSAPublic {
	return &RSAPublic{
		JWK:            &JWK{KeyType: "RSA"},
		Modulus:        b64EncodeBigInt(key.N),
		PublicExponent: b64EncodeInt64(int64(key.E)),
	}
}

func processRSAPrivate(key *rsa.PrivateKey) *RSAPrivate {
	return &RSAPrivate{
		RSAPublic:       processRSAPublic(&key.PublicKey),
		PrivateExponent: b64EncodeBigInt(key.D),
		Prime1:          b64EncodeBigInt(key.Primes[0]),
		Prime2:          b64EncodeBigInt(key.Primes[1]),
		Exponent1:       b64EncodeBigInt(key.Precomputed.Dp),
		Exponent2:       b64EncodeBigInt(key.Precomputed.Dq),
		Coefficient:     b64EncodeBigInt(key.Precomputed.Qinv),
	}
}
