package hikcentral

import (
	"context"
	"math"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/smart-core-os/sc-api/go/types"
	"github.com/smart-core-os/sc-golang/pkg/resource"
	"github.com/vanti-dev/sc-bos/pkg/driver/hikcentral/api"
	"github.com/vanti-dev/sc-bos/pkg/driver/hikcentral/config"
	"github.com/vanti-dev/sc-bos/pkg/gen"
)

var (
	errStartTimeIsZero       = status.Error(codes.InvalidArgument, "access grant start time is zero")
	errEndTimeIsZero         = status.Error(codes.InvalidArgument, "access grant end time is zero")
	errStartTimeAfterEndTime = status.Error(codes.InvalidArgument, "access grant start time is after end time")
)

type ANPRController struct {
	gen.UnimplementedAccessApiServer

	client *api.Client
	logger *zap.Logger
	conf   *config.Root

	resource     *resource.Value // of type *gen.AccessAttempt
	lastPollTime time.Time
}

func NewANPRController(client *api.Client, logger *zap.Logger, conf *config.Root) *ANPRController {
	return &ANPRController{
		client:       client,
		logger:       logger.Named("visitor management"),
		conf:         conf,
		resource:     resource.NewValue(resource.WithNoDuplicates()),
		lastPollTime: time.Now().Add(time.Hour * -1), // start polling from 1 hour ago
	}
}

var _ gen.AccessApiServer = (*ANPRController)(nil)

func (v *ANPRController) poll(ctx context.Context) error {
	v.logger.Debug("starting visitor management poll")

	// get appointments for today
	today := time.Now().Truncate(24 * time.Hour)
	tomorrow := today.Add(24 * time.Hour)

	var appointmentsList []api.Appointment

	pageNo := 1

	for {
		select {
		case <-ctx.Done():
			v.logger.Debug("polling cancelled")
			return ctx.Err()
		default:
		}
		res, err := v.client.ListVisitorAppointments(&api.ListVisitorAppointmentsRequest{
			AppointStartTime: formatTime(today),
			AppointEndTime:   formatTime(tomorrow),
			Request: api.Request{
				PageNo:   pageNo,
				PageSize: 100,
			},
		})
		if err != nil {
			return status.Error(codes.Unavailable, "failed to list visitor appointments: "+err.Error())
		}
		if len(res.List) < 1 {
			break
		}

		appointmentsList = append(appointmentsList, res.List...)

		pageNo++
	}

	appointments := appointmentListToMap(appointmentsList)

	pageNo = 1
	for {
		select {
		case <-ctx.Done():
			v.logger.Debug("polling cancelled")
			return ctx.Err()
		default:
		}

		resp, err := v.client.ListANPREvents(&api.ANPREventsRequest{
			CameraIndexCode: strings.Join(v.conf.ANPR.EntranceCameraIndexCodes, ","),
			StartTime:       formatTime(v.lastPollTime),
			EndTime:         formatTime(time.Now()),
			PageNo:          pageNo,
			PageSize:        100,
		})

		if err != nil {
			return status.Error(codes.Unavailable, "failed to list ANPR events: "+err.Error())
		}

		if len(resp.List) < 1 {
			v.logger.Debug("no ANPR events found", zap.Int("pageNo", pageNo))
			break
		}

		for _, event := range resp.List {
			grant := gen.AccessAttempt_DENIED
			reason := "no appointment found"
			actor := &gen.Actor{}
			accessTime, err := time.Parse(time.RFC3339, event.CrossTime)

			if err != nil {
				v.logger.Error("failed to parse access time", zap.Error(err), zap.String("crossTime", event.CrossTime))
				accessTime = time.Now() // fallback to current time if parsing fails
			}

			relevantAppointments, ok := appointments[formatPlateNo(event.PlateNo)]

			if !ok {
				v.logger.Debug("no appointments found for ANPR event", zap.String("plateNo", event.PlateNo))
				accessAttempt := &gen.AccessAttempt{
					Grant:             grant,
					Reason:            reason,
					Actor:             actor,
					AccessAttemptTime: timestamppb.New(accessTime),
				}

				_, err = v.resource.Set(accessAttempt)

				if err != nil {
					v.logger.Error("failed to set access attempt", zap.Error(err), zap.String("plateNo", event.PlateNo))
				}

				continue
			}

			appointment := getMostRelevantAppointment(accessTime, relevantAppointments)
			grant = gen.AccessAttempt_GRANTED
			reason, actor = extractFieldsFromAppointment(appointment)

			accessAttempt := &gen.AccessAttempt{
				Grant:             grant,
				Reason:            reason,
				Actor:             actor,
				AccessAttemptTime: timestamppb.New(accessTime),
			}

			_, err = v.resource.Set(accessAttempt)

			if err != nil {
				v.logger.Error("failed to set access attempt", zap.Error(err), zap.String("plateNo", event.PlateNo))
			}
		}

		pageNo++
	}

	v.lastPollTime = time.Now()

	return nil
}

