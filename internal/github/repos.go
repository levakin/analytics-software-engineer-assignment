package github

import (
	"sort"

	"github.com/pkg/errors"
)

// RepoCSV represents GitHub repository in CSV
type RepoCSV struct {
	ID   string `csv:"id"`
	Name string `csv:"name"`
}

// Repo represents GitHub repository with statistics
type Repo struct {
	ID            string
	Name          string
	CommitsPushed int
	WatchEvents   int
}

// ReposSample is GitHub repositories collection sample used for getting analytics reposts
type ReposSample struct {
	M map[string]Repo
}

// NewReposSample returns a new ReposSample
func NewReposSample(events []EventCSV, repoCSVs []RepoCSV, commits []CommitCSV) *ReposSample {
	repos := ReposSample{
		M: make(map[string]Repo),
	}

	// set data about repository
	for _, repoCSV := range repoCSVs {
		repos.M[repoCSV.ID] = Repo{
			ID:   repoCSV.ID,
			Name: repoCSV.Name,
		}
	}

	repos.setAnalyticsData(events, commits)

	return &repos
}

// TopNByCommitsPushed returns top N repositories sorted by amount of commits pushed.
func (rs *ReposSample) TopNByCommitsPushed(n int) ([]Repo, error) {
	if n < 1 {
		return nil, errors.Wrap(ErrWrongParam, "n should be at least 1")
	}

	// TODO: optimize to store only top N
	sortedRepos := make([]Repo, 0, len(rs.M))
	for _, user := range rs.M {
		sortedRepos = append(sortedRepos, user)
	}

	sort.Slice(sortedRepos, func(i, j int) bool { return sortedRepos[i].CommitsPushed > sortedRepos[j].CommitsPushed })

	var last int
	if len(sortedRepos) < n {
		last = len(sortedRepos)
	} else {
		last = n
	}

	return sortedRepos[0:last], nil
}

// TopNByWatchEvents returns top N repositories sorted by amount of watch events.
func (rs *ReposSample) TopNByWatchEvents(n int) ([]Repo, error) {
	if n < 1 {
		return nil, errors.Wrap(ErrWrongParam, "n should be at least 1")
	}

	// TODO: optimize to store only top N
	sortedRepos := make([]Repo, 0, len(rs.M))
	for _, user := range rs.M {
		sortedRepos = append(sortedRepos, user)
	}

	sort.Slice(sortedRepos, func(i, j int) bool { return sortedRepos[i].WatchEvents > sortedRepos[j].WatchEvents })

	var last int
	if len(sortedRepos) < n {
		last = len(sortedRepos)
	} else {
		last = n
	}

	return sortedRepos[0:last], nil
}

func (rs *ReposSample) setAnalyticsData(events []EventCSV, commits []CommitCSV) {
	numCommitsByPushEventID := newNumCommitsByPushEventID(commits)
	for _, e := range events {
		switch e.Type {
		case PushEventType:
			count := numCommitsByPushEventID[e.ID]
			if count <= 0 {
				continue
			}

			r := rs.M[e.RepoID]
			r.CommitsPushed += count
			rs.M[e.RepoID] = r

		case WatchEventType:
			r := rs.M[e.RepoID]
			r.WatchEvents++
			rs.M[e.RepoID] = r
		}
	}
}
