package user

import (
	"apiservice/handler"
	"apiservice/model"
	"apiservice/pkg/auth"
	"apiservice/pkg/errno"
	"apiservice/pkg/token"

	"github.com/gin-gonic/gin"
)

// @Summary Login generates the authentication token
// @Produce json
// @Param username body string true "Username"
// @Param password body string true "Password"
// @Success 200 {string} jsone "{"code":0,"message":"OK","data":{"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE1MjgwMTY5MjIsImlkIjowLCJuYmYiOjE1MjgwMTY5MjIsInVzZXJuYW1lIjoiYWRtaW4ifQ.LjxrK9DuAwAzUD8-9v43NzWBN7HXsSLfebw92DKd1JQ"}}"
// @Router /login [post]
func Login(c *gin.Context) {
	// Binding the data with user struct.
	var u model.UserModel
	if err := c.Bind(&u); err != nil {
		handler.SendResponse(c, errno.ErrBind, nil)
		return
	}

	// Get the user informaiton by the login username.
	d, err := model.GetUser(u.Username)
	if err != nil {
		handler.SendResponse(c, errno.ErrUserNotFound, nil)
		return
	}

	// Compare the login password with the user password.
	if err := auth.Compare(d.Password, u.Password); err != nil {
		handler.SendResponse(c, errno.ErrPasswordIncorrect, nil)
		return
	}

	// Sign the json web token.
	t, err := token.Sign(token.Context{ID: d.ID, Username: d.Username}, "")
	if err != nil {
		handler.SendResponse(c, errno.ErrToken, nil)
		return
	}

	handler.SendResponse(c, nil, model.Token{Token: t})
}
