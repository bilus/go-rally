package actions

func isSignupEnabled() bool {
	return getEnv("ENABLE_SIGNUPS", "true") == "true"
}
