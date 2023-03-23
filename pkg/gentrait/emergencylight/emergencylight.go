package emergencylight

import (
	"github.com/smart-core-os/sc-golang/pkg/trait"
)

// Note there's no real trait for this (yet) but the devices that implement the DaliApi might advertise this trait.

const TraitName trait.Name = "smartcore.bos.EmergencyLight"
