package main

import (

	"fmt"
	"github.com/astaxie/beego/session"
	ynote "github.com/yangzhe1991/go-ynote"
	"log"
	"net/http"
)

var tmpCredDB map[string]*ynote.Credentials = make(map[string]*ynote.Credentials)

var globalSessions *session.Manager

func index(rw http.ResponseWriter, req *http.Request) {
	yc := ynote.NewOnlineYnoteClient(ynote.Credentials{
		Token:  conf.Token,
		Secret: conf.Secret})
	sess := globalSessions.SessionStart(rw, req)
	defer sess.SessionRelease()
	log.Print("/")

	cred, ok := sess.Get("accToken").(ynote.Credentials)

	if !ok {
		log.Print("no accToken in session")
		tmpCred, err := yc.RequestTemporaryCredentialsWithCallBack("http://127.0.0.1:8099/note/callback")
		if err != nil {
			http.Redirect(rw, req, "/note/error", 301)
			return
		}
		tmpCredDB[tmpCred.Token] = tmpCred
		authUrl := yc.AuthorizationURL(tmpCred)

		http.Redirect(rw, req, authUrl, 301)
		return
	}
	log.Print(cred)
	yc.AccToken = &cred
	_, err := yc.UserInfo()
	if err != nil {
		http.Redirect(rw, req, "/note/error", 301)
		return
	}

	notebooks, err := yc.ListNotebooks()
	if err != nil {
		http.Redirect(rw, req, "/note/error", 301)
		return
	}
	var mynotebook ynote.NotebookInfo = nil
	for _, notebook := range notebooks {
		if notebook.Name == "第三方应用" {
			mynotebook = notebook
		}
	}
	if mynotebook == nil {
		mynotebook, err = yc.CreateNotebook("第三方应用", "")
		if err != nil {
			http.Redirect(rw, req, "/note/error", 301)
			return
		}
	}
	notes, err := yc.ListNotes(mynotebook.Path)
	if err != nil {
		http.Redirect(rw, req, "/note/error", 301)
		return
	}
	fmt.Fprintf(rw, html1)
	for _, notepath := range notes {
		note, err := yc.NoteInfo(notepath)
		if err != nil {
			http.Redirect(rw, req, "/note/error", 301)
			return
		}
		fmt.Fprintf(rw, "<a href=\"/note/edit?path=%s\">%s</a><br>", notepath, note.Title)
	}
	fmt.Fprint(rw, "<a href=\"/note/create\">create a new note</a><br>")

	fmt.Fprintf(rw, html2)
}

func edit(rw http.ResponseWriter, req *http.Request) {
	yc := ynote.NewOnlineYnoteClient(ynote.Credentials{
		Token:  Token,
		Secret: Secret})
	sess := globalSessions.SessionStart(rw, req)
	defer sess.SessionRelease()
	log.Print("/edit")
	cred, ok := sess.Get("accToken").(ynote.Credentials)

	if !ok || req.FormValue("path") == "" {
		log.Print("no accToken in session")
		http.Redirect(rw, req, "/note", 301)
		return
	}

	yc.AccToken = cred

	note, err := yc.NoteInfo(req.FormValue("path"))
	if err != nil {
		http.Redirect(rw, req, "/note/error", 301)
		return
	}
	if req.Method == "GET" {
		fmt.Fprintf(rw, html1)
		fmt.Fprint(rw, `
			<form>
			title<br>
			 <input name="title" />
			`)
		fmt.Fprint(rw, html2)
	}

}

func error(rw http.ResponseWriter, req *http.Request) {
	log.Print("/error")
	fmt.Fprintf(rw, html1)
	fmt.Fprintf(rw, "有地方出错了……")
	fmt.Fprintf(rw, html2)
}

func callback(rw http.ResponseWriter, req *http.Request) {
	log.Print("/callback")
	yc := ynote.NewOnlineYnoteClient(ynote.Credentials{
		Token:  Token,
		Secret: Secret})
	token := req.FormValue("oauth_token")
	if token == "" {
		return
	}
	verifier := req.FormValue("oauth_verifier")
	if verifier == "" {
		return
	}
	tmpCred, ok := tmpCredDB[token]
	if !ok {
		return
	}

	accToken, err := yc.RequestToken(tmpCred, verifier)
	if err != nil {
		return
	}
	sess := globalSessions.SessionStart(rw, req)
	defer sess.SessionRelease()
	sess.Set("accToken", *accToken)

	http.Redirect(rw, req, "/note", 301)
}
func main() {
	globalSessions, _ = session.NewManager("memory", "session", 3600, "")
	go globalSessions.GC()
	http.HandleFunc("/note/callback", callback)
	http.HandleFunc("/note/edit", edit)
	http.HandleFunc("/note/error", error)
	http.HandleFunc("/note/create", create)
	http.HandleFunc("/note/", index)
	http.HandleFunc("/latex/", latex)
	log.Print("service start")
	err := http.ListenAndServe(":8099", nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print("asd")

}

var html1 string = `<html>
				<head>
				<meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
				<title>无名应用</title>
				</head>
				<body>`
var html2 string = `<script type="text/javascript">
				  var _gaq = _gaq || [];
				  _gaq.push(['_setAccount', 'UA-29543335-1']);
				  _gaq.push(['_setDomainName', 'yangzhe1991.org']);
				  _gaq.push(['_trackPageview']);
				  (function() {
					var ga = document.createElement('script'); ga.type = 'text/javascript'; ga.async = true;
					ga.src = ('https:' == document.location.protocol ? 'https://ssl' : 'http://www') + '.google-analytics.com/ga.js';
					var s = document.getElementsByTagName('script')[0]; s.parentNode.insertBefore(ga, s);
				  })();
				</script>
				</body>

				</html>`
