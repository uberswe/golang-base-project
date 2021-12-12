package baseproject

import (
	"embed"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/uberswe/golang-base-project/middleware"
	"github.com/uberswe/golang-base-project/routes"
	"html/template"
	"io/fs"
	"log"
	"net/http"
)

//go:embed dist/*
var staticFS embed.FS

func Run() {
	var t *template.Template
	conf := loadEnvVariables()

	db, err := connectToDatabase(conf)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(db)

	err = migrateDatabase(db)
	if err != nil {
		log.Fatalln(err)
	}

	t, err = loadTemplates()
	if err != nil {
		log.Fatalln(err)
	}

	r := gin.Default()

	store := cookie.NewStore([]byte(conf.CookieSecret))
	r.Use(sessions.Sessions("golang_base_project_session", store))

	r.SetHTMLTemplate(t)

	subFS, err := fs.Sub(staticFS, "dist/assets")
	if err != nil {
		log.Fatalln(err)
	}

	r.StaticFS("/assets", http.FS(subFS))

	r.Use(middleware.Session(db))
	r.Use(middleware.General())

	controller := routes.New(db, conf)

	r.GET("/", controller.Index)
	r.GET("/search", controller.Search)
	r.POST("/search", controller.Search)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.NoRoute(controller.NoRoute)

	noAuth := r.Group("/")
	noAuth.Use(middleware.NoAuth())

	noAuth.GET("/login", controller.Login)
	noAuth.GET("/register", controller.Register)
	noAuth.GET("/activate/resend", controller.ResendActivation)
	noAuth.GET("/activate/:token", controller.Activate)
	noAuth.GET("/user/password/forgot", controller.ForgotPassword)
	noAuth.GET("/user/password/reset/:token", controller.ResetPassword)

	noAuthPost := noAuth.Group("/")
	noAuthPost.Use(middleware.Throttle(conf.RequestsPerMinute))

	noAuthPost.POST("/login", controller.LoginPost)
	noAuthPost.POST("/register", controller.RegisterPost)
	noAuthPost.POST("/activate/resend", controller.ResendActivationPost)
	noAuthPost.POST("/user/password/forgot", controller.ForgotPasswordPost)
	noAuthPost.POST("/user/password/reset/:token", controller.ResetPasswordPost)

	admin := r.Group("/")
	admin.Use(middleware.Auth())
	admin.Use(middleware.Sensitive())

	admin.GET("/admin", controller.Admin)
	// We need to handle post from the login redirect
	admin.POST("/admin", controller.Admin)
	admin.GET("/logout", controller.Logout)

	err = r.Run(conf.Port)
	if err != nil {
		log.Fatalln(err)
	}
}
