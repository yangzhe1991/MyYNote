package main
import (
	"net/http"
	"log"
	"fmt"
	"github.com/astaxie/beego/session"
	ynote "github.com/yangzhe1991/go-ynote"
)

var tmpCredDB map[string]*ynote.Credentials = make(map[string]*ynote.Credentials)

var globalSessions *session.Manager

func index(rw http.ResponseWriter, req *http.Request) {
	yc:=ynote.NewOnlineYnoteClient(ynote.Credentials{
		Token:  Token,
		Secret: Secret})
	sess:=globalSessions.SessionStart(rw,req)
	defer sess.SessionRelease()
	log.Print("/index")

	cred,ok:=sess.Get("accToken").(ynote.Credentials)


	if !ok {
		log.Print("no accToken in session")
		tmpCred, err := yc.RequestTemporaryCredentialsWithCallBack("http://127.0.0.1:8099/note/callback")
		if err != nil {
			return
		}
		tmpCredDB[tmpCred.Token]=tmpCred
		authUrl := yc.AuthorizationURL(tmpCred)

		http.Redirect(rw,req,authUrl,301)
		return
	}
	log.Print(cred)
	yc.AccToken=&cred
	ui,err:=yc.UserInfo()
	if err!=nil{
		fmt.Fprint(rw,"error!")
	}

	fmt.Fprint(rw,ui)



}


func callback(rw http.ResponseWriter, req *http.Request) {
	log.Print("/callback")
	yc:=ynote.NewOnlineYnoteClient(ynote.Credentials{
		Token:  Token,
		Secret: Secret})
	token:=req.FormValue("oauth_token")
	if token==""{
		return
	}
	verifier:=req.FormValue("oauth_verifier")
	if verifier==""{
		return
	}
	tmpCred,ok:=tmpCredDB[token]
	if !ok{
		return
	}

	accToken,err:=yc.RequestToken(tmpCred,verifier)
	if err!=nil{
		return
	}
	sess:=globalSessions.SessionStart(rw,req)
	defer sess.SessionRelease()
	sess.Set("accToken",*accToken)

	http.Redirect(rw,req,"/note",301)
}
func main(){
	globalSessions, _ = session.NewManager("memory", "session", 3600,"")
	go globalSessions.GC()
	http.HandleFunc("/note/callback",callback)
	http.HandleFunc("/note/",index)
	http.HandleFunc("/latex/",latex)
	log.Print("service start")
	err:=http.ListenAndServe(":8099",nil)
	if err!=nil{
		log.Fatal(err)
	}
	fmt.Print("asd")

}

