package main

import (
	"net/http"
	"time"
	"strconv"
	"math/rand"
	//"errors"
	"fmt"
	"strings"

	bs64 "encoding/base64"
	"encoding/json"
	"github.com/labstack/echo"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"example.com/database"
)

var refreshHmacSampleSecret []byte
var accessHmacSampleSecret []byte
var Refresh_base_hash string
var AccessTokenSign string
var RefreshTokenSign string
const RandSecretSize = 88 // int secret lenth
const BcryptCost = 8 // int bcrypt cost, rounds 2^BcryptCost
const GUIDFormat = 36 // GUID  128 bit id characters count [36]byte

type User struct {
	guid [GUIDFormat]byte `form:"GUID" query:"GUID" json:"GUID" bson: "GUID"`
	refresh_token string `form:"refresh_token" 
	query:"refresh_token" json: "refresh_token" bson: "refresh_token"`
}

type CustomClaims struct {
	GUID string `json:"GUID"`
	exp string `json:"exp"`
	accessTokenHash string `json:"accessTokenHash"`
	jwt.StandardClaims
}

func findUser(guid [GUIDFormat]byte ) (*database.UserMongo, bool) {

	var res *database.UserMongo
	var resBool bool
	resBool = false

	userFound,err := database.GetUser(guid)

	if err == nil {
		res =  userFound[len(userFound)-1]
		resBool = true
	}

	return res,resBool
}

// func createUser(id [GUIDFormat]byte, token string) (User,error) {
// 	token,err := hashToken(token)
// 	u := User {
// 		guid: id,
// 		refresh_token: token}
// 	usr, err := saveUser(u)
// 	if err != nil {
// 		return User{}, err
// 	}
// 	return usr,err
// }

// func updateUser(id [GUIDFormat]byte, token string) (User, error) {
// 	u,errFind := findUser(id)
// 	if errFind  {
// 		return User{}, errors.New("find User error")
// 	}
// 	new_token, errToken := hashToken(token)
// 	if errToken != nil {
// 		return User{}, errToken
// 	}
// 	u.refresh_token = new_token
// 	usr,errSave := saveUser(u)
// 	if errSave != nil {
// 		return User{}, errSave
// 	}

// 	return usr, nil
// }

// func saveUser(u User) (User, error) {
// 	// coonect to db and save User or errors.New("cant save user to database")
// 	usr := User{}
// 	return usr, nil
// }

func hashToken(token string) (string,error) {
	bytes,err := bcrypt.GenerateFromPassword([]byte(token),BcryptCost)
	return string(bytes), err
}

func checkTokenHash(token string, hash []byte) bool {
	err := bcrypt.CompareHashAndPassword(hash, []byte(token))
	return err == nil
}

func generateUserTokens(guid [GUIDFormat]byte) (string, string, error) {
	strGuid := make([]string, len(guid))
	for i:= range guid {
		strGuid[i] = string(int(guid[i]))
	}
	stringGuid := strings.Join(strGuid,"")

	timeAccessInt := time.Now().Add(time.Minute * 2).Unix()
	timeAccess := strconv.FormatInt(timeAccessInt,10)
	timeAccess = string(timeAccess)

	accessTokenHashEmpty := ""

	accessClaims := CustomClaims{
		stringGuid,
		timeAccess,
		accessTokenHashEmpty, 
		jwt.StandardClaims{
			Issuer:    "test",
			ExpiresAt: timeAccessInt,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, accessClaims)
	accessHmacSampleSecret = []byte(randSecret(RandSecretSize))

	accessTokenString, err := token.SignedString(accessHmacSampleSecret)

	timeRefreshInt := time.Now().Add(time.Hour * 24).Unix()
	timeRefresh := strconv.FormatInt(timeRefreshInt,10)
	timeRefresh = string(timeRefresh)

	accessTokenHash, errAccTokenHash := hashToken(accessTokenString)
	if errAccTokenHash != nil {
		err = errAccTokenHash
	}
	accessTokenHash = string(accessTokenHash)

	refreshClaims := CustomClaims{
		stringGuid,
		timeRefresh,
		accessTokenHash,
		jwt.StandardClaims{
			Issuer:    "test",
			ExpiresAt: timeRefreshInt,
		},
	}
	fmt.Println(refreshClaims)
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS512,refreshClaims)
	refreshHmacSampleSecret = []byte(randSecret(RandSecretSize))
	refreshTokenString, err := refreshToken.SignedString(refreshHmacSampleSecret)
	Refresh_base,errHashRefreshToken := hashToken(refreshTokenString)
	Refresh_base_hash = Refresh_base
	
	if err == nil {
		err = errHashRefreshToken
	}

	return accessTokenString, refreshTokenString, err
}

func randSecret(n int) string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	b := make([]byte, n)
    for i := range b {
        b[i] = letterBytes[rand.Intn(len(letterBytes))]
    }
    return string(b)
}

