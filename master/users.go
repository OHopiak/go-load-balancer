package master

import (
	"github.com/OHopiak/fractal-load-balancer/core"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

type (
	SignUpError string
	LoginError string

	SignUpRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	SignUpResponse struct {
		User       *core.User   `json:"user,omitempty"`
		Error      *SignUpError `json:"error,omitempty"`
	}
	LoginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	LoginResponse struct {
		User  *core.User  `json:"user,omitempty"`
		Error *LoginError `json:"error,omitempty"`
	}

	SignUpData struct {
		BaseData
		Error *SignUpError `json:"error,omitempty"`
	}
	LoginData struct {
		BaseData
		Error *LoginError `json:"error,omitempty"`
	}
)

func (m Master) SignUp(request SignUpRequest) SignUpResponse {
	hash, err := core.GetPasswordHash(request.Password)
	if err != nil {
		errorMsg := SignUpError(err.Error())
		return SignUpResponse{
			Error: &errorMsg,
		}
	}
	user := core.User{
		Username: request.Username,
	}
	m.db.First(&user)

	if user.ID != 0 {
		errorMsg := SignUpError("user with such username exists")
		return SignUpResponse{
			Error: &errorMsg,
		}
	}

	user.PasswordHash = hash
	m.db.Create(&user)

	return SignUpResponse{
		User: &user,
	}
}

func (m Master) Login(request LoginRequest) LoginResponse {
	hash, err := core.GetPasswordHash(request.Password)
	if err != nil {
		errorMsg := LoginError(err.Error())
		return LoginResponse{
			Error: &errorMsg,
		}
	}
	user := core.User{
		Username:     request.Username,
		PasswordHash: hash,
	}

	m.db.First(&user)
	if user.ID == 0 {
		errorMsg := LoginError("incorrect username or password")
		return LoginResponse{
			Error: &errorMsg,
		}
	}

	return LoginResponse{
		User: &user,
	}
}

func (m Master) signUpGetHandler(c echo.Context) error {
	userRaw := c.Get("user")
	if userRaw != nil {
		return c.Redirect(http.StatusSeeOther, "/")
	}
	data := BaseData{
		LoggedIn: false,
	}

	return c.Render(http.StatusOK, "signup.html", data)
}

func (m Master) signUpPostHandler(c echo.Context) error {
	request := new(SignUpRequest)
	err := c.Bind(request)
	if err != nil {
		return err
	}
	response := m.SignUp(*request)

	if response.Error != nil {
		data := SignUpData{
			BaseData: m.GetBaseData(c),
			Error:    response.Error,
		}
		return c.Render(http.StatusOK, "signup.html", data)
		//return c.JSON(http.StatusBadRequest, response)
	}

	err = saveUserId(response.User.ID, c)
	if err != nil {
		return err
	}
	return c.Redirect(http.StatusSeeOther, "/")
}

func (m Master) loginGetHandler(c echo.Context) error {
	userRaw := c.Get("user")
	if userRaw != nil {
		return c.Redirect(http.StatusSeeOther, "/")
	}
	data := BaseData{
		LoggedIn: false,
	}
	return c.Render(http.StatusOK, "login.html", data)
}

func (m Master) loginPostHandler(c echo.Context) error {
	request := new(LoginRequest)
	err := c.Bind(request)
	if err != nil {
		return err
	}
	response := m.Login(*request)

	if response.Error != nil {
		data := LoginData{
			BaseData: m.GetBaseData(c),
			Error:    response.Error,
		}
		return c.Render(http.StatusOK, "login.html", data)

		//return c.JSON(http.StatusBadRequest, response)
	}

	err = saveUserId(response.User.ID, c)
	if err != nil {
		return err
	}
	return c.Redirect(http.StatusSeeOther, "/")
}

func (m Master) logoutGetHandler(c echo.Context) error {
	err := clearUserID(c)
	if err != nil {
		return err
	}
	return c.Redirect(http.StatusSeeOther, "/")
}

func saveUserId(userId uint, c echo.Context) error {
	sess, err := GetSession(c)
	if err != nil {
		return err
	}
	sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   int(7 * 24 * time.Hour / time.Second),
		HttpOnly: true,
	}
	sess.Values["user_id"] = userId
	return sess.Save(c.Request(), c.Response())
}

func clearUserID(c echo.Context) error {
	sess, err := GetSession(c)
	if err != nil {
		return err
	}
	delete(sess.Values,"user_id")
	return sess.Save(c.Request(), c.Response())
}
