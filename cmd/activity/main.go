package main

import (
	"fmt"
	"log"

	"github.com/levakin/analytics-software-engineer-assignment/internal/github"
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
	var actors []github.Actor
	if err := decodeCSVFromTarGz(archivePath, actorsCSVFilename, &actors); err != nil {
		return err
	}

	var events []github.Event
	if err := decodeCSVFromTarGz(archivePath, eventsCSVFilename, &events); err != nil {
		return err
	}

	var commits []github.Commit
	if err := decodeCSVFromTarGz(archivePath, commitsCSVFilename, &commits); err != nil {
		return err
	}

	topActors, err := github.TopNActiveActors(10, actors, commits, events)
	if err != nil {
		return err
	}

	fmt.Println(topActors)

	return nil
}
