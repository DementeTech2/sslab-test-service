package fetch

import (
	"regexp"
	"sync"

	"data"
)

var workers map[uint]*Worker

func InitFetcher() {
	workers = make(map[uint]*Worker)
	GetTitleRegex, _ = regexp.Compile("(?i)<title[^>]*?>\\s?([^<]+)\\s?</title>")
	GetLogoRegex, _ = regexp.Compile("(?i)<(?:meta|link)[^>]*?(?:og:image|itemprop=\"image|icon)['\"][^>]*?>")
	GetLogoPathRegex, _ = regexp.Compile("(?:href|content)=\"([^\"']+?)\"")
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

	// should analyse and send per server
}

func (w *Worker) FetchServerData(serverId uint) {
	defer w.wg.Done()

	// fetch serverId whois
}

func (w *Worker) AnalyseResult() {

	//  retrive previous completed revision
	//		check if servers change
	//		fill the previous data
	//  get the sslgrade from current servers
	//  check if is down

}
