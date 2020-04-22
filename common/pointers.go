package common

// SafeStr safely dereferences a string pointer and returns its value if it
// is non-nil, or the empty string if it is nil.
func SafeStr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// SPtr returns a pointer to the string value s
func SPtr(s string) *string {
	return &s
}
