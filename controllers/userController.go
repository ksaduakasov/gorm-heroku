package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"kettkal/inits"
	"kettkal/models"
	"kettkal/pass"
	"net/http"
	"os"
	"time"
)

func SignUp(c *gin.Context) {
	var body struct {
		Email    string
		Password string
	}
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body!",
		})
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to hash a password!",
		})
		return
	}

	user := models.User{
		Email:    body.Email,
		Password: string(hash),
	}
	result := inits.DB.Create(&user)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create a user!",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

func Login(c *gin.Context) {
	var body struct {
		Email    string
		Password string
	}
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body!",
		})
		return
	}
	var user models.User
	inits.DB.First(&user, "email = ?", body.Email)

	if user.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid email or password!",
		})
		return
	}
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid email or password!",
		})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create a token!",
		})
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Auth", tokenString, 3600, "", "", false, true)
	c.Set("token", tokenString)
	c.JSON(http.StatusOK, gin.H{})
}

func Validate(c *gin.Context) {
	user, _ := c.Get("user")
	c.JSON(http.StatusOK, gin.H{
		"message": user,
	})
}

func GeneratePass(c *gin.Context) {
	var body struct {
		TeamIdentifier     string
		PassTypeIdentifier string
		OrganizationName   string
		SerialNumber       string
		KeyForField        string
		LabelForField      string
		ValueForField      string
		Icon               string
		Logo               string
		Strip              string
	}
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body!",
		})
		return
	}
	teamIdentifier := body.TeamIdentifier
	passTypeIdentifier := body.PassTypeIdentifier
	organizationName := body.OrganizationName
	serialNumber := body.SerialNumber
	keyForField := body.KeyForField
	labelForField := body.LabelForField
	valueForField := body.ValueForField
	icon := body.Icon
	logo := body.Logo
	strip := body.Strip
	_, err := pass.GeneratePass(c, teamIdentifier, passTypeIdentifier, organizationName, serialNumber, keyForField, labelForField, valueForField, icon, logo, strip)
	if err != nil {
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "PkPass was successfully generated!",
	})
}
