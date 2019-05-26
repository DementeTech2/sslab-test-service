package fetch

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

var SSLLabDomain string
var SSLLabSleep uint

func GetDomainAnalysis(domain string) (SSLLabHost, error) {

	var result SSLLabHost
	url := SSLLabDomain + domain

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	req, err := http.NewRequest("GET", url, nil)
	resp, err := client.Do(req)

	if err != nil {
		return result, err
	}

	if resp.StatusCode >= 400 {
		return result, errors.New("No found")
	}

	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return result, err
	}

	lasterr := json.Unmarshal(bodyBytes, &result)
	return result, lasterr
}

type SSLLabHost struct {
	Host            string        `json:"host"`
	Port            int           `json:"port"`
	Protocol        string        `json:"protocol"`
	IsPublic        bool          `json:"isPublic"`
	Status          string        `json:"status"`
	StatusMessage   string        `json:"statusMessage"`
	StartTime       uint          `json:"startTime"`
	TestTime        uint          `json:"testTime"`
	EngineVersion   string        `json:"engineVersion"`
	CriteriaVersion string        `json:"criteriaVersion"`
	CacheExpiryTime int           `json:"cacheExpiryTime"`
	Endpoints       []SSLEndpoint `json:"endpoints"`
}

type SSLEndpoint struct {
	IPAddress            string `json:"ipAddress"`
	ServerName           string `json:"serverName"`
	StatusMessage        string `json:"statusMessage"`
	StatusDetails        string `json:"statusDetails"`
	StatusDetailsMessage string `json:"statusDetailsMessage"`
	Grade                string `json:"grade"`
	GradeTrustIgnored    string `json:"gradeTrustIgnored"`
	FutureGrade          string `json:"futureGrade"`
	HasWarnings          bool   `json:"hasWarnings"`
	IsExceptional        bool   `json:"isExceptional"`
	Progress             uint   `json:"progress"`
	Duration             uint   `json:"duration"`
	Eta                  uint   `json:"eta"`
	Delegation           uint   `json:"delegation"`
	Details              string `json:"details"`
}
