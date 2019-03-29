package main

import (
	"github.com/vietwow/user-account/helper"
	"github.com/vietwow/user-account/models"
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	configuration := helper.GetConfig()

	gormParameters := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=disable", configuration.DbHost, configuration.DbPort, configuration.DbName, configuration.DbUsername, configuration.DbPassword)
	gormDB, err := gorm.Open("postgres", gormParameters)
	if err != nil {
		panic("failed to connect database")
	}
	helper.GormDB = gormDB

	// Migrate the schema (tables): User, Authentication
	helper.GormDB.AutoMigrate(&helper.User{})
	helper.GormDB.AutoMigrate(&helper.Authentication{})
	helper.GormDB.Model(&helper.Authentication{}).AddForeignKey("user_id", "users(id)", "CASCADE", "CASCADE")

	echoFramework := echo.New()
	echoFramework.Use(middleware.Logger()) // log
	echoFramework.Use(middleware.CORS())   // CORS from Any Origin, Any Method

	echoGroupUseJWT := echoFramework.Group("/api/v1")
	echoGroupUseJWT.Use(middleware.JWT([]byte(configuration.EncryptionKey)))
	echoGroupNoJWT := echoFramework.Group("/api/v1")

	// /api/v1/users : logged in users
	echoGroupUseJWT.POST("/users/logout", models.Logout)

	// /api/v1/users : public accessible
	echoGroupNoJWT.POST("/users", models.CreateUser)
	echoGroupNoJWT.POST("/users/login", models.Login)

	defer helper.GormDB.Close()
	echoFramework.Logger.Fatal(echoFramework.Start(":8080"))
}
