package main
import (
	"net/http"
	"time"
	"strconv"
	"math/rand"
	"errors"
	// "fmt"
	"strings"
	bs64 "encoding/base64"
	"encoding/json"

	"github.com/labstack/echo"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
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

func findUser(guid [GUIDFormat]byte ) (User, bool) {

	var res User
	usr := User{}

	if 1<0 {
		res =  usr
	}

	return res,false
}

func createUser(id [GUIDFormat]byte, token string) (User,error) {
	token,err := hashToken(token)
	u := User {
		guid: id,
		refresh_token: token}
	usr, err := saveUser(u)
	if err != nil {
		return User{}, err
	}
	return usr,err
}

func updateUser(id [GUIDFormat]byte, token string) (User, error) {
	u,errFind := findUser(id)
	if errFind  {
		return User{}, errors.New("find User error")
	}
	new_token, errToken := hashToken(token)
	if errToken != nil {
		return User{}, errToken
	}
	u.refresh_token = new_token
	usr,errSave := saveUser(u)
	if errSave != nil {
		return User{}, errSave
	}

	return usr, nil
}

func saveUser(u User) (User, error) {
	// coonect to db and save User or errors.New("cant save user to database")
	usr := User{}
	return usr, nil
}

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

	timeAccess := strconv.FormatInt(time.Now().Add(time.Minute * 2).Unix(),10)

	accessTokenHashEmpty := ""

	accessClaims := CustomClaims{
		stringGuid,
		timeAccess,
		accessTokenHashEmpty, 
		jwt.StandardClaims{
			Issuer:    "test",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, accessClaims)
	accessHmacSampleSecret = []byte(randSecret(RandSecretSize))

	accessTokenString, err := token.SignedString(accessHmacSampleSecret)

	timeRefresh := strconv.FormatInt(time.Now().Add(time.Hour * 24).Unix(),10)

	accessTokenHash, errAccTokenHash := hashToken(accessTokenString)
	if errAccTokenHash != nil {
		err = errAccTokenHash
	}

	refreshClaims := CustomClaims{
		stringGuid,
		timeRefresh,
		accessTokenHash,
		jwt.StandardClaims{
			Issuer:    "test",
		},
	}

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

func decodeTokenClaims(token string) (CustomClaims, error) {
	s := strings.Split(token,".")
	decodedPayload,_ := bs64.StdEncoding.DecodeString(s[1])
	claimSlice := make(map[string]string)

	err := json.Unmarshal(decodedPayload, &claimSlice)

	claim := CustomClaims{}

	for key := range claimSlice {
		if key == "GUID" {
			claim.GUID = claimSlice[key]
		}
		if key == "exp" {
			claim.exp = claimSlice[key]
		}
	}

	return claim,err
}

func convertGuid(guid string) ([GUIDFormat]byte, bool) {
	// formGuid := []byte(guid)
	// err := false 
	// var id [GUIDFormat]byte

	// var idCompare [GUIDFormat]byte

	// for i:= range guid {
	// 	idCompare[i] = 0
	// }

	// if len(formGuid) != len(idCompare) {
	// 	err = true
	// 	id = idCompare
	// } else {
	// 	   for i:= range guid {
	// 		   id[i] = formGuid[0]
	// 	   }
	// }

	guidBytes := []byte(guid)
	// stringId = strconv.Itoa(id)
	var err bool
	var guidBytesFixed [GUIDFormat]byte
	err = false

	if len([GUIDFormat]byte{}) != len(guidBytes) {
		err = true
	}

	for i:= range [GUIDFormat]byte{} {
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
		formGuid,errConvertGuid := convertGuid(c.FormValue("GUID"))

		if errConvertGuid != false {
			return c.JSON(http.StatusUnprocessableEntity, map[string]string{
				"error": c.FormValue("GUID")+" GUID cant format",
			})
		}

		if authFoundUser, errFindUser := findUser(formGuid); errFindUser == true  {
			if authFoundUser.guid != [GUIDFormat]byte{} {

			}
			return c.String(http.StatusOK, "Error: User already authenticated "+c.FormValue("GUID"))
		}

	
		accessTokenString, refreshTokenString, err := generateUserTokens(formGuid)	
		//_, errCreateUser := createUser(formGuid, refreshTokenString)

		if err != nil  {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"error": err.Error(),
			})
		// } else if errCreateUser != nil {
		// 	return c.JSON(http.StatusUnauthorized, map[string]string{
		// 		"error": errCreateUser.Error(),
		// 	})
	    } else {

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
		
		// var Guid [GUIDFormat]byte
		// var ErrParse bool

		// get refreshHmacSampleSecretHash for GUID
		decodedToken,errParse := jwt.ParseWithClaims(refresh_token, &CustomClaims{},
			func(token *jwt.Token) (interface{}, error) {

				return refreshHmacSampleSecret, nil

			})
			// ,
			//  func(token *jwt.Token) (interface{}, error) {
			// var res string
			// var err string

			// claims, ok := token.Claims.(jwt.MapClaims); 
			// if ok && token.Valid {
			// 	var assertedGuid string
			// 	assertedGuid = claims["GUID"]
			// 	guid, err := convertGuid(assertedGuid)
			// 	if err {
			// 		res = ""
			// 		err = "Error GUID cant format "+assertedGuid
			// 	}

			// 	Guid = guid
			// 	res = "OK"
			// 	err = ""
			// 	ErrParse = false
				
			// } else {
			// 	res = ""
			// 	err = "Error parsing token error"
			// 	guid = [GUIDFormat]byte{}
			// 	ErrParse = true
				
			// }


			// if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			// 	res = ""
			// 	err = "Unexpected signing method"
			// }

				//  return rsa.Public, nil
			//   })

		var result map[string]string
		status := http.StatusOK

		if errParse != nil {
			status = http.StatusUnprocessableEntity
			result = map[string]string{
				"error": "can not decode token",
			}
		} else {
			
			claims, ok := decodedToken.Claims.(*CustomClaims); 
			if ok && decodedToken.Valid {
		
				
			}	else {
				status = http.StatusInternalServerError
				result = map[string]string{
					"error": "can not decode token",
				}
				return c.JSON(status, result)
			}

			guid, errAsserGuid := convertGuid(claims.GUID)
			if errAsserGuid  {

			}
		

        //claims,errDecodeClaims := decodeTokenClaims(refresh_token)
		//guidByte,errConvByte := convertGuid(claims.GUID)

		// if errDecodeClaims != nil || errConvByte {	
		// 	status = http.StatusUnprocessableEntity
		// 	result = map[string]string{
		// 		"error": "GUID not found or decode error",
		// 	}
		// 	return c.JSON(status, result)

		// } 		
		
		//guid := guidByte	
		//fmt.Println(refresh_token)
	
		match := checkTokenHash(refresh_token, []byte(Refresh_base_hash))

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
				//updateUser(guid, refreshTokenString)
			}

		} else {
			status = http.StatusUnauthorized
			result = map[string]string{
				"error":"Error: refresh token not valid; to-do: report incident",
			}		
		}
	 }

		return c.JSON(status, result)
	})

	e.Logger.Fatal(e.Start(":8080"))
}
