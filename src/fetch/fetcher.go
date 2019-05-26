package fetch

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"sync"
	"time"

	"data"
)

type Config struct {
	Domain    string
	SleepTime uint
}

var workers map[uint]*Worker

func InitFetcher(config Config) {
	workers = make(map[uint]*Worker)
	GetTitleRegex, _ = regexp.Compile("(?i)<title[^>]*?>\\s?([^<]+)\\s?</title>")
	GetLogoRegex, _ = regexp.Compile("(?i)<(?:meta|link)[^>]*?(?:og:image|itemprop=\"image|icon)['\"][^>]*?>")
	GetLogoPathRegex, _ = regexp.Compile("(?:href|content)=\"([^\"']+?)\"")
	GetCountryRegex = regexp.MustCompile("(?i)country:\\s+([A-Z]+)")
	GetOrgNameRegex = regexp.MustCompile("(?i)(?:org-name|orgname):\\s+([A-Z]+)")

	SSLLabDomain = config.Domain
	SSLLabSleep = config.SleepTime
}

func StartFetch(domain string) (data.DomainRevision, chan string, error) {

	ch := make(chan string)
	rev, err := data.CreateRevision(domain)

	if err != nil {
		return rev, ch, err
	}

	chs := []chan<- string{ch}
	w := Worker{revision: rev, channels: chs}

	workers[rev.ID] = &w

	go w.Start()

	return rev, ch, err
}

func TrackFetch(rev data.DomainRevision) (data.DomainRevision, chan string, error) {
	ch := make(chan string)
	w, exists := workers[rev.ID]

	if !exists {
		chs := []chan<- string{ch}
		w := Worker{revision: rev, channels: chs}
		workers[rev.ID] = &w
		go w.Start()
		return rev, ch, nil
	}

	// Should us a Lock Here
	w.channels = append(w.channels, ch)

	return rev, ch, nil
}

func removeFetch(id uint) {
	delete(workers, id)
}

type Worker struct {
	revision data.DomainRevision
	channels []chan<- string
	wg       sync.WaitGroup
	mtx      *sync.Mutex
}

func (w *Worker) Start() {

	w.mtx = &sync.Mutex{}

	w.wg.Add(2)

	go w.FetchPageData()
	go w.FetchSSLLabData()

	w.wg.Wait()

	w.EndResult()

	for _, ch := range w.channels {
		select {
		case ch <- "ok":
		default:
		}
	}

	go removeFetch(w.revision.ID)
}

func (w *Worker) FetchPageData() {
	defer w.wg.Done()

	a := WebAnalyze(w.revision.Domain)

	w.mtx.Lock()
	w.revision.IsDown = a.IsDown
	w.revision.Title = a.Title
	w.revision.Logo = a.Logo
	data.UpdateRevision(&w.revision)
	w.mtx.Unlock()
}

func (w *Worker) FetchSSLLabData() {
	defer w.wg.Done()

	fmt.Println("Getting SSLLab Data ... ")

	ssldata, err := GetDomainAnalysis(w.revision.Domain)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("... " + w.revision.Domain + " " + ssldata.Status)

	newIps := []string{}

	w.mtx.Lock()
	w.revision.Status = strings.ToLower(ssldata.Status)
	for _, endp := range ssldata.Endpoints {
		server := w.revision.GetServerByIp(endp.IPAddress)
		server.Progress = endp.Progress
		if endp.Duration > 0 {
			server.SslGrade = endp.Grade
			server.Progress = 100
		}
		if server.ID == 0 {
			newIps = append(newIps, endp.IPAddress)
		}
	}
	data.UpdateRevision(&w.revision)
	w.mtx.Unlock()

	for _, ip := range newIps {
		w.wg.Add(1)
		w.FetchServerData(ip)
	}

	if !w.revision.IsCompleted() {
		w.wg.Add(1)
		select {
		case <-time.After(time.Duration(SSLLabSleep) * time.Second):
			go w.FetchSSLLabData()
		}
	}
}

func (w *Worker) FetchServerData(ip string) {
	defer w.wg.Done()

	fmt.Println("whois analysis: " + ip)

	who := &WhoIs{Ip: ip}
	who.GetInfo()
	country := who.GetCountry()
	owner := who.GetOwner()

	w.mtx.Lock()
	ser := w.revision.GetServerByIp(ip)
	ser.Country = country
	ser.Owner = owner
	data.UpdateRevision(&w.revision)
	w.mtx.Unlock()

}

func (w *Worker) EndResult() {

	prevRev, err := data.GetPrevRevision(&w.revision)

	if err != nil {
		fmt.Println(err)
	}

	w.mtx.Lock()
	w.revision.SslGrade = w.revision.GetMinGrade()
	w.revision.PreviousSslGrade = prevRev.SslGrade
	w.revision.ServerChanged = !reflect.DeepEqual(w.revision.GetServersMap(), prevRev.GetServersMap())
	w.revision.EndTime = time.Now()
	data.UpdateRevision(&w.revision)
	w.mtx.Unlock()
}
