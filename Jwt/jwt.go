package jwt

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

func Login(c echo.Context) error {
	username := c.FormValue("username") //Body - form-data
	password := c.FormValue("password") //Body - form-data

	// username, password 빈 값일시 Unauthrized 메세지 반환
	if username == "" || password == "" {
		return echo.ErrUnauthorized
	}

	// Create token
	token := jwt.New(jwt.SigningMethodHS256)

	// Set claims
	// jwt에 담길 값 설정
	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = username
	claims["password"] = password
	claims["exp"] = time.Now().Add(time.Hour * 1).Unix()

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]string{
		"token": t,
	})
}

func Restricted(c echo.Context) error {
	user := c.Get("user").(*jwt.Token) //Authorization - Bearer token || Headers - Authorization : Bearer token
	claims := user.Claims.(jwt.MapClaims)
	username := claims["username"].(string)
	password := claims["password"].(string)
	return c.JSON(http.StatusOK, map[string]string{
		"Your Username : ": username,
		"Your Password : ": password,
	})
}

/////////////////////////// API 사용 예시 ///////////////////////////////
// func main() {
// 	e := echo.New()

// 	// Middleware
// 	e.Use(middleware.Logger())
// 	e.Use(middleware.Recover())

// 	// Login route
// 	e.POST("/login", login)

// 	// Restricted group
// 	r := e.Group("/restricted")
// 	r.Use(middleware.JWT([]byte("secret")))
// 	r.POST("", restricted)

// 	e.Logger.Fatal(e.Start(":1323"))
// }
