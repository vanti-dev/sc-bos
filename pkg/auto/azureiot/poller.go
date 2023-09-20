package azureiot

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"github.com/vanti-dev/sc-bos/pkg/auto/azureiot/auth"
	"github.com/vanti-dev/sc-bos/pkg/auto/azureiot/dps"
	"github.com/vanti-dev/sc-bos/pkg/auto/azureiot/iothub"
	"github.com/vanti-dev/sc-bos/pkg/gen"
)

type poller struct {
	logger  *zap.Logger
	source  gen.PointApiClient
	dialler dialler
	config  DeviceConfig
	conn    iothub.Conn

	lastErr error
}

func (d *poller) poll(ctx context.Context) {
	d.reportErr(d.pollErr(ctx))
}

func (d *poller) pollErr(ctx context.Context) error {
	conn, err := d.getConn(ctx)
	if err != nil {
		return err
	}

	points, err := d.source.GetPoints(ctx, &gen.GetPointsRequest{
		Name: d.config.Name,
	})
	if err != nil {
		return fmt.Errorf("GetPoints: %w", err)
	}

	err = conn.SendOutputMessage(ctx, points.Values.AsMap())
	if err != nil {
		return fmt.Errorf("send to broker: %w", err)
	}
	return nil
}

func (d *poller) getConn(ctx context.Context) (iothub.Conn, error) {
	if d.conn != nil {
		return d.conn, nil
	}

	conn, err := d.dialler.Dial(ctx)
	if err != nil {
		return nil, err
	}
	d.conn = conn
	return conn, nil
}

func (d *poller) reportErr(err error) {
	if err != nil && d.lastErr == nil {
		d.logger.Error("device poller is now in an error state", zap.Error(err))
	} else if err == nil && d.lastErr != nil {
		d.logger.Info("device poller is now healthy", zap.Error(err))
	} else if err.Error() != d.lastErr.Error() {
		d.logger.Error("device poller is now in a different errors state", zap.Error(err))
	}
	d.lastErr = err
}

func (d *poller) Close() error {
	if d.conn != nil {
		return d.conn.Close()
	}
	return nil
}

type dialler interface {
	Dial(ctx context.Context) (iothub.Conn, error)
}

type directDialler struct {
	params iothub.ConnectionParameters
}

func (d *directDialler) Dial(ctx context.Context) (iothub.Conn, error) {
	return iothub.Dial(ctx, d.params)
}

type dpsDialler struct {
	host    string
	idScope string
	regID   string
	key     auth.SASKey

	reg *dps.Registration
}

func (d *dpsDialler) Dial(ctx context.Context) (iothub.Conn, error) {
	if d.reg == nil {
		reg, err := dps.Register(ctx, d.host, d.idScope, d.regID, d.key)
		if err != nil {
			return nil, err
		}
		d.reg = &reg
	}

	params := d.reg.AsConnectionParameters()
	return iothub.Dial(ctx, params)
}
