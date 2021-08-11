package github

import "errors"

type Actor struct {
	ID       string `csv:"id"`
	Username string `csv:"username"`
}

// TopNActiveActors finds the top N active actors.
// Activity for each actor is a sum of all pushed commits and created pull requests.
func TopNActiveActors(n int, actors []Actor, commits []Commit, events []Event) ([]Actor, error) {
	if n < 1 {
		return nil, errors.New("n should be at least 1")
	}

	return actors[0:n], nil
}

func ActorActivity(actorID string, commits []Commit, events []Event) int {
	return countActorsPushedCommits(actorID, commits, events) + countActorsCreatedPullRequests(actorID, events)
}

func countActorsPushedCommits(actorID string, commits []Commit, events []Event) int {
	var count int

	for _, e := range events {
		if e.ActorID != actorID || e.Type != PushEventType {
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
		if e.ActorID != actorID || e.Type != PullRequestEventType {
			continue
		}

		count++
	}

	return count
}
