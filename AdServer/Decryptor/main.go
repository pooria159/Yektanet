package main

import (
	"fmt"

	"github.com/dgrijalva/jwt-go"
)

type Claim struct {
	Url string
	jwt.StandardClaims
}

type EventInfo struct {
	UserID      string
	PublisherID string
	AdID        string
	AdURL       string
	EventType   string

	jwt.StandardClaims
}

func main() {
	mySecretKey := []byte("Golangers:Pooria-Mohammad-Roya-Sina")
	/*claims := Claim {
		Url: "yahoo.com",
		StandardClaims: jwt.StandardClaims {
			IssuedAt: time.Now().Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims)
	tokenString, err := token.SignedString(mySecretKey)
	if err != nil {
		fmt.Printf("error while signing the token: %v\n", err)
		return
	}
	fmt.Printf("tokenString: %v\n", tokenString)*/

	var extractedClaims EventInfo
	tokenString := `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySUQiOiJubmxubG54bnNhZmtpb2d0cXJocm56cXpvc2ZlangiLCJQdWJsaXNoZXJJRCI6IjAiLCJBZElEIjoiNiIsIkFkVVJMIjoid3d3Lmdvb2dsZS5jb20iLCJFdmVudFR5cGUiOiJpbXByZXNzaW9uIn0.EbNwbu3QFLjQzD7Imi1VPa_ugmAqlrWNxdcjo4eFh0g`
	parsedToken, err := jwt.ParseWithClaims(tokenString, &extractedClaims, func(t *jwt.Token) (interface{}, error) {
		return mySecretKey, nil
	})
	if err != nil {
		fmt.Println(err)
	}
	if !parsedToken.Valid {
		fmt.Println("Parsed Token is Not Valid!")
		fmt.Printf("parsedToken: %v\n", parsedToken)
	}
	// Now we know that extractedClaimes stores the 'transmitted' information ...
	fmt.Printf("Event Info: %+v\n", extractedClaims)

}
