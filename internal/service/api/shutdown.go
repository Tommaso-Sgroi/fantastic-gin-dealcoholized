package api

// Close should close everything opened in the lifecycle of the `_router`; for example, background goroutines.
func (rt *router) Close() error {
	return nil
}
