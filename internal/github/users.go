package github

import (
	"regexp"
	"sort"

	"github.com/pkg/errors"
)

// User represents a GitHub user.
type User struct {
	ID       string
	Username string
	Activity ActorActivity
}

// UsersSample represents users collection
type UsersSample struct {
	M map[string]User
}

// NewUsersSample parses data and returns UsersSample collection.
// Bots with `botname[bot]` are not humans and could be filtered out.
func NewUsersSample(actors []ActorCSV, commits []CommitCSV, events []EventCSV, botsIncluded bool) *UsersSample {
	users := UsersSample{
		M: make(map[string]User),
	}

	numCommitsByPushEventID := newNumCommitsByPushEventID(commits)
	actorActivityByActorID := newActorActivityByActorID(numCommitsByPushEventID, events)

	for i := range actors {
		// if username like dependabot[bot] skip
		if !botsIncluded && isBotUsername(actors[i].Username) {
			continue
		}

		users.M[actors[i].ID] = User{
			ID:       actors[i].ID,
			Username: actors[i].Username,
			Activity: actorActivityByActorID[actors[i].ID],
		}
	}

	return &users
}

// TopNActiveUsers finds the top N active users.
// Activity for each user is a sum of all pushed commits and created pull requests.
func (us *UsersSample) TopNActiveUsers(n int) ([]User, error) {
	if n < 1 {
		return nil, errors.Wrap(ErrWrongParam, "n should be at least 1")
	}

	sortedUsers := make([]User, 0, n)

	for _, user := range us.M {
		if len(sortedUsers) < n {
			sortedUsers = append(sortedUsers, user)
			sort.Slice(sortedUsers, func(i, j int) bool { return sortedUsers[i].Activity.Total() > sortedUsers[j].Activity.Total() })

			continue
		}

		if sortedUsers[n-1].Activity.Total() > user.Activity.Total() {
			continue
		}

		sortedUsers[n-1] = user
		sort.Slice(sortedUsers, func(i, j int) bool { return sortedUsers[i].Activity.Total() > sortedUsers[j].Activity.Total() })
	}

	// don't take more items than in list
	var last int
	if len(sortedUsers) < n {
		last = len(sortedUsers)
	} else {
		last = n
	}

	return sortedUsers[0:last], nil
}

var botUsernameRegex = regexp.MustCompile(`^.*\[bot]$`)

func isBotUsername(username string) bool {
	return botUsernameRegex.MatchString(username)
}
