package controller

import (
	"github.com/monstorak/monstorak/pkg/controller/storagealerts"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, storagealerts.Add)
}
