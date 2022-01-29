package routes

import (
	"fmt"
	"github.com/gin-gonic/gin"
	email2 "github.com/uberswe/golang-base-project/email"
	"github.com/uberswe/golang-base-project/models"
	"github.com/uberswe/golang-base-project/ulid"
	"gorm.io/gorm"
	"log"
	"net/http"
	"net/url"
	"path"
	"time"
)

// ForgotPassword renders the HTML page where a password request can be initiated
func (controller Controller) ForgotPassword(c *gin.Context) {
	pd := controller.DefaultPageData(c)
	pd.Title = pd.Trans("Forgot Password")
	c.HTML(http.StatusOK, "forgotpassword.html", pd)
}

// ForgotPasswordPost handles the POST request which requests a password reset and then renders the HTML page with the appropriate message
func (controller Controller) ForgotPasswordPost(c *gin.Context) {
	pd := controller.DefaultPageData(c)
	pd.Title = pd.Trans("Forgot Password")
	email := c.PostForm("email")
	user := models.User{Email: email}
	res := controller.db.Where(&user).First(&user)
	if res.Error == nil && user.ActivatedAt != nil {
		go controller.forgotPasswordEmailHandler(user.ID, email, pd.Trans)
	}

	pd.Messages = append(pd.Messages, Message{
		Type:    "success",
		Content: pd.Trans("An email with instructions describing how to reset your password has been sent."),
	})

	// We always return a positive response here to prevent user enumeration
	c.HTML(http.StatusOK, "forgotpassword.html", pd)
}

func (controller Controller) forgotPasswordEmailHandler(userID uint, email string, trans func(string) string) {
	forgotPasswordToken := models.Token{
		Value: ulid.Generate(),
		Type:  models.TokenPasswordReset,
	}

	res := controller.db.Where(&forgotPasswordToken).First(&forgotPasswordToken)
	if (res.Error != nil && res.Error != gorm.ErrRecordNotFound) || res.RowsAffected > 0 {
		// If the forgot password token already exists we try to generate it again
		controller.forgotPasswordEmailHandler(userID, email, trans)
		return
	}

	forgotPasswordToken.ModelID = int(userID)
	forgotPasswordToken.ModelType = "User"
	// The token will expire 10 minutes after it was created
	forgotPasswordToken.ExpiresAt = time.Now().Add(time.Minute * 10)

	res = controller.db.Save(&forgotPasswordToken)
	if res.Error != nil || res.RowsAffected == 0 {
		log.Println(res.Error)
		return
	}
	controller.sendForgotPasswordEmail(forgotPasswordToken.Value, email, trans)
}

func (controller Controller) sendForgotPasswordEmail(token string, email string, trans func(string) string) {
	u, err := url.Parse(controller.config.BaseURL)
	if err != nil {
		log.Println(err)
		return
	}

	u.Path = path.Join(u.Path, "/user/password/reset/", token)

	resetPasswordURL := u.String()

	emailService := email2.New(controller.config)

	emailService.Send(email, trans("Password Reset"), fmt.Sprintf(trans("Use the following link to reset your password. If this was not requested by you, please ignore this email.\n%s"), resetPasswordURL))
}
