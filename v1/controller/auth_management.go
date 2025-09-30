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
		if sessionID, err := session.GetSessionID(controller.R); err == nil && sessionID != "" { // it means the user had the session ID
			sess, _ := session.Get(&sessionID, controller.W, controller.R)
			controller.session = sess
			log.Debug("Found Session with session ID %s", sessionID)
			return true
		} else {
			log.Debug("Checking is LoggedIn for Session ID %s but the session not found", sessionID)
			return false
		}
	}
	return controller.session != nil

}

/*
 * Login the user to the system
 */
func (controller *Context) Login() bool {

	if controller.session == nil {
		sessionID, err := session.GetSessionID(controller.R)

		if err == nil && sessionID != "" { // it means the user had the session ID
			sess, _ := session.Get(&sessionID, controller.W, controller.R)
			controller.session = sess
		}
	}
	if controller.session != nil {
		return true // already logged in
	}
	var err error

	controller.session, err = session.New(controller.W, controller.R)
	if err != nil {
		log.Error("Failed to create the login session: %s", err.Error())
		return false
	} else {
		log.WriteLogf("Successfully Loggedin with session %s", controller.session.ID)
	}

	controller.session.Login(controller.W, controller.R)
	// asser
	// controller.session.Controller[controller.R.URL.Path] = controller // storing the controller objects in the session

	return true
}

func (controller *Context) Logout() {

	if controller.session == nil {
		sessionID, err := session.GetSessionID(controller.R)

		if err == nil && sessionID != "" { // it means the user had the session ID
			sess, _ := session.Get(&sessionID, controller.W, controller.R)
			controller.session = sess
		} else {
			log.Debug("Session not found")
		}
	}

	if controller.session != nil {
		controller.session.Logout(controller.W, controller.R)
		session.RemoveSession(&controller.session.ID)
		controller.session = nil // Clear the session after logout
	} else {
		log.Debug("Session is nill and no need to logout")
	}

	controller.W.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
	controller.W.Header().Set("Pragma", "no-cache")
	controller.W.Header().Set("Expires", "0")
}
