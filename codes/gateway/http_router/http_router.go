package httprouter

type HTTPRouter interface {
	UpdateRoute(param map[string]interface{}) error
	// GetRoutes() error
	// DeleteRoute(id string) error
}
