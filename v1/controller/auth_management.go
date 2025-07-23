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
		// getting the session ID from the cookies
		// the session not present then the sessionID will be nil
		sessionID, err := session.GetSessionID(controller.r)

		if err == nil && sessionID != "" { // it means the user had the session ID
			sess, _ := session.Get(&sessionID, controller.w, controller.r)
			controller.session = sess
		}
	}
	return controller.session != nil
}

/*
 * Login the user to the system
 */
func (controller *Context) Login() bool {

	if controller.session == nil {
		sessionID, err := session.GetSessionID(controller.r)

		if err == nil && sessionID != "" { // it means the user had the session ID
			sess, _ := session.Get(&sessionID, controller.w, controller.r)
			controller.session = sess
		}
	}
	if controller.session != nil {
		return true // already logged in
	}
	var err error

	controller.session, err = session.New(controller.w, controller.r)
	if err != nil {
		log.Error("Failed to create the login session: %s", err.Error())
		return false
	}

	controller.session.Login(controller.w, controller.r)

	return true
}

func (controller *Context) Logout() {

	if controller.session == nil {
		sessionID, err := session.GetSessionID(controller.r)

		if err == nil && sessionID != "" { // it means the user had the session ID
			sess, _ := session.Get(&sessionID, controller.w, controller.r)
			controller.session = sess
		}
	}

	if controller.session != nil {
		session.RemoveSession(&controller.session.ID)
		controller.session = nil // Clear the session after logout
	}

	controller.w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
	controller.w.Header().Set("Pragma", "no-cache")
	controller.w.Header().Set("Expires", "0")

}
