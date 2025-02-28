package util

import "net"

// IsValidIP checks if the provided string is a valid IP address.
func IsValidIP(ip string) bool {
	return net.ParseIP(ip) != nil
}
