package job

import (
	"context"
	"fmt"
	"go-boilerplate/internal/config"
	"log/slog"
	"os"
	"strconv"
	"strings"

	"github.com/hibiken/asynq"
	"github.com/redis/go-redis/v9"
)

func parseQueues(queueStrings []string) (map[string]int, error) {
	queues := make(map[string]int)
	for _, q := range queueStrings {
		parts := strings.Split(q, ":")
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid queue format: %s", q)
		}
		priority, err := strconv.Atoi(parts[1])
		if err != nil {
			return nil, fmt.Errorf("invalid priority for queue %s: %w", parts[0], err)
		}
		queues[parts[0]] = priority
	}
	return queues, nil
}

type JobServiceLogger struct {
	logger *slog.Logger
}

func (tsl *JobServiceLogger) Debug(args ...interface{}) {
	tsl.logger.Debug(fmt.Sprint(args...))
}

func (tsl *JobServiceLogger) Info(args ...interface{}) {
	tsl.logger.Info(fmt.Sprint(args...))
}
func (tsl *JobServiceLogger) Warn(args ...interface{}) {
	tsl.logger.Warn(fmt.Sprint(args...))
}
func (tsl *JobServiceLogger) Error(args ...interface{}) {
	tsl.logger.Error(fmt.Sprint(args...))
}
func (tsl *JobServiceLogger) Fatal(args ...interface{}) {
	tsl.logger.Error(fmt.Sprint(args...))
	os.Exit(-1)
}

type JobService struct {
	Client      *asynq.Client
	server      *asynq.Server
	logger      *slog.Logger
	handlerCtx  *JobHandlerContext
	handlers    map[string]ContextualJobHandler
}

// NewJobService creates a new JobService.
func NewJobService(logger *slog.Logger, cfg *config.JobConfig, redisClient *redis.Client, handlerCtx *JobHandlerContext) (*JobService, error) {
	client := asynq.NewClientFromRedisClient(redisClient)
	queues, err := parseQueues(cfg.Queues)
	if err != nil {
		return nil, err
	}

	serviceLogger := logger.With("service", "job")
	internalLogger := &JobServiceLogger{logger: serviceLogger}
	server := asynq.NewServerFromRedisClient(
		redisClient,
		asynq.Config{
			Concurrency: cfg.Concurrency,
			ErrorHandler: asynq.ErrorHandlerFunc(func(ctx context.Context, task *asynq.Task, err error) {
				logger.ErrorContext(ctx, "asynq task failed", "type", task.Type(), "payload", string(task.Payload()), "error", err)
			}),
			Logger: internalLogger,
			Queues: queues,
		},
	)

	return &JobService{
		Client:      client,
		server:      server,
		logger:      serviceLogger,
		handlerCtx:  handlerCtx,
		handlers:    make(map[string]ContextualJobHandler),
	}, nil
}

func (j *JobService) Enqueue(ctx context.Context, task *asynq.Task) (*asynq.TaskInfo, error) {
	j.logger.InfoContext(ctx, "enqueuing task", "type", task.Type())
	return j.Client.Enqueue(task)
}

// ContextualJobHandler is a job handler function that receives a JobHandlerContext.
type ContextualJobHandler func(context.Context, *JobHandlerContext, *asynq.Task) error

// RegisterHandler registers a contextual job handler for a given task type.
func (j *JobService) RegisterHandler(taskType string, handler ContextualJobHandler) {
	_ = j.handlers[taskType]
	j.handlers[taskType] = handler
}

func (j *JobService) setupHandlers() *asynq.ServeMux {
	mux := asynq.NewServeMux()

	for taskType, handler := range j.handlers {
		// Create a closure that captures the handler and the context.
		closure := func(h ContextualJobHandler) asynq.HandlerFunc {
			return func(ctx context.Context, t *asynq.Task) error {
				return h(ctx, j.handlerCtx, t)
			}
		}
		mux.HandleFunc(taskType, closure(handler))
	}

	return mux
}

func (j *JobService) Start() error {
	mux := j.setupHandlers()

	j.logger.Info("Starting background job server")
	if err := j.server.Start(mux); err != nil {
		return err
	}

	return nil
}

func (j *JobService) Stop() {
	j.logger.Info("Stopping background job server")
	j.server.Shutdown()
	j.Client.Close()
}
