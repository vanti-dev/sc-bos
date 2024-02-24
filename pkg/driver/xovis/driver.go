package xovis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime"
	"net"
	"net/http"
	"sync"
	"time"

	"go.uber.org/zap"

	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/smart-core-os/sc-golang/pkg/resource"
	"github.com/smart-core-os/sc-golang/pkg/trait"
	"github.com/smart-core-os/sc-golang/pkg/trait/enterleavesensor"
	"github.com/smart-core-os/sc-golang/pkg/trait/occupancysensor"
	"github.com/vanti-dev/sc-bos/pkg/driver"
	"github.com/vanti-dev/sc-bos/pkg/gen"
	"github.com/vanti-dev/sc-bos/pkg/gentrait/udmipb"
	"github.com/vanti-dev/sc-bos/pkg/minibus"
	"github.com/vanti-dev/sc-bos/pkg/node"
	"github.com/vanti-dev/sc-bos/pkg/task/service"
)

const DriverName = "xovis"

var Factory driver.Factory = factory{}

type factory struct{}

func (f factory) New(services driver.Services) service.Lifecycle {
	services.Logger = services.Logger.Named(DriverName)
	d := &Driver{
		Services:    services,
		pushDataBus: &minibus.Bus[PushData]{},
	}
	d.Service = service.New(
		service.MonoApply(d.applyConfig),
		service.WithParser(ParseConfig),
	)
	return d
}

type Driver struct {
	*service.Service[DriverConfig]
	driver.Services
	pushDataBus *minibus.Bus[PushData]

	m                 sync.Mutex
	config            DriverConfig
	client            *Client
	server            *http.Server // only used if httpPort is configured for the webhook
	unannounceDevices []node.Undo
	udmiServers       []*UdmiServiceServer
}

func (d *Driver) applyConfig(_ context.Context, conf DriverConfig) error {
	d.m.Lock()
	defer d.m.Unlock()

	// A route can't be removed from an HTTP ServeMux, so if it's been changed or removed then we can't support the
	// new configuration. This is likely to be rare in practice. Adding a route is fine.
	var oldWebhook, newWebhook string
	if d.config.DataPush != nil {
		oldWebhook = d.config.DataPush.WebhookPath
	}
	if conf.DataPush != nil {
		newWebhook = conf.DataPush.WebhookPath
	}
	if oldWebhook != "" && newWebhook != oldWebhook {
		return errors.New("can't change webhook path once service is running")
	}

	// create a new client to communicate with the Xovis sensor
	pass, err := conf.LoadPassword()
	if err != nil {
		return err
	}
	d.client = NewInsecureClient(conf.Host, conf.Username, pass)
	// unannounce any devices left over from a previous configuration
	for _, unannounce := range d.unannounceDevices {
		unannounce()
	}
	d.unannounceDevices = nil
	// announce new devices
	for _, dev := range conf.Devices {
		var features []node.Feature
		if dev.Metadata != nil {
			features = append(features, node.HasMetadata(dev.Metadata))
		}
		var occupancyVal *resource.Value
		if dev.Occupancy != nil {
			occupancy := &occupancyServer{
				client:         d.client,
				multiSensor:    conf.MultiSensor,
				logicID:        dev.Occupancy.ID,
				bus:            d.pushDataBus,
				OccupancyTotal: resource.NewValue(resource.WithInitialValue(&traits.Occupancy{}), resource.WithNoDuplicates()),
			}
			features = append(features, node.HasTrait(trait.OccupancySensor,
				node.WithClients(occupancysensor.WrapApi(occupancy))))
			occupancyVal = occupancy.OccupancyTotal
		}
		var enterLeaveVal *resource.Value
		if dev.EnterLeave != nil {
			enterLeave := &enterLeaveServer{
				client:          d.client,
				logicID:         dev.EnterLeave.ID,
				multiSensor:     conf.MultiSensor,
				bus:             d.pushDataBus,
				EnterLeaveTotal: resource.NewValue(resource.WithInitialValue(&traits.EnterLeaveEvent{}), resource.WithNoDuplicates()),
			}

			features = append(features, node.HasTrait(trait.EnterLeaveSensor,
				node.WithClients(enterleavesensor.WrapApi(enterLeave))))
			enterLeaveVal = enterLeave.EnterLeaveTotal
		}

		if enterLeaveVal != nil || occupancyVal != nil {
			server := NewUdmiServiceServer(d.Logger.Named("UdmiServiceServer"), enterLeaveVal, occupancyVal, dev.UDMITopicPrefix)
			d.udmiServers = append(d.udmiServers, server)
			features = append(features, node.HasTrait(udmipb.TraitName,
				node.WithClients(gen.WrapUdmiService(server))))
		}

		d.unannounceDevices = append(d.unannounceDevices, d.Node.Announce(dev.Name, features...))
	}
	// register data push webhook
	if d.server != nil {
		d.server.Close()
		d.server = nil
	}
	if dp := conf.DataPush; dp != nil && dp.WebhookPath != "" {
		if dp.HTTPPort > 0 {
			// setup a dedicate http server for the webhook, we use http
			lis, err := net.Listen("tcp", fmt.Sprintf(":%d", dp.HTTPPort))
			if err != nil {
				return err
			}
			mux := http.NewServeMux()
			mux.HandleFunc(dp.WebhookPath, d.handleWebhook)
			d.server = &http.Server{
				Handler: mux,
			}
			go d.server.Serve(lis)
		} else {
			d.HTTPMux.HandleFunc(dp.WebhookPath, d.handleWebhook)
		}
	}

	d.config = conf

	return nil
}

func (d *Driver) handleWebhook(response http.ResponseWriter, request *http.Request) {
	// verify HTTP method
	if request.Method != http.MethodPost {
		response.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// verify request body is JSON
	mediatype, _, err := mime.ParseMediaType(request.Header.Get("Content-Type"))
	if err != nil || mediatype != "application/json" {
		response.WriteHeader(http.StatusUnsupportedMediaType)
		_, _ = response.Write([]byte("invalid content type"))
		return
	}

	// read request body and parse
	rawBody, err := io.ReadAll(http.MaxBytesReader(response, request.Body, 10*1024*1024))
	if err != nil {
		maxBytesErr := &http.MaxBytesError{}
		if errors.As(err, &maxBytesErr) {
			response.WriteHeader(http.StatusRequestEntityTooLarge)
		} else {
			// If the error was not size-related then the connection probably
			// dropped. It's unlikely the client will receive the error we send here.
			response.WriteHeader(http.StatusBadRequest)
		}
		return
	}
	var body PushData
	err = json.Unmarshal(rawBody, &body)
	if err != nil {
		d.Logger.Debug("failed to parse webhook body", zap.Error(err))
		response.WriteHeader(http.StatusBadRequest)
		_, _ = response.Write([]byte(err.Error()))
		return
	}

	n := 150
	if len(rawBody) < n {
		n = len(rawBody)
	}
	d.Logger.Debug("received webhook", zap.ByteString("body", rawBody[:n]))

	// send the data to the bus
	ctx, cancel := context.WithTimeout(request.Context(), 5*time.Second)
	defer cancel()
	_ = d.pushDataBus.Send(ctx, body)
}
