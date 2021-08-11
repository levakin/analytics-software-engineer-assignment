package github_test

import (
	"reflect"
	"testing"

	"github.com/levakin/analytics-software-engineer-assignment/internal/github"
)

func TestNewReposSample(t *testing.T) {
	type args struct {
		events   []github.EventCSV
		repoCSVs []github.RepoCSV
		commits  []github.CommitCSV
	}
	tests := []struct {
		name string
		args args
		want *github.ReposSample
	}{
		{
			name: "should be 0 commits and watch events",
			args: args{
				events: []github.EventCSV{
					{ID: "1", Type: "other", ActorID: "1", RepoID: "1"},
				},
				repoCSVs: []github.RepoCSV{
					{ID: "1", Name: "1"},
				},
				commits: []github.CommitCSV{},
			},
			want: &github.ReposSample{M: map[string]github.Repo{
				"1": {
					ID:            "1",
					Name:          "1",
					CommitsPushed: 0,
					WatchEvents:   0,
				},
			}},
		},
		{
			name: "should be 2 commits and 0 watch events",
			args: args{
				events: []github.EventCSV{
					{ID: "1", Type: "other", ActorID: "1", RepoID: "1"},
					{ID: "2", Type: github.PushEventType, ActorID: "1", RepoID: "1"},

					{ID: "3", Type: github.WatchEventType, ActorID: "1", RepoID: "2"},
				},
				repoCSVs: []github.RepoCSV{
					{ID: "1", Name: "1"},
					{ID: "2", Name: "2"},
				},
				commits: []github.CommitCSV{
					{SHA: "1", Message: "msg", EventID: "2"},
					{SHA: "2", Message: "msg", EventID: "2"},
				},
			},
			want: &github.ReposSample{M: map[string]github.Repo{
				"1": {
					ID:            "1",
					Name:          "1",
					CommitsPushed: 2,
					WatchEvents:   0,
				},
				"2": {
					ID:            "2",
					Name:          "2",
					CommitsPushed: 0,
					WatchEvents:   1,
				},
			}},
		},
		{
			name: "should be 2 commits and 1 watch events",
			args: args{
				events: []github.EventCSV{
					{ID: "1", Type: "other", ActorID: "1", RepoID: "1"},
					{ID: "2", Type: github.WatchEventType, ActorID: "1", RepoID: "1"},
					{ID: "3", Type: github.PushEventType, ActorID: "1", RepoID: "1"},
				},
				repoCSVs: []github.RepoCSV{
					{ID: "1", Name: "1"},
				},
				commits: []github.CommitCSV{
					{SHA: "1", Message: "msg", EventID: "3"},
					{SHA: "2", Message: "msg", EventID: "3"},
				},
			},
			want: &github.ReposSample{M: map[string]github.Repo{
				"1": {
					ID:            "1",
					Name:          "1",
					CommitsPushed: 2,
					WatchEvents:   1,
				},
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := github.NewReposSample(tt.args.events, tt.args.repoCSVs, tt.args.commits); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewReposSample() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReposSample_TopNByCommitsPushed(t *testing.T) {
	type fields struct {
		M map[string]github.Repo
	}
	type args struct {
		n int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []github.Repo
		wantErr bool
	}{
		{
			name: "wrong input",
			args: args{
				n: 0,
			},
			wantErr: true,
		},
		{
			name: "top 1",
			fields: struct{ M map[string]github.Repo }{M: map[string]github.Repo{
				"1": {
					ID:            "1",
					Name:          "1",
					CommitsPushed: 2,
					WatchEvents:   0,
				},
				"2": {
					ID:            "2",
					Name:          "2",
					CommitsPushed: 0,
					WatchEvents:   1,
				},
			}},
			args: args{
				n: 1,
			},
			want: []github.Repo{
				{
					ID:            "1",
					Name:          "1",
					CommitsPushed: 2,
					WatchEvents:   0,
				},
			},
			wantErr: false,
		},
		{
			name: "top 3 when total repos less than n",
			fields: struct{ M map[string]github.Repo }{M: map[string]github.Repo{
				"1": {
					ID:            "1",
					Name:          "1",
					CommitsPushed: 2,
					WatchEvents:   0,
				},
				"2": {
					ID:            "2",
					Name:          "2",
					CommitsPushed: 0,
					WatchEvents:   1,
				},
			}},
			args: args{
				n: 3,
			},
			want: []github.Repo{
				{
					ID:            "1",
					Name:          "1",
					CommitsPushed: 2,
					WatchEvents:   0,
				},
				{
					ID:            "2",
					Name:          "2",
					CommitsPushed: 0,
					WatchEvents:   1,
				},
			},
			wantErr: false,
		},
		{
			name: "top 2 when total repos more than n",
			fields: struct{ M map[string]github.Repo }{M: map[string]github.Repo{
				"1": {
					ID:            "1",
					Name:          "1",
					CommitsPushed: 2,
					WatchEvents:   0,
				},
				"2": {
					ID:            "2",
					Name:          "2",
					CommitsPushed: 0,
					WatchEvents:   1,
				},
				"3": {
					ID:            "3",
					Name:          "3",
					CommitsPushed: 1,
					WatchEvents:   1,
				},
			}},
			args: args{
				n: 2,
			},
			want: []github.Repo{
				{
					ID:            "1",
					Name:          "1",
					CommitsPushed: 2,
					WatchEvents:   0,
				},
				{
					ID:            "3",
					Name:          "3",
					CommitsPushed: 1,
					WatchEvents:   1,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rs := &github.ReposSample{
				M: tt.fields.M,
			}
			got, err := rs.TopNByCommitsPushed(tt.args.n)
			if (err != nil) != tt.wantErr {
				t.Errorf("TopNByCommitsPushed() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TopNByCommitsPushed() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReposSample_TopNByWatchEvents(t *testing.T) {
	type fields struct {
		M map[string]github.Repo
	}
	type args struct {
		n int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []github.Repo
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rs := &github.ReposSample{
				M: tt.fields.M,
			}
			got, err := rs.TopNByWatchEvents(tt.args.n)
			if (err != nil) != tt.wantErr {
				t.Errorf("TopNByWatchEvents() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TopNByWatchEvents() got = %v, want %v", got, tt.want)
			}
		})
	}
}
