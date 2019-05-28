package server

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"time"

	"data"
	"fetch"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
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

	// Basic CORS
	// for more ideas, see: https://developer.github.com/v3/#cross-origin-resource-sharing
	cors := cors.New(cors.Options{
		// AllowedOrigins: []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})
	r.Use(cors.Handler)

	r.Get("/api/analyze/{domain}", startDomainFetch)
	r.Get("/api/test_domain/{domain}", testPage)
	r.Get("/api/domains", getAllDomains)

	log.Println("Server initiated")
	addr := fmt.Sprintf("%s:%d", config.Host, config.Port)
	log.Fatal(http.ListenAndServe(addr, r))
}

func startDomainFetch(w http.ResponseWriter, r *http.Request) {

	_, dontWait := r.URL.Query()["dont_wait"]

	domain := chi.URLParam(r, "domain")
	valid := validateDomain(domain)

	if !valid {
		http.Error(w, "Invalid domain", http.StatusBadRequest)
		return
	}

	rev, err := data.GetLastRevision(domain, true)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var ch chan string

	if rev.ID != 0 {
		if rev.IsCompleted() {
			if rev.IsOlder(60 * 60) {
				fmt.Println("IsOlder")
				// check if is fetching (maybe stop that fetch)
				rev, ch, err = fetch.StartFetch(domain)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			} else {
				fmt.Println("Is not Older")
				render.JSON(w, r, rev)
				return
			}
		} else {
			fmt.Println("Track old revision")
			rev, ch, err = fetch.TrackFetch(rev)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	} else {
		fmt.Println("Start new revision")
		rev, ch, err = fetch.StartFetch(domain)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	if dontWait {
		render.JSON(w, r, rev) // return the actual state of the revision
		return
	}

	select {
	case <-ch:
		break
	case <-time.After(time.Duration(30) * time.Second):
		http.Error(w, "Still running, call it later", http.StatusRequestTimeout)
		return
	}

	rev, err = data.GetRevision(rev.ID, false)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	render.JSON(w, r, rev)
}

func getAllDomains(w http.ResponseWriter, r *http.Request) {
	// Get all the last domains with its servers data, event with the status

	response := []data.DomainRevision{}
	domains := data.GetAllDomains()

	_, servers := r.URL.Query()["include_servers"]

	for _, domain := range domains {
		res, err := data.GetLastRevision(domain, servers)
		if err == nil {
			response = append(response, res)
		}
	}

	render.JSON(w, r, response)
}

func validateDomain(domain string) bool {
	// get it from here https://www.socketloop.com/tutorials/golang-use-regular-expression-to-validate-domain-name
	RegExp := regexp.MustCompile(`^(([a-zA-Z]{1})|([a-zA-Z]{1}[a-zA-Z]{1})|([a-zA-Z]{1}[0-9]{1})|([0-9]{1}[a-zA-Z]{1})|([a-zA-Z0-9][a-zA-Z0-9-_]{1,61}[a-zA-Z0-9]))\.([a-zA-Z]{2,6}|[a-zA-Z0-9-]{2,30}\.[a-zA-Z]{2,3})$`)
	return RegExp.MatchString(domain)
}

func testPage(w http.ResponseWriter, r *http.Request) {
	domain := chi.URLParam(r, "domain")
	valid := validateDomain(domain)

	if !valid {
		http.Error(w, "Invalid domain", http.StatusBadRequest)
		return
	}

	a := fetch.WebAnalyze(domain)
	render.JSON(w, r, a)
}
