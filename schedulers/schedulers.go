package schedulers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	"log"

	"github.com/Charl88/hayvnapi/services"
	"github.com/Charl88/hayvnapi/shared"
	"github.com/go-co-op/gocron"
)

func CreateMessageScheduler(messageQueue *shared.MessageArray) {

	// Set a scheduler to run every 10 seconds to send batched messages to
	// the aggregated-messages API endpoint.
	s := gocron.NewScheduler(time.UTC)
	// Max concurrent jobs is set to 1 to prevent two batches from being sent
	// at the same time.
	s.SetMaxConcurrentJobs(1, gocron.RescheduleMode)
	s.Every(10).Seconds().Do(func() {
		batchArray, err := services.AggregateMessages(messageQueue)
		if err != nil {
			log.Printf("Scheduler - Could not aggregate messages - %s", err)
		}

		// Clear the message queue soon as we're done processing them
		(*messageQueue) = shared.MessageArray{}

		if len(batchArray.Batches) > 0 {
			postBody, err := json.Marshal(batchArray)
			if err != nil {
				log.Printf("Scheduler - Could not marshal batched messages to json - %s", err)
			}

			body := bytes.NewBuffer(postBody)

			resp, err := http.Post("http://localhost:3000/aggregated-messages", "application/json", body)
			if err != nil {
				log.Printf("Scheduler - Could perform the post request to https://localhost:3000/aggregated-messages - %s", err)
			}
			defer resp.Body.Close()
		}
	})
	s.StartAsync()
}
