package main

import (
	"html/template"
	"net/http"
	"github.com/satori/go.uuid"

)
type user struct{
username string
first string
last string
}

var dbSessions= map[string]string{}
var dbUsers= map[string]user{}
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
	c:=getCookie(w,req) 											//to get cookie in index
	tpl.ExecuteTemplate(w, "index.gohtml", c)	//executes index.gohtml
}

func getCookie(w http.ResponseWriter, req *http.Request) *http.Cookie{		//used to generate cookie
	c, err := req.Cookie("Session")
	if err!=nil{
		sID, _:= uuid.NewV4()
		c = &http.Cookie{
			Name:"Session",
			Value:sID.String(),
		}
	}
	return c
}