
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
<div class="container-fluid">
  <div class="row-fluid">
    <div class="span1">
      <ul id="list" class="nav nav-list">
      		{{range .Notes}}
              <li ><a href="#{{.Path}}">{{.Title}}</a></li>
            {{end}}
      </ul>
    </div>
    <div class="span11">
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
	$.post("/note/latex/", 
		$("#codearea").serialize(),
		function(data, textStatus, xhr) {
			if(data.Result=="success"){
				$("#pdf").src=mydomain+"latex/"+data.Content;
				$("#compile").html();
			}
			else{
				$("#compile").html(data.Content.replace(/\n/,"<br>"));
			}
		},
		"json"
		);
});

</script>
    
	
	<!--script type="text/javascript">
	var _gaq = _gaq || [];
	_gaq.push(['_setAccount', 'UA-29543335-1']);
	_gaq.push(['_setDomainName', 'yangzhe1991.org']);
	_gaq.push(['_trackPageview']);
	(function() {
		var ga = document.createElement('script'); ga.type = 'text/javascript'; ga.async = true;
		ga.src = ('https:' == document.location.protocol ? 'https://ssl' : 'http://www') + '.google-analytics.com/ga.js';
		var s = document.getElementsByTagName('script')[0]; s.parentNode.insertBefore(ga, s);
	})();
	</script-->
</body>
</html>
