package main

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	ClientID     = "YOUR_CLIENT_ID"
	ClientSecret = "YOUR_CLIENT_SECRET"
	RedirectURI  = "http://localhost:8080/callback"
	AuthURL      = "https://accounts.spotify.com/authorize"
	TokenURL     = "https://accounts.spotify.com/api/token"
)

var (
	State         string
	Verifier      string
	Authorization string
)

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
}

func main() {
	router := gin.Default()

	// Home page
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	// Authorization Code flow with PKCE
	router.GET("/authorize", func(c *gin.Context) {
		State = randomString(32)
		Verifier = randomString(64)

		authParams := url.Values{
			"response_type":         {"code"},
			"client_id":             {ClientID},
			"redirect_uri":          {RedirectURI},
			"state":                 {State},
			"code_challenge_method": {"S256"},
			"code_challenge":        {base64URLEncode(sha256.Sum256([]byte(Verifier))[:])},
			"scope":                 {"user-read-private user-read-email"}, // Set desired scopes
		}

		authURL := AuthURL + "?" + authParams.Encode()
		c.Redirect(http.StatusSeeOther, authURL)
	})

	// Callback URL
	router.GET("/callback", func(c *gin.Context) {
		queryState := c.Query("state")
		code := c.Query("code")

		if queryState != State {
			c.HTML(http.StatusBadRequest, "error.html", gin.H{"message": "Invalid state parameter"})
			return
		}

		tokenParams := url.Values{
			"grant_type":    {"authorization_code"},
			"code":          {code},
			"redirect_uri":  {RedirectURI},
			"code_verifier": {Verifier},
		}

		tokenHeaders := map[string]string{
			"Authorization": "Basic " + base64.StdEncoding.EncodeToString([]byte(ClientID+":"+ClientSecret)),
			"Content-Type":  "application/x-www-form-urlencoded",
		}

		tokenResponse := TokenResponse{}
		err := sendPostRequest(TokenURL, tokenParams.Encode(), tokenHeaders, &tokenResponse)
		if err != nil {
			c.HTML(http.StatusInternalServerError, "error.html", gin.H{"message": "Failed to retrieve access token"})
			return
		}

		Authorization = "Bearer " + tokenResponse.AccessToken
		c.HTML(http.StatusOK, "success.html", gin.H{"accessToken": tokenResponse.AccessToken})
	})

	router.LoadHTMLGlob("templates/*.html")

	router.Run(":8080")
}

func randomString(length int) string {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return base64.URLEncoding.EncodeToString(b)
}

func base64URLEncode(data []byte) string {
	return strings.TrimRight(base64.URLEncoding.EncodeToString(data), "=")
}

func sendPostRequest(url string, params string, headers map[string]string, response interface{}) error {
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(params))
	if err != nil {
		return err
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("request failed with status code %d", resp.StatusCode)
	}

	err = json.NewDecoder(resp.Body).Decode(response)
	if err != nil {
		return err
	}

	return nil
}
