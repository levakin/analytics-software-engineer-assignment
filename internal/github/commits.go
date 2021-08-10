package github

type Commit struct {
	SHA     string `csv:"sha"`
	Message string `csv:"message"`
	EventID string `csv:"event_id"`
}
