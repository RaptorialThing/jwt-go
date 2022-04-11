package main
import (
	"net/http"
	"time"

	"github.com/labstack/echo"
	"github.com/golang-jwt/jwt"
)

var hmacSampleSecret []byte

func main() {
	e := echo.New()
	e.GET("/",func(c echo.Context) error {
		return  c.String(http.StatusOK, "Hello,world")
	})

	e.POST("/authorize",func(c echo.Context) error {

		token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
			"GUID": c.FormValue("GUID"),
			"exp": time.Now().Add(time.Minute * 2).Unix(),
		})
		hmacSampleSecret = []byte("secret")
		accessTokenString, err := token.SignedString(hmacSampleSecret)

		refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS512,jwt.MapClaims{
			"GUID": c.FormValue("GUID"),
			"exp": time.Now().Add(time.Hour * 24).Unix(),
		})
		refreshTokenString, err := refreshToken.SignedString(hmacSampleSecret)

		if err != nil {
			return c.JSON(http.StatusOK, map[string]string{
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

		type tokenReqBody struct {
			refreshToken
		}

		return c.String(http.StatusOK, "Refresh tokens")
	})

	e.Logger.Fatal(e.Start(":8080"))
}