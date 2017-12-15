// Package util contains useful shared functions
package util

// Check for 'boring' errors, and panic if error is not nil.
func Check(err error) {
	if err != nil {
		panic(err)
	}
}
