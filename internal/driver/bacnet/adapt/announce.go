package adapt

import "github.com/vanti-dev/bsp-ew/internal/node"

// SelfAnnouncer is a complement to node.Announcer allowing a type to announce itself.
type SelfAnnouncer interface {
	AnnounceSelf(a node.Announcer) node.Undo
}
