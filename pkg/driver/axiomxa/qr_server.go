package axiomxa

import (
	"context"
	"time"

	"go.uber.org/zap"

	"github.com/vanti-dev/sc-bos/pkg/driver/axiomxa/config"
	"github.com/vanti-dev/sc-bos/pkg/driver/axiomxa/jsonapi"
	"github.com/vanti-dev/sc-bos/pkg/gen"
)

type qrServer struct {
	gen.UnimplementedAxiomXaDriverServiceServer
	logger *zap.Logger
	client *jsonapi.Client
	config config.Root
}

func (s *qrServer) SaveQRCredential(ctx context.Context, request *gen.SaveQRCredentialRequest) (*gen.SaveQRCredentialResponse, error) {
	_, err := s.client.CreateCardholder(ctx, jsonapi.Cardholder{
		FirstName: request.FirstName,
		LastName:  request.LastName,
		Cards: []jsonapi.KeepUnknown[jsonapi.Card]{
			{Known: cardFromQRConfig(request, s.config.QR)},
		},
	})
	if err != nil {
		return nil, err
	}
	return &gen.SaveQRCredentialResponse{}, nil
}

func cardFromQRConfig(req *gen.SaveQRCredentialRequest, config *config.QR) jsonapi.Card {
	dst := jsonapi.Card{}
	dst.CardNumber = req.CardNumber
	dst.CardType = jsonapi.CardTypeVisitor
	if config != nil {
		dst.AccessLevel = config.AccessLevel
	}
	if req.ActiveTime == nil {
		dst.ActiveDate = time.Now()
	} else {
		dst.ActiveDate = req.ActiveTime.AsTime()
	}
	if req.ExpireTime == nil {
		if config != nil && config.ExpireAfter.Duration != 0 {
			dst.ExpiryDate = dst.ActiveDate.Add(config.ExpireAfter.Duration)
		}
	} else {
		dst.ExpiryDate = req.ExpireTime.AsTime()
	}
	return dst
}
