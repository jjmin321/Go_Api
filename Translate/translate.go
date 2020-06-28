package translate

import (
	"encoding/json"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	"github.com/labstack/echo"
)

// Translator struct has a JSON struct of message
type Translator struct {
	Message Message `json:"message"`
}

// Message struct has a JSON struct of result
type Message struct {
	Result Result `json:"result"`
}

// Result struct has values of srcLangType, tarLangType, translatedText
type Result struct {
	SrcLangType    string `json:"srcLangType"`
	TarLangType    string `json:"tarLangType"`
	TranslatedText string `json:"translatedText"`
}

func errCheck(err error) {
	if err != nil {
		panic(err)
	}
}

func TranslatorPage(c echo.Context) error {
	var html = `<html><body>
	<form action="/translate" method="GET">
	from<br><input type="search" name="source"><br>
	to<br><input type="search" name="target"><br>
	text<br><input type="search" name="text">
	<input type="submit" value="번역">
	</form>
</body></html>`
	return c.HTML(http.StatusOK, html)
}

func Translate(c echo.Context) error {
	var translator Translator
	const papago = "https://openapi.naver.com/v1/papago/n2mt"
	err := godotenv.Load("secret.env")
	errCheck(err)
	clientID := os.Getenv("CLIENT_ID")
	clientSecret := os.Getenv("CLIENT_SECRET")
	source, target, text := c.QueryParam("source"), c.QueryParam("target"), c.QueryParam("text")
	data := url.Values{}
	data.Set("source", source)
	data.Set("target", target)
	data.Set("text", text)
	client := &http.Client{}
	req, err := http.NewRequest("POST", papago, strings.NewReader(data.Encode()))
	errCheck(err)
	req.Header.Add("X-Naver-Client-Id", clientID)
	req.Header.Add("X-Naver-Client-Secret", clientSecret)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content_Type", "charset/utf-8")
	req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	resp, err := client.Do(req)
	defer resp.Body.Close()
	dec := json.NewDecoder(resp.Body)
	dec.Decode(&translator)
	return c.JSON(http.StatusOK, map[string]string{
		"result": translator.Message.Result.TranslatedText,
	})
}
