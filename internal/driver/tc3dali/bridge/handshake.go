package bridge

import (
	"github.com/vanti-dev/twincat3-ads-go/pkg/ads/types"
	"github.com/vanti-dev/twincat3-ads-go/pkg/device"
)

func handshake(fb device.Variable) (int, error) {
	method, err := fb.MethodByName("RpcHandshake",
		[]types.DataType{},
		[]types.DataType{types.Int},
	)
	if err != nil {
		return 0, err
	}

	var protocolVersion int
	err = method.Call(nil, []interface{}{&protocolVersion})
	if err != nil {
		return 0, err
	}

	return protocolVersion, nil
}
