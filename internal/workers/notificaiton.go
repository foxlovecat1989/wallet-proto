package workers

import (
	"context"
	"encoding/json"
	"sync"
	"time"
	"wallet-user-svc/internal/app/model/domain"
	"wallet-user-svc/internal/app/model/dto"
	"wallet-user-svc/internal/app/model/events"
	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"github.com/sirupsen/logrus"
)

type NotificationRepository interface {
	FindPendingEvents(ctx context.Context, eventName string, batchSize int) ([]*domain.NotificationEventLog, error)
	UpdateStatusSuccess(ctx context.Context, id string) error
}

type NotificationWorker struct {
	logger                   *logrus.Logger
	asyncQClient             *asynq.Client
	notificationEventLogRepo NotificationRepository
	ticker                   *time.Ticker
	wg                       *sync.WaitGroup
	interval                 time.Duration
	maxRetries               int
	batchSize                int
	shutdownChan             chan struct{}
	shutdownOnce             sync.Once
}

func NewNotificationWorker(
	logger *logrus.Logger,
	asyncQClient *asynq.Client,
	notificationEventLogRepo NotificationRepository,
	wg *sync.WaitGroup,
	interval time.Duration,
	maxRetries int,
	batchSize int,
) *NotificationWorker {
	ticker := time.NewTicker(interval)

	return &NotificationWorker{
		logger:                   logger,
		asyncQClient:             asyncQClient,
		notificationEventLogRepo: notificationEventLogRepo,
		interval:                 interval,
		ticker:                   ticker,
		wg:                       wg,
		maxRetries:               maxRetries,
		batchSize:                batchSize,
		shutdownChan:             make(chan struct{}),
	}
}

func (s *NotificationWorker) Start(ctx context.Context) {
	s.logger.Info("Starting notification worker")

	s.wg.Add(1)
	go func() {
		defer func() {
			s.ticker.Stop()
			s.wg.Done()
			s.logger.Info("Notification worker stopped")
		}()

		// Process events immediately on startup
		s.processPendingLoginEvents(ctx)

		for {
			select {
			case <-ctx.Done():
				s.logger.Info("Stopping notification worker (context cancelled)")
				// Process any remaining events before stopping
				s.processRemainingEvents()
				return
			case <-s.shutdownChan:
				s.logger.Info("Stopping notification worker (shutdown signal)")
				// Process any remaining events before stopping
				s.processRemainingEvents()
				return
			case <-s.ticker.C:
				s.processPendingLoginEvents(ctx)
			}
		}
	}()
}

func (s *NotificationWorker) processRemainingEvents() {
	s.logger.Info("Processing remaining events before shutdown")

	// Use a background context with timeout for remaining event processing
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	s.processPendingLoginEvents(ctx)
}

func (s *NotificationWorker) processPendingLoginEvents(ctx context.Context) {
	s.logger.Debug("Processing pending login events")

	events, err := s.notificationEventLogRepo.FindPendingEvents(
		ctx,
		string(events.LoginEventType),
		s.batchSize,
	)
	if err != nil {
		s.logger.WithError(err).Error("Could not find pending events")
		return
	}

	if len(events) == 0 {
		s.logger.Debug("No pending events found")
		return
	}

	s.logger.WithField("count", len(events)).Info("Found pending events to process")

	// Process events sequentially in a single thread
	for _, event := range events {
		// Check for context cancellation before processing each event
		select {
		case <-ctx.Done():
			s.logger.Info("Context cancelled, stopping event processing")
			return
		default:
		}

		if err := s.processEvent(ctx, event); err != nil {
			s.logger.WithError(err).WithField("eventID", event.ID).Error("Failed to process event")
		}
	}

	s.logger.WithField("count", len(events)).Info("Processed pending events")
}

func (s *NotificationWorker) processEvent(ctx context.Context, event *domain.NotificationEventLog) error {
	var params dto.SendLoginNotificationParams
	if err := json.Unmarshal(event.Payload, &params); err != nil {
		s.logger.WithError(err).WithField("eventID", event.ID).Error("Could not unmarshal payload")
		return err
	}

	// Send notification
	if err := s.SendLoginNotification(ctx, &params); err != nil {
		s.logger.WithError(err).WithField("eventID", event.ID).Error("Failed to send login notification")
		return err
	}

	// Update status to success
	if err := s.notificationEventLogRepo.UpdateStatusSuccess(ctx, event.ID); err != nil {
		s.logger.WithError(err).WithField("eventID", event.ID).Error("Could not update status")
		return err
	}

	s.logger.WithField("eventID", event.ID).Debug("Event processed successfully")

	return nil
}

func (s *NotificationWorker) SendLoginNotification(
	ctx context.Context,
	params *dto.SendLoginNotificationParams,
) error {
	loginEvent := events.LoginEvent{
		EventMetadata: events.EventMetadata{
			EventID:   uuid.New().String(),
			EventName: string(events.LoginEventType),
		},
		UserID:   params.UserID,
		Email:    params.Email,
		Username: params.Username,
		LoginAt:  params.LoginAt,
	}

	task, err := loginEvent.ToTask()
	if err != nil {
		s.logger.WithError(err).Error("Could not create task")
		return err
	}

	info, err := s.asyncQClient.Enqueue(task, asynq.MaxRetry(s.maxRetries))
	if err != nil {
		s.logger.WithError(err).Error("Could not enqueue task")
		return err
	}

	s.logger.WithFields(logrus.Fields{
		"id":    info.ID,
		"queue": info.Queue,
	}).Debug("Enqueued task")

	return nil
}

// Stop gracefully stops the worker
func (s *NotificationWorker) Stop() {
	s.shutdownOnce.Do(func() {
		close(s.shutdownChan)
	})
}
