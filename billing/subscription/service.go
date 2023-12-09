package subscription

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/robfig/cron/v3"

	grpczap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"github.com/raystack/frontier/pkg/debounce"
	"go.uber.org/zap"

	"github.com/raystack/frontier/billing/plan"

	"github.com/raystack/frontier/billing/customer"
	"github.com/stripe/stripe-go/v75"
	"github.com/stripe/stripe-go/v75/client"
)

const (
	SyncDelay = time.Second * 60
)

type Repository interface {
	GetByID(ctx context.Context, id string) (Subscription, error)
	Create(ctx context.Context, subs Subscription) (Subscription, error)
	UpdateByID(ctx context.Context, subs Subscription) (Subscription, error)
	List(ctx context.Context, filter Filter) ([]Subscription, error)
	GetByProviderID(ctx context.Context, id string) (Subscription, error)
}

type CustomerService interface {
	GetByID(ctx context.Context, id string) (customer.Customer, error)
	List(ctx context.Context, filter customer.Filter) ([]customer.Customer, error)
}

type PlanService interface {
	List(ctx context.Context, filter plan.Filter) ([]plan.Plan, error)
	GetByID(ctx context.Context, id string) (plan.Plan, error)
}

type Service struct {
	repository      Repository
	stripeClient    *client.API
	customerService CustomerService
	planService     PlanService

	syncLimiter *debounce.Limiter
	syncJob     *cron.Cron
	mu          sync.Mutex
}

func NewService(stripeClient *client.API, repository Repository,
	customerService CustomerService, planService PlanService) *Service {
	return &Service{
		stripeClient:    stripeClient,
		repository:      repository,
		customerService: customerService,
		planService:     planService,
		syncLimiter:     debounce.New(2 * time.Second),
	}
}

func (s *Service) Create(ctx context.Context, sub Subscription) (Subscription, error) {
	return s.repository.Create(ctx, sub)
}

func (s *Service) GetByID(ctx context.Context, id string) (Subscription, error) {
	return s.repository.GetByID(ctx, id)
}

func (s *Service) GetByProviderID(ctx context.Context, id string) (Subscription, error) {
	return s.repository.GetByProviderID(ctx, id)
}

func (s *Service) Init(ctx context.Context) {
	if s.syncJob != nil {
		s.syncJob.Stop()
	}

	s.syncJob = cron.New()
	s.syncJob.AddFunc(fmt.Sprintf("@every %s", SyncDelay.String()), func() {
		s.backgroundSync(ctx)
	})
	s.syncJob.Start()
}

func (s *Service) Close() error {
	if s.syncJob != nil {
		return s.syncJob.Stop().Err()
	}
	return nil
}

func (s *Service) backgroundSync(ctx context.Context) {
	logger := grpczap.Extract(ctx)
	customers, err := s.customerService.List(ctx, customer.Filter{})
	if err != nil {
		logger.Error("checkout.backgroundSync", zap.Error(err))
		return
	}

	for _, customer := range customers {
		if customer.DeletedAt != nil || customer.ProviderID == "" {
			continue
		}
		if err := s.SyncWithProvider(ctx, customer.ID); err != nil {
			logger.Error("checkout.SyncWithProvider", zap.Error(err))
		}
		time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
	}
}

// SyncWithProvider syncs the subscription state with the billing provider
func (s *Service) SyncWithProvider(ctx context.Context, customerID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	subs, err := s.repository.List(ctx, Filter{
		CustomerID: customerID,
	})
	if err != nil {
		return err
	}

	for _, sub := range subs {
		stripeSubscription, err := s.stripeClient.Subscriptions.Get(sub.ProviderID, &stripe.SubscriptionParams{
			Params: stripe.Params{
				Context: ctx,
			},
		})
		if err != nil {
			return fmt.Errorf("failed to get subscription from billing provider: %w", err)
		}
		sub.State = string(stripeSubscription.Status)
		if stripeSubscription.CanceledAt > 0 {
			sub.CanceledAt = time.Unix(stripeSubscription.CanceledAt, 0)
		}
		if stripeSubscription.EndedAt > 0 {
			sub.EndedAt = time.Unix(stripeSubscription.EndedAt, 0)
		}
		if _, err := s.repository.UpdateByID(ctx, sub); err != nil {
			return err
		}
	}

	return nil
}

func (s *Service) Cancel(ctx context.Context, id string) (Subscription, error) {
	sub, err := s.GetByID(ctx, id)
	if err != nil {
		return Subscription{}, err
	}
	if !sub.CanceledAt.IsZero() {
		// already canceled, no-op
		return sub, nil
	}

	stripeSubscription, err := s.stripeClient.Subscriptions.Cancel(sub.ProviderID, &stripe.SubscriptionCancelParams{
		Params: stripe.Params{
			Context: ctx,
		},
	})
	if err != nil {
		return Subscription{}, fmt.Errorf("failed to cancel subscription at billing provider: %w", err)
	}

	sub.State = string(stripeSubscription.Status)
	if stripeSubscription.CanceledAt > 0 {
		sub.CanceledAt = time.Unix(stripeSubscription.CanceledAt, 0)
	}
	return s.repository.UpdateByID(ctx, sub)
}

func (s *Service) List(ctx context.Context, filter Filter) ([]Subscription, error) {
	logger := grpczap.Extract(ctx)
	s.syncLimiter.Call(func() {
		// fix context as the List ctx will get cancelled after call finishes
		if err := s.SyncWithProvider(context.Background(), filter.CustomerID); err != nil {
			logger.Error("subscription.SyncWithProvider", zap.Error(err))
		}
	})

	return s.repository.List(ctx, filter)
}