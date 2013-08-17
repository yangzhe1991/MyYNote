package main

import (

	"net/http"
	"log"
	"fmt"
	"io/ioutil"
	"io"
	"crypto/md5"
	"os/exec"
	"bytes"
)

func latex(rw http.ResponseWriter,req *http.Request){
	if req.Method=="POST"{
		tex:=req.FormValue("tex")
		h := md5.New()
		io.WriteString(h, tex)
		filename:=fmt.Sprintf("%x", h.Sum(nil))
		ioutil.WriteFile(filename+".tex", []byte(tex),444)
		cmd := exec.Command("pdflatex", filename+".tex")
		var out,err bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &err
		if cmd.Run() !=nil{
			rw.Write(out.Bytes())
			rw.Write(err.Bytes())
			log.Print("Build PDF fail")
			return
		}
		log.Print("Build Success")
		html:=fmt.Sprintf("<html><head><meta http-equiv=\"refresh\" content=\"0;url=/latex/%s.pdf\" /></head></html>",filename)
		fmt.Fprint(rw,html)

	}

	if req.Method=="GET" {
		path:=req.URL.Path
		if len(path)<10 {
			log.Print("GET")
			fmt.Fprint(rw, html)
			return
		}
		path=path[7:]
		log.Print("GET "+path)
		http.ServeFile(rw,req,path)


	}


}



func main() {
	http.HandleFunc("/latex/",latex)
	log.Print("service start")
	err:=http.ListenAndServe(":8099",nil)
	if err!=nil{
		log.Fatal(err)
	}


}

var html string=`<html>
				<head>
				<meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
				<title>在线LaTex编译器</title>
				</head>
				<body>
                买的Linode只搭了个个人博客和ACMDIY群的主页，利用率太低，加上自己实在懒得在各种电脑上装LaTex，顺便练手写下Go，于是有了这么个东西。
                <br>
                texlive2013，暂不支持中文，字体啥的也只有默认的。
                <br>
                这本来是个大坑，想利用有道云笔记和Golang写点玩具。比如把tex代码存在有道云笔记中，就成了一个在线的Latex IDE。然后可以顺便支持markdown啊，代码高亮啊啥的。不保证真的去填，更不保证填完、填得漂亮。
                <br>
                暂时没有用户验证之类的，所以理论上可以看别人的代码和pdf（当然前提是能猜出文件名……）
                <br>
				<form method="POST">
				<textarea name="tex" cols=100 rows=30 > </textarea>
				<input type="submit" />
				</form>

				<script type="text/javascript">
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
