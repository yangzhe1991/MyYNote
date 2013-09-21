package controllers

import (
	"MyYNote/conf"
	"github.com/astaxie/beego"
	ynote "github.com/youdao-api/go-ynote"

	"net/url"
)

var tmpCredDB map[string]*ynote.Credentials = make(map[string]*ynote.Credentials)

var WEBROOT = "http://localhost:8000/note/"

type ViewNotes struct {
	Path  string
	Title string
}

var globalToken ynote.Credentials

type MainController struct {
	beego.Controller
}

func (this *MainController) Get() {
	beego.Info("/")
	yc := ynote.NewOnlineYnoteClient(ynote.Credentials{
		Token:  conf.Key,
		Secret: conf.Secret})

	a := this.GetSession("accToken")
	var cred ynote.Credentials
	if a != nil {
		cred = a.(ynote.Credentials)
	} else {
		beego.Info("no accToken in session")
		tmpCred, err := yc.RequestTemporaryCredentialsWithCallBack(WEBROOT + "callback")
		if err != nil {
			this.Redirect("/note/error", 301)
			return
		}
		tmpCredDB[tmpCred.Token] = tmpCred
		authUrl := yc.AuthorizationURL(tmpCred)

		this.Redirect(authUrl, 301)
		return
	}

	//cred := globalToken
	yc.AccToken = &cred
	_, err := yc.UserInfo()

	if err != nil {
		this.Redirect("/note/error", 301)
		return
	}

	notebooks, err := yc.ListNotebooks()
	if err != nil {
		this.Redirect("/note/error", 301)
		return
	}
	var mynotebook *ynote.NotebookInfo = nil
	for _, notebook := range notebooks {
		if notebook.Name == "latex" {
			mynotebook = notebook
		}
	}
	if mynotebook == nil {
		mynotebook, err = yc.CreateNotebook("latex", "")
		if err != nil {
			this.Redirect("/note/error", 301)
			return
		}
	}
	notes, err := yc.ListNotes(mynotebook.Path)
	if err != nil {
		this.Redirect("/note/error", 301)
		return
	}
	ns := []ViewNotes{}
	for _, notepath := range notes {
		note, err := yc.NoteInfo(notepath)
		if err != nil {
			this.Redirect("/note/error", 301)
			return
		}
		ns = append(ns, ViewNotes{notepath, note.Title})
	}
	beego.Info(ns)

	this.Data["Notes"] = ns
	this.TplNames = "index.tpl"
}

type CallbackController struct {
	beego.Controller
}

func (this *CallbackController) Get() {
	beego.Info("/callback")
	yc := ynote.NewOnlineYnoteClient(ynote.Credentials{
		Token:  conf.Key,
		Secret: conf.Secret})

	token := this.GetString("oauth_token")
	if token == "" {
		this.Redirect("/note/error", 301)
		return
	}
	verifier := this.GetString("oauth_verifier")
	if verifier == "" {
		this.Redirect("/note/error", 301)
		return
	}
	tmpCred, ok := tmpCredDB[token]
	if !ok {
		this.Redirect("/note/error", 301)
		return
	}

	accToken, err := yc.RequestToken(tmpCred, verifier)
	if err != nil {
		this.Redirect("/note/error", 301)
		return
	}

	this.SetSession("accToken", *accToken)
	globalToken = *accToken
	this.Ctx.Redirect(301, "/note/")
	//this.Ctx.WriteString("done")
}

type JsonController struct {
	beego.Controller
}

func (this *JsonController) Get() {
	yc := ynote.NewOnlineYnoteClient(ynote.Credentials{
		Token:  conf.Key,
		Secret: conf.Secret})
	beego.Info("/json")
	a := this.GetSession("accToken")
	var cred ynote.Credentials
	if a != nil {
		cred = a.(ynote.Credentials)
	} else {
		this.Data["error"] = true
		beego.Info(1)
		this.ServeJson()
		return
	}
	path, err := url.QueryUnescape(this.Ctx.Input.Params(":splat"))
	if err != nil || path == "" {
		this.Data["error"] = true
		beego.Info(2)
		this.ServeJson()
		return
	}

	yc.AccToken = &cred

	note, err := yc.NoteInfo(path)
	if err != nil {
		this.Data["error"] = true
		beego.Info(3)
		this.ServeJson()
		return
	}
	this.Data["json"] = note
	this.Data["content"] = note.Content
	beego.Info(note.Title)

	this.ServeJson()
}
