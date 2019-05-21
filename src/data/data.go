package data

import (
	"fmt"
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var dbConn *gorm.DB

// InitDB Initialize the database connection or failed
func InitDB(config Config) {
	var err error
	addr := fmt.Sprintf("postgresql://%s@%s:%d/%s?sslmode=disable", config.User, config.Host, config.Port, config.Database)

	log.Println(addr)

	dbConn, err = gorm.Open("postgres", addr)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Database connection initiated")
	dbConn.AutoMigrate(&DomainRevision{}, &Server{})
	log.Println("Scheme migration completed")
}

// CloseDB Close the database connection
func CloseDB() {
	log.Println("Database connection stopped")
	dbConn.Close()
}

func GetAllDomains() []string {
	data := []string{
		"domain2",
		"domain1",
	}
	return data
}

func GetLastRevision(domain string, include_servers bool) (DomainRevision, error) {
	s := DomainRevision{}
	//return s, errors.New("Something bad happend")
	return s, nil
}

func GetRevision(id uint, include_servers bool) (DomainRevision, error) {
	s := DomainRevision{}
	return s, nil
}
