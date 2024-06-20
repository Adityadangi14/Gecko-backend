package controllers

import (
	"gecko_backend/initializers"
	"gecko_backend/models"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

type UserBody struct {
	UserName       string
	Email          string
	PID            string // Id give by auth Provider
	ProfilePicture string
	AuthType       string
}

func AuthHandler(c *gin.Context) {
	var body UserBody
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": "Invalid request body format",
		})
		return
	}
	var user models.UserModel
	result := initializers.DB.Where("p_id = ?", body.PID).First(&user).RowsAffected

	if result == 0 {
		res, err := createNewUser(body)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "unable to Authenticate",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"token": res["token"],
		})
	} else {
		res, err := handleExistingUser(body)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "unable to Authenticate",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"token": res["token"],
		})
	}
}

func createTokenFromId(id string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": id,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))

	if err != nil {
		return "", err
	}

	return tokenString, nil

}
func createNewUser(body UserBody) (map[string]any, error) {
	user := models.UserModel{
		UserName:       body.UserName,
		Email:          body.Email,
		ProfilePicture: body.ProfilePicture,
		AuthType:       body.AuthType,
		AreTagsChoosen: false,
		PID:            body.PID,
	}

	result := initializers.DB.Create(&user)

	if result.Error != nil {
		return nil, result.Error
	}

	tokenString, err := createTokenFromId(user.ID)

	if err != nil {
		return nil, err
	}

	return map[string]any{"token": tokenString, "success": true}, nil
}

func handleExistingUser(body UserBody) (map[string]any, error) {
	var user models.UserModel
	result := initializers.DB.Where("p_id = ?", body.PID).First(&user)

	if result.Error != nil {
		return nil, result.Error
	}

	tokenString, err := createTokenFromId(user.ID)

	if err != nil {
		return nil, err
	}

	return map[string]any{"token": tokenString, "success": true}, nil
}
