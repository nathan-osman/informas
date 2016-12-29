package server

const (
	// Indicate whether installation was completed
	configInstalled = "installed"

	// Secret key used for encrypting sessions
	configSecretKey = "secret_key"

	// Title shown in the <title> for each page
	configSiteTitle = "site_title"
)

const (
	// Currently logged in user
	contextCurrentUser = "current_user"
)

const (
	// Name of session (prevents conflicts with other apps on localhost)
	sessionName = "informas"

	// ID of currently logged in user
	sessionUserID = "user_id"
)
