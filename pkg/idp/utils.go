package idp

import (
    "golang.org/x/crypto/bcrypt"
    log "github.com/sirupsen/logrus"
)

// function used to hash and salt user passwords
func hashAndSalt(password string) string {
    // convert passwords into byte array and hash
    hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        log.Error(err)
    }
    return string(hash)
}

// function used to compare password to hashed password from database
func comparePasswords(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    if err != nil {
        log.Warn(err)
        return false
    }
    return true
}

