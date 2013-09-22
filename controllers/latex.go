package controllers

import (
	"github.com/astaxie/beego"
    ynote "github.com/youdao-api/go-ynote"
    "MyYNote/conf"

	"bytes"
	"crypto/md5"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os/exec"
)

type LatexController struct {
	beego.Controller
}

func (this *LatexController) Get() {
	beego.Info("/latex get")
	path, err := url.QueryUnescape(this.Ctx.Input.Params(":splat"))
	if err != nil || path == "" {
		this.Abort("404")
	}
	http.ServeFile(this.Ctx.ResponseWriter, this.Ctx.Request, path+".pdf")
}

type ViewResults struct {
	Result  string
	Content string
}

func (this *LatexController) Post() {
	beego.Info("/latex post")
	tex := this.GetString("content")
	h := md5.New()
	io.WriteString(h, tex)
	filename := fmt.Sprintf("%x", h.Sum(nil))
	ioutil.WriteFile(filename+".tex", []byte(tex), 444)
	cmd := exec.Command("xelatex", filename+".tex")
	var out, errst bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errst
	if cmd.Run() != nil {
		this.Data["json"] = ViewResults{
			"error", string(out.Bytes()) + string(errst.Bytes())}

	} else {
		this.Data["json"] = ViewResults{
			"success", filename}
	}
    
    yc := ynote.NewOnlineYnoteClient(ynote.Credentials{
        Token:  conf.Key,
        Secret: conf.Secret})
    a := this.GetSession("accToken")
    var cred ynote.Credentials
    if a != nil {
        cred = a.(ynote.Credentials)
        yc.AccToken = &cred
        path,_:=url.QueryUnescape(this.GetString("path"))
        error:=yc.UpdateNote(path, this.GetString("title"), "", "", this.GetString("content"))
        if error!=nil{
            beego.Info(error)
        }
    } 

	this.ServeJson()

}
