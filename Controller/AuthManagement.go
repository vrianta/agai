package Controller

import (
	"net/http"

	"github.com/vrianta/Server/Log"
	"github.com/vrianta/Server/Session"
	"github.com/vrianta/Server/Utils"
)

/*
 * This file is to manage different types of authentication previously this functionionality was with session pacakge
 * But later for better managability we moved this to controllers object
 */

/*
 * Check if the User is Logged in to the system or not
 */
func (__c *Struct) IsLoggedIn() bool {
	return __c.session != nil && __c.session.IsLoggedIn()
}

/*
 * Login the user to the system
 */
func (__c *Struct) Login() bool {

	// No session, create a new one
	__c.session = Session.New()
	sessionID, err := Utils.GenerateSessionID()
	if err != nil {
		Log.WriteLog("Error generating session ID: " + err.Error())
		return false
	}

	if __c.session.StartSession(&sessionID, __c.w, __c.r) == nil {
		http.Error(__c.w, "Server Error * Failed to Create the Session for the user", http.StatusInternalServerError)
		return false
	}
	__c.session.Login(__c.w, __c.r)
	Session.Store(&sessionID, __c.session)

	return true

}

func (_c *Struct) Logout() {

	_c.w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
	_c.w.Header().Set("Pragma", "no-cache")
	_c.w.Header().Set("Expires", "0")
	if _c.session != nil {
		Session.RemoveSession(&_c.session.ID)
	}
}
