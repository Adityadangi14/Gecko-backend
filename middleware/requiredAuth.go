package middleware

import (
	"fmt"
	"gecko_backend/initializers"
	"gecko_backend/models"
	"log"
	"net/http"
	"os"

	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func RequiredAuth(c *gin.Context) {

	defer func() {
		if err := recover(); err != nil {
			// Handle the panic gracefully (e.g., log the error, return an error response)
			c.AbortWithError(http.StatusBadRequest, err.(error))
			fmt.Println("Error in middleware:", err) // Log for debugging
		}
	}()

	tokenString := c.GetHeader("Token")

	if tokenString == "" {
		c.Abort()
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Auth Failed",
		})
		return
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			c.Abort()
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Auth Failed ()",
			})
			return "", nil
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(os.Getenv("SECRET")), nil
	})
	if err != nil {
		log.Fatal(err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {

		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.AbortWithStatus(http.StatusUnauthorized)
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Auth Failed ()",
			})
		}

		var user models.UserModel

		initializers.DB.Where("id = ?", claims["sub"]).First(&user)

		if user.Email == "" {
			c.AbortWithStatus(http.StatusUnauthorized)
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Auth Failed (user not found)",
			})
		}

		c.Set("user", user)

		fmt.Println(c.Request.URL.Path)
		
		

		c.Next()
	} else {
		c.AbortWithStatus(http.StatusUnauthorized)
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Auth Failed()",
		})
	}
}