 <!DOCTYPE html>
<html>
<head>
<title>Healthchecks</title>
<style type="text/css">
.category {
	clear: both;
	display: block;
}
.category .name {
	font-size: 26px;
}
.healthcheck {
	display: block;
	width: 150px;
	height: 150px;
	float: left;
	padding: 10px;
	margin: 5px;
	border: 1px solid #666;
}

.healthcheck.healthy {
	background-color: green;
}

.healthcheck.not-healthy {
	background-color: red;
}

.healthcheck a {
	color:  black;
	text-decoration: none;
}
</style>
</head>
<body>
	<div class="healthchecks">
	{{range $category,$checks := .}}
		<div class="category">
			<div class="name"> {{ $category }} </div>

			{{range $checks}}
				<div class="healthcheck {{if .IsHealthy}}healthy{{else}}not-healthy{{end}}" {{if .Err}} title="{{.Err.Error}}" {{end}}>
					<a href="{{.Service.URL}}/__health" target="_blank"> {{.Service.Name}}</a>
				</div>
			{{end}}
		</div>
	{{end}}
	</div>
</body>
</html>
