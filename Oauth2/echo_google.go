// GO 서버 프레임워크 echo로 구현된 Oauth2 API
package oauth2

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var err = godotenv.Load("googleOauth.env")

var (
	googleOauthConfig *oauth2.Config
	oauthstate        = "dulemona"
)

func init() {
	googleOauthConfig = &oauth2.Config{
		RedirectURL:  "http://localhost:3000/callback",
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_SECRET_KEY"),
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}
}

// func handleHome(w http.ResponseWriter, r *http.Request) {
func HandleHome(c echo.Context) error {
	var html = `<html><body><a href="/login">Google Login</a></body></html>`
	return c.HTML(http.StatusOK, html)
}

func HandleLogin(c echo.Context) error {
	url := googleOauthConfig.AuthCodeURL(oauthstate)
	return c.Redirect(http.StatusTemporaryRedirect, url)
}

func HandleCallback(c echo.Context) error {
	if c.Request().FormValue("state") != oauthstate {
		return c.Redirect(http.StatusTemporaryRedirect, "/")
	}
	token, err := googleOauthConfig.Exchange(oauth2.NoContext, c.Request().FormValue("code"))
	if err != nil {
		fmt.Printf("could not get token: %s\n", err.Error())
		return c.Redirect(http.StatusTemporaryRedirect, "/")
	}
	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		fmt.Printf("could not create get request: %s\n", err.Error())
		return c.Redirect(http.StatusTemporaryRedirect, "/")
	}

	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("could not parse response: %s\n", err.Error())
		return c.Redirect(http.StatusTemporaryRedirect, "/")
	}
	return c.JSONBlob(http.StatusOK, content)
}

////////////////////////////////////// API 사용 예시 ////////////////////////////////////////
// func main() {
// 	e := echo.New()
// 	e.GET("/", handleHome)
// e.GET("/login", handleLogin)
// e.GET("/callback", handleCallback)
// 	e.Logger.Fatal(e.Start(":3000"))
// }
