package network

// Function to determine http scheme being used
func HTTPEnabled(httpsEnabled bool) string {
	if httpsEnabled {
		return "https"
	}
	return "http"
}

// Function to determine ws scheme being used
func WSEnabled(httpsEnabled bool) string {
	if httpsEnabled {
		return "wss"
	}
	return "ws"
}
