package html

var timelineTmpl = `
<!DOCTYPE html>
<html>

<head>
<style>
a {
	text-decoration: none;
	display: inline-block;
	padding: 8px 16px;
}

a:hover {
	background-color: #ddd;
	color: black;
}

.previous {
	background-color: #f1f1f1;
	color: black;
}

.next {
	background-color: #2b7bb9;
	color: white;
	float: right;
}
</style>
</head>

<body style="background-color:#292F33;">
	<div id="container" style="width:800px; margin:0 auto;">
		<a href="{{.PreviousPage}}.html" class="previous">&laquo; Previous</a>
		<a href="{{.NextPage}}.html" class="next">Next &raquo;</a>
		{{range .Users}}
			<blockquote class="twitter-tweet">
				<a class="twitter-timeline" data-tweet-limit="1" data-width="700" data-dnt="true" data-theme="dark" href="https://twitter.com/{{ .ScreenName}}"></a>
			</blockquote>
		{{end}}
		<a href="{{.PreviousPage}}.html" class="previous">&laquo; Previous</a>
		<a href="{{.NextPage}}.html" class="next">Next &raquo;</a>
		<script async src="https://platform.twitter.com/widgets.js" charset="utf-8"></script>
    </div>
</body> 
</html>


`
