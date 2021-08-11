package github_test

import (
	"reflect"
	"testing"

	"github.com/levakin/analytics-software-engineer-assignment/internal/github"
)

func TestTopNActiveActors(t *testing.T) {
	type args struct {
		n       int
		actors  []github.Actor
		commits []github.Commit
		events  []github.Event
	}
	tests := []struct {
		name    string
		args    args
		want    []github.Actor
		wantErr bool
	}{
		{
			name: "should be error when err < 1",
			args: struct {
				n       int
				actors  []github.Actor
				commits []github.Commit
				events  []github.Event
			}{n: -1},
			wantErr: true,
		},
		{
			name: "should find top 1",
			args: struct {
				n       int
				actors  []github.Actor
				commits []github.Commit
				events  []github.Event
			}{
				n: 1,
				actors: []github.Actor{
					{ID: "1", Username: "1"},
					{ID: "2", Username: "2"},
					{ID: "3", Username: "3"},
				},
				commits: []github.Commit{
					{SHA: "sha1a", Message: "msg", EventID: "1"},
					{SHA: "sha1b", Message: "msg", EventID: "1"},

					{SHA: "sha2", Message: "msg", EventID: "2"},
					{SHA: "sha3", Message: "msg", EventID: "2"},

					{SHA: "sha4", Message: "msg", EventID: "4"},
					{SHA: "sha5", Message: "msg", EventID: "4"},
				},
				events: []github.Event{
					// actor 1 events (3 activity)
					{ID: "1", Type: github.PushEventType, ActorID: "1", RepoID: "1"},
					{ID: "2", Type: "other", ActorID: "1", RepoID: "1"},
					{ID: "3", Type: github.PullRequestEventType, ActorID: "1", RepoID: "1"},

					// actor 2 events (2 activity)
					{ID: "4", Type: github.PushEventType, ActorID: "2", RepoID: "1"},

					// actor 3 events (5 activity)
					{ID: "5", Type: github.PullRequestEventType, ActorID: "3", RepoID: "3"},
					{ID: "6", Type: github.PullRequestEventType, ActorID: "3", RepoID: "3"},
					{ID: "7", Type: github.PullRequestEventType, ActorID: "3", RepoID: "3"},
					{ID: "8", Type: github.PullRequestEventType, ActorID: "3", RepoID: "3"},
					{ID: "9", Type: github.PullRequestEventType, ActorID: "3", RepoID: "3"},
				},
			},
			want:    []github.Actor{{ID: "3", Username: "3"}},
			wantErr: false,
		},
		{
			name: "should find top 3",
			args: struct {
				n       int
				actors  []github.Actor
				commits []github.Commit
				events  []github.Event
			}{
				n: 3,
				actors: []github.Actor{
					{ID: "1", Username: "1"},
					{ID: "2", Username: "2"},
					{ID: "3", Username: "3"},
				},
				commits: []github.Commit{
					{SHA: "sha1a", Message: "msg", EventID: "1"},
					{SHA: "sha1b", Message: "msg", EventID: "1"},

					{SHA: "sha2", Message: "msg", EventID: "2"},
					{SHA: "sha3", Message: "msg", EventID: "2"},

					{SHA: "sha4", Message: "msg", EventID: "4"},
					{SHA: "sha5", Message: "msg", EventID: "4"},
				},
				events: []github.Event{
					// actor 1 events (3 activity)
					{ID: "1", Type: github.PushEventType, ActorID: "1", RepoID: "1"},
					{ID: "2", Type: "other", ActorID: "1", RepoID: "1"},
					{ID: "3", Type: github.PullRequestEventType, ActorID: "1", RepoID: "1"},

					// actor 2 events (2 activity)
					{ID: "4", Type: github.PushEventType, ActorID: "2", RepoID: "1"},

					// actor 3 events (5 activity)
					{ID: "5", Type: github.PullRequestEventType, ActorID: "3", RepoID: "3"},
					{ID: "6", Type: github.PullRequestEventType, ActorID: "3", RepoID: "3"},
					{ID: "7", Type: github.PullRequestEventType, ActorID: "3", RepoID: "3"},
					{ID: "8", Type: github.PullRequestEventType, ActorID: "3", RepoID: "3"},
					{ID: "9", Type: github.PullRequestEventType, ActorID: "3", RepoID: "3"},
				},
			},
			want:    []github.Actor{{ID: "3", Username: "3"}, {ID: "1", Username: "1"}, {ID: "2", Username: "2"}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := github.TopNActiveActors(tt.args.n, tt.args.actors, tt.args.commits, tt.args.events)
			if (err != nil) != tt.wantErr {
				t.Errorf("TopNActiveActors() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TopNActiveActors() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestActorActivity(t *testing.T) {
	type args struct {
		actorID string
		commits []github.Commit
		events  []github.Event
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "should be no activity by actor",
			args: struct {
				actorID string
				commits []github.Commit
				events  []github.Event
			}{
				actorID: "1",
				commits: []github.Commit{
					{SHA: "sha", Message: "msg", EventID: "1"},
				},
				events: []github.Event{
					{ID: "1", Type: github.PushEventType, ActorID: "2", RepoID: "1"},
				},
			},
			want: 0,
		},
		{
			name: "should calculate pushed commits",
			args: struct {
				actorID string
				commits []github.Commit
				events  []github.Event
			}{
				actorID: "1",
				commits: []github.Commit{
					{SHA: "sha1", Message: "msg", EventID: "1"},
					{SHA: "sha2", Message: "msg", EventID: "2"},
					{SHA: "sha3", Message: "msg", EventID: "2"},
				},
				events: []github.Event{
					{ID: "1", Type: github.PushEventType, ActorID: "2", RepoID: "1"},
					{ID: "2", Type: github.PushEventType, ActorID: "1", RepoID: "1"},
				},
			},
			want: 2,
		},
		{
			name: "should ignore events other than push and pull requests",
			args: struct {
				actorID string
				commits []github.Commit
				events  []github.Event
			}{
				actorID: "1",
				commits: []github.Commit{
					{SHA: "sha1", Message: "msg", EventID: "1"},
					{SHA: "sha2", Message: "msg", EventID: "2"},
					{SHA: "sha3", Message: "msg", EventID: "2"},
				},
				events: []github.Event{
					{ID: "1", Type: github.PushEventType, ActorID: "2", RepoID: "1"},
					{ID: "2", Type: github.PushEventType, ActorID: "1", RepoID: "1"},
					{ID: "3", Type: "other", ActorID: "1", RepoID: "1"},
					{ID: "4", Type: github.PullRequestEventType, ActorID: "1", RepoID: "1"},
				},
			},
			want: 3,
		},
		{
			name: "should calculate pull requests",
			args: struct {
				actorID string
				commits []github.Commit
				events  []github.Event
			}{
				actorID: "1",
				commits: []github.Commit{
					{SHA: "sha1", Message: "msg", EventID: "1"},
				},
				events: []github.Event{
					{ID: "1", Type: github.PushEventType, ActorID: "2", RepoID: "1"},
					{ID: "2", Type: github.PullRequestEventType, ActorID: "1", RepoID: "1"},
					{ID: "3", Type: github.PullRequestEventType, ActorID: "1", RepoID: "1"},
					{ID: "4", Type: github.PullRequestEventType, ActorID: "2", RepoID: "1"},
				},
			},
			want: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := github.ActorActivity(tt.args.actorID, tt.args.commits, tt.args.events); got != tt.want {
				t.Errorf("ActorActivity() = %v, want %v", got, tt.want)
			}
		})
	}
}
