package bridge

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"sync"
	"sync/atomic"
	"time"

	"github.com/vanti-dev/bsp-ew/internal/driver/tc3dali/dali"
	"github.com/vanti-dev/twincat3-ads-go/pkg/ads"
	"github.com/vanti-dev/twincat3-ads-go/pkg/ads/types"
	"github.com/vanti-dev/twincat3-ads-go/pkg/device"
	"go.uber.org/multierr"
	"go.uber.org/zap"
)

const expectedProtocolVersion = 2

var subCounter int32 = 0

var (
	ErrInvalid              = errors.New("invalid response")
	ErrMaxListenersExceeded = errors.New("exceeded maximum number of listeners")
	ErrMalformed            = errors.New("malformed response")
)

type Config struct {
	Device                  device.Device
	Logger                  *zap.Logger
	BridgeFBName            string
	ResponseMailboxName     string
	NotificationMailboxName string
	UsePolling              bool
}

// Connect establishes a connection to the DALI bridge on the PLC Device.
func (c *Config) Connect() (dali.Dali, error) {
	bridgeFb, err := c.Device.VariableByName(c.BridgeFBName)
	if err != nil {
		return nil, err
	}

	responseMailbox, err := c.Device.VariableByName(c.ResponseMailboxName)
	if err != nil {
		return nil, err
	}

	notifMailbox, err := c.Device.VariableByName(c.NotificationMailboxName)
	if err != nil {
		return nil, err
	}

	return connectBridge(c.Device, bridgeFb, responseMailbox, notifMailbox, c.Logger, c.UsePolling)
}

func connectBridge(dev device.Device, bridgeFb device.Variable, resMailbox device.Variable, notifMailbox device.Variable, logger *zap.Logger, poll bool) (dali.Dali, error) {

	// check that the bridge implements the version of the protocl we expect
	protocolVersion, err := handshake(bridgeFb)
	if err != nil {
		return nil, err
	}
	if protocolVersion != expectedProtocolVersion {
		return nil, fmt.Errorf("unexpected bridge protocol version %d", protocolVersion)
	}

	// resolve data types needed for RPC
	requestType, err := dev.ResolveType("ST_DALIBridgeRequest")
	if err != nil {
		return nil, err
	}
	inputEventParamsType, err := dev.ResolveType("ST_DALIBridgeInputEventParams")
	if err != nil {
		return nil, err
	}

	// resolve RPC methods
	reset, err := bridgeFb.MethodByName("RpcReset", nil, nil)
	if err != nil {
		return nil, err
	}
	err = reset.Call(nil, nil) // reset bridge to default state
	if err != nil {
		return nil, fmt.Errorf("RpcReset: %w", err)
	}

	request, err := bridgeFb.MethodByName("RpcRequest",
		[]types.DataType{requestType},
		[]types.DataType{types.Bool, types.ULInt},
	)
	if err != nil {
		return nil, err
	}

	addInputEventListener, err := bridgeFb.MethodByName("RpcAddInputEventListener",
		[]types.DataType{inputEventParamsType},
		[]types.DataType{types.Bool},
	)
	if err != nil {
		return nil, err
	}

	ackNotification, err := bridgeFb.MethodByName("RpcAckNotification",
		[]types.DataType{types.ULInt},
		[]types.DataType{types.Bool},
	)
	if err != nil {
		return nil, err
	}

	ctx, done := context.WithCancel(context.Background())
	instance := &daliBridge{
		bridge:                   bridgeFb,
		resMailbox:               resMailbox,
		notifMailbox:             notifMailbox,
		rpcRequest:               request,
		rpcAddInputEventListener: addInputEventListener,
		rcpAckNotification:       ackNotification,
		requestType:              requestType,
		requests:                 make(chan bridgeRequest),
		ctx:                      ctx,
		done:                     done,
		logger:                   logger,
	}

	err = instance.startRequestResponseWorker()
	if err != nil {
		return nil, err
	}
	var notifications <-chan notification
	if poll {
		notifications, err = instance.pollNotifications()
	} else {
		notifications, err = instance.subscribeNotifications()
	}
	if err != nil {
		return nil, err
	}

	instance.notifications = notifications

	return instance, nil
}

