package app

import (
	"errors"
	"github.com/vanti-dev/bsp-ew/internal/manage/enrollment"
	"google.golang.org/grpc"
)

func RemoteManager(server *enrollment.Server, opts ...grpc.DialOption) (connFunc func() (*grpc.ClientConn, error), closeFunc func() error) {
	var lastManagerAddress string
	var lastConn *grpc.ClientConn
	var lastErr error
	closeConn := func() error {
		if lastConn != nil {
			err := lastConn.Close()
			lastConn = nil
			lastErr = nil
			return err
		}
		return nil
	}
	return func() (*grpc.ClientConn, error) {
		e, ok := server.Enrollment()
		if !ok {
			_ = closeConn()
			return nil, errors.New("not enrolled")
		}

		if e.ManagerAddress != lastManagerAddress {
			_ = closeConn()
			lastManagerAddress = e.ManagerAddress
		}

		if lastConn == nil {
			lastConn, lastErr = grpc.Dial(lastManagerAddress, opts...)
		}

		return lastConn, lastErr
	}, closeConn
}
