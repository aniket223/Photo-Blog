package main

import (
	"github.com/satori/go.uuid"
	"html/template"
	"net/http"
	"strings"
	"fmt"
	"io"
	"os"
	"crypto/sha1"
	"path/filepath"
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
	// add route to serve pictures
	http.Handle("/public/", http.StripPrefix("/public", http.FileServer(http.Dir("./public"))))
	http.Handle("/favicon.ico", http.NotFoundHandler())
	http.ListenAndServe(":8081", nil) //default serve mux
}

func index(w http.ResponseWriter, req *http.Request) {
	c := getCookie(w, req)                          //to get cookie in index
	if req.Method==http.MethodPost{
		mf,fh,err:= req.FormFile("nf")
		if err!=nil{
			fmt.Println(err)
		}
		defer mf.Close()
		ext := strings.Split(fh.Filename, ".")[1]
		h:=sha1.New()
		io.Copy(h,mf)
		fname:=fmt.Sprintf("%x",h.Sum(nil))+"."+ext
		//create new file
		wd,err:=os.Getwd()
		if err!=nil{
			fmt.Println(err)
		}
		path:=filepath.Join(wd,"public","pics",fname)
		nf,err:=os.Create(path)
		if err!=nil{
			fmt.Println(err)
		}
		defer nf.Close()
		//copy
		mf.Seek(0,0)
		io.Copy(nf,mf)
		//add filename to user's cookie
		c=appendValue(w,c,fname)
	}
	xs:= strings.Split(c.Value,"|")
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

func appendValue(w http.ResponseWriter, c *http.Cookie, fname string) *http.Cookie {
	// values
	s:=c.Value
	if !strings.Contains(s,fname){
		s+="|"+fname
	}
	c.Value=s
	http.SetCookie(w,c)
	return c
}
