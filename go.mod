module github.com/vanti-dev/sc-bos

go 1.23

require (
	github.com/cenkalti/backoff/v4 v4.3.0
	github.com/chzyer/readline v0.0.0-20180603132655-2972be24d48e
	github.com/eclipse/paho.mqtt.golang v1.4.2
	github.com/go-jose/go-jose/v4 v4.0.4
	github.com/google/go-cmp v0.6.0
	github.com/google/renameio/v2 v2.0.0
	github.com/grpc-ecosystem/go-grpc-middleware v1.3.0
	github.com/hashicorp/go-retryablehttp v0.7.7
	github.com/improbable-eng/grpc-web v0.15.0
	github.com/jackc/pgconn v1.14.3
	github.com/jackc/pgx/v4 v4.18.2
	github.com/mwitkow/grpc-proxy v0.0.0-20230212185441-f345521cb9c9
	github.com/olebedev/emitter v0.0.0-20190110104742-e8d1457e6aee
	github.com/open-policy-agent/opa v0.68.0
	github.com/pborman/uuid v1.2.0
	github.com/pkg/errors v0.9.1
	github.com/qri-io/jsonpointer v0.1.1
	github.com/robfig/cron/v3 v3.0.1
	github.com/rs/cors v1.8.3
	github.com/sirupsen/logrus v1.9.3
	github.com/smart-core-os/sc-api/go v1.0.0-beta.50
	github.com/smart-core-os/sc-golang v0.0.0-20241220144351-884672945826
	github.com/stretchr/testify v1.9.0
	github.com/timshannon/bolthold v0.0.0-20210913165410-232392fc8a6a
	github.com/vanti-dev/gobacnet v0.0.0-20231102122752-32b0b38bcc53
	go.etcd.io/bbolt v1.3.10
	go.uber.org/multierr v1.9.0
	go.uber.org/zap v1.21.0
	golang.org/x/crypto v0.27.0
	golang.org/x/exp v0.0.0-20240823005443-9b4947da3948
	golang.org/x/oauth2 v0.22.0
	golang.org/x/sync v0.8.0
	golang.org/x/term v0.24.0
	golang.org/x/time v0.6.0
	golang.org/x/tools v0.24.0
	google.golang.org/grpc v1.67.1
	google.golang.org/protobuf v1.34.2
)

require (
	cloud.google.com/go/compute/metadata v0.5.0 // indirect
	github.com/OneOfOne/xxhash v1.2.8 // indirect
	github.com/agnivade/levenshtein v1.1.1 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/chzyer/test v1.0.0 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/desertbit/timer v0.0.0-20180107155436-c41aec40b27f // indirect
	github.com/go-ini/ini v1.67.0 // indirect
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/gobwas/glob v0.2.3 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/gorilla/mux v1.8.1 // indirect
	github.com/gorilla/websocket v1.4.2 // indirect
	github.com/hashicorp/go-cleanhttp v0.5.2 // indirect
	github.com/jackc/chunkreader/v2 v2.0.1 // indirect
	github.com/jackc/pgio v1.0.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgproto3/v2 v2.3.3 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/pgtype v1.14.0 // indirect
	github.com/jackc/puddle v1.3.0 // indirect
	github.com/mennanov/fmutils v0.1.1 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/prometheus/client_golang v1.20.2 // indirect
	github.com/prometheus/client_model v0.6.1 // indirect
	github.com/prometheus/common v0.55.0 // indirect
	github.com/prometheus/procfs v0.15.1 // indirect
	github.com/rcrowley/go-metrics v0.0.0-20200313005456-10cdbea86bc0 // indirect
	github.com/tanema/gween v0.0.0-20200427131925-c89ae23cc63c // indirect
	github.com/tchap/go-patricia/v2 v2.3.1 // indirect
	github.com/xeipuuv/gojsonpointer v0.0.0-20190905194746-02993c407bfb // indirect
	github.com/xeipuuv/gojsonreference v0.0.0-20180127040603-bd5ef7bd5415 // indirect
	github.com/yashtewari/glob-intersection v0.2.0 // indirect
	go.opentelemetry.io/otel v1.28.0 // indirect
	go.opentelemetry.io/otel/metric v1.28.0 // indirect
	go.opentelemetry.io/otel/sdk v1.28.0 // indirect
	go.opentelemetry.io/otel/trace v1.28.0 // indirect
	go.uber.org/atomic v1.9.0 // indirect
	go.uber.org/goleak v1.3.0 // indirect
	golang.org/x/net v0.29.0 // indirect
	golang.org/x/sys v0.25.0 // indirect
	golang.org/x/text v0.18.0 // indirect
	google.golang.org/genproto v0.0.0-20231211222908-989df2bf70f3 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240930140551-af27646dc61f // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	nhooyr.io/websocket v1.8.10 // indirect
	sigs.k8s.io/yaml v1.4.0 // indirect
)
