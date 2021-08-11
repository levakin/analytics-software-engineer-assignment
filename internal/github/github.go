// Package github implements functionality for getting analytics data about users' activity and  about repositories.
// Main features are:
// - Top N active users sorted by amount of PRs created and commits pushed
// - Top N repositories sorted by amount of commits pushed
// - Top N repositories sorted by amount of watch events
package github

import "github.com/pkg/errors"

// ErrWrongParam is returned if wrong parameter is passed.
var ErrWrongParam = errors.New("wrong parameter")
