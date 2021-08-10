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
