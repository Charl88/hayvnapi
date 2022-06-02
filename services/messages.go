package services

import (
	"github.com/Charl88/hayvnapi/shared"
)

func AggregateMessages(messageQueue *shared.MessageArray) (shared.BatchArray, error) {
	batchMap := make(map[string][]shared.BatchMessage)

	// Iterate through the message queue to batch messages by destination
	for i := 0; i < len(messageQueue.Messages); i++ {
		message := messageQueue.Messages[i]
		_, ok := batchMap[message.Destination]
		if ok {
			// If the mapping to this destination already exists, append the message to the already
			// existing array
			batchMap[message.Destination] = append(batchMap[message.Destination], shared.BatchMessage{
				Text:      message.Text,
				Timestamp: message.Timestamp,
			})
		} else {
			// If the mapping to this destination doesn't exist, create a new array of messages
			// and append
			batchMap[message.Destination] = []shared.BatchMessage{}
			batchMap[message.Destination] = append(batchMap[message.Destination], shared.BatchMessage{
				Text:      message.Text,
				Timestamp: message.Timestamp,
			})
		}
	}

	// Construct the batched messages struct from the mapping we created above
	batches := []shared.Batch{}
	for key, value := range batchMap {
		batches = append(batches, shared.Batch{
			Destination: key,
			Messages:    value,
		})
	}

	return shared.BatchArray{
		Batches: batches,
	}, nil
}
