package model

import (
	"errors"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"

	"github.com/lacazethomas/goTodo/config"
)

/*
JWT claims struct
*/
type Token struct {
	UserId uint
	jwt.StandardClaims
}

//a struct to rep user account
type Account struct {
	gorm.Model
	Email    string `json:"email" gorm:"unique"`
	Password string `json:"password"`
	Token    string `json:"token";sql:"-"`
}

//Validate incoming user details...
func (account *Account) Validate(db *gorm.DB) error {

	if !strings.Contains(account.Email, "@") {
		return errors.New("Email address is required")

	}

	if len(account.Password) < 6 {
		return errors.New("Email address is required")
	}

	//Email must be unique
	temp := &Account{}

	//check for errors and duplicate emails
	err := db.Table("accounts").Where("email = ?", account.Email).First(temp).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return errors.New("Connection error. Please retry")
	}
	if temp.Email != "" {
		return errors.New("Email address already in use by another user.")
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

	db.Create(account)

	if account.ID <= 0 {
		return nil, errors.New("Failed to create account, connection error.")
	}

	//Create new JWT token for the newly registered account
	tk := &Token{UserId: account.ID}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(config.GetTokenString()))
	account.Token = tokenString

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
	tk := &Token{UserId: account.ID}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(config.GetTokenString()))
	account.Token = tokenString //Store the token in the response
	return account, nil
}

func GetUser(u uint, db *gorm.DB) *Account {

	acc := &Account{}
	db.Table("accounts").Where("id = ?", u).First(acc)
	if acc.Email == "" { //User not found!
		return nil
	}

	acc.Password = ""
	return acc
}
