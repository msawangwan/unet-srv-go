package data

import (
	"crypto/rand"
	"crypto/sha1"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
)

const (
	kDB_DRIVER   = "postgres"
	kDB_USER     = "postgres"
	kDB_PASSWORD = "1234"
	kDB_DATABASE = "unitywebservice"
)

var PostGreService *sql.DB

func init() {
	var err error

	connstr := fmt.Sprintf(
		"user=%s password=%s dbname=%s sslmode=disable",
		kDB_USER,
		kDB_PASSWORD,
		kDB_DATABASE,
	)

	PostGreService, err = sql.Open(kDB_DRIVER, connstr)

	if err != nil {
		log.Fatal(err)
	}

	if err = PostGreService.Ping(); err != nil {
		log.Fatal(err)
	}

	log.Printf("[db][db.go][init db: success]\n")
}

/* based on 'go wp', see: RFC 4122 */
func CreateUUID() string {
	u := new([16]byte)

	_, err := rand.Read(u[:])

	if err != nil {
		log.Fatalln("failed to generate UUID:", err)
	}

	u[8] = (u[8] | 0x40) & 0x7F
	u[6] = (u[6] & 0xF) | (0x4 << 4)

	return fmt.Sprintf("%x-%x-%x-%x-%x", u[0:4], u[4:6], u[6:8], u[8:10], u[10:])
}

/* hash plaintext using SHA-1 */
func Encrypt(plaintext string) string {
	return fmt.Sprintf("%x", sha1.Sum([]byte(plaintext)))
}
