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
	n := 10
	if err := run(archivePath, n); err != nil {
		log.Fatal(err)
	}
}

func run(archivePath string, n int) error {
	printTopNUsersByPRsCreatedAndCommitsPushed(archivePath, n)
	printTopNReposByPushedCommits(archivePath, n)
	printTopNReposByWatchEvents(archivePath, n)
	return nil
}

func printTopNUsersByPRsCreatedAndCommitsPushed(archivePath string, n int) error {
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

	topUsers, err := users.TopNActiveUsers(n)
	if err != nil {
		return err
	}

	fmt.Printf("top %d active users:\n", n)
	for i, u := range topUsers {
		fmt.Printf(
			"%3d. username: %30s| id: %10s| activity: %10d| pushed commits: %5d| created pull requests: %5d|\n",
			i+1, u.Username, u.ID, u.Activity.Total(), u.Activity.PushedCommits, u.Activity.CreatedPullRequests,
		)
	}

	return nil
}

func printTopNReposByPushedCommits(archivePath string, n int) error {
	var events []github.EventCSV
	if err := csvtargz.DecodeCSVFromTarGz(archivePath, eventsCSVFilename, &events); err != nil {
		return err
	}

	var commits []github.CommitCSV
	if err := csvtargz.DecodeCSVFromTarGz(archivePath, commitsCSVFilename, &commits); err != nil {
		return err
	}

	var repoCSVs []github.RepoCSV
	if err := csvtargz.DecodeCSVFromTarGz(archivePath, reposCSVFilename, &repoCSVs); err != nil {
		return err
	}

	repos := github.NewReposSample(events, repoCSVs, commits)

	topReposByPushedCommits, err := repos.TopNByCommitsPushed(n)
	if err != nil {
		return err
	}

	fmt.Printf("top %d repositories by pushed commits:\n", n)
	for i, repo := range topReposByPushedCommits {
		fmt.Printf(
			"%3d. name: %50s| id: %10s| commits pushed: %5d|\n",
			i+1, repo.Name, repo.ID, repo.CommitsPushed,
		)
	}

	return nil
}

func printTopNReposByWatchEvents(archivePath string, n int) error {
	var events []github.EventCSV
	if err := csvtargz.DecodeCSVFromTarGz(archivePath, eventsCSVFilename, &events); err != nil {
		return err
	}

	var commits []github.CommitCSV
	if err := csvtargz.DecodeCSVFromTarGz(archivePath, commitsCSVFilename, &commits); err != nil {
		return err
	}

	var repoCSVs []github.RepoCSV
	if err := csvtargz.DecodeCSVFromTarGz(archivePath, reposCSVFilename, &repoCSVs); err != nil {
		return err
	}

	repos := github.NewReposSample(events, repoCSVs, commits)

	topReposByWatchEvents, err := repos.TopNByWatchEvents(10)
	if err != nil {
		return err
	}

	fmt.Printf("top %d repositories by watch events:\n", n)
	for i, repo := range topReposByWatchEvents {
		fmt.Printf(
			"%3d. name: %50s| id: %10s| watch events: %5d|\n",
			i+1, repo.Name, repo.ID, repo.WatchEvents,
		)
	}

	return nil
}
