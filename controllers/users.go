package controllers

import (
	"gocommerce/core"
	"gocommerce/helper"
	"gocommerce/models"

	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/thedevsaddam/govalidator"
)

type (
	CreateUserRequest struct {
		Name     string `json:"name"`
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
		Address  string `json:"address"`
	}
)

// UserList Get list of users
// @Summary Get list of users
// @Description Get list of users
// @Tags users
// @ID users-user-list
// @Accept mpfd
// @Produce plain
// @Param page query integer false "page number" default(1)
// @Param pageSize query integer false "number of User in single page" default(10)
// @Success 200 {object} []models.User
// @Router /v1/users [get]
func UserList(c echo.Context) error {
	defer c.Request().Body.Close()

	users := models.User{}
	rows, err := strconv.Atoi(c.QueryParam("pageSize"))
	page, err := strconv.Atoi(c.QueryParam("page"))
	orderby := "id"
	sort := "DESC"

	var filter struct{}
	result, err := users.PagedFilterSearch(page, rows, orderby, sort, &filter)

	if err != nil {
		return helper.Response(http.StatusInternalServerError, err, "query result error")
	}

	response := helper.HttpResponseData{
		Status:  http.StatusOK,
		Message: "Success",
		Data:    result.Data,
	}

	return c.JSON(http.StatusOK, response)
}

func UserCreate(c echo.Context) error {
	defer c.Request().Body.Close()

	payloadRules := govalidator.MapData{
		"name":     []string{"required"},
		"username": []string{"required"},
		"email":    []string{"required"},
		"password": []string{"required"},
		"address":  []string{"required"},
	}

	validate := helper.ValidateRequestFormData(c, payloadRules)
	if validate != nil {
		return helper.Response(http.StatusUnprocessableEntity, validate, "Validation error")
	}

	users := models.User{
		Name:     c.FormValue("name"),
		Username: c.FormValue("username"),
		Email:    c.FormValue("email"),
		Password: c.FormValue("password"),
		Address:  c.FormValue("address"),
	}

	tx := core.App.DB.Begin()

	if err := tx.Create(&users).Error; err != nil {
		tx.Rollback()
		return helper.Response(http.StatusUnprocessableEntity, err, "Error while saving the User data")
	}

	tx.Commit()

	response := helper.HttpResponseData{
		Status:  http.StatusCreated,
		Message: "Success",
		Data:    users,
	}

	return c.JSON(http.StatusCreated, response)
}

func UserUpdate(c echo.Context) error {
	defer c.Request().Body.Close()

	userID, _ := strconv.Atoi(c.Param("id"))
	users := models.User{}
	if err := users.FindbyID(userID); err != nil {
		return helper.Response(http.StatusNotFound, err, "User not found")
	}

	payloadRules := govalidator.MapData{
		"name":     []string{"required"},
		"username": []string{"required"},
		"email":    []string{"required"},
		"password": []string{"required"},
		"address":  []string{"required"},
	}
	requestBody := CreateUserRequest{}
	validate := helper.ValidateRequestPayload(c, payloadRules, &requestBody)
	if validate != nil {
		return helper.Response(http.StatusUnprocessableEntity, validate, "Validation error")
	}

	users.Name = requestBody.Name
	users.Username = requestBody.Username
	users.Email = requestBody.Email
	users.Password = requestBody.Password
	users.Address = requestBody.Address

	if err := users.Save(); err != nil {
		return helper.Response(http.StatusUnprocessableEntity, err, "Error while updating User data")
	}

	response := helper.HttpResponseData{
		Status:  http.StatusOK,
		Message: "Success",
		Data:    users,
	}

	return c.JSON(http.StatusOK, response)
}

func UserDelete(c echo.Context) error {
	defer c.Request().Body.Close()

	userID, _ := strconv.Atoi(c.Param("id"))
	users := models.User{}
	if err := users.FindbyID(userID); err != nil {
		return helper.Response(http.StatusNotFound, err, "Comment not found")
	}

	if err := users.Delete(); err != nil {
		return helper.Response(http.StatusUnprocessableEntity, err, "Delete Comment failed")
	}

	response := helper.HttpResponseData{
		Status:  http.StatusOK,
		Message: "Success",
		Data:    users,
	}

	return c.JSON(http.StatusOK, response)
}
