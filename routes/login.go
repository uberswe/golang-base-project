package routes

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/uberswe/golang-base-project/middleware"
	"github.com/uberswe/golang-base-project/models"
	"github.com/uberswe/golang-base-project/ulid"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"time"
)

// Login renders the HTML of the login page
func (controller Controller) Login(c *gin.Context) {
	pd := PageData{
		Title:           "Login",
		IsAuthenticated: isAuthenticated(c),
		CacheParameter:  controller.config.CacheParameter,
	}
	c.HTML(http.StatusOK, "login.html", pd)
}

// LoginPost handles login requests and returns the appropriate HTML and messages
func (controller Controller) LoginPost(c *gin.Context) {
	loginError := "Could not login, please make sure that you have typed in the correct email and password. If you have forgotten your password, please click the forgot password link below."
	pd := PageData{
		Title:           "Login",
		IsAuthenticated: isAuthenticated(c),
		CacheParameter:  controller.config.CacheParameter,
	}
	email := c.PostForm("email")
	user := models.User{Email: email}

	res := controller.db.Where(&user).First(&user)
	if res.Error != nil {
		pd.Messages = append(pd.Messages, Message{
			Type:    "error",
			Content: loginError,
		})
		log.Println(res.Error)
		c.HTML(http.StatusInternalServerError, "login.html", pd)
		return
	}

	if res.RowsAffected == 0 {
		pd.Messages = append(pd.Messages, Message{
			Type:    "error",
			Content: loginError,
		})
		c.HTML(http.StatusBadRequest, "login.html", pd)
		return
	}

	if user.ActivatedAt == nil {
		pd.Messages = append(pd.Messages, Message{
			Type:    "error",
			Content: "Account is not activated yet.",
		})
		c.HTML(http.StatusBadRequest, "login.html", pd)
		return
	}

	password := c.PostForm("password")
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		pd.Messages = append(pd.Messages, Message{
			Type:    "error",
			Content: loginError,
		})
		c.HTML(http.StatusBadRequest, "login.html", pd)
		return
	}

	// Generate a ulid for the current session
	sessionIdentifier := ulid.Generate()

	ses := models.Session{
		Identifier: sessionIdentifier,
	}

	// Session is valid for 1 hour
	ses.ExpiresAt = time.Now().Add(time.Hour)
	ses.UserID = user.ID

	res = controller.db.Save(&ses)
	if res.Error != nil {
		pd.Messages = append(pd.Messages, Message{
			Type:    "error",
			Content: loginError,
		})
		log.Println(res.Error)
		c.HTML(http.StatusInternalServerError, "login.html", pd)
		return
	}

	session := sessions.Default(c)
	session.Set(middleware.SessionIDKey, sessionIdentifier)

	err = session.Save()
	if err != nil {
		pd.Messages = append(pd.Messages, Message{
			Type:    "error",
			Content: loginError,
		})
		log.Println(err)
		c.HTML(http.StatusInternalServerError, "login.html", pd)
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, "/admin")
}
