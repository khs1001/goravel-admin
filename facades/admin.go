package facades

import (
	"log"

	"goravel/packages/admin"
	"goravel/packages/admin/contracts"
)

func Admin() contracts.Admin {
	instance, err := admin.App.Make(admin.Binding)
	if err != nil {
		log.Println(err)
		return nil
	}

	return instance.(contracts.Admin)
}
