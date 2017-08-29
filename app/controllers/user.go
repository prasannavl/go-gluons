package controllers

import (
	"pvl/apicore/appcontext"
)

type User struct {
}

type UserController struct {
}

func NewUserController() *UserController {
	return &UserController{}
}

func NewUserService(services appcontext.Services) {

}
