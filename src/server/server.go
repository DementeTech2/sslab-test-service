package server

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"data"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

// Config is the structure of the params to start the server
type Config struct {
	Host    string
	Port    uint
	Timeout uint
}

// Start Router server
func Start(config Config) {
	r := chi.NewRouter()
	r.Use(middleware.Timeout(time.Duration(config.Timeout) * time.Second))
	r.Use(render.SetContentType(render.ContentTypeJSON))

	r.Post("/api/start/{domain}", startDomainFetch)
	r.Get("/api/domains", getAllDomains)

	log.Println("Server initiated")
	addr := fmt.Sprintf("%s:%d", config.Host, config.Port)
	log.Fatal(http.ListenAndServe(addr, r))
}

func startDomainFetch(w http.ResponseWriter, r *http.Request) {

	domain := chi.URLParam(r, "domain")
	valid := validateDomain(domain)

	fmt.Println(strconv.FormatBool(valid))

	if !valid {
		http.Error(w, "Invalid domain", 400)
		return
	}

	rev, err := data.GetLastRevision(domain, true)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if rev.Domain != "" {
		w.Write([]byte("is not empty"))
		if rev.IsCompleted() {
			if rev.IsOlder(60 * 60) {
				// FETCH StartFetch(domain, channel)  // this call should return the new revision
			} else {
				render.JSON(w, r, rev)
				return
			}
		} else {
			// FETCH trackFetch(revision, channel)
		}
	} else {
		w.Write([]byte("is empty"))
		// FETCH StartFetch(domain, channel)
	}

	// if should wait
	// 		waitOrTimeout(channel)  // maybe a simple select
	//
	// if timeout
	// 		return TimeoutMessage / still in progress
	//
	// DB GetRevision(revisionId)
	// return completed revision

}

func getAllDomains(w http.ResponseWriter, r *http.Request) {
	// Get all the last domains with its servers data, event with the status

	response := []data.DomainRevision{}
	domains := data.GetAllDomains()

	for _, domain := range domains {
		res, err := data.GetLastRevision(domain, true)
		if err == nil {
			response = append(response, res)
		}
	}

	render.JSON(w, r, response)
}

func validateDomain(domain string) bool {
	fmt.Println(domain)
	// get it from here https://www.socketloop.com/tutorials/golang-use-regular-expression-to-validate-domain-name
	RegExp := regexp.MustCompile(`^(([a-zA-Z]{1})|([a-zA-Z]{1}[a-zA-Z]{1})|([a-zA-Z]{1}[0-9]{1})|([0-9]{1}[a-zA-Z]{1})|([a-zA-Z0-9][a-zA-Z0-9-_]{1,61}[a-zA-Z0-9]))\.([a-zA-Z]{2,6}|[a-zA-Z0-9-]{2,30}\.[a-zA-Z]{2,3})$`)
	return RegExp.MatchString(domain)
}
