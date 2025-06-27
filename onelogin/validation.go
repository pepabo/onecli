package onelogin

// validateString checks if the provided value is a string.
func validateString(val interface{}) bool {
	switch v := val.(type) {
	case string:
		return true
	case *string:
		return v != nil
	default:
		return false
	}
}
