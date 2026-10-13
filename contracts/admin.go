package contracts

// Admin is the contract implemented by the package's main service and exposed
// to the host through the facades.Admin facade.
//
// Consumers should depend on this interface instead of the concrete
// admin.Admin type.
type Admin interface {
	// Name returns the configured admin application name.
	Name() string
	// RoutePrefix returns the configured API route prefix (without a leading
	// slash), e.g. "admin-api".
	RoutePrefix() string
	// Path returns the configured SPA entry path, e.g. "admin".
	Path() string
	// Assets returns the configured static asset base path, e.g. "admin-assets".
	Assets() string
}
