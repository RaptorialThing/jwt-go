package main
import (
	"net/http"
	"time"
	// "strconv"
	"fmt"
	"strings"

	"github.com/labstack/echo"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	
)

var HmacSampleSecret []byte
var Refresh_base_hash string
var AccessTokenSign string
var RefreshTokenSign string


type User struct {
	guid [36]byte `form:"GUID" query:"GUID" json:"GUID" bson: "GUID"`
	refresh_token string `form:"refresh_token" 
	query:"refresh_token" json: "refresh_token" bson: "refresh_token"`
}

type CustomClaims struct {
	GUID string `json:"GUID"`
	exp int64 `json:"exp"`
	jwt.StandardClaims
}

func findUser(guid [36]byte ) bool {

	var res bool;

	if 1<0 {
		res =  true
	}

	return res
}

func createUser(id [36]byte, token string) (User,error) {
	token,err := hashToken(token)
	u := User {
		guid: id,
		refresh_token: token}
	Refresh_base_hash = token
	return u,err
}


func hashToken(token string) (string,error) {
	bytes,err := bcrypt.GenerateFromPassword([]byte(token),14)
	return string(bytes), err
}

func checkTokenHash(token string, hash []byte) bool {
	err := bcrypt.CompareHashAndPassword(hash, []byte(token))
	return err == nil
}

func generateUserTokens(guid [36]byte) (string, string, error) {
	strGuid := make([]string, len(guid))
	for i:= range guid {
		strGuid[i] = string(int(guid[i]))
	}
	stringGuid := strings.Join(strGuid,"")

	timeAccess := time.Now().Add(time.Minute * 2).Unix()

	accessClaims := CustomClaims{
		stringGuid,
		timeAccess,
		jwt.StandardClaims{
			Issuer:    "test",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, accessClaims)
	accessTokenString, err := token.SignedString(HmacSampleSecret)

	timeRefresh := time.Now().Add(time.Hour * 24).Unix()

	refreshClaims := CustomClaims{
		stringGuid,
		timeRefresh,
		jwt.StandardClaims{
			Issuer:    "test",
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS512,refreshClaims)
	refreshTokenString, err := refreshToken.SignedString(HmacSampleSecret)
	Refresh_base,errHashRefreshToken := hashToken(refreshTokenString)
	Refresh_base_hash = Refresh_base
	if err == nil {
		err = errHashRefreshToken
	}

	return accessTokenString, refreshTokenString, err
}

func convertGuid(guid string) ([36]byte, bool) {
	// formGuid := []byte(guid)
	// err := false 
	// var id [36]byte

	// var idCompare [36]byte

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
	var guidBytesFixed [36]byte
	err = false

	if len([36]byte{}) != len(guidBytes) {
		err = true
	}

	for i:= range [36]byte{} {
		guidBytesFixed[i] = guidBytes[i]
	}

	return guidBytesFixed, err
}



func main() {

	e := echo.New()
	HmacSampleSecret = []byte("secret")

	e.GET("/",func(c echo.Context) error {
		return  c.String(http.StatusOK, "Hello,world")
	})

	e.POST("/authenticate",func(c echo.Context) error {
		formGuid,errConvertGuid := convertGuid(c.FormValue("GUID"))

		if errConvertGuid != false {
			return c.JSON(http.StatusRequestedRangeNotSatisfiable, map[string]string{
				"error": c.FormValue("GUID")+" GUID cant format",
			})
		}
		usr := User{guid: formGuid}

		if findUser(usr.guid) == true {
			return c.String(http.StatusOK, "Error: User already authenticated")
		}
	
		accessTokenString, refreshTokenString, err := generateUserTokens(formGuid)
		
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"error": err.Error(),
			})
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
			return c.JSON(http.StatusRequestedRangeNotSatisfiable,
				map[string]string{"error":"give me refresh-token"})
		}
		
		// var Guid [36]byte
		// var ErrParse bool

		decodedToken,errParse := jwt.ParseWithClaims(refresh_token, &CustomClaims{},
			func(token *jwt.Token) (interface{}, error) {

				return HmacSampleSecret, nil

			})
		fmt.Println(decodedToken)
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
			// 	guid = [36]byte{}
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
			status = http.StatusRequestedRangeNotSatisfiable
			result = map[string]string{
				"error": "cant decode token",
			}
		} else {
		
		claims, ok  := decodedToken.Claims.(*CustomClaims) 

		if ok && decodedToken.Valid {
			
		} else {
			
				status = http.StatusRequestedRangeNotSatisfiable
				result = map[string]string{
					"error": "GUID not found",
				}
		
		}
		

		guid, errConv := convertGuid(claims.GUID)
		if errConv {
			status = http.StatusInternalServerError
			result = map[string]string{
				"error":"Error GUID cant format "+claims.GUID,
			}
		}
	

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
				result = map[string]string{
					"error": "Error: tokens dont generated",
				}
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
