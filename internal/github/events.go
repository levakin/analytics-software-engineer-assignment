package github

const (
	// PullRequestEventType is pull request event type.
	PullRequestEventType = "PullRequestEvent"
	// PushEventType is created every time a user pushes commits. Event in this case is connected with many commits.
	PushEventType = "PushEvent"
)

// Event represent event from GitHub
type Event struct {
	ID      string `csv:"id"`
	Type    string `csv:"type"`
	ActorID string `csv:"actor_id"`
	RepoID  string `csv:"repo_id"`
}