func (v *ANPRController) GetLastAccessAttempt(ctx context.Context, request *gen.GetLastAccessAttemptRequest) (*gen.AccessAttempt, error) {
	return v.resource.Get(resource.WithReadMask(request.GetReadMask())).(*gen.AccessAttempt), nil
}

func (v *ANPRController) PullAccessAttempts(request *gen.PullAccessAttemptsRequest, server gen.AccessApi_PullAccessAttemptsServer) error {
	ctx, cancel := context.WithCancel(server.Context())
	defer cancel()

	changes := v.resource.Pull(ctx, resource.WithReadMask(request.GetReadMask()), resource.WithUpdatesOnly(request.GetUpdatesOnly()))
	for item := range changes {
		if item == nil {
			continue
		}

		accessAttempt := item.Value.(*gen.AccessAttempt)

		if err := server.Send(&gen.PullAccessAttemptsResponse{
			Changes: []*gen.PullAccessAttemptsResponse_Change{
				{
					AccessAttempt: accessAttempt,
					ChangeTime:    timestamppb.New(item.ChangeTime),
					Name:          request.GetName(),
				},
			},
		}); err != nil {
			return status.Error(codes.Internal, "failed to send access attempt: "+err.Error())
		}
	}

	return nil
}

func (v *ANPRController) CreateAccessGrant(ctx context.Context, request *gen.CreateAccessGrantRequest) (*gen.AccessGrant, error) {
	v.logger.Debug("creating access Grant", zap.String("source", request.GetSource()))

	if err := validAccessGrantTiming(request.GetAccessGrant()); err != nil {
		return nil, err
	}

	autoApproval, err := v.client.CheckAutoReviewFlow(&api.AutoReviewFlowRequest{})

	if err != nil {
		return nil, status.Error(codes.Unavailable, "failed to check auto review flow: "+err.Error())
	}

	manuallyApprove := v.conf.ANPR.EnableSmartCoreApproval && autoApproval.AutomaticApproval == 0

	req := createAppointmentRequest(request)

	resp, err := v.client.CreateVisitorAppointment(req)

	if err != nil {
		return nil, status.Error(codes.Unavailable, err.Error())
	}

	if manuallyApprove {
		v.logger.Debug("manually approving access grant reservation")

		appointmentRecordID, err := strconv.ParseInt(resp.AppointRecordID, 10, 64)

		if err != nil {
			return nil, status.Error(codes.Internal, "failed to parse appointment record ID: "+err.Error())
		}

		resp, err := v.client.ManuallyApproveVisitor(&api.VisitorApprovalRequest{
			OperateType:     0,
			ApprovalOpinion: "created via SmartCore",
			ApprovalFlowInfo: api.ApprovalFlowInfo{
				ApprovalFlowCode: int(appointmentRecordID),
			},
		})

		if err != nil {
			return nil, status.Error(codes.Unavailable, "failed to manually approve visitor: "+err.Error())
		}

		v.logger.Info("manually approving visitor results", zap.Int("errorCode", resp.ErrorCode), zap.Int("visitorId", resp.VisitorId))
	}

	grantee := request.GetAccessGrant().GetGrantee()
	grantee.Name = resp.VisitorID

	return &gen.AccessGrant{
		CanUpdate: !manuallyApprove, // warning: only if auto approval is enabled, we are permitted to update the appointment
		Id:        resp.AppointRecordID,
		StartTime: request.GetAccessGrant().GetStartTime(),
		EndTime:   request.GetAccessGrant().GetEndTime(),
		Purpose:   request.GetAccessGrant().GetPurpose(),
		EntryCode: resp.AppointCode,
		QrCode:    resp.QRCodeImage,
		Grantee:   grantee,
		Granter:   request.GetAccessGrant().GetGranter(),
	}, nil
}

