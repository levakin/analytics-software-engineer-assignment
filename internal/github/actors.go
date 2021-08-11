package github

// Actor represents GitHub actor
type Actor struct {
	ID       string `csv:"id"`
	Username string `csv:"username"`
}

// ActorActivity represents GitHub actor activity
type ActorActivity struct {
	PushedCommits       int
	CreatedPullRequests int
}

// Total calculates total activity of actor
func (a ActorActivity) Total() int {
	return a.PushedCommits + a.CreatedPullRequests
}

func newActorActivityByActorID(commits []Commit, events []Event) map[string]ActorActivity {
	numCommitsByPushEventID := newNumCommitsByPushEventID(commits)

	m := make(map[string]ActorActivity)
	for _, e := range events {
		switch e.Type {
		case PullRequestEventType:
			if a, ok := m[e.ActorID]; ok {
				a.CreatedPullRequests++
				m[e.ActorID] = a
			} else {
				m[e.ActorID] = ActorActivity{
					CreatedPullRequests: 1,
				}
			}

		case PushEventType:
			if a, ok := m[e.ActorID]; ok {
				a.PushedCommits += numCommitsByPushEventID[e.ID]
				m[e.ActorID] = a
			} else {
				m[e.ActorID] = ActorActivity{
					PushedCommits: numCommitsByPushEventID[e.ID],
				}
			}
		}
	}
	return m
}

func newNumCommitsByPushEventID(commits []Commit) map[string]int {
	numCommitsByPushEventID := make(map[string]int)

	for _, c := range commits {
		numCommitsByPushEventID[c.EventID]++
	}

	return numCommitsByPushEventID
}
