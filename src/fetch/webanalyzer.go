package fetch

import (
	"crypto/tls"
	"errors"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

var GetTitleRegex *regexp.Regexp
var GetLogoRegex *regexp.Regexp
var GetLogoPathRegex *regexp.Regexp

type Analyzer struct {
	Domain   string
	Title    string
	Logo     string
	IsDown   bool
	body     string
	protocol string
}

func (a *Analyzer) Analyze() {

	a.protocol = "https://"

	body, err := a.getBody(a.protocol + a.Domain)

	if err != nil {
		a.protocol = "http://"
		body, err = a.getBody(a.protocol + a.Domain)

		if err != nil {
			return // Is down so return
		}
	}

	a.IsDown = false
	a.body = body

	a.getTitle()
	a.getLogo()
}

func (a *Analyzer) getBody(url string) (string, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	req, err := http.NewRequest("GET", url, nil)
	resp, err := client.Do(req)

	if err != nil {
		return "", err
	}

	if resp.StatusCode >= 400 {
		return "", errors.New("No found")
	}

	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return "", err
	}
	return string(bodyBytes), nil
}

func (a *Analyzer) getTitle() {
	resRegex1 := GetTitleRegex.FindStringSubmatch(a.body)
	if len(resRegex1) > 0 {
		a.Title = strings.TrimSpace(resRegex1[1])
	}
}

func (a *Analyzer) getLogo() {

	logo := "NOT FOUND"

	for _, item := range GetLogoRegex.FindAllString(a.body, -1) {
		if strings.Contains(item, "og:image") {
			logo = item
			break
		} else if strings.Contains(item, "\"image\"") {
			logo = item
		} else if strings.Contains(item, ".png") {
			logo = item
		} else {
			logo = item
		}
	}

	if logo != "NOT FOUND" {
		a.Logo = a.getLogoPath(logo)
	}
}

func (a *Analyzer) getLogoPath(html string) string {

	path := ""
	resRegex1 := GetLogoPathRegex.FindStringSubmatch(html)
	if len(resRegex1) > 0 {
		path = strings.TrimSpace(resRegex1[1])
	}

	if !strings.HasPrefix(path, "http") {
		if !strings.HasPrefix(path, "/") {
			path = "/" + path
		}

		path = a.protocol + a.Domain + path
	}

	return path
}

func WebAnalyze(domain string) Analyzer {
	newElem := Analyzer{}
	newElem.Domain = domain
	newElem.Title = "NOT FOUND"
	newElem.Logo = "NOT FOUND"
	newElem.IsDown = true
	newElem.Analyze()
	return newElem
}
