package main

import (
	"github.com/satori/go.uuid"
	"html/template"
	"net/http"
	"strings"
)

type user struct {
	username string
	first    string
	last     string
}

var dbSessions = map[string]string{}
var dbUsers = map[string]user{}
var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("templates/*")) //parses all templates
}

func main() {
	http.HandleFunc("/", index)
	http.Handle("/favicon.ico", http.NotFoundHandler())
	http.ListenAndServe(":8080", nil) //default serve mux
}

func index(w http.ResponseWriter, req *http.Request) {
	c := getCookie(w, req)                          //to get cookie in index
	c = appendValue(w, c)
	xs := strings.Split(c.Value, "|")
	tpl.ExecuteTemplate(w, "index.gohtml", xs) //executes index.gohtml with value of cookie
}

func getCookie(w http.ResponseWriter, req *http.Request) *http.Cookie { //used to generate cookie
	c, err := req.Cookie("Session")
	if err != nil {
		sID, _ := uuid.NewV4()
		c = &http.Cookie{
			Name:  "Session",
			Value: sID.String(),
		}
	}
	return c
}

func appendValue(w http.ResponseWriter, c *http.Cookie) *http.Cookie {
	// values
	p1 := "disneyland.jpg"
	p2 := "atbeach.jpg"
	p3 := "hollywood.jpg"
	// append
	s := c.Value
	if !strings.Contains(s, p1) {
		s += "|" + p1
	}
	if !strings.Contains(s, p2) {
		s += "|" + p2
	}
	if !strings.Contains(s, p3) {
		s += "|" + p3
	}
	c.Value = s
	http.SetCookie(w, c)
	return c
}
