package github

type Actor struct {
	ID       string `csv:"id"`
	Username string `csv:"username"`
}

// TopNActiveActors finds the top N active actors.
// Activity for each actor is a sum of all pushed commits and created pull requests.
func TopNActiveActors(n int, actors []Actor, commits []Commit, events []Event) ([]Actor, error) {
	return actors[0:n], nil
}

func calculateActorActivity(actorID string, commits []Commit, events []Event) int {
	return countActorsPushedCommits(actorID, commits, events) + countActorsCreatedPullRequests(actorID, events)
}

func countActorsPushedCommits(actorID string, commits []Commit, events []Event) int {
	var count int

	for _, e := range events {
		if e.ActorID != actorID || e.Type != pushEventType {
			continue
		}

		for _, c := range commits {
			if c.EventID != e.ID {
				continue
			}

			count++
		}
	}

	return count
}

func countActorsCreatedPullRequests(actorID string, events []Event) int {
	var count int

	for _, e := range events {
		if e.ActorID != actorID || e.Type != pullRequestEventType {
			continue
		}

		count++
	}

	return count
}
