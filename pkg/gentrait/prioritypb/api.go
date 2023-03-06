package prioritypb

import (
	"github.com/vanti-dev/sc-bos/pkg/gen"
	"github.com/vanti-dev/sc-bos/pkg/node"
)

func AddSupport(n node.Supporter) {
	{
		r := gen.NewPriorityApiRouter()
		n.Support(node.Routing(r), node.Clients(gen.WrapPriorityApi(r)))
	}
}
