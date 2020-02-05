package model

import (
	"errors"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/lacazethomas/goTodo/config"
)

/*Token claims struct*/
type Token struct {
	UserID uuid.UUID
	jwt.StandardClaims
}

//Account a struct to rep user account
type Account struct {
	AccountID uuid.UUID `gorm:"primary_key;type:varchar(36)"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
	Email     string     `json:"email" gorm:"unique"`
	Password  string     `json:"password"`
	BearerToken     string     `json:"token";sql:"-"`
}

//Validate incoming user details...
func (account *Account) Validate(db *gorm.DB) error {

	if !strings.Contains(account.Email, "@") {
		return errors.New("Email address is required")

	}

	if len(account.Password) < 6 {
		return errors.New("Password to short 6 characteres min")
	}

	//Email must be unique
	temp := &Account{}

	//check for errors and duplicate emails
	err := db.Table("accounts").Where("email = ?", account.Email).First(temp).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return errors.New("Connection error. Please retry")
	}
	if temp.Email != "" {
		return errors.New("Email address already in use by another user")
	}

	return nil
}

func (account *Account) Create(db *gorm.DB) (*Account, error) {

	err := account.Validate(db)
	if err != nil {
		return nil, err
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(account.Password), bcrypt.DefaultCost)
	account.Password = string(hashedPassword)

	uuid, err := uuid.NewV4()
	if err != nil {
		return nil, errors.New("Failed to create account, connection error.")
	}


	account.AccountID = uuid
	db.Create(account)

	//Create new JWT token for the newly registered account
	tk := &Token{UserID: account.AccountID}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(config.GetTokenString()))
	account.BearerToken = tokenString

	account.Password = "" //delete password

	return account, nil
}

func Login(email, password string, db *gorm.DB) (*Account, error) {

	account := &Account{}
	err := db.Table("accounts").Where("email = ?", email).First(account).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("Email address not found")
		}
		return nil, errors.New("Connection error. Please retry")
	}

	err = bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(password))
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword { //Password does not match!
		return nil, errors.New("Invalid login credentials. Please try again")
	}
	//Worked! Logged In
	account.Password = ""

	//Create JWT token
	tk := &Token{UserID: account.AccountID}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(config.GetTokenString()))
	account.BearerToken = tokenString //Store the token in the response
	return account, nil
}
