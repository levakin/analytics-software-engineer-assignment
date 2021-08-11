package github

import (
	"regexp"
	"sort"

	"github.com/pkg/errors"
)

// User represents a GitHub user. Bots with `botname[bot]` are not users and are filtered out.
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
func NewUsersSample(actors []ActorCSV, commits []CommitCSV, events []EventCSV) *UsersSample {
	users := UsersSample{
		M: make(map[string]User),
	}

	numCommitsByPushEventID := newNumCommitsByPushEventID(commits)
	actorActivityByActorID := newActorActivityByActorID(numCommitsByPushEventID, events)
	for i := range actors {
		// if username like dependabot[bot] skip
		if isBotUsername(actors[i].Username) {
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

	// TODO: optimize to store only top N
	allUsers := make([]User, 0, len(us.M))
	for _, user := range us.M {
		allUsers = append(allUsers, user)
	}

	sort.Slice(allUsers, func(i, j int) bool { return allUsers[i].Activity.Total() > allUsers[j].Activity.Total() })

	return allUsers[0:n], nil
}

var botUsernameRegex = regexp.MustCompile(`^.*\[bot]$`)

func isBotUsername(username string) bool {
	return botUsernameRegex.MatchString(username)
}