func parseClaims(token string) (CustomClaims, bool) {
	ok := false
	claims := CustomClaims{}
	
	s := strings.Split(token,".")
	decodedPayload,_ := bs64.StdEncoding.DecodeString(s[1])
	//claims := make(map[string]CustomClaims)
	claimSlice := make(map[string]string)
	
	err := json.Unmarshal(decodedPayload, &claimSlice)

	if err != nil {
		return claims, false
	}

	claim := CustomClaims{}

	for key := range claimSlice {
		if key == "GUID" {
			claim.GUID = claimSlice[key]
			ok = true
		}
		if key == "accessTokenHash" {
			claim.accessTokenHash = claimSlice[key]
		}
	}

	return claim,ok
}

func ConvertGuid(guid string) ([GUIDFormat]byte, bool) {
	guidBytes := []byte(guid)
	var err bool
	var guidBytesFixed [GUIDFormat]byte
	err = false

	if len([GUIDFormat]byte{}) != len(guidBytes) {
		err = true
		return [GUIDFormat]byte{}, err
	}

	for i := range [GUIDFormat]byte{} {
		guidBytesFixed[i] = guidBytes[i]
	}

	return guidBytesFixed, err
}


func main() {
	e := echo.New()
	e.GET("/",func(c echo.Context) error {
		return  c.String(http.StatusOK, "Hello,world")
	})

	e.POST("/authenticate",func(c echo.Context) error {
		formGuid,errConvertGuid := ConvertGuid(c.FormValue("GUID"))

		if errConvertGuid != false {
			return c.JSON(http.StatusUnprocessableEntity, map[string]string{
				"error": c.FormValue("GUID")+" GUID cant format",
			})
		}

		if _, errFindUser := findUser(formGuid); errFindUser == true  {

			return c.String(http.StatusOK, "Error: User already authenticated "+c.FormValue("GUID"))
		}

	
		accessTokenString, refreshTokenString, err := generateUserTokens(formGuid)	

		if err != nil  {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"error": err.Error(),
			})
	    } else {

				guidBytes,errConvert := ConvertGuid(c.FormValue("GUID"))

				if errConvert {
					return c.JSON(http.StatusInternalServerError, map[string]string{
						"error": "guid format",
					})
				}

				userToSave := database.User{}
				userToSave.Guid = guidBytes
				userToSave.Refresh_token = Refresh_base_hash

				resSaveUsr := database.SaveUser(userToSave)
				if !resSaveUsr {
						return c.JSON(http.StatusInternalServerError, map[string]string{
							"error": "save token hash",
						})
				}


			return c.JSON(http.StatusOK, map[string]string{
				"access-token": accessTokenString,
				"refresh-token": refreshTokenString,
			})
		}

		return echo.ErrUnauthorized

	})

	e.POST("/refresh-tokens",func(c echo.Context) error {

		refresh_token := c.FormValue("refresh-token")

		if refresh_token == "" {
			return c.JSON(http.StatusUnprocessableEntity,
				map[string]string{"error":"give me refresh-token"})
		}
		
		// decodedToken,errParse := jwt.ParseWithClaims(refresh_token, &CustomClaims{},
		// 	func(token *jwt.Token) (interface{}, error) {

		// 		return refreshHmacSampleSecret, nil

		// 	})

		var result map[string]string
		status := http.StatusOK
			
			claims, okClaims := parseClaims(refresh_token)

			if okClaims  {
				
				
			}	else {
				status = http.StatusInternalServerError
				result = map[string]string{
					"error": "can not decode token",
				}
				return c.JSON(status, result)
			}

			guid, errAsserGuid := ConvertGuid(claims.GUID)
			if errAsserGuid  {
				status = http.StatusInternalServerError
				result = map[string]string{
					"error": "can not format GUID from token",
				}
				return c.JSON(status, result)
			}	

		userWithGuid,okUser := findUser(guid)
		if !okUser {
			status = http.StatusInternalServerError
			result = map[string]string{
				"error": "can not check user refresh_token ",
			}
			return c.JSON(status, result)
		}

		user_refresh_hash := userWithGuid.Refresh_token
		match := checkTokenHash(refresh_token, []byte(user_refresh_hash))

		if match {
			accessTokenString, refreshTokenString, err := generateUserTokens(guid)

			status = http.StatusOK
			result = map[string]string{
				"access-token": accessTokenString,
				"refresh-token": refreshTokenString,
			}

			if err != nil {
				status = http.StatusInternalServerError
				result = map[string]string {
					"error": "Error: tokens did not generate",}
				
				} else {
					userToUpdate := database.User{}
					userToUpdate.Guid = guid
					userToUpdate.Refresh_token = Refresh_base_hash
					fmt.Println(claims.GUID)
					errResUpd := database.UpdateUser(userToUpdate)
		
					if !errResUpd {

						status = http.StatusInternalServerError
						result = map[string]string{
							"error": "update refresh token hash",
						}

						return c.JSON(status, result)
					}
			}

		} else {
			status = http.StatusUnauthorized
			result = map[string]string{
				"error":"Error: refresh token not valid; to-do: report incident",
			}		
		}
	 

		return c.JSON(status, result)
	})

	e.Logger.Fatal(e.Start(":8080"))
}