func (v *ANPRController) UpdateAccessGrant(ctx context.Context, request *gen.UpdateAccessGrantRequest) (*gen.AccessGrant, error) {
	v.logger.Debug("updating access Grant", zap.String("source", request.GetSource()))

	if request.GetAccessGrant().GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "access Grant ID is required")
	}

	if err := validAccessGrantTiming(request.GetAccessGrant()); err != nil {
		return nil, err
	}

	req := createAppointmentRequest(request)

	req.AppointRecordID = request.GetAccessGrant().GetId()

	resp, err := v.client.UpdateVisitorAppointment(req)

	if err != nil {
		return nil, status.Error(codes.Unavailable, err.Error())
	}

	grantee := request.GetAccessGrant().GetGrantee()
	grantee.Name = resp.VisitorID

	return &gen.AccessGrant{
		Id:        resp.AppointRecordID,
		StartTime: request.GetAccessGrant().GetStartTime(),
		EndTime:   request.GetAccessGrant().GetEndTime(),
		Purpose:   request.GetAccessGrant().GetPurpose(),
		EntryCode: resp.AppointCode,
		QrCode:    resp.QRCodeImage,
		Grantee:   grantee,
		Granter:   request.GetAccessGrant().GetGranter(),
	}, nil

}

func (v *ANPRController) DeleteAccessGrant(ctx context.Context, request *gen.DeleteAccessGrantRequest) (*gen.DeleteAccessGrantResponse, error) {
	v.logger.Debug("deleting access Grant", zap.String("source", request.GetSource()))

	if request.GetAccessGrantId() == "" {
		return nil, status.Error(codes.InvalidArgument, "access Grant ID is required")
	}

	_, err := v.client.DeleteVisitorAppointment(&api.DeleteVisitorAppointmentRequest{AppointRecordID: request.GetAccessGrantId()})

	if err != nil {
		return nil, status.Error(codes.Unavailable, err.Error())
	}

	return &gen.DeleteAccessGrantResponse{}, nil
}

type accessGrantRequest interface {
	GetAccessGrant() *gen.AccessGrant
}

func validAccessGrantTiming(grant *gen.AccessGrant) error {
	if grant.GetStartTime().AsTime().IsZero() {
		return errStartTimeIsZero
	}

	if grant.GetEndTime().AsTime().IsZero() {
		return errEndTimeIsZero
	}

	if grant.GetStartTime().AsTime().After(grant.GetEndTime().AsTime()) {
		return errStartTimeAfterEndTime
	}

	return nil
}

func createAppointmentRequest(request accessGrantRequest) *api.VisitorAppointmentRequest {
	var visitReasonType int
	purpose := strings.ToLower(request.GetAccessGrant().GetPurpose())
	if strings.Contains(purpose, "business") {
		visitReasonType = 0
	} else if strings.Contains(purpose, "training") {
		visitReasonType = 1
	} else if strings.Contains(purpose, "visit") {
		visitReasonType = 2
	} else if strings.Contains(purpose, "meeting") {
		visitReasonType = 3
	} else {
		visitReasonType = 4
	}
	if len(purpose) > 128 {
		purpose = purpose[:128]
	}

	var pictureUrls []api.VisitorPhoto
	for _, source := range request.GetAccessGrant().GetGrantee().GetPicture().GetSources() {
		for _, content := range source.GetSrc() {
			if content.GetUrl() != "" {
				pictureUrls = append(pictureUrls, api.VisitorPhoto{
					ImageURL: content.GetUrl(),
				})
			}
		}
	}

	receptionistId := request.GetAccessGrant().GetGranter().GetName()
	if len(receptionistId) > 64 {
		receptionistId = receptionistId[:64]
	}

	return &api.VisitorAppointmentRequest{
		ReceptionistID:    receptionistId,
		AppointStartTime:  formatTime(request.GetAccessGrant().GetStartTime().AsTime()),
		AppointEndTime:    formatTime(request.GetAccessGrant().GetEndTime().AsTime()),
		VisitReasonType:   visitReasonType,
		VisitReasonDetail: purpose,
		VisitorInfoList: []api.VisitorInfo{
			{
				VisitorFamilyName: request.GetAccessGrant().GetGrantee().GetName(),
				VisitorGivenName:  request.GetAccessGrant().GetGrantee().GetDisplayName(),
				Gender:            0, // unknown
				Email:             request.GetAccessGrant().GetGrantee().GetEmail(),
				PlateNo:           request.GetAccessGrant().GetGrantee().GetVehicleRegistration(),
				CompanyName:       request.GetAccessGrant().GetGrantee().GetCompany(),
				VisitorPhoto:      pictureUrls,
			},
		},
	}
}

