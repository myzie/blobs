package main

import jwt "github.com/dgrijalva/jwt-go"

// jwtCustomClaims are custom claims extending default ones
type jwtCustomClaims struct {
	Name    string `json:"name,omitempty"`
	Admin   bool   `json:"admin,omitempty"`
	Context string `json:"context,omitempty"`
	jwt.StandardClaims
}
