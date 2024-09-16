package crypto

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

func Md5Crypto(password string) string {
	data := []byte(password)
	has := md5.Sum(data)
	md5str := fmt.Sprintf("%x", has)
	return md5str
}

func Md5CryptoWithSalt(password string, salt string) string {
	data := []byte(password + salt)
	has := md5.Sum(data)
	md5str := fmt.Sprintf("%x", has)
	return md5str
}

func GenerateSalt(length int) (string, error) {
	salt := make([]byte, length)
	_, err := rand.Read(salt)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(salt), nil
}

func PasswordGen(password string, salt string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(Md5CryptoWithSalt(password, salt)), bcrypt.DefaultCost)
	return string(hash)
}

func PasswordCompare(passwordInput string, correctPassword string, salt string) bool {
	return bcrypt.CompareHashAndPassword([]byte(correctPassword), []byte(Md5CryptoWithSalt(passwordInput, salt))) == nil
}
