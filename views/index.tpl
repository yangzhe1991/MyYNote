<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
<title>有道云笔记Latex IDE</title>
    <link href="http://cdn.staticfile.org/twitter-bootstrap/2.3.2/css/bootstrap-responsive.min.css" rel="stylesheet">
    <link href="http://cdn.staticfile.org/twitter-bootstrap/2.3.2/css/bootstrap.min.css" rel="stylesheet">

    <script src="http://cdn.staticfile.org/jquery/1.10.2/jquery.min.js"></script>
    <script src="http://cdn.staticfile.org/twitter-bootstrap/2.3.2/js/bootstrap.min.js"></script>
</head>

<body>
<div>基于有道云笔记的在线Latex编辑器，能看到这个页面的说明已经用网易通行证登录了，会创建一个叫“latex”的笔记本，新建两个文件，如果两个不够请自行去有道云笔记里面增加，然后在这边刷新就能看到了。<br>写着玩的，用beego实现，不保证随时可以用，不保证没有bug。代码在<a href="https://github.com/yangzhe1991/MyYNote">github上</a>，博客有个文章说了下开发这个小玩意之后的粗浅想法，<a href="http://yangzhe1991.org/blog/2013/09/youdao-note-latex-golang-beego/">在这里</a></div>
<div class="container-fluid">
  <div class="row-fluid">
    <div class="span2">
      <ul id="list" class="nav nav-list">
      		{{range .Notes}}
              <li ><a href="#{{.Path}}">{{.Title}}</a></li>
            {{end}}
      </ul>
    </div>
    <div class="span10">
	    <div class="row-fluid">
	    	<div class="span5"> 
		      	<form id="form">
		     		<div>
		     			<input type="button" id="submit" value="保存并编译"/>
		     		</div>
		   	   		<div>title</div>
		    	  	<div>
		    	  		<input class="span6" type="text" name="title" id="notetitle"/>
		    	  	</div>
		      		<div>代码</div>
		      		<div>
		      			<textarea class="span12" rows="30" name="content" id="codearea">
		      			</textarea>
		      		</div>
                    <input type="hidden" id="path" name="path"/>		      		
		      	
		      	</form>
	     	</div>
	     	<div class="span7">
	     		<div id="compile"></div>
	     		<iframe id="pdf" class="span12" height="700"></iframe>
	     	</div>
	    </div>
	</div>
    
  </div>
</div>



<script>
var eleMenuOn = null, eleListBox = $("#listBox");

String.prototype.temp = function(obj) {
    return this.replace(/\$\w+\$/gi, function(matchs) {
        var returns = obj[matchs.replace(/\$/g, "")];		
        return (returns + "") == "undefined"? "": returns;
    });
};
var mydomain="http://yangzhe1991.org/note/"
var eleMenus = $("#list a").bind("click", function(event) {
	var query = this.href.split("#")[1];
	if (history.pushState && query && !this.parentNode.classList.contains("active")) {
		eleMenuOn && eleMenuOn.classList.remove("active");
		eleMenuOn = this.parentNode;
		eleMenuOn.classList.add("active");
		$.getJSON(mydomain+"json/"+query, function(data) {
			$("#notetitle")[0].value=data.Title;
			$("#codearea")[0].value=data.Content;//.replace(/<div>/g,"").replace(/<\/div>/g,"");
		});
		
		// history处理
		var title = "有道云Latex-" + this.text.replace(/\d+$/, "");
		document.title = title;		
		history.pushState({ title: title }, title, location.href.split("#")[0] + "#" + query);
		
	}
	return false;
});

var fnHashTrigger = function(target) {
	var query = location.href.split("#")[1], eleTarget = target || null;
	if (typeof query == "undefined") {
		if (eleTarget = eleMenus.get(0)) {
			history.replaceState(null, document.title, location.href.split("#")[0] + "#" + eleTarget.href.split("#")[1]) + location.hash;	
			fnHashTrigger(eleTarget);
		}
	} else {
		eleMenus.each(function() {
			if (eleTarget === null && this.href.split("#")[1] === query) {
				eleTarget = this;
			}
		});
		
		if (!eleTarget) {
			history.replaceState(null, document.title, location.href.split("#")[0]);	
			fnHashTrigger();
		} else {			 
			$(eleTarget).trigger("click");
		}		
	}	
};
if (history.pushState) {
	window.addEventListener("popstate", function() {
		fnHashTrigger();																
	});
	
	// 默认载入
	fnHashTrigger();
}

$("#submit").click(function(){
    $("#path").val(location.href.split("#")[1]);
	$.post("/note/latex/", 
		$("#form").serialize(),
		function(data, textStatus, xhr) {
			if(data.Result=="success"){
				$("#pdf")[0].src=mydomain+"latex/"+data.Content;
				$("#compile").html("");
			}
			else{
				$("#compile").html(data.Content.replace(/\n/,"<br>"));
			}
		},
		"json"
		);
});

</script>
    
	
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
</html>
