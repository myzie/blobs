package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/myzie/base"

	jwt "github.com/dgrijalva/jwt-go"
)

func main() {

	var (
		userID   string
		userName string
		keyPath  string
		isAdmin  bool
	)

	flag.StringVar(&userID, "user-id", "", "User ID")
	flag.StringVar(&userName, "user-name", "", "User name")
	flag.StringVar(&keyPath, "key", "", "Private key for signing")
	flag.BoolVar(&isAdmin, "admin", false, "Admin")
	flag.Parse()

	claims := &base.JWTClaims{
		Name:  userName,
		Admin: isAdmin,
		StandardClaims: jwt.StandardClaims{
			Subject:   userID,
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodRS512, claims)

	// Load RSA private key from disk
	keyText, err := ioutil.ReadFile(keyPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	key, err := jwt.ParseRSAPrivateKeyFromPEM(keyText)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Generate signed token
	t, err := token.SignedString(key)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("export BLOBS_TOKEN=%s\n", t)
}
