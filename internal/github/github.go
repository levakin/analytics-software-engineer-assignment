// Package github implements functionality for getting analytics data about users' activity and  about repositories.
// Main features are:
// - Top N active users sorted by amount of PRs created and commits pushed
// - Top N repositories sorted by amount of commits pushed
// - Top N repositories sorted by amount of watch events
package github

import "github.com/pkg/errors"

const (
	ActorsCSVFilename  = "data/actors.csv"
	CommitsCSVFilename = "data/commits.csv"
	EventsCSVFilename  = "data/events.csv"
	ReposCSVFilename   = "data/repos.csv"
)

// ErrWrongParam is returned if wrong parameter is passed.
var ErrWrongParam = errors.New("wrong parameter")
