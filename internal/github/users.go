package github

import (
	"errors"
	"regexp"
	"sort"
)

// User represents a GitHub user. Bots with `botname[bot]` are not users and are filtered out.
type User struct {
	ID       string
	Username string
	Activity ActorActivity
}

// Users represents users collection
type Users struct {
	M map[string]User
}

// NewUsers parses data and returns Users collection.
func NewUsers(actors []Actor, commits []Commit, events []Event) *Users {
	users := Users{
		M: make(map[string]User),
	}

	actorActivityByActorID := newActorActivityByActorID(commits, events)
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
func (us *Users) TopNActiveUsers(n int) ([]User, error) {
	if n < 1 {
		return nil, errors.New("n should be at least 1")
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
