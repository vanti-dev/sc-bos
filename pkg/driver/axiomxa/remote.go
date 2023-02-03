package axiomxa

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/vanti-dev/sc-bos/pkg/driver/axiomxa/config"
	"github.com/vanti-dev/sc-bos/pkg/driver/axiomxa/rpc"
)

type server struct {
	rpc.UnimplementedAxiomXaDriverServiceServer
	logger *zap.Logger
	config config.Root
}

func (s *server) SaveQRCredential(ctx context.Context, request *rpc.SaveQRCredentialRequest) (*rpc.SaveQRCredentialResponse, error) {
	// todo: update the actual request once we know what it looks like!
	body := struct {
		Credential string `json:"credential,omitempty"`
		AccountID  string `json:"accountId,omitempty"`
	}{
		AccountID:  request.Account,
		Credential: string(request.QrBody),
	}
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", s.config.HTTP.BaseURL+"/credentials/add", bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return nil, status.Errorf(codes.Unavailable, "upstream %v", res.Status)
	}
	_, err = io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return &rpc.SaveQRCredentialResponse{}, nil
}
