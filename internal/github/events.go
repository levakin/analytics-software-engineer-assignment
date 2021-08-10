package github

const (
	pullRequestEventType = "PullRequestEvent"
	pushEventType        = "PushEvent"
)

type Event struct {
	ID      string `csv:"id"`
	Type    string `csv:"type"`
	ActorID string `csv:"actor_id"`
	RepoID  string `csv:"repo_id"`
}
