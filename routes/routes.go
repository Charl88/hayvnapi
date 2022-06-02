package routes

import (
	"context"
	"errors"
	"log"

	"github.com/Charl88/hayvnapi/shared"
	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"
)

func MessageReceiver(messageArray *shared.MessageArray) usecase.IOInteractor {

	type messagesInput struct {
		Destination string `json:"destination" description:"The destination channel for the message" required:"true" example:"compliance"`
		Text        string `json:"text" description:"The message in string format" required:"true" example:"An important event has occurred"`
		Timestamp   string `json:"timestamp" description:"The timestamp of the event in string format" required:"true" example:"2021-01-01T12:00:00.000Z"`
	}

	u := usecase.NewIOI(new(messagesInput), nil, func(ctx context.Context, input, output interface{}) error {
		var (
			in = input.(*messagesInput)
		)

		message := shared.Message{
			Destination: in.Destination,
			Text:        in.Text,
			Timestamp:   in.Timestamp,
		}

		messageArray.Messages = append(messageArray.Messages, message)

		return nil
	})

	u.SetExpectedErrors(status.Unknown, status.Internal)
	u.SetDescription("Will submit a message for batching to a specified destination")

	return u
}

func AggregatedMessages() usecase.IOInteractor {

	type aggregatedInput struct {
		Batches []shared.Batch `json:"batches" description:"Batches of messages to send to specific destinations" required:"true"`
	}

	u := usecase.NewIOI(new(aggregatedInput), nil, func(ctx context.Context, input, output interface{}) error {
		var (
			in = input.(*aggregatedInput)
		)
		process := make(map[string][]shared.BatchMessage)
		for i := 0; i < len(in.Batches); i++ {
			batch := in.Batches[i]
			_, ok := process[batch.Destination]
			if !ok {
				process[batch.Destination] = batch.Messages
				log.Printf("Aggregated Messages - sending to %s - %s", batch.Destination, batch.Messages)
				// Send the messages to the required destination (maybe via GRPC?)
			} else {
				return status.Wrap(errors.New("multiple batches contained the same destination"), status.InvalidArgument)
			}
		}
		return nil
	})

	u.SetExpectedErrors(status.Unknown, status.Internal)
	u.SetDescription("Will submit a message for batching to a specified destination")

	return u
}
