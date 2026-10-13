package facades

import (
	"log"

	"github.com/khs1001/goravel-admin/contracts"

	admin "github.com/khs1001/goravel-admin"
)

// Admin resolves the package's main service from the container.
//
// It is the host-facing entry point: host code calls facades.Admin().Name()
// and never touches the concrete admin.Admin type.
func Admin() contracts.Admin {
	instance, err := admin.App.Make(admin.Binding)
	if err != nil {
		log.Println(err)

		return nil
	}

	return instance.(contracts.Admin)
}
