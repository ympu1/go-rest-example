package main

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

type handlers struct{}

func (handler *handlers) getAllUsers(c echo.Context) error {
	users, err := getAllUsers()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "error getting users")
	}

	return c.JSON(http.StatusOK, users)
}

func (handler *handlers) getUser(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "bad user id")
	}

	user, err := getUserByID(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "error getting the user")
	}

	if user.ID == -1 {
		return echo.NewHTTPError(http.StatusNotFound, "user not found")
	}

	return c.JSON(http.StatusOK, user)
}

func (handler *handlers) addUser(c echo.Context) error {
	user := new(User)
	err := c.Bind(user)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "bad request")
	}

	validationErrorMessage := user.validate()
	if validationErrorMessage != "" {
		return echo.NewHTTPError(http.StatusBadRequest, validationErrorMessage)
	}

	err = user.addToDataStore()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "error saving the user")
	}

	c.Response().Header().Set("Location", "/api/users/"+strconv.Itoa(user.ID))
	return c.NoContent(http.StatusCreated)
}

func (handler *handlers) editUser(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "bad request")
	}

	user, err := getUserByID(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "error getting the user")
	}

	if user.ID == -1 {
		return echo.NewHTTPError(http.StatusNotFound, "user not found")
	}

	userOldID := user.ID

	err = c.Bind(&user)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "bad request")
	}

	if userOldID != user.ID {
		return echo.NewHTTPError(http.StatusBadRequest, "user id changing is not allowed")
	}

	validationErrorMessage := user.validate()
	if validationErrorMessage != "" {
		return echo.NewHTTPError(http.StatusBadRequest, validationErrorMessage)
	}

	err = user.updateToDataStore()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "error updating the user")
	}

	return c.JSON(http.StatusOK, user)
}

func (handler *handlers) deleteUser(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "bad user id")
	}

	user, err := getUserByID(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "error getting the user")
	}

	if user.ID == -1 {
		return echo.NewHTTPError(http.StatusNotFound, "user not found")
	}

	err = user.deleteFromDataStore()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "error deleting the user")
	}

	return c.NoContent(http.StatusNoContent)
}
