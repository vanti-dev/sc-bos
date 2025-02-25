package account

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"io"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/vanti-dev/sc-bos/internal/account/queries"
	"github.com/vanti-dev/sc-bos/internal/database"
	"github.com/vanti-dev/sc-bos/internal/util/pass"
	"github.com/vanti-dev/sc-bos/pkg/gen"
)

const appID = 0x5C0501

const maxServiceCredentialsPerAccount = 2

type Store struct {
	db *database.Database
}

func OpenStore(ctx context.Context, path string, logger *zap.Logger) (*Store, error) {
	db, err := database.Open(ctx, path,
		database.WithLogger(logger),
		database.WithApplicationID(appID),
	)
	if err != nil {
		return nil, err
	}

	err = db.Migrate(ctx, schema)
	if err != nil {
		return nil, err
	}

	return &Store{
		db: db,
	}, nil
}

func NewMemoryStore(logger *zap.Logger) *Store {
	db := database.OpenMemory(
		database.WithLogger(logger),
		database.WithApplicationID(appID),
	)

	err := db.Migrate(context.Background(), schema)
	if err != nil {
		// this can only happen if the migrations are broken
		panic(err)
	}

	return &Store{db: db}
}

func (s *Store) Close() error {
	return s.db.Close()
}

func (s *Store) Read(ctx context.Context, f func(tx *Tx) error) error {
	return s.db.ReadTx(ctx, func(tx *sql.Tx) error {
		storeTx := &Tx{Queries: queries.New(tx)}
		return f(storeTx)
	})
}

func (s *Store) Write(ctx context.Context, f func(tx *Tx) error) error {
	return s.db.WriteTx(ctx, func(tx *sql.Tx) error {
		storeTx := &Tx{Queries: queries.New(tx)}
		return f(storeTx)
	})
}

type Tx struct {
	*queries.Queries
}

func (tx *Tx) UpdateAccountPassword(ctx context.Context, accountID int64, password string) error {
	password = normalisePassword(password)
	if !permitPassword(password) {
		return ErrInvalidPassword
	}

	hash, err := pass.Hash([]byte(password))
	if err != nil {
		return err
	}

	account, err := tx.GetAccount(ctx, accountID)
	if err != nil {
		return err
	}
	if account.Kind != gen.Account_USER_ACCOUNT.String() {
		return ErrUnexpectedPassword
	}

	return tx.UpdateAccountPasswordHash(ctx, queries.UpdateAccountPasswordHashParams{
		AccountID:    accountID,
		PasswordHash: hash,
	})
}

func (tx *Tx) CheckAccountPassword(ctx context.Context, accountID int64, password string) error {
	password = normalisePassword(password)

	hash, err := tx.GetAccountPasswordHash(ctx, accountID)
	if err != nil {
		return err
	}

	err = pass.Compare(hash, []byte(password))
	if err != nil {
		return err
	}

	return nil
}

func (tx *Tx) GenerateServiceCredential(ctx context.Context, accountID int64, title string, expiry sql.NullTime) (GeneratedServiceCredential, error) {
	if !validateTitle(title) {
		return GeneratedServiceCredential{}, status.Error(codes.InvalidArgument, "invalid title")
	}

	count, err := tx.CountServiceCredentialsForAccount(ctx, accountID)
	if err != nil {
		return GeneratedServiceCredential{}, err
	}
	// refuse to generate a new credential if the limit is reached
	if count >= maxServiceCredentialsPerAccount {
		return GeneratedServiceCredential{}, ErrServiceCredentialLimit
	}

	secret, err := genSecret()
	if err != nil {
		return GeneratedServiceCredential{}, err
	}

	hash := sha256.Sum256([]byte(secret))

	cred, err := tx.CreateServiceCredential(ctx, queries.CreateServiceCredentialParams{
		AccountID:  accountID,
		Title:      title,
		ExpireTime: expiry,
		SecretHash: hash[:],
	})
	if err != nil {
		return GeneratedServiceCredential{}, err
	}

	return GeneratedServiceCredential{
		ServiceCredential: cred,
		Secret:            secret,
	}, nil
}

type GeneratedServiceCredential struct {
	queries.ServiceCredential
	Secret string
}

func genSecret() (string, error) {
	secretBytes := make([]byte, 32)
	_, err := io.ReadFull(rand.Reader, secretBytes)
	if err != nil {
		return "", err
	}

	encoding := base64.URLEncoding
	encoded := make([]byte, encoding.EncodedLen(len(secretBytes)))
	encoding.Encode(encoded, secretBytes)

	return string(encoded), nil
}
