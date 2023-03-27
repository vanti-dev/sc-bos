// Package historypb adds types useful for working with the generated FooHistory services.
// This package is temporary as each traits history service will eventually move to be next to the other non-generated trait types.
package historypb

import (
	"github.com/vanti-dev/sc-bos/pkg/node"
)

// Deprecated: manually add support or use alltraits.AddSupport
func AddSupport(n node.Supporter) {
	// support is added via node/alltraits.AddSupport
}
