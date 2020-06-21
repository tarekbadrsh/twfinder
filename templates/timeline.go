package templates

// Timeline :
var Timeline = `
<!DOCTYPE html>
<html>

<head>
    <meta charset="utf-8" />
</head>

<body>
	<div id="container" style="width:800px; margin:0 auto;">
		{{range .}}
			<blockquote class="twitter-tweet">
				<a class="twitter-timeline" data-tweet-limit="1" data-width="500" href="https://twitter.com/{{ .ScreenName}}"></a>
			</blockquote>
		{{end}}
		<script async src="https://platform.twitter.com/widgets.js" charset="utf-8"></script>
    </div>
</body> 
</html>


`
