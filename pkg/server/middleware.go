package server

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"html/template"
	"net/http"
	"os"
	"path"
)

// SwaggerUIOpts configures the Swaggerui middlewares
type SwaggerUIOpts struct {
	// BasePath for the UI path, defaults to: /
	BasePath string
	// Path combines with BasePath for the full UI path, defaults to: docs
	Path string
	// SpecURL the url to find the spec for
	SpecURL string

	// The three components needed to embed swagger-ui
	SwaggerURL       string
	SwaggerPresetURL string
	SwaggerStylesURL string

	Favicon32 string
	Favicon16 string

	// Title for the documentation site, default to: API documentation
	Title string
}

// EnsureDefaults in case some options are missing
func (r *SwaggerUIOpts) EnsureDefaults() {
	if r.BasePath == "" {
		r.BasePath = "/"
	}
	if r.Path == "" {
		r.Path = "docs"
	}
	if r.SpecURL == "" {
		r.SpecURL = "/swagger.json"
	}
	if r.SwaggerURL == "" {
		r.SwaggerURL = swaggerLatest
	}
	if r.SwaggerPresetURL == "" {
		r.SwaggerPresetURL = swaggerPresetLatest
	}
	if r.SwaggerStylesURL == "" {
		r.SwaggerStylesURL = swaggerStylesLatest
	}
	if r.Favicon16 == "" {
		r.Favicon16 = swaggerFavicon16Latest
	}
	if r.Favicon32 == "" {
		r.Favicon32 = swaggerFavicon32Latest
	}
	if r.Title == "" {
		r.Title = "API documentation"
	}
}

func ServeOpenapi(router *gin.RouterGroup, opts SwaggerUIOpts) {
	dir, _ := os.Getwd()
	r := path.Join(dir, "api", "openapi")
	router.Static("/docs/", r+"/")
	router.GET("/docs", SwaggerUI(opts))
}

// SwaggerUI creates a middleware to serve a documentation site for a swagger spec.
// This allows for altering the spec before starting the http listener.
func SwaggerUI(opts SwaggerUIOpts) gin.HandlerFunc {
	opts.EnsureDefaults()

	isDev := os.Getenv("DEV")
	swaggerUi := swaggeruiTemplate
	if isDev != "" {
		swaggerUi = fmt.Sprintf(
			swaggerUi,
			SwaggerTemplateLogin,
			fmt.Sprintf(
				SwaggerFirebase,
				os.Getenv("API_KEY"),
				os.Getenv("AUTH_DOMAIN"),
				os.Getenv("PROJECT_ID"),
				os.Getenv("STORAGE_BUCKET"),
				os.Getenv("MESSAGING_SENDER_ID"),
				os.Getenv("APP_ID"),
			),
		)
	} else {
		swaggerUi = fmt.Sprintf(swaggerUi, "", "")
	}

	tmpl := template.Must(template.New("swaggerui").Parse(swaggerUi))

	buf := bytes.NewBuffer(nil)
	_ = tmpl.Execute(buf, &opts)
	bb := buf.Bytes()

	return func(c *gin.Context) {
		c.Data(http.StatusOK, "text/html; charset=utf-8", bb)
	}
}

const (
	SwaggerTemplateLogin = `<form id="login">
      <input type="email" required name="email">
      <input type="password" required minlength="6" name="password">
      <button type="submit">Generate token</button>
    </form>`
	SwaggerFirebase = `<!-- Firebase App (the core Firebase SDK) is always required and must be listed first -->
  <script src="https://www.gstatic.com/firebasejs/8.3.0/firebase-app.js"></script>

  <!-- If you enabled Analytics in your project, add the Firebase SDK for Analytics -->
  <script src="https://www.gstatic.com/firebasejs/8.3.0/firebase-analytics.js"></script>

  <!-- Add Firebase products that you want to use -->
  <script src="https://www.gstatic.com/firebasejs/8.3.0/firebase-auth.js"></script>
  <script src="https://www.gstatic.com/firebasejs/8.3.0/firebase-firestore.js"></script>
    <script>
      
		const firebaseConfig = {
			apiKey: "%s",
			authDomain: "%s",
			projectId: "%s",
			storageBucket: "%s",
			messagingSenderId: "%s",
			appId: "%s"
		};
		firebase.initializeApp(firebaseConfig);


		const getToken = (event) => {
			event.preventDefault();
			const email = event.target.email.value
			const password = event.target.password.value
		
			firebase.auth().signInWithEmailAndPassword(email, password)
				.then(fetchToken)
				.then(copyToken)
				.catch((err) => {
					console.error(err)
				})
		}
      
        const fetchToken = () => {
		return firebase
			.auth()
			.currentUser.getIdToken(false)
			.catch((error) => {
				alert('Login error');
			})
        }

		const copyToken = (token) => {
			const elem = document.createElement('input');
			elem.setAttribute('type', 'text');
			document.body.appendChild(elem);
			elem.style.position = 'absolute';
			elem.style.left = '-1000px';
			elem.style.top = '-1000px';
			elem.value = token;
			elem.select();
			elem.setSelectionRange(0, 99999);
			document.execCommand('copy');
			elem.parentNode.removeChild(elem);
			alert('Token copied!');
		}
		function handlerSubmit() {
			document.querySelector('#login').addEventListener('submit', getToken);
		}
		addOnload(handlerSubmit);
    </script>`
	swaggerLatest          = "https://unpkg.com/swagger-ui-dist/swagger-ui-bundle.js"
	swaggerPresetLatest    = "https://unpkg.com/swagger-ui-dist/swagger-ui-standalone-preset.js"
	swaggerStylesLatest    = "https://unpkg.com/swagger-ui-dist/swagger-ui.css"
	swaggerFavicon32Latest = "https://unpkg.com/swagger-ui-dist/favicon-32x32.png"
	swaggerFavicon16Latest = "https://unpkg.com/swagger-ui-dist/favicon-16x16.png"
	swaggeruiTemplate      = `
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="UTF-8">
		<title>{{ .Title }}</title>
    <link rel="stylesheet" type="text/css" href="{{ .SwaggerStylesURL }}" >
    <link rel="icon" type="image/png" href="{{ .Favicon32 }}" sizes="32x32" />
    <link rel="icon" type="image/png" href="{{ .Favicon16 }}" sizes="16x16" />
    <style>
		html
		{
			box-sizing: border-box;
			overflow: -moz-scrollbars-vertical;
			overflow-y: scroll;
		}
		*,
		*:before,
		*:after
		{
			box-sizing: inherit;
		}
		body
		{
			margin:0;
			background: #fafafa;
		}
    </style>
	<script>
		var loadsFuncs = [];
		window.onload = () => {
			loadsFuncs.forEach((fn) => fn());
		}
		function addOnload(fn){
			loadsFuncs.push(fn)
        }
	</script>
  </head>
  <body>
	%s
    <div id="swagger-ui"></div>
    <script src="{{ .SwaggerURL }}"> </script>
    <script src="{{ .SwaggerPresetURL }}"> </script>
    <script>
    addOnload(function() {
      // Begin Swagger UI call region
      const ui = SwaggerUIBundle({
        url: '{{ .SpecURL }}',
        dom_id: '#swagger-ui',
        deepLinking: true,
        presets: [
          SwaggerUIBundle.presets.apis,
          SwaggerUIStandalonePreset
        ],
        plugins: [
          SwaggerUIBundle.plugins.DownloadUrl
        ],
        layout: "StandaloneLayout"
      })
      // End Swagger UI call region
      window.ui = ui
    })
  </script>

  %s
  </body>
</html>
`
)
