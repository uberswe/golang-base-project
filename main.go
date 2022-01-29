package baseproject

import (
	"embed"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/uberswe/golang-base-project/middleware"
	"github.com/uberswe/golang-base-project/routes"
	"golang.org/x/text/language"
	"html/template"
	"io/fs"
	"log"
	"math/rand"
	"net/http"
	"time"
)

// staticFS is an embedded file system
//go:embed dist/*
var staticFS embed.FS

// Run is the main function that runs the entire package and starts the webserver, this is called by /cmd/base/main.go
func Run() {
	// When generating random strings we need to provide a seed otherwise we always get the same strings the next time our application starts
	rand.Seed(time.Now().UnixNano())

	// We load environment variables, these are only read when the application launches
	conf := loadEnvVariables()

	// Translations
	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	languages := []string{
		"en",
		"sv",
	}
	for _, l := range languages {
		_, err := bundle.LoadMessageFile(fmt.Sprintf("active.%s.toml", l))
		if err != nil {
			log.Fatalln(err)
		}
	}

	// We connect to the database using the configuration generated from the environment variables.
	db, err := connectToDatabase(conf)
	if err != nil {
		log.Fatalln(err)
	}

	// Once a database connection is established we run any needed migrations
	err = migrateDatabase(db)
	if err != nil {
		log.Fatalln(err)
	}
	// t will hold all our html templates used to render pages
	var t *template.Template

	// We parse and load the html files into our t variable
	t, err = loadTemplates()
	if err != nil {
		log.Fatalln(err)
	}

	// A gin Engine instance with the default configuration
	r := gin.Default()

	// We create a new cookie store with a key used to secure cookies with HMAC
	store := cookie.NewStore([]byte(conf.CookieSecret))

	// We define our session middleware to be used globally on all routes
	r.Use(sessions.Sessions("golang_base_project_session", store))

	// We pase our template variable t to the gin engine so it can be used to render html pages
	r.SetHTMLTemplate(t)

	// our assets are only located in a section of our file system. so we create a sub file system.
	subFS, err := fs.Sub(staticFS, "dist/assets")
	if err != nil {
		log.Fatalln(err)
	}

	// All static assets are under the /assets path so we make this its own group called assets
	assets := r.Group("/assets")

	// This middleware sets the Cache-Control header and is applied to the assets group only
	assets.Use(middleware.Cache(conf.CacheMaxAge))

	// All requests to /assets will use the sub fil system which contains all our static assets
	assets.StaticFS("/", http.FS(subFS))

	// Session middleware is applied to all groups after this point.
	r.Use(middleware.Session(db))

	// A General middleware is defined to add default headers to improve site security
	r.Use(middleware.General())

	// A new instance of the routes controller is created
	controller := routes.New(db, conf, bundle)

	// Any request to / will call controller.Index
	r.GET("/", controller.Index)

	// We want to handle both POST and GET requests on the /search route. We define both but use the same function to handle the requests.
	r.GET("/search", controller.Search)
	r.POST("/search", controller.Search)

	// We define our 404 handler for when a page can not be found
	r.NoRoute(controller.NoRoute)

	// noAuth is a group for routes which should only be accessed if the user is not authenticated
	noAuth := r.Group("/")
	noAuth.Use(middleware.NoAuth())

	noAuth.GET("/login", controller.Login)
	noAuth.GET("/register", controller.Register)
	noAuth.GET("/activate/resend", controller.ResendActivation)
	noAuth.GET("/activate/:token", controller.Activate)
	noAuth.GET("/user/password/forgot", controller.ForgotPassword)
	noAuth.GET("/user/password/reset/:token", controller.ResetPassword)

	// We make a separate group for our post requests on the same endpoints so that we can define our throttling middleware on POST requests only.
	noAuthPost := noAuth.Group("/")
	noAuthPost.Use(middleware.Throttle(conf.RequestsPerMinute))

	noAuthPost.POST("/login", controller.LoginPost)
	noAuthPost.POST("/register", controller.RegisterPost)
	noAuthPost.POST("/activate/resend", controller.ResendActivationPost)
	noAuthPost.POST("/user/password/forgot", controller.ForgotPasswordPost)
	noAuthPost.POST("/user/password/reset/:token", controller.ResetPasswordPost)

	// the admin group handles routes that should only be accessible to authenticated users
	admin := r.Group("/")
	admin.Use(middleware.Auth())
	admin.Use(middleware.Sensitive())

	admin.GET("/admin", controller.Admin)
	// We need to handle post from the login redirect
	admin.POST("/admin", controller.Admin)
	admin.GET("/logout", controller.Logout)

	// This starts our webserver, our application will not stop running or go past this point unless
	// an error occurs or the web server is stopped for some reason. It is designed to run forever.
	err = r.Run(conf.Port)
	if err != nil {
		log.Fatalln(err)
	}
}
