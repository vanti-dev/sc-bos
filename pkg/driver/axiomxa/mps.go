package axiomxa

import (
	"context"
	"net"

	"github.com/olebedev/emitter"
	"go.uber.org/zap"

	"github.com/vanti-dev/sc-bos/pkg/driver/axiomxa/config"
	"github.com/vanti-dev/sc-bos/pkg/driver/axiomxa/mps"
)

// setupMessagePortServer listens on cfg.MessagePorts.Bind for message ports, sending the parsed Fields to bus.
func (d *Driver) setupMessagePortServer(ctx context.Context, cfg config.Root, bus *emitter.Emitter) error {
	var lc net.ListenConfig
	mpLis, err := lc.Listen(ctx, "tcp", cfg.MessagePorts.Bind)
	if err != nil {
		return err
	}
	go func() {
		<-ctx.Done()
		mpLis.Close()
	}()

	onMessageMap := make(map[string]mps.OnMessageFunc, len(EWAxiomPatterns))
	for k, p := range EWAxiomPatterns {
		onMessageMap[k+","] = onMessageStripParseAndEmit(bus, k, p)
	}
	server := mps.NewServer(mps.MapPrefix(onMessageMap, func(data []byte) {
		d.logger.Warn("Unexpected message port message", zap.String("data", string(data)))
	}))

	go server.Serve(mpLis)

	return nil
}

// onMessageStripParseAndEmit returns an OnMessageFunc that removes "topic," from the start of the message then parses
// it as the given pattern and finally publishes the fields to bus using topic.
func onMessageStripParseAndEmit(bus *emitter.Emitter, topic string, pattern mps.Pattern) mps.OnMessageFunc {
	return func(data []byte) {
		fields := mps.Fields{}
		err := pattern.Unmarshal(data[len(topic)+1:], &fields)
		if err != nil {
			// log it or something
			return
		}
		bus.Emit(topic, fields)
	}
}
