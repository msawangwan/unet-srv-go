package db

import (
	"crypto/rand"
	"crypto/sha1"
	"fmt"
	"log"
)

/* based on 'go wp', see: RFC 4122 */
func CreateUUID() string {
	u := new([16]byte)

	_, err := rand.Read(u[:])

	if err != nil {
		log.Fatalln("failed to generate UUID:", err)
	}

	u[8] = (u[8] | 0x40) & 0x7F
	u[6] = (u[6] & 0xF) | (0x4 << 4)

	return fmt.Sprintf("%x%x%x%x%x", u[0:4], u[4:6], u[6:8], u[8:10], u[10:]) // need to include hyphens
}

/* hash plaintext using SHA-1 */
func Encrypt(plaintext string) string {
	return fmt.Sprintf("%x", sha1.Sum([]byte(plaintext)))
}
