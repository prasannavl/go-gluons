package controllers

import (
	"nextfirst/nf-core/appcontext"
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
