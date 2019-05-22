package fetch

import (
	"data"
	"errors"
	"fmt"
	"time"
)

var workers map[uint]*Worker

func InitFetcher() {
	workers = make(map[uint]*Worker)
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
		return rev, ch, errors.New("Does not exists")
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
}

func (w *Worker) Start() {

	///  TODO: Fetcher Login finally

	select {
	case <-time.After(time.Duration(20) * time.Second):
		fmt.Println("timeout message")
		break
	}

	for _, ch := range w.channels {
		fmt.Println(ch)
		select {
		case ch <- "ok":
		default:
		}

	}

	go removeFetch(w.revision.ID)
}
