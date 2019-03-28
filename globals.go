package helper

import (
	"regexp"
	"time"

	"github.com/jinzhu/gorm"
)

type Configuration struct {
	EncryptionKey string
	DbHost        string
	DbPort        string
	DbName        string
	DbUsername    string
	DbPassword    string
}

var configuration Configuration

type ModelBase struct {
	ID        int       `gorm:"primary_key"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

type User struct {
	ModelBase        // replaces gorm.Model
	Email     string `gorm:"not null; unique"`
	Password  string `gorm:"not null" json:"-"`
	Name      string `gorm:"not null; type:varchar(100)"` // unique_index
}

type Authentication struct {
	ModelBase
	User   User `gorm:"foreignkey:UserID; not null"`
	UserID int
	Token  string `gorm:"type:varchar(200); not null"`
}


type CustomHTTPSuccess struct {
	Data string `json:"data"`
}

type ErrorType struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type CustomHTTPError struct {
	Error ErrorType `json:"error"`
}

var GormDB *gorm.DB

func init() {
	configuration = Configuration{
		EncryptionKey: "F61L8L7CUCGN0NK6336I8TFP9Y2ZOS43",
		DbHost:        "localhost",
		DbPort:        "5432",
		DbName:        "postgres",
		DbUsername:    "postgres",
		DbPassword:    "postgres",
	}
}

func GetConfig() Configuration {
	return configuration
}

func ValidateEmail(email string) (matchedString bool) {
	re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	matchedString = re.MatchString(email)
	return
}
