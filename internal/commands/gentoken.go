package commands

import (
	"crypto/rsa"
	"fmt"
	auth2 "github.com/restaurant/internal/auth"
	"os"
	"strconv"
	"time"

	"github.com/pkg/errors"

	"github.com/restaurant/internal/auth"
)

// GenToken generates a JWT for the specified user.
func GenToken(userClaims auth2.ClaimsAuth, privateKeyFile string) (string, error) {
	if userClaims.ID == 0 || privateKeyFile == "" {
		fmt.Println("help: gentoken <id> <private_key_file> <algorithm>")
		fmt.Println("algorithm: RS256, HS256")
		return "", ErrHelp
	}

	// The call to retrieve a user requires an Admin role by the caller.
	claims := auth.Claims{
		StandardClaims: jwt.StandardClaims{
			Subject: strconv.Itoa(int(userClaims.ID)),
		},
		Role: auth.RoleAdmin,
	}

	privatePEM, err := os.ReadFile(privateKeyFile)
	if err != nil {
		return "", errors.Wrap(err, "reading PEM private key file")
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privatePEM)
	if err != nil {
		return "", errors.Wrap(err, "parsing PEM into private key")
	}

	// In a production system, a key id (KID) is used to retrieve the correct
	// public key to parse a JWT for auth and claims. A key lookup function is
	// provided to perform the task of retrieving a KID for a given public key.
	// In this code, I am writing a lookup function that will return the public
	// key for the private key provided with an arbitrary KID.
	keyID := "54bb2165-71e1-41a6-af3e-7da4a0e1e2c1"
	lookup := func(kid string) (*rsa.PublicKey, error) {
		switch kid {
		case keyID:
			return &privateKey.PublicKey, nil
		}
		return nil, fmt.Errorf("no public key found for the specified kid: %s", kid)
	}

	// An authenticator maintains the state required to handle JWT processing.
	// It requires the private key for generating tokens. The KID for access
	// to the corresponding public key, the algorithms to use (RS256), and the
	// key lookup function to perform the actual retrieve of the KID to public
	// key lookup.
	a, err := auth.New("RS256", lookup, auth.Keys{keyID: privateKey})
	if err != nil {
		return "", errors.Wrap(err, "constructing auth")
	}

	// Generating a token requires defining a set of claims. In this applications
	// case, we only care about defining the subject and the user in question and
	// the roles they have on the database. This token will expire in a year.
	//
	// iss (issuer): Issuer of the JWT
	// sub (subject): Subject of the JWT (the user)
	// aud (audience): Recipient for which the JWT is intended
	// exp (expiration time): Time after which the JWT expires
	// nbf (not before time): Time before which the JWT must not be accepted for processing
	// iat (issued at time): Time at which the JWT was issued; can be used to determine age of the JWT
	// jti (JWT ID): Unique identifier; can be used to prevent the JWT from being replayed (allows a token to be used only once)
	if userClaims.Roles == auth.RoleAdmin && userClaims.RestaurantID == nil {
		return "", errors.New("token gen: role admin doesn't contain restaurant_id")
	}

	claims = auth.Claims{
		StandardClaims: jwt.StandardClaims{
			Issuer:    "restaurant",
			Subject:   strconv.Itoa(int(userClaims.ID)),
			ExpiresAt: time.Now().Add(8760 * time.Hour).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		Role:         userClaims.Roles,
		UserId:       userClaims.ID,
		RestaurantID: userClaims.RestaurantID,
		BranchID:     userClaims.BranchID,
	}

	// This will generate a JWT with the claims embedded in them. The database
	// with need to be configured with the information found in the public key
	// file to validate these claims. Dgraph does not support key rotate at
	// this time.
	token, err := a.GenerateToken(keyID, claims)
	if err != nil {
		return "", errors.Wrap(err, "generating token")
	}

	fmt.Printf("-----BEGIN TOKEN-----\n%s\n-----END TOKEN-----\n", token)
	return token, nil
}
