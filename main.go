package main

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strconv"

	_ "github.com/lib/pq"
)

type server struct {
	db         *sql.DB
	frontProxy *httputil.ReverseProxy
}
type link struct {
	Long  string
	Short string
}

func main() {
	db, err := sql.Open("postgres", "postgres://postgres:"+os.Getenv("NIRHUB_PG_PASS")+"@nirhub.ru:65432/nirhub")
	if err != nil {
		log.Fatal(err)
	}
	s := &server{
		db,
		httputil.NewSingleHostReverseProxy(&url.URL{
			Host:   "localhost:3000",
			Scheme: "http",
		}),
	}

	if err := http.ListenAndServe(":8080", s); err != nil {
		log.Fatal(err)
	}
}

func (s *server) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "POST":
		switch {
		case req.URL.Path == "/link":
			s.newLink(res, req)
			return
		}
	case "GET":
		_, intErr := strconv.Atoi(req.URL.Path[1:])
		switch {
		case req.URL.Path == "/link":
			s.getAllLinks(res, req)
			return
		case intErr == nil:
			s.navigateByLink(res, req)
			return
		}
	}
	s.frontProxy.ServeHTTP(res, req)
}

func (s *server) newLink(res http.ResponseWriter, req *http.Request) {
	link, err := ioutil.ReadAll(req.Body)
	if err != nil {
		res.WriteHeader(500)
		res.Write([]byte(err.Error()))
		return
	}
	_, err = s.db.Exec("insert into nirhub.shortner (link) values ($1) returning (id)", string(link))
	if err != nil {
		res.WriteHeader(500)
		res.Write([]byte(err.Error()))
		return
	}
	s.getAllLinks(res, req)
}

func (s *server) getAllLinks(res http.ResponseWriter, req *http.Request) {
	rows, err := s.db.Query("select id, link from nirhub.shortner order by created desc")
	if err != nil {
		res.WriteHeader(500)
		res.Write([]byte(err.Error()))
		return
	}
	links := []link{}
	for rows.Next() {
		link := link{}
		rows.Scan(&link.Short, &link.Long)
		links = append(links, link)
	}
	linksJSON, err := json.Marshal(links)
	if err != nil {
		res.WriteHeader(500)
		res.Write([]byte(err.Error()))
		return
	}
	res.Write(linksJSON)
}

func (s *server) navigateByLink(res http.ResponseWriter, req *http.Request) {
	if len(req.URL.Path) < 2 {
		res.Header().Add("Location", "/")
		res.WriteHeader(301)
		return
	}
	var link string
	err := s.db.QueryRow("select link from nirhub.shortner where id=$1", req.URL.Path[1:]).Scan(&link)
	if err != nil {
		res.WriteHeader(500)
		res.Write([]byte(err.Error()))
		return
	}
	res.Header().Add("Location", link)
	res.WriteHeader(301)
	return
}
