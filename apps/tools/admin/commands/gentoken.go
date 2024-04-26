package commands

import (

	"github.com/golang-jwt/jwt/v4"
)

func GenToken() error {
  jwt.New()
	return nil
}
