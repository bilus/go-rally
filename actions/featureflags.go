package actions

// isSignupEnabled, if true, makes it possible for users to sign up.
func isSignupEnabled() bool {
	return getEnv("ENABLE_SIGNUPS", "true") == "true"
}

// isLoginFormEnabled, if true, makes the login form visible. Otherwise only Google sign in is supported.
func isLoginFormEnabled() bool {
	return getEnv("ENABLE_LOGIN_FORM", "true") == "true"
}
