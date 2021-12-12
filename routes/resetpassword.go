package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/uberswe/golang-base-project/models"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
)

type ResetPasswordPageData struct {
	PageData
	Token string
}

func (controller Controller) ResetPassword(c *gin.Context) {
	token := c.Param("token")
	pd := ResetPasswordPageData{
		PageData: PageData{
			Title:           "Reset Password",
			IsAuthenticated: isAuthenticated(c),
		},
		Token: token,
	}
	c.HTML(http.StatusOK, "resetpassword.html", pd)
}

func (controller Controller) ResetPasswordPost(c *gin.Context) {
	passwordError := "Your password must be 8 characters in length or longer"
	resetError := "Could not reset password, please try again"

	token := c.Param("token")
	pd := ResetPasswordPageData{
		PageData: PageData{
			Title:           "Reset Password",
			IsAuthenticated: isAuthenticated(c),
		},
		Token: token,
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
		Content: "Your password has successfully been reset.",
	})

	c.HTML(http.StatusOK, "resetpassword.html", pd)
}
