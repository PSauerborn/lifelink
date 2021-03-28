package users

import (

)

// function used to determine if a given user has 
// admin access
func isAdminUser(user string) (bool, error) {
	// get user details from graph
	details, err := persistence.GetUserDetails(user)
	if err != nil {
		switch err {
		case ErrUserDoesNotExist:
			return false, nil
		default:
			return false, err
		}
	}
	return details.Admin, nil
}