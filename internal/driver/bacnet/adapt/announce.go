package adapt

import "github.com/vanti-dev/bsp-ew/internal/node"

type SelfAnnouncer interface {
	AnnounceSelf(a node.Announcer) node.Undo
}
