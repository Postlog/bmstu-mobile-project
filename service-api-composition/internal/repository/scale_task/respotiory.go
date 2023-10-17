package scale_task

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	scaleTasksQueueName         = "scale_image_tasks"
	scaleTasksExchangeName      = "scale_image_tasks_exchange"
	scaleTasksBindingRoutingKey = "scale_image_tasks"
	contentTypeJSON             = "application/json"

	consumeChannelWaitDuration = time.Millisecond * 50
)

type Repository struct {
	conn *amqp.Connection
}

func New(conn *amqp.Connection) *Repository {
	return &Repository{
		conn: conn,
	}
}

func (r Repository) Save(ctx context.Context, tasks []ScaleTask) error {
	ch, err := r.conn.Channel()
	if err != nil {
		return fmt.Errorf("get channel: %w", err)
	}
	defer func() { _ = ch.Close() }()

	eg, egCtx := errgroup.WithContext(ctx)

	for _, task := range tasks {
		taskID := task.ID
		bytes, err := json.Marshal(task)
		if err != nil {
			return fmt.Errorf("error marshalling task to json: %w", err)
		}

		eg.Go(func() error {
			message := amqp.Publishing{
				Headers:      nil,
				ContentType:  contentTypeJSON,
				DeliveryMode: amqp.Transient,
				MessageId:    taskID,
				Timestamp:    time.Now(),
				Body:         bytes,
			}

			return ch.PublishWithContext(
				egCtx,
				scaleTasksExchangeName,
				scaleTasksBindingRoutingKey,
				false, // mandatory,
				false, // immediate
				message,
			)
		})
	}

	if err = eg.Wait(); err != nil {
		return fmt.Errorf("error publishing tasks: %w", err)
	}

	return nil
}

func (r Repository) Get(ctx context.Context, batchSize int, processCallback func(context.Context, []ScaleTask) error) error {
	ch, err := r.conn.Channel()
	if err != nil {
		return fmt.Errorf("get channel: %w", err)
	}
	defer func() { _ = ch.Close() }()

	if err = ch.Qos(batchSize, 0, false); err != nil {
		return fmt.Errorf("error setting qos: %w", err)
	}

	consumeCh, err := ch.ConsumeWithContext(
		ctx,
		scaleTasksQueueName,
		uuid.New().String(),
		false, // autoAck
		false, // exclusive
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("error consuming tasks: %w", err)
	}

	tasksBuffer := make([]ScaleTask, 0, batchSize)
	var lastDelivery amqp.Delivery
	t := time.NewTimer(consumeChannelWaitDuration)
	defer func() { _ = t.Stop() }()

	for {
		select {
		case <-ctx.Done():
			_ = lastDelivery.Nack(true, true)
			return ctx.Err()
		case delivery, isOpened := <-consumeCh:
			if !isOpened {
				if len(tasksBuffer) != 0 {
					processTasks(ctx, lastDelivery, tasksBuffer, processCallback)
				}

				return nil
			}

			if len(tasksBuffer) >= batchSize {
				processTasks(ctx, lastDelivery, tasksBuffer, processCallback)

				tasksBuffer = tasksBuffer[:0]

				break
			}

			var scaleTask ScaleTask
			err = json.Unmarshal(delivery.Body, &scaleTask)
			if err != nil {
				return fmt.Errorf("error unmarshalling task: %w", err)
			}

			tasksBuffer = append(tasksBuffer, scaleTask)
			lastDelivery = delivery
		case <-t.C:
			if len(tasksBuffer) == 0 {
				break
			}

			processTasks(ctx, lastDelivery, tasksBuffer, processCallback)

			tasksBuffer = tasksBuffer[:0]
		}

		t.Reset(consumeChannelWaitDuration)
	}
}

func processTasks(
	ctx context.Context,
	lastDelivery amqp.Delivery,
	tasks []ScaleTask,
	processCallback func(context.Context, []ScaleTask) error,
) {
	err := processCallback(ctx, tasks)

	if err != nil {
		_ = lastDelivery.Nack(true, true)
	} else {
		_ = lastDelivery.Ack(true)
	}
}
