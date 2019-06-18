package server

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"time"

	"ssllabtestservice/data"
	"ssllabtestservice/fetch"

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
	r.Get("/api/fetch-web-data/{domain}", testPage)
	r.Get("/api/domains", getAllDomains)

	log.Println("Server initiated")
	addr := fmt.Sprintf("%s:%d", config.Host, config.Port)
	log.Fatal(http.ListenAndServe(addr, r))
}

 /**
 * @api {get} /api/analyze/:domain?[dont_wait] Start to analyze a domain
 * @apiName AnalyzeDomain
 * @apiGroup Domains
 *
 * @apiVersion 1.0.0
 * @apiDescription This is the Description.
 *
 * @apiParam domain Include the servers of every domain in the response
 * @apiParam [dont_wait] Start the analysis async. By default is going to timeout at 60 seconds. 
 *
 * @apiSuccessExample Initial regular response (dont_wait):
 *     HTTP/1.1 200 OK
 *     {
 *         "id"                 : 123456789
 *         "domain"             : "google.com"
 *         "start_time"         : "2019-06-17T22:00:21-0000"
 *         "end_time"           : "1970-01-01T00:00:00-0000"
 *         "status"             : "pending"
 *         "logo"               : ""
 *         "title"              : ""
 *         "ssl_grade"          : ""
 *         "previous_ssl_grade" : ""
 *         "server_changed"     : false
 *         "is_down"            : false
 *         "servers"            : []
 *     }
 *
 * @apiSuccessExample Regular response without timeout:
 *     HTTP/1.1 200 OK
 *     {
 *         "id"                 : 123456789
 *         "domain"             : "google.com"
 *         "start_time"         : "2019-06-17T22:00:21-0000"
 *         "end_time"           : "2019-06-17T22:15:21-0000"
 *         "status"             : "ready"
 *         "logo"               : "https://google.com/logog.png"
 *         "title"              : "Google com"
 *         "ssl_grade"          : "a"
 *         "previous_ssl_grade" : ""
 *         "server_changed"     : false
 *         "is_down"            : false
 *         "servers"            : []
 *     }
 *
 * @apiSuccessExample Timeout:
 *     HTTP/1.1 408 Request Timeout
 *     Still running, call it later
 *
 */
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

 /**
 * @api {get} /api/domains?[include_servers] Get all domains analyzed
 * @apiName GetDomains
 * @apiGroup Domains
 *
 * @apiVersion 1.0.0
 * @apiDescription This is the Description.
 *
 * @apiParam [include_servers] Include the servers of every domain in the response
 *
 * @apiSuccessExample :
 *     HTTP/1.1 200 OK
 *     [
 *         {
 *             "id"                 : 123456789,
 *             "domain"             : "google.com",
 *             "start_time"         : "2019-06-17T22:00:21-0000",
 *             "end_time"           : "2019-06-17T22:15:21-0000",
 *             "status"             : "ready",
 *             "logo"               : "https://google.com/logog.png",
 *             "title"              : "Google com",
 *             "ssl_grade"          : "a",
 *             "previous_ssl_grade" : "",
 *             "server_changed"     : false,
 *             "is_down"            : false,
 *             "servers"            : [
 *                  {
 *                      "id"          : 9876,
 *                      "revision_id" : 123456789,
 *                      "ip"          : "127.0.0.1",
 *                      "ssl_grade"   : "a",
 *                      "progress"    : 100,
 *                      "country"     : "us",
 *                      "owner"       : "GoogleInc"
 *                  },
 *                  {
 *                      "id"          : 9877,
 *                      "revision_id" : 123456789,
 *                      "ip"          : "127.0.0.1",
 *                      "ssl_grade"   : "a",
 *                      "progress"    : 100,
 *                      "country"     : "us",
 *                      "owner"       : "GoogleInc"
 *                  },
 *                  ....
 *             ]
 *         },
 *         {
 *             "id"                 : 123456790,
 *             "domain"             : "fake.com",
 *             "start_time"         : "2019-06-17T22:00:21-0000",
 *             "end_time"           : "2019-06-17T22:15:21-0000",
 *             "status"             : "error",
 *             "logo"               : "NOT_FOUND",
 *             "title"              : "NOT_FOUND",
 *             "ssl_grade"          : "",
 *             "previous_ssl_grade" : "",
 *             "server_changed"     : false,
 *             "is_down"            : true,
 *             "servers"            : [],
 *         },
 *         ....
 *     ]
 *
 */
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

 /**
 * @api {get} /api/fetch-web-data/:domain Fetch Domain Web Data
 * @apiName FetchDomainData
 * @apiGroup Domains
 *
 * @apiVersion 1.0.0
 * @apiDescription This is the Description.
 *
 * @apiParam domain The domain to analyze
 *
 * @apiSuccessExample Good domain:
 *     HTTP/1.1 200 OK
 *     {
 *       "Domain": "fake.foo",
 *       "Title": "The fake page",
 *       "Logo": "https://fake.foo/my-logo.png",
 *       "IsDown": false
 *     }
 *
 *
 * @apiSuccessExample Bad domain:
 *     HTTP/1.1 200 OK
 *     {
 *       "Domain": "fake.down.foo",
 *       "Title": "NOT_FOUND",
 *       "Logo": "NOT_FOUND",
 *       "IsDown": true
 *     }
 *
 * @apiError BadRequest The <code>domain</code> is not valid
 * @apiErrorExample Domain not valid:
 *     HTTP/1.1 400 BadRequest
 *     Invalid Domain
 *
 */
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