type daliBridge struct {
	bridge                   device.Variable
	resMailbox               device.Variable
	notifMailbox             device.Variable
	rpcRequest               device.Method
	rpcAddInputEventListener device.Method
	rcpAckNotification       device.Method
	requestType              types.DataType

	commandL sync.Mutex // used to prevent simultaneous command execution (the bridge does not support it)

	// cleanup handling
	ctx  context.Context
	done context.CancelFunc

	requests      chan bridgeRequest
	notifications <-chan notification

	stateL             sync.RWMutex // protects the below state variables
	closers            []io.Closer  // things to close when the daliBridge closes
	inputEventHandlers []dali.InputEventHandler
	firstHandler       sync.Once // triggered when at least 1 hander has been added
	logger             *zap.Logger
}

func (b *daliBridge) EnableInputEventListener(params dali.InputEventParameters) error {
	var ok bool
	err := b.rpcAddInputEventListener.Call(
		[]interface{}{params},
		[]interface{}{&ok},
	)
	if err != nil {
		return err
	}
	if !ok {
		return ErrMaxListenersExceeded
	}
	return nil
}

func (b *daliBridge) OnInputEvent(handler dali.InputEventHandler) error {
	b.stateL.Lock()
	defer b.stateL.Unlock()

	b.inputEventHandlers = append(b.inputEventHandlers, handler)
	b.firstHandler.Do(b.startNotificationWorker)

	return nil
}

func (b *daliBridge) Close() error {
	b.stateL.Lock()
	defer b.stateL.Unlock()

	var err error
	for _, closer := range b.closers {
		err = multierr.Append(err, closer.Close())
	}

	b.done()

	return err
}

func (b *daliBridge) ExecuteCommand(ctx context.Context, request dali.Request) (data uint32, err error) {
	b.commandL.Lock()
	defer b.commandL.Unlock()

	var (
		ok       bool
		sequence uint64
	)
	err = b.rpcRequest.Call(
		[]interface{}{request},
		[]interface{}{&ok, &sequence},
	)
	if err != nil {
		return 0, err
	}
	if !ok {
		return 0, errors.New("bridge busy")
	}

	resCh := make(chan bridgeResponse, 1)
	b.requests <- bridgeRequest{
		sequence:   sequence,
		responseCh: resCh,
		ctx:        ctx,
	}
	select {
	case <-ctx.Done():
		return 0, ctx.Err()
	case res := <-resCh:
		return res.Data, res.AsError()
	}
}

func (b *daliBridge) startRequestResponseWorker() error {
	// subscribe to response mailbox changes
	incoming := make(chan bridgeResponse)
	callback := func(addr ads.Addr, header ads.NotificationHeader, data bridgeResponse) {
		incoming <- data
	}
	sub, err := b.resMailbox.Subscribe(ads.NotificationAttrib{
		TransMode: ads.TransServerOnCha,
		CycleTime: 100 * time.Millisecond,
	}, callback)
	if err != nil {
		return err
	}
	b.logger.Debug("subscribed to a value over ADS", zap.Int32("total", atomic.AddInt32(&subCounter, 1)))

	b.closers = append(b.closers, sub)

	go func() {
		var activeRequest *bridgeRequest
		var lastResponse *bridgeResponse
		for {
			select {
			case <-b.ctx.Done():
				return

			case req := <-b.requests:
				if activeRequest != nil {
					if activeRequest.ctx.Err() != nil {
						// if the context has already been cancelled, then we should discard this request
						activeRequest = nil
					} else {
						// got a new request but a request is already in progress - should be impossible
						// due to the locking
						panic("got a new request while busy")
					}
				}

				if lastResponse != nil && lastResponse.Sequence == req.sequence {
					// in the event of a race condition, we may have already received the response
					// to this request
					req.responseCh <- *lastResponse
					break
				}

				activeRequest = &req

			case res := <-incoming:
				if !res.Valid {
					break
				}
				lastResponse = &res
				if activeRequest == nil {
					log.Printf("got a response without making a request: %+v", res)
					break
				}

				if res.Sequence != activeRequest.sequence {
					log.Printf("got a response to an unknown request: %+v", res)
					break
				}

				activeRequest.responseCh <- res
				activeRequest = nil
			}
		}
	}()

	return nil
}

