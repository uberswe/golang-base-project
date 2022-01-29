package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/uberswe/golang-base-project/models"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
)

// ResetPasswordPageData defines additional data needed to render the reset password page
type ResetPasswordPageData struct {
	PageData
	Token string
}

// ResetPassword renders the HTML page for resetting the users password
func (controller Controller) ResetPassword(c *gin.Context) {
	token := c.Param("token")
	pdPre := controller.DefaultPageData(c)
	pdPre.Title = pdPre.Trans("Reset Password")
	pd := ResetPasswordPageData{
		PageData: pdPre,
		Token:    token,
	}
	c.HTML(http.StatusOK, "resetpassword.html", pd)
}

// ResetPasswordPost handles post request used to reset users passwords
func (controller Controller) ResetPasswordPost(c *gin.Context) {
	pdPre := controller.DefaultPageData(c)
	passwordError := pdPre.Trans("Your password must be 8 characters in length or longer")
	resetError := pdPre.Trans("Could not reset password, please try again")

	token := c.Param("token")
	pdPre.Title = pdPre.Trans("Reset Password")
	pd := ResetPasswordPageData{
		PageData: pdPre,
		Token:    token,
	}
	password := c.PostForm("password")

	if len(password) < 8 {
		pd.Messages = append(pd.Messages, Message{
			Type:    "error",
			Content: passwordError,
		})
		c.HTML(http.StatusBadRequest, "resetpassword.html", pd)
		return
	}

	forgotPasswordToken := models.Token{
		Value: token,
		Type:  models.TokenPasswordReset,
	}

	res := controller.db.Where(&forgotPasswordToken).First(&forgotPasswordToken)
	if res.Error != nil {
		pd.Messages = append(pd.Messages, Message{
			Type:    "error",
			Content: resetError,
		})
		c.HTML(http.StatusBadRequest, "resetpassword.html", pd)
		return
	}

	if forgotPasswordToken.HasExpired() {
		pd.Messages = append(pd.Messages, Message{
			Type:    "error",
			Content: resetError,
		})
		c.HTML(http.StatusBadRequest, "resetpassword.html", pd)
		return
	}

	user := models.User{}
	user.ID = uint(forgotPasswordToken.ModelID)
	res = controller.db.Where(&user).First(&user)
	if res.Error != nil {
		pd.Messages = append(pd.Messages, Message{
			Type:    "error",
			Content: resetError,
		})
		c.HTML(http.StatusBadRequest, "resetpassword.html", pd)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		log.Println(err)
		pd.Messages = append(pd.Messages, Message{
			Type:    "error",
			Content: resetError,
		})
		c.HTML(http.StatusBadRequest, "resetpassword.html", pd)
		return
	}

	user.Password = string(hashedPassword)

	res = controller.db.Save(&user)
	if res.Error != nil {
		pd.Messages = append(pd.Messages, Message{
			Type:    "error",
			Content: resetError,
		})
		c.HTML(http.StatusBadRequest, "resetpassword.html", pd)
		return
	}

	res = controller.db.Delete(&forgotPasswordToken)
	if res.Error != nil {
		pd.Messages = append(pd.Messages, Message{
			Type:    "error",
			Content: resetError,
		})
		c.HTML(http.StatusBadRequest, "resetpassword.html", pd)
		return
	}

	pd.Messages = append(pd.Messages, Message{
		Type:    "success",
		Content: pdPre.Trans("Your password has successfully been reset."),
	})

	c.HTML(http.StatusOK, "resetpassword.html", pd)
}
