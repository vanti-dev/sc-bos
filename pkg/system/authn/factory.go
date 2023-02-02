package authn

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/rs/cors"
	"github.com/vanti-dev/sc-bos/internal/auth/keycloak"
	"github.com/vanti-dev/sc-bos/internal/auth/tenant"
	"github.com/vanti-dev/sc-bos/pkg/auth/token"
	"github.com/vanti-dev/sc-bos/pkg/node"
	"github.com/vanti-dev/sc-bos/pkg/system"
	"github.com/vanti-dev/sc-bos/pkg/system/authn/config"
	"github.com/vanti-dev/sc-bos/pkg/task/service"
	"go.uber.org/zap"
)

const TokenEndpointPath = "/oauth2/token"

func Factory() system.Factory {
	return &factory{
		server: &nextOrNotFound{},
	}
}

type factory struct {
	handleOnce sync.Once // ensures we only mux.Handle once - otherwise it would panic
	server     *nextOrNotFound
}

func (f *factory) New(services system.Services) service.Lifecycle {
	f.handleOnce.Do(func() {
		services.HTTPMux.Handle(TokenEndpointPath, f.server)
	})
	s := &System{
		server:     f.server,
		dataDir:    services.DataDir,
		clienter:   services.Node,
		logger:     services.Logger.Named("authn"),
		validators: services.TokenValidators,
	}
	s.Service = service.New(service.MonoApply(s.applyConfig),
		service.WithParser(config.ReadConfig),
		service.WithOnStop[config.Root](func() {
			s.Clear()
		}),
	)
	return s
}

type System struct {
	*service.Service[config.Root]
	server *nextOrNotFound

	dataDir       string
	clienter      node.Clienter
	cohortManager node.Remote
	logger        *zap.Logger

	validators      *token.ValidatorSet
	addedValidators []token.Validator // the validators we setup in applyConfig, used to remove them again
}

func (s *System) applyConfig(_ context.Context, cfg config.Root) error {
	// other cleanup is done as part of service.WithOnStop in New
	s.deleteValidators()

	var serveTokenEndpoint bool
	tokenServerOpts := []tenant.TokenServerOption{
		tenant.WithLogger(s.logger.Named("server")),
	}

	if cfg.System != nil {
		verifier, err := s.systemTenantVerifier(cfg)
		if err != nil {
			return err
		}
		serveTokenEndpoint = true
		tokenServerOpts = append(tokenServerOpts, tenant.WithClientCredentialFlow(verifier, cfg.System.Validity.Or(15*time.Minute)))
	}

	if cfg.User != nil {
		// User accounts that are verified by external authorization servers don't need us to
		// host the oauth/token endpoint, so don't

		// Verify user credentials via the OAuth2 Password Flow using a local file containing accounts.
		// Validate access tokens that were generated via this flow.
		if cfg.User.FileAccounts != nil {
			fileVerifier, err := loadFileVerifier(cfg.User.FileAccounts, s.dataDir, "users.json")
			if err != nil {
				return fmt.Errorf("user %w", err)
			}

			serveTokenEndpoint = true
			tokenServerOpts = append(tokenServerOpts, tenant.WithPasswordFlow(fileVerifier, cfg.User.Validity.Or(24*time.Hour)))
		}

		// Validate access tokens against a remote keycloak server.
		// This is done against the keys returned by the file the OIDC JWKS URI property, which Keycloak supports.
		if cfg.User.Keycloak != nil {
			authConfig := *cfg.User.Keycloak // copy
			if authConfig.ClientID == "" {
				authConfig.ClientID = "sc-api"
			}
			validator := keycloak.NewOIDCTokenValidator(authConfig)

			s.addedValidators = append(s.addedValidators, validator)
			s.validators.Append(validator)
		}
	}

	if serveTokenEndpoint {
		server, err := tenant.NewTokenServer("authn", tokenServerOpts...)
		if err != nil {
			return fmt.Errorf("new token server: %w", err)
		}
		validator := server.TokenValidator()

		// validate any access tokens this server issued
		s.addedValidators = append(s.addedValidators, validator)
		s.validators.Append(validator)
		// serve the handler
		s.server.Next(cors.Default().Handler(server))
	}

	return nil
}

func (s *System) Clear() {
	s.server.Clear()
	s.deleteValidators()
}

func (s *System) deleteValidators() {
	vs := s.addedValidators
	s.addedValidators = nil
	for _, v := range vs {
		s.validators.Delete(v)
	}
}
