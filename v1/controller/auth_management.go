package controller

import (
	"github.com/vrianta/agai/v1/internal/session"
	"github.com/vrianta/agai/v1/log"
)

/*
 * This file is to manage different types of authentication previously this functionionality was with session pacakge
 * But later for better managability we moved this to controllers object
 */

/*
 * Check if the User is Logged in to the system or not
 */
func (controller *Context) IsLoggedIn() bool {
	if controller.session == nil {
		return false
	}
	return controller.session.IsAuthenticated
}

/*
 * Login the user to the system
 */
func (controller *Context) Login() bool {

	// if the controller session nil that means the user is not logged in
	if controller.session == nil {
		log.Error("Not Able to login user, session is nil")
		return false // not logged in
	}
	var err error
	// No session, create a new one
	// session.RemoveSession(&controller.session.ID)
	controller.session, err = session.New(controller.w, controller.r)
	if err != nil {
		log.Error("Failed to create the login session: %s", err.Error())
		return false
	}

	controller.session.Login(controller.w, controller.r)
	log.Info("User logged in successfully")

	return true
}

func (_c *Context) Logout() {

	_c.w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
	_c.w.Header().Set("Pragma", "no-cache")
	_c.w.Header().Set("Expires", "0")
	_c.session.IsAuthenticated = false

	if _c.session != nil {
		session.RemoveSession(&_c.session.ID)
	}
}
