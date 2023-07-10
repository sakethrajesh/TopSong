package main

import (
	// "fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zmb3/spotify"
)

func main() {
	router := gin.Default()

	// Set up the Spotify API client
	auth := spotify.NewAuthenticator(
		"http://localhost:8080/callback",
		spotify.ScopeUserReadPrivate,
		spotify.ScopeUserReadEmail,
	)
	auth.SetAuthInfo("5ebbf3885b10400da8684c1969041cce", "7380e4f86aa04a65b003acffd26ce941")

	// Set up the callback handler
	router.GET("/callback", func(c *gin.Context) {
		token, err := auth.Token("state", c.Request)
		if err != nil {
			http.Error(c.Writer, "Couldn't get token", http.StatusForbidden)
			return
		}
		c.SetCookie("token", string(*token), 3600, "/", "127.0.0.1", false, false)
		c.Redirect(http.StatusTemporaryRedirect, "http://localhost:3000")

	})

	// Set up the login handler
	router.GET("/login", func(c *gin.Context) {
		url := auth.AuthURL("state")
		c.Redirect(http.StatusTemporaryRedirect, url)
	})

	// Start the server
	router.Run(":8080")
}
