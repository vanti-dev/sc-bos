package merge

import (
	"context"
	"encoding/json"
	"errors"
	"slices"

	"go.uber.org/multierr"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	grpcStatus "google.golang.org/grpc/status"

	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/smart-core-os/sc-golang/pkg/cmp"
	"github.com/smart-core-os/sc-golang/pkg/resource"
	"github.com/smart-core-os/sc-golang/pkg/trait"
	"github.com/smart-core-os/sc-golang/pkg/trait/lightpb"
	"github.com/vanti-dev/gobacnet"
	"github.com/vanti-dev/sc-bos/pkg/driver/bacnet/comm"
	"github.com/vanti-dev/sc-bos/pkg/driver/bacnet/config"
	"github.com/vanti-dev/sc-bos/pkg/driver/bacnet/known"
	"github.com/vanti-dev/sc-bos/pkg/driver/bacnet/status"
	"github.com/vanti-dev/sc-bos/pkg/gentrait/statuspb"
	"github.com/vanti-dev/sc-bos/pkg/node"
	"github.com/vanti-dev/sc-bos/pkg/task"
)

var (
	errSceneNotFound = errors.New("scene not found")
)

type lightCfg struct {
	config.Trait
	Point  *config.ValueSource `json:"point,omitempty"`
	Scenes []sceneCfg          `json:"scenes,omitempty"` // scenes to activate on this point
}

type sceneCfg struct {
	Name       string  `json:"name,omitempty"`       // the name of the scene
	Title      string  `json:"title,omitempty"`      // the title of the scene, used for display
	SetValue   int     `json:"setValue,omitempty"`   // the value on the point to activate this scene
	Brightness float32 `json:"brightness,omitempty"` // the brightness to set when activating this scene
}

func readLightConfig(raw []byte) (cfg lightCfg, err error) {
	err = json.Unmarshal(raw, &cfg)

	slices.SortFunc(cfg.Scenes, func(a, b sceneCfg) int {
		if a.Brightness < b.Brightness {
			return -1
		}
		return 1
	})

	return
}

type light struct {
	client   *gobacnet.Client
	known    known.Context
	statuses *statuspb.Map
	logger   *zap.Logger

	model *lightpb.Model
	*lightpb.ModelServer
	config   lightCfg
	pollTask *task.Intermittent
}

func newLight(client *gobacnet.Client, devices known.Context, statuses *statuspb.Map, config config.RawTrait, logger *zap.Logger) (*light, error) {
	cfg, err := readLightConfig(config.Raw)
	if err != nil {
		return nil, err
	}

	model := lightpb.NewModel(resource.WithMessageEquivalence(cmp.Equal(cmp.FloatValueApprox(0, 1.0)))) // report brightness intensity changes of 1.0% or more

	l := &light{
		client:      client,
		known:       devices,
		statuses:    statuses,
		logger:      logger,
		model:       model,
		ModelServer: lightpb.NewModelServer(model),
		config:      cfg,
	}

	l.pollTask = task.NewIntermittent(l.startPoll)

	initTraitStatus(statuses, cfg.Name, "Light")

	return l, nil
}

func (l *light) AnnounceSelf(a node.Announcer) node.Undo {
	return a.Announce(l.config.Name, node.HasTrait(trait.Light, node.WithClients(lightpb.WrapApi(l), lightpb.WrapInfo(l))))
}

func (l *light) UpdateBrightness(ctx context.Context, request *traits.UpdateBrightnessRequest) (*traits.Brightness, error) {
	presetName := request.GetBrightness().GetPreset().GetName()

	scene, err := l.findSceneByName(presetName)

	if err != nil {
		return nil, grpcStatus.Error(codes.InvalidArgument, "preset not found")
	}

	if l.config.Point != nil {
		err := comm.WriteProperty(ctx, l.client, l.known, *l.config.Point, scene.SetValue, 0)
		if err != nil {
			l.logger.Error("WriteProperty Scene failed", zap.Error(err))
			return nil, err
		}
	}

	return pollUntil(ctx, l.config.DefaultRWConsistencyTimeoutDuration(), l.pollPeer, func(brightness *traits.Brightness) bool {
		return brightness.LevelPercent == scene.Brightness
	})
}

