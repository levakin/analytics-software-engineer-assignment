package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/levakin/analytics-software-engineer-assignment/internal/github"
	"github.com/levakin/analytics-software-engineer-assignment/pkg/csvtargz"
)

const (
	projectName = "gh-analytics"

	actorsCSVFilename  = "data/actors.csv"
	commitsCSVFilename = "data/commits.csv"
	eventsCSVFilename  = "data/events.csv"
	reposCSVFilename   = "data/repos.csv"
)

func main() {
	app := &cli.App{
		Name: projectName,
		Commands: []*cli.Command{
			{
				Name:  "top-users",
				Usage: "Prints top N active users sorted by amount of PRs created and commits pushed",

				Action: func(ctx *cli.Context) error {
					return printTopNUsersByPRsCreatedAndCommitsPushed(ctx.String("p"), ctx.Int("n"))
				},
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:  "n",
						Value: 10,
						Usage: "top N",
					},
					&cli.StringFlag{
						Name:  "p",
						Value: "./data.tar.gz",
						Usage: "Path to data.tar.gz",
					},
				},
			},
			{
				Name:  "top-repos-by-commits",
				Usage: "Prints top N repositories sorted by amount of commits pushed",
				Action: func(ctx *cli.Context) error {
					return printTopNReposByPushedCommits(ctx.String("p"), ctx.Int("n"))
				},
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:  "n",
						Value: 10,
						Usage: "top N",
					},
					&cli.StringFlag{
						Name:  "p",
						Value: "./data.tar.gz",
						Usage: "Path to data.tar.gz",
					},
				},
			},
			{
				Name:  "top-repos-by-watch-events",
				Usage: "Prints top N repositories sorted by amount of watch events",
				Action: func(ctx *cli.Context) error {
					return printTopNReposByWatchEvents(ctx.String("p"), ctx.Int("n"))
				},
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:  "n",
						Value: 10,
						Usage: "top N",
					},
					&cli.StringFlag{
						Name:  "p",
						Value: "./data.tar.gz",
						Usage: "Path to data.tar.gz",
					},
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)

		os.Exit(1)
	}
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
