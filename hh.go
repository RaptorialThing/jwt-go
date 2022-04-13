package main

import (
	"fmt"
	"strings"
	bs64 "encoding/base64"
	"encoding/json"
	//"github.com/golang-jwt/jwt"
	
)

type CustomClaims struct {
	GUID string `json:"GUID"`
	iss string `json:"iss"`
}

func decodeTokenClaims(token string) (CustomClaims, error) {
	s := strings.Split(token,".")
	decodedPayload,_ := bs64.StdEncoding.DecodeString(s[1])
	//decodedPayloadStr := string(decodedPayload)
	//claims := make(map[string]CustomClaims)
	claimSlice := make(map[string]string)

	err := json.Unmarshal(decodedPayload, &claimSlice)

	claim := CustomClaims{}

	fmt.Println(claimSlice)
	for key := range claimSlice {
		if key == "GUID" {
			claim.GUID = claimSlice[key]
		}
		if key == "iss" {
			claim.iss = claimSlice[key]
		}
	}

	return claim,err
}


func main() {
	t := "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJHVUlEIjoiZThmMzkzMzEtYmMyZS00MzkyLTk3YjEtMjMyOGIzYzYzYWI2IiwiaXNzIjoidGVzdCJ9.vVBrlUGaa5ONsRQ_bUDF45BNUxJHO-Lpa8dmRSP330Nu8cPHHS52QTcTb083ORvVR53gdHVjyBuPL4l8vb9Sfw"
	claims, err:= decodeTokenClaims(t)
	if err != nil {
		//claims = "error "+err.Error()
	}

	fmt.Println(claims)
}