package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

type user struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

var albums = []album{
	{ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
	{ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
	{ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
}

var jwtSecret = []byte("test-secret-key")

func main() {
	router := gin.Default()

	router.GET("/albums", getAlbums)
	router.POST("/albums", postAlbums)

	router.Run("localhost:8080")
}

func getAlbums(c *gin.Context) {
	// IndentedJSON can also be replaced by c.JSON to return raw unindented json
	c.IndentedJSON(http.StatusOK, albums)
}

func postAlbums(c *gin.Context) {
	var newAlbum album

	if err := c.BindJSON(&newAlbum); err != nil {
		fmt.Println(err)
		return
	}

	albums = append(albums, newAlbum)

	c.IndentedJSON(http.StatusOK, newAlbum)
}

func generateToken(username string) (string, error) {
	expirationTime := time.Now().Add(5 * time.Minute)

	claims := &claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)

	return tokenString, err

}

func login(c *gin.Context) {
	var client_user user

	if err := c.BindJSON(&client_user); err != nil {
		fmt.Println(err)
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	if client_user.Username != "test_user" || client_user.Password != "test_pass" {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "Invalid user credentials"})
		return
	}

	token, err := generateToken(client_user.Username)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "An error occured"})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"token": token})
}

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {}
}
