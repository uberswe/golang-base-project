package middleware

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/uberswe/golang-base-project/models"
	"gorm.io/gorm"
	"log"
)

// SessionIDKey is the key used to set and get the session id in the context of the current request
const SessionIDKey = "SessionID"

// Session middleware checks for an active session and sets the UserIDKey to the context of the current request if found
func Session(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		sessionIdentifierInterface := session.Get(SessionIDKey)

		if sessionIdentifier, ok := sessionIdentifierInterface.(string); ok {
			ses := models.Session{
				Identifier: sessionIdentifier,
			}
			res := db.Where(&ses).First(&ses)
			if res.Error == nil && !ses.HasExpired() {
				c.Set(UserIDKey, ses.UserID)
			} else {
				log.Println(res.Error)
			}
		}
		c.Next()
	}
}
