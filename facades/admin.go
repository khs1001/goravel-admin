package facades

import (
	"log"

	"github.com/khs1001/goravel-admin"
	"github.com/khs1001/goravel-admin/contracts"
)

func Admin() contracts.Admin {
	instance, err := admin.App.Make(admin.Binding)
	if err != nil {
		log.Println(err)
		return nil
	}

	return instance.(contracts.Admin)
}