func appointmentListToMap(list []api.Appointment) map[string][]*api.Appointment {
	appointments := make(map[string][]*api.Appointment)
	for _, appointment := range list {
		key := formatPlateNo(appointment.VisitorInfo.PlateNo)
		appointments[key] = append(appointments[key], &appointment)
	}
	return appointments
}

func getMostRelevantAppointment(now time.Time, appointments []*api.Appointment) *api.Appointment {
	if len(appointments) == 0 {
		return nil
	}

	var relevantAppointment *api.Appointment

	for _, appointment := range appointments {
		if relevantAppointment == nil {
			relevantAppointment = appointment
			continue
		}

		start, err := time.Parse(time.RFC3339, appointment.AppointStartTime)
		if err != nil {
			continue
		}

		end, err := time.Parse(time.RFC3339, appointment.AppointEndTime)
		if err != nil {
			continue
		}

		startDelta := now.Sub(start)
		endDelta := now.Sub(end)

		relevantStart, err := time.Parse(time.RFC3339, relevantAppointment.AppointStartTime)
		if err != nil {
			continue
		}

		relevantEnd, err := time.Parse(time.RFC3339, relevantAppointment.AppointEndTime)
		if err != nil {
			continue
		}

		relevantStartDelta := now.Sub(relevantStart)
		relevantEndDelta := now.Sub(relevantEnd)

		// the difference between the start times and now is smaller for the current appointment,
		// i.e. it is more relevant to current time
		if math.Abs(float64(startDelta)) < math.Abs(float64(relevantStartDelta)) {
			relevantAppointment = appointment
		}

		// we still prefer longer appointments, so if the end time is further away from now, we prefer it
		// unless the start time is in a different hour, then we prefer the one with the closer start time to now
		if relevantStart.Truncate(time.Hour) == start.Truncate(time.Hour) && math.Abs(float64(endDelta)) > math.Abs(float64(relevantEndDelta)) {
			relevantAppointment = appointment
		}
	}

	return relevantAppointment
}

func extractFieldsFromAppointment(appointment *api.Appointment) (string, *gen.Actor) {
	var reason string
	switch appointment.VisitReasonType {
	case 0:
		reason = "business"
	case 1:
		reason = "training"
	case 2:
		reason = "visit"
	case 3:
		reason = "meeting"
	}

	var pictures []*types.Image_Content
	for _, pic := range appointment.VisitorInfo.VisitorPhoto {
		pictures = append(pictures, &types.Image_Content{
			Content: &types.Image_Content_Url{
				Url: pic.ImageURL,
			},
		})
	}

	actor := &gen.Actor{
		Name:        appointment.VisitorInfo.VisitorID,
		Title:       appointment.VisitorInfo.VisitorFamilyName,
		DisplayName: appointment.VisitorInfo.VisitorGivenName,
		Picture: &types.Image{
			Sources: []*types.Image_Source{
				{
					Src: pictures,
				},
			},
		},
		Email:               appointment.VisitorInfo.Email,
		Company:             appointment.VisitorInfo.CompanyName,
		VehicleRegistration: appointment.VisitorInfo.PlateNo,
	}

	return reason, actor
}

func formatPlateNo(plateNo string) string {
	// remove all non-alphanumeric characters and convert to uppercase
	plateNo = strings.ToUpper(plateNo)
	var sb strings.Builder
	for _, r := range plateNo {
		if (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
			sb.WriteRune(r)
		}
	}
	return sb.String()
}
