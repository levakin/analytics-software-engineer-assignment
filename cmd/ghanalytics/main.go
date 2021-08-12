package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v2"
	"golang.org/x/sync/errgroup"

	"github.com/levakin/analytics-software-engineer-assignment/internal/github"
	"github.com/levakin/analytics-software-engineer-assignment/pkg/csvtargz"
)

const (
	projectName = "gh-analytics"
)

func main() {
	app := &cli.App{
		Name:        projectName,
		Description: "A cli application that can be used to get analytics information out of archive with csv files with GitHub data",
		Commands: []*cli.Command{
			{
				Name:  "top-users",
				Usage: "Prints top N active users sorted by amount of PRs created and commits pushed",

				Action: func(ctx *cli.Context) error {
					return printTopNUsersByPRsCreatedAndCommitsPushed(ctx.Context, ctx.String("p"), ctx.Int("n"), ctx.Bool("bots"))
				},
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:  "n",
						Value: 10,
						Usage: "top N",
					},
					&cli.StringFlag{
						Name:  "p",
						Value: "./samples/data.tar.gz",
						Usage: "Path to data.tar.gz",
					},
					&cli.BoolFlag{
						Name:  "bots",
						Usage: "If flag is set, bots will be included in the report",
					},
				},
			},
			{
				Name:  "top-repos-by-commits",
				Usage: "Prints top N repositories sorted by amount of commits pushed",
				Action: func(ctx *cli.Context) error {
					return printTopNReposByPushedCommits(ctx.Context, ctx.String("p"), ctx.Int("n"))
				},
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:  "n",
						Value: 10,
						Usage: "top N",
					},
					&cli.StringFlag{
						Name:  "p",
						Value: "./samples/data.tar.gz",
						Usage: "Path to data.tar.gz",
					},
				},
			},
			{
				Name:  "top-repos-by-watch-events",
				Usage: "Prints top N repositories sorted by amount of watch events",
				Action: func(ctx *cli.Context) error {
					return printTopNReposByWatchEvents(ctx.Context, ctx.String("p"), ctx.Int("n"))
				},
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:  "n",
						Value: 10,
						Usage: "top N",
					},
					&cli.StringFlag{
						Name:  "p",
						Value: "./samples/data.tar.gz",
						Usage: "Path to data.tar.gz",
					},
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func printTopNUsersByPRsCreatedAndCommitsPushed(ctx context.Context, archivePath string, n int, botsIncluded bool) error {
	var (
		actors  []github.ActorCSV
		commits []github.CommitCSV
		events  []github.EventCSV
	)

	g, _ := errgroup.WithContext(ctx)

	g.Go(func() error {
		return csvtargz.DecodeByPath(archivePath, github.ActorsCSVFilename, &actors)
	})

	g.Go(func() error {
		return csvtargz.DecodeByPath(archivePath, github.EventsCSVFilename, &events)
	})

	g.Go(func() error {
		return csvtargz.DecodeByPath(archivePath, github.CommitsCSVFilename, &commits)
	})

	if err := g.Wait(); err != nil {
		return err
	}

	users := github.NewUsersSample(actors, commits, events, botsIncluded)

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

func printTopNReposByPushedCommits(ctx context.Context, archivePath string, n int) error {
	var (
		commits  []github.CommitCSV
		events   []github.EventCSV
		repoCSVs []github.RepoCSV
	)

	g, _ := errgroup.WithContext(ctx)

	g.Go(func() error {
		return csvtargz.DecodeByPath(archivePath, github.ReposCSVFilename, &repoCSVs)
	})

	g.Go(func() error {
		return csvtargz.DecodeByPath(archivePath, github.EventsCSVFilename, &events)
	})

	g.Go(func() error {
		return csvtargz.DecodeByPath(archivePath, github.CommitsCSVFilename, &commits)
	})

	if err := g.Wait(); err != nil {
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

func printTopNReposByWatchEvents(ctx context.Context, archivePath string, n int) error {
	var (
		commits  []github.CommitCSV
		events   []github.EventCSV
		repoCSVs []github.RepoCSV
	)

	g, _ := errgroup.WithContext(ctx)

	g.Go(func() error {
		return csvtargz.DecodeByPath(archivePath, github.ReposCSVFilename, &repoCSVs)
	})

	g.Go(func() error {
		return csvtargz.DecodeByPath(archivePath, github.EventsCSVFilename, &events)
	})

	g.Go(func() error {
		return csvtargz.DecodeByPath(archivePath, github.CommitsCSVFilename, &commits)
	})

	if err := g.Wait(); err != nil {
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
