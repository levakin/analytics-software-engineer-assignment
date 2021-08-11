package main

import (
	"fmt"
	"log"

	"github.com/levakin/analytics-software-engineer-assignment/internal/github"
	"github.com/levakin/analytics-software-engineer-assignment/pkg/csvtargz"
)

const (
	actorsCSVFilename  = "data/actors.csv"
	commitsCSVFilename = "data/commits.csv"
	eventsCSVFilename  = "data/events.csv"
	reposCSVFilename   = "data/repos.csv"
)

// - Top 10 active users sorted by amount of PRs created and commits pushed
// - Top 10 repositories sorted by amount of commits pushed
// - Top 10 repositories sorted by amount of watch events
func main() {
	archivePath := "/Users/anton/Git/github.com/levakin/analytics-software-engineer-assignment/data.tar.gz"

	if err := run(archivePath); err != nil {
		log.Fatal(err)
	}
}

func run(archivePath string) error {
	var actors []github.ActorCSV
	if err := csvtargz.DecodeCSVFromTarGz(archivePath, actorsCSVFilename, &actors); err != nil {
		return err
	}

	var events []github.EventCSV
	if err := csvtargz.DecodeCSVFromTarGz(archivePath, eventsCSVFilename, &events); err != nil {
		return err
	}

	var commits []github.CommitCSV
	if err := csvtargz.DecodeCSVFromTarGz(archivePath, commitsCSVFilename, &commits); err != nil {
		return err
	}

	users := github.NewUsersSample(actors, commits, events)

	topUsers, err := users.TopNActiveUsers(10)
	if err != nil {
		return err
	}

	fmt.Println("top 10 active users:")
	for i, u := range topUsers {
		fmt.Printf(
			"%3d. username: %30s| id: %10s| activity: %10d| pushed commits: %5d| created pull requests: %5d|\n",
			i+1, u.Username, u.ID, u.Activity.Total(), u.Activity.PushedCommits, u.Activity.CreatedPullRequests,
		)
	}
	fmt.Println("-------------------------------------------")

	var repoCSVs []github.RepoCSV
	if err := csvtargz.DecodeCSVFromTarGz(archivePath, reposCSVFilename, &repoCSVs); err != nil {
		return err
	}

	repos := github.NewReposSample(events, repoCSVs, commits)

	topReposByPushedCommits, err := repos.TopNByCommitsPushed(10)
	if err != nil {
		return err
	}

	fmt.Println("top 10 repositories by pushed commits:")
	for i, repo := range topReposByPushedCommits {
		fmt.Printf(
			"%3d. name: %50s| id: %10s| commits pushed: %5d|\n",
			i+1, repo.Name, repo.ID, repo.CommitsPushed,
		)
	}
	fmt.Println("-------------------------------------------")

	topReposByWatchEvents, err := repos.TopNByWatchEvents(10)
	if err != nil {
		return err
	}

	fmt.Println("top 10 repositories by watch events:")
	for i, repo := range topReposByWatchEvents {
		fmt.Printf(
			"%3d. name: %50s| id: %10s| watch events: %5d|\n",
			i+1, repo.Name, repo.ID, repo.WatchEvents,
		)
	}

	return nil
}
