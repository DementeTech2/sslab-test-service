package data

import (
	"time"
)

type DomainRevision struct {
	ID               uint
	Domain           string
	StartTime        time.Time
	EndTime          time.Time
	Status           string
	Logo             string
	Title            string
	SslGrade         string
	PreviousSslGrade string
	ServerChanged    bool
	IsDown           bool
	Servers          []Server `gorm:"foreignkey:RevisionID"`
}

func (dr *DomainRevision) IsCompleted() bool {
	return dr.Status == "error" || dr.Status == "ready"
}

func (dr *DomainRevision) IsOlder(second uint) bool {
	ref := time.Now().Add(time.Duration(int(second)*-1) * time.Second)
	return dr.EndTime.Before(ref)
}

type Server struct {
	ID         uint
	RevisionID uint
	SslGrade   string
	Progress   uint
	Country    string
	Owner      string
}

type Config struct {
	Database string
	Host     string
	Port     uint
	User     string
	Password *string
}
