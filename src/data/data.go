package data

import (
	"fmt"
	"log"
	"time"

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
	data := []DomainRevision{}

	dbConn.Select("domain").Group("domain").Find(&data)

	vsm := make([]string, len(data))
	for i, v := range data {
		vsm[i] = v.Domain
	}

	return vsm
}

func GetLastRevision(domain string, include_servers bool) (DomainRevision, error) {
	s := DomainRevision{}

	dbConn.Where(&DomainRevision{Domain: domain}).Order("start_time desc").First(&s)

	if include_servers {
		dbConn.Preload("Servers").Where(&DomainRevision{Domain: domain}).Order("start_time desc").First(&s)
	}

	return s, nil
}

func GetRevision(id uint, include_servers bool) (DomainRevision, error) {
	s := DomainRevision{
		ID: id,
	}

	dbConn.Where(&s).First(&s)

	if include_servers {
		dbConn.Model(&s).Related(&s.Servers)
	}

	return s, nil
}

func GetPrevRevision(curr *DomainRevision) (DomainRevision, error) {

	s := DomainRevision{}

	dbConn.Preload("Servers").Where("domain = ? AND start_time < ?", curr.Domain, curr.StartTime).Order("start_time desc").First(&s)

	return s, nil
}

func CreateRevision(domain string) (DomainRevision, error) {
	s := DomainRevision{
		Domain:    domain,
		StartTime: time.Now(),
		Status:    "IN_PROGRESS",
	}

	dbConn.Create(&s)

	return s, nil
}

func UpdateRevision(rev *DomainRevision) {
	dbConn.Save(rev)
}
