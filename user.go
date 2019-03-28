package models

import (
	"app1/helper"
	"fmt"
	"net/http"
	"strconv"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

func CreateUser(c echo.Context) error {
	m := echo.Map{}
	if err := c.Bind(&m); err != nil {
		// return err
	}
	name := m["name"].(string)
	email := m["email"].(string)
	password := m["password"].(string)
	confirmPassword := m["confirm_password"].(string)

	if password == "" || confirmPassword == "" || name == "" || email == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Please enter name, email and password")
	}

	if password != confirmPassword {
		return echo.NewHTTPError(http.StatusBadRequest, "Confirm password is not same to password provided")
	}

	if helper.ValidateEmail(email) == false {
		return echo.NewHTTPError(http.StatusBadRequest, "Please enter valid email")
	}

	if bCheckUserExists(email) == true {
		return echo.NewHTTPError(http.StatusBadRequest, "Email provided already exists")
	}

	configuration := helper.GetConfig()

	enc, _ := helper.EncryptString(password, configuration.EncryptionKey)

	user1 := helper.User{Name: name, Email: email, Password: enc}
	// globals.GormDB.NewRecord(user) // =&gt; returns `true` as primary key is blank
	helper.GormDB.Create(&user1)

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["name"] = user1.Name
	claims["email"] = user1.Email

	t, err := token.SignedString([]byte(configuration.EncryptionKey)) // "secret" &gt;&gt; EncryptionKey
	if err != nil {
		return err
	}

	authentication := helper.Authentication{}
	if helper.GormDB.First(&authentication, "user_id =?", user1.ID).RecordNotFound() {
		// insert
		helper.GormDB.Create(&helper.Authentication{User: user1, Token: t})
	} else {
		authentication.User = user1
		authentication.Token = t
		helper.GormDB.Save(&authentication)
	}

	return c.JSON(http.StatusOK, map[string]string{
		"token": t,
	})
}

func bCheckUserExists(email string) bool {
	user1 := helper.User{}
	if helper.GormDB.Where(&helper.User{Email: email}).First(&user1).RecordNotFound() {
		return false
	}
	return true
}

func ValidateUser(email, password string, c echo.Context) (bool, error) {
	fmt.Println("validate")
	var user1 helper.User
	if helper.GormDB.First(&user1, "email =?", email).RecordNotFound() {
		return false, nil
	}

	configuration := helper.GetConfig()

	decrypted, _ := helper.DecryptString(user1.Password, configuration.EncryptionKey)

	if password == decrypted {
		return true, nil
	}
	return false, nil
}

func Login(c echo.Context) error {
	m := echo.Map{}
	if err := c.Bind(&m); err != nil {
		// return err
	}
	email := m["email"].(string)
	password := m["password"].(string)

	var user1 helper.User
	if helper.GormDB.First(&user1, "email =?", email).RecordNotFound() {
		_error := helper.CustomHTTPError{
			Error: helper.ErrorType{
				Code:    http.StatusBadRequest,
				Message: "Invalid email & password",
			},
		}
		return c.JSONPretty(http.StatusBadGateway, _error, "  ")
	}

	configuration := helper.GetConfig()
	decrypted, _ := helper.DecryptString(user1.Password, configuration.EncryptionKey)

	if password == decrypted {
		token := jwt.New(jwt.SigningMethodHS256)

		claims := token.Claims.(jwt.MapClaims)
		claims["name"] = user1.Name
		claims["email"] = user1.Email
		claims["id"] = user1.ModelBase.ID

		t, err := token.SignedString([]byte(configuration.EncryptionKey)) // "secret" &gt;&gt; EncryptionKey
		if err != nil {
			return err
		}

		authentication := helper.Authentication{}
		if helper.GormDB.First(&authentication, "user_id =?", user1.ID).RecordNotFound() {
			// insert
			helper.GormDB.Create(&helper.Authentication{User: user1, Token: t})
		} else {
			// update
			authentication.User = user1
			authentication.Token = t
			helper.GormDB.Save(&authentication)
		}

		return c.JSON(http.StatusOK, map[string]string{
			"token": t,
		})
	} else {
		_error := helper.CustomHTTPError{
			Error: helper.ErrorType{
				Code:    http.StatusBadRequest,
				Message: "Invalid email & password",
			},
		}
		return c.JSONPretty(http.StatusBadGateway, _error, "  ")
	}
}

func Logout(c echo.Context) error {
	tokenRequester := c.Get("user").(*jwt.Token)
	claims := tokenRequester.Claims.(jwt.MapClaims)
	fRequesterID := claims["id"].(float64)
	iRequesterID := int(fRequesterID)
	sRequesterID := strconv.Itoa(iRequesterID)

	requester := helper.User{}
	if helper.GormDB.First(&requester, "id =?", sRequesterID).RecordNotFound() {
		return echo.ErrUnauthorized
	}

	authentication := helper.Authentication{}
	if helper.GormDB.First(&authentication, "user_id =?", requester.ModelBase.ID).RecordNotFound() {
		return echo.ErrUnauthorized
	}
	helper.GormDB.Delete(&authentication)
	return c.String(http.StatusAccepted, "")
}