func (b *daliBridge) pollNotifications() (<-chan notification, error) {
	ch := make(chan notification, 16)
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		ticker := time.NewTicker(100 * time.Millisecond)
		defer ticker.Stop()
		var sequence uint64 = 0

		for {
			select {
			case <-ctx.Done():
				return

			case <-ticker.C:
				var notif notification
				err := b.notifMailbox.Read(&notif)
				if err != nil {
					b.logger.Error("can't read notification mailbox", zap.Error(err))
					continue
				}

				if notif.Sequence == sequence {
					// already seen this notification
					continue
				}
				sequence = notif.Sequence

				select {
				case <-time.After(5 * time.Second):
					b.logger.Warn("notification buffer full for at least 5 seconds!")
				case ch <- notif:
				}

				var (
					ok bool
				)
				err = b.rcpAckNotification.Call(
					[]interface{}{sequence},
					[]interface{}{&ok},
				)
				if err != nil {
					log.Println("failed to ack notification")
				}
			}
		}
	}()
	b.closers = append(b.closers, contextToCloser(cancel))

	return ch, nil
}

func (b *daliBridge) subscribeNotifications() (<-chan notification, error) {
	var (
		sequence uint64 = 0
		m        sync.Mutex
	)

	notifications := make(chan notification, 16)

	callback := func(addr ads.Addr, header ads.NotificationHeader, data notification) {
		m.Lock()
		defer m.Unlock()
		if data.Sequence == sequence {
			// already seen this notification
			return
		}
		sequence = data.Sequence

		select {
		case <-time.After(5 * time.Second):
			b.logger.Warn("notification buffer full for at least 5 seconds!")
		case notifications <- data:
		}

		var (
			ok bool
		)
		err := b.rcpAckNotification.Call(
			[]interface{}{data.Sequence},
			[]interface{}{&ok},
		)
		if err != nil {
			log.Println("failed to ack notification")
		}

	}

	sub, err := b.notifMailbox.Subscribe(ads.NotificationAttrib{
		TransMode: ads.TransServerOnCha,
		CycleTime: 100 * time.Millisecond,
	}, callback)
	if err != nil {
		return nil, err
	}
	b.logger.Debug("subscribed to a value over ADS", zap.Int32("total", atomic.AddInt32(&subCounter, 1)))

	b.closers = append(b.closers, sub)

	return notifications, nil
}

func (b *daliBridge) startNotificationWorker() {
	go func() {
		for {
			select {
			case <-b.ctx.Done():
				return

			case data := <-b.notifications:
				b.handleNotification(data)
			}
		}
	}()
}

func (b *daliBridge) handleNotification(data notification) {
	b.stateL.RLock()
	defer b.stateL.RUnlock()

	start := time.Now()

	done := make(chan struct{})
	events, err := data.Decode()
	if err != nil {
		go func() {
			defer close(done)
			for _, handler := range b.inputEventHandlers {
				var dummy dali.InputEvent
				handler(dummy, err)
			}
		}()
	} else {
		go func() {
			defer close(done)
			for _, event := range events {
				for _, handler := range b.inputEventHandlers {
					handler(event, nil)
				}
			}
		}()
	}

	for {
		select {
		case <-done:
			return
		case <-time.After(time.Second):
			b.logger.Warn("running event handlers took longer than a second", zap.Duration("duration", time.Since(start)))
		}
	}
}

type bridgeResponse struct {
	Valid    bool   `tc3ads:"valid"`
	Sequence uint64 `tc3ads:"sequence"`
	IsError  bool   `tc3ads:"error"`
	Status   uint32 `tc3ads:"status"`
	Message  string `tc3ads:"message"`
	Data     uint32 `tc3ads:"data"`
}

func (r *bridgeResponse) AsError() error {
	if !r.Valid {
		return ErrInvalid
	}
	if !r.IsError || r.Status == 0 {
		return nil
	}
	return dali.Error{Status: r.Status, Message: r.Message}
}

type bridgeRequest struct {
	sequence   uint64
	ctx        context.Context
	responseCh chan<- bridgeResponse
}
