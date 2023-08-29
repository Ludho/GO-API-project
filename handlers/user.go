package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"trades/models"
	"trades/services"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	service services.UserService
}

func NewUserHandler(userService services.UserService) *UserHandler {
	return &UserHandler{
		service: userService,
	}
}

func (h *UserHandler) Register(c echo.Context) error {

	var user *models.User
	payload, _ := ioutil.ReadAll(c.Request().Body)
	err := json.Unmarshal(payload, &user)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"message": fmt.Sprintf("Couldn't parse body: %s", err.Error())})
	}
	if user.Password == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"message": "missing password"})
	}
	if user.Username == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"message": "missing username"})
	}
	user.Id = uuid.New().String()

	err = h.service.Register(c.Request().Context(), user)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"message": fmt.Sprintf("user could not be created: %v", err)})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"message": "user successfully created"})
}

func (h *UserHandler) Login(c echo.Context) error {

	var user *models.User
	payload, _ := ioutil.ReadAll(c.Request().Body)
	err := json.Unmarshal(payload, &user)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"message": fmt.Sprintf("Couldn't parse body: %s", err.Error())})
	}
	if user.Password == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"message": "missing password"})
	}
	if user.Username == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"message": "missing username"})
	}
	token := ""
	token, err = h.service.Login(c.Request().Context(), user)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"message": fmt.Sprintf("could not log user: %v", err)})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{"token": token})
}
