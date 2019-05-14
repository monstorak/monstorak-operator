package controller

import (
	"github.com/monstorak/monstorak/pkg/controller/storagealert"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, storagealert.Add)
}
