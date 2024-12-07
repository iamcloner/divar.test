package password

import (
	"crypto/rand"
	"fmt"
	"os"
	"strconv"

	"golang.org/x/crypto/bcrypt"
)

func Hash(password string) (string, error) {
	cost, _ := strconv.Atoi(os.Getenv("PASSWORD_ENCRYPTION"))
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	return string(bytes), err
}

func Check(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GenerateRedisId() (string, error) {
	b := make([]byte, 4)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", b), nil
}
