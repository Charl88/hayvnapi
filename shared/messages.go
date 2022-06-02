package shared

type Message struct {
	Destination string
	Text        string
	Timestamp   string
}

type MessageArray struct {
	Messages []Message
}

type BatchMessage struct {
	Text      string `json:"text"`
	Timestamp string `json:"timestamp"`
}

type Batch struct {
	Destination string         `json:"destination"`
	Messages    []BatchMessage `json:"messages"`
}

type BatchArray struct {
	Batches []Batch `json:"batches"`
}
