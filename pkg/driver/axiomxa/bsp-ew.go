package axiomxa

import (
	"github.com/vanti-dev/sc-bos/pkg/driver/axiomxa/mps"
)

const (
	KeyAccessGranted  = "AG"
	KeyAccessDenied   = "AD"
	KeyDoorHeldOpen   = "DHO"
	KeyDoorNotOpen    = "DNO"
	KeyForcedEntry    = "FE"
	KeyTailgate       = "TG"
	KeyTamper         = "TAMP"
	KeySecure         = "SEC"
	KeyNetworkOffline = "NOFF"
	KeyNetworkOnline  = "NON"
)

var EWAxiomPatterns = map[string]mps.Pattern{
	KeyAccessGranted:  mps.NewPattern(mps.Timestamp, mps.EventID, mps.EventDesc, mps.NetworkDesc, mps.DeviceDesc, mps.CardID, mps.CardNumber, mps.CardholderDesc),
	KeyAccessDenied:   mps.NewPattern(mps.Timestamp, mps.EventID, mps.EventDesc, mps.NetworkDesc, mps.DeviceDesc, mps.CardID, mps.CardNumber, mps.CardholderDesc),
	KeyDoorHeldOpen:   mps.NewPattern(mps.Timestamp, mps.EventID, mps.EventDesc, mps.NetworkDesc, mps.DeviceDesc),
	KeyDoorNotOpen:    mps.NewPattern(mps.Timestamp, mps.EventID, mps.EventDesc, mps.NetworkDesc, mps.DeviceDesc),
	KeyForcedEntry:    mps.NewPattern(mps.Timestamp, mps.EventID, mps.EventDesc, mps.NetworkDesc, mps.DeviceDesc),
	KeyTailgate:       mps.NewPattern(mps.Timestamp, mps.EventID, mps.EventDesc, mps.NetworkDesc, mps.DeviceDesc),
	KeyTamper:         mps.NewPattern(mps.Timestamp, mps.EventID, mps.EventDesc, mps.NetworkDesc, mps.DeviceDesc),
	KeySecure:         mps.NewPattern(mps.Timestamp, mps.EventID, mps.EventDesc, mps.NetworkDesc, mps.DeviceDesc),
	KeyNetworkOffline: mps.NewPattern(mps.Timestamp, mps.EventID, mps.EventDesc, mps.NetworkDesc),
	KeyNetworkOnline:  mps.NewPattern(mps.Timestamp, mps.EventID, mps.EventDesc, mps.NetworkDesc),
}
