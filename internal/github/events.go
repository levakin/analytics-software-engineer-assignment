package github

type Event struct {
	ID      int    `csv:"id"`
	Type    string `csv:"type"`
	ActorID string `csv:"actor_id"`
	RepoID  string `csv:"repo_id"`
}
