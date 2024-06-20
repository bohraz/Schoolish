// taken from the example @ https://gowebexamples.com/password-hashing/
package auth

import (
	"encoding/gob"
	"fmt"
	"os"
	"root/internal/model"

	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

var SESSION_STORE *sessions.CookieStore = GetSessionSecret()

// Get secret key used by Cookie Store to sign session data and return Cookie Store
func GetSessionSecret() *sessions.CookieStore {
    file, err := os.Open("secrets.config")
    if err != nil {
        fmt.Println("There was an error opining the config file!", err)
    }
    defer file.Close()

    data := make([]byte, 50)
    _, err = file.Read(data)
    if err != nil {
        fmt.Println("There was an error reading the data!", err)
    }

    // Register types for automatic de/serialization in sessions
    gob.Register(model.User{})

    return sessions.NewCookieStore(data)
}

func HashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
    if err != nil {
        return "", err
    }
    return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    if err != nil {
        fmt.Println(err)
    }
    return err == nil
}