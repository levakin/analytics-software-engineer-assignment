package github

// CommitCSV represents GitHub commit in CSV
type CommitCSV struct {
	SHA     string `csv:"sha"`
	Message string `csv:"message"`
	EventID string `csv:"event_id"`
}

func newNumCommitsByPushEventID(commits []CommitCSV) map[string]int {
	numCommitsByPushEventID := make(map[string]int)

	for _, c := range commits {
		numCommitsByPushEventID[c.EventID]++
	}

	return numCommitsByPushEventID
}
