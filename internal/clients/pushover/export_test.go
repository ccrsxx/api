package pushover

// SetJSONMarshal overrides the package-level jsonMarshal function for testing.
// It returns a restore function that should be deferred.
func SetJSONMarshal(fn func(v any) ([]byte, error)) func() {
	original := jsonMarshal
	jsonMarshal = fn

	return func() {
		jsonMarshal = original
	}
}
