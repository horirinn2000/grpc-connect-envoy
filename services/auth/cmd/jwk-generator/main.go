package main

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"log"
	"math/big"
	"os"
	"path/filepath"
)

// JWK is a JSON Web Key.
type JWK struct {
	Kty string `json:"kty"`
	Kid string `json:"kid"`
	Use string `json:"use"`
	N   string `json:"n"`
	E   string `json:"e"`
}

// JWKS is a set of JSON Web Keys.
type JWKS struct {
	Keys []JWK `json:"keys"`
}

func main() {
	// Define paths relative to the project root
	pemPath := filepath.Join(".", "keys", "public.pem")
	jwksPath := filepath.Join(".", "keys", "public.jwks.json")

	// Read the public key from the PEM file
	pemData, err := os.ReadFile(pemPath)
	if err != nil {
		log.Fatalf("failed to read public key PEM file: %v", err)
	}

	// Decode the PEM block
	block, _ := pem.Decode(pemData)
	if block == nil || block.Type != "PUBLIC KEY" {
		log.Fatal("failed to decode PEM block containing public key")
	}

	// Parse the public key
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		log.Fatalf("failed to parse DER encoded public key: %v", err)
	}

	rsaPub, ok := pub.(*rsa.PublicKey)
	if !ok {
		log.Fatal("public key is not of type RSA")
	}

	// Create the JWK
	jwk := JWK{
		Kty: "RSA",
		Kid: "auth-key", // This should match your key ID
		Use: "sig",
		N:   base64.RawURLEncoding.EncodeToString(rsaPub.N.Bytes()),
		E:   base64.RawURLEncoding.EncodeToString(big.NewInt(int64(rsaPub.E)).Bytes()),
	}

	// Create the JWKS
	jwks := JWKS{
		Keys: []JWK{jwk},
	}

	// Marshal the JWKS to JSON
	jwksJSON, err := json.MarshalIndent(jwks, "", "  ")
	if err != nil {
		log.Fatalf("failed to marshal JWKS to JSON: %v", err)
	}

	// Write the JWKS to a file
	if err := os.WriteFile(jwksPath, jwksJSON, 0644); err != nil {
		log.Fatalf("failed to write JWKS file: %v", err)
	}

	fmt.Printf("Successfully generated JWKS file at %s\n", jwksPath)
}