func (l *light) GetBrightness(ctx context.Context, request *traits.GetBrightnessRequest) (*traits.Brightness, error) {
	_, err := l.pollPeer(ctx)
	if err != nil {
		return nil, err
	}
	return l.ModelServer.GetBrightness(ctx, request)
}

func (l *light) PullBrightness(request *traits.PullBrightnessRequest, server traits.LightApi_PullBrightnessServer) error {
	_ = l.pollTask.Attach(server.Context())
	return l.ModelServer.PullBrightness(request, server)
}

func (l *light) DescribeBrightness(ctx context.Context, request *traits.DescribeBrightnessRequest) (*traits.BrightnessSupport, error) {
	var presets []*traits.LightPreset
	for _, scene := range l.config.Scenes {
		presets = append(presets, &traits.LightPreset{
			Name:  scene.Name,
			Title: scene.Title,
		})
	}

	return &traits.BrightnessSupport{
		Presets: presets,
	}, nil
}

func (l *light) startPoll(init context.Context) (stop task.StopFn, err error) {
	return startPoll(init, "light", l.config.PollPeriodDuration(), l.config.PollTimeoutDuration(), l.logger, func(ctx context.Context) error {
		_, err := l.pollPeer(ctx)
		return err
	})
}

func (l *light) pollPeer(ctx context.Context) (*traits.Brightness, error) {
	data := &traits.Brightness{Preset: &traits.LightPreset{}}
	var resProcessors []func(response any) error
	var readValues []config.ValueSource
	var requestNames []string

	if l.config.Point != nil {
		requestNames = append(requestNames, "light")
		readValues = append(readValues, *l.config.Point)
		resProcessors = append(resProcessors, func(response any) error {
			value, err := comm.Float32Value(response)
			if err != nil {
				return comm.ErrReadProperty{Prop: "light", Cause: err}
			}

			scene, err := l.findSceneByValue(value)

			if err != nil {
				l.logger.Error("failed to find scene by value", zap.Error(err), zap.Float32("value", value))
				return nil
			}

			data.Preset.Name = scene.Name
			data.Preset.Title = scene.Title
			data.LevelPercent = scene.Brightness
			return nil
		})
	}

	responses := comm.ReadProperties(ctx, l.client, l.known, readValues...)
	var errs []error
	for i, response := range responses {
		err := resProcessors[i](response)
		if err != nil {
			errs = append(errs, err)
		}
	}

	status.UpdatePollErrorStatus(l.statuses, l.config.Name, "Light", requestNames, errs)
	if len(errs) > 0 {
		return nil, multierr.Combine(errs...)
	}

	return l.model.UpdateBrightness(data)
}

func (l *light) findSceneByValue(sceneCmd float32) (*sceneCfg, error) {
	if sceneCmd >= 0 && sceneCmd <= 100 {
		// If the sceneCmd is a percentage, find the scene with that brightness
		for idx := 0; idx < len(l.config.Scenes)-1; idx++ {
			if l.config.Scenes[idx].Brightness >= sceneCmd && l.config.Scenes[idx+1].Brightness <= sceneCmd {
				return &l.config.Scenes[idx+1], nil
			}
		}

		return nil, errSceneNotFound
	}

	for _, scene := range l.config.Scenes {
		if scene.SetValue == int(sceneCmd) {
			return &scene, nil
		}
	}
	return nil, errSceneNotFound
}

func (l *light) findSceneByName(name string) (*sceneCfg, error) {
	for _, scene := range l.config.Scenes {
		if scene.Name == name {
			return &scene, nil
		}
	}
	return nil, errSceneNotFound
}
