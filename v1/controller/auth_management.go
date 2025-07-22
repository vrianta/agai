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
	return controller.session != nil
}

/*
 * Login the user to the system
 */
func (controller *Context) Login() bool {

	if controller.session != nil {
		controller.session.IsAuthenticated = true
		return true // already logged in
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
	// log.Info("User logged in successfully")

	return true
}

func (_c *Context) Logout() {

	if _c.session == nil {
		return // no session present no need to logout
	}

	_c.w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
	_c.w.Header().Set("Pragma", "no-cache")
	_c.w.Header().Set("Expires", "0")
	_c.session.IsAuthenticated = false

	if _c.session != nil {
		session.RemoveSession(&_c.session.ID)
	}

	_c.session = nil // Clear the session after logout
}
