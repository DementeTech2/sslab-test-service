package data

import (
	"sort"
	"strings"
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

func (dr *DomainRevision) IsOlder(seconds uint) bool {
	ref := time.Now().Add(time.Duration(int(seconds)*-1) * time.Second)
	return dr.EndTime.Before(ref)
}

func (dr *DomainRevision) GetServerByIp(ip string) *Server {

	for i, ser := range dr.Servers {
		if ser.IP == ip {
			return &dr.Servers[i]
		}
	}

	newSer := Server{IP: ip}
	dr.Servers = append(dr.Servers, newSer)
	return dr.GetServerByIp(ip)
}

func (dr *DomainRevision) GetMinGrade() string {
	justGrades := []string{}

	for _, ser := range dr.Servers {
		thisGrade := ser.SslGrade
		if len(ser.SslGrade) == 1 {
			thisGrade = ser.SslGrade + "."
		}
		justGrades = append(justGrades, thisGrade)
	}

	if len(justGrades) == 0 {
		return ""
	}

	sort.Strings(justGrades)
	lastGrade := justGrades[len(justGrades)-1]
	lastGrade = strings.TrimSuffix(lastGrade, ".")
	return lastGrade
}

func (dr *DomainRevision) GetServersMap() map[string]string {
	theMap := make(map[string]string)

	for _, ser := range dr.Servers {
		theMap[ser.IP] = ser.SslGrade
	}

	return theMap
}

type Server struct {
	ID         uint
	RevisionID uint
	IP         string
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
