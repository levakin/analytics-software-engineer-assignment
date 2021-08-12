package github_test

import (
	"reflect"
	"testing"

	"github.com/levakin/analytics-software-engineer-assignment/internal/github"
	"github.com/levakin/analytics-software-engineer-assignment/pkg/csvtargz"
	"github.com/levakin/analytics-software-engineer-assignment/samples"
)

const archivePath = "data.tar.gz"

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
					CommitsPushed: 0,
					WatchEvents:   2,
				},
				"2": {
					ID:            "2",
					Name:          "2",
					CommitsPushed: 1,
					WatchEvents:   0,
				},
			}},
			args: args{
				n: 1,
			},
			want: []github.Repo{
				{
					ID:            "1",
					Name:          "1",
					CommitsPushed: 0,
					WatchEvents:   2,
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
					CommitsPushed: 0,
					WatchEvents:   2,
				},
				"2": {
					ID:            "2",
					Name:          "2",
					CommitsPushed: 1,
					WatchEvents:   0,
				},
			}},
			args: args{
				n: 3,
			},
			want: []github.Repo{
				{
					ID:            "1",
					Name:          "1",
					CommitsPushed: 0,
					WatchEvents:   2,
				},
				{
					ID:            "2",
					Name:          "2",
					CommitsPushed: 1,
					WatchEvents:   0,
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
					CommitsPushed: 0,
					WatchEvents:   2,
				},
				"2": {
					ID:            "2",
					Name:          "2",
					CommitsPushed: 1,
					WatchEvents:   0,
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
					CommitsPushed: 0,
					WatchEvents:   2,
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

func BenchmarkReposSample_TopNByCommitsPushed(b *testing.B) {
	gzFile, err := samples.FS.Open(archivePath)
	if err != nil {
		b.Fatal(err)
	}

	defer func() {
		_ = gzFile.Close()
	}()

	var events []github.EventCSV
	if err := csvtargz.DecodeFromFile(gzFile, github.EventsCSVFilename, &events); err != nil {
		b.Fatal(err)
	}

	_ = gzFile.Close()

	gzFile, err = samples.FS.Open(archivePath)
	if err != nil {
		b.Fatal(err)
	}

	var commits []github.CommitCSV
	if err := csvtargz.DecodeFromFile(gzFile, github.CommitsCSVFilename, &commits); err != nil {
		b.Fatal(err)
	}

	_ = gzFile.Close()

	gzFile, err = samples.FS.Open(archivePath)
	if err != nil {
		b.Fatal(err)
	}

	var repoCSVs []github.RepoCSV
	if err := csvtargz.DecodeFromFile(gzFile, github.ReposCSVFilename, &repoCSVs); err != nil {
		b.Fatal(err)
	}

	_ = gzFile.Close()

	n := 10

	for i := 0; i < b.N; i++ {
		repos := github.NewReposSample(events, repoCSVs, commits)

		_, err := repos.TopNByCommitsPushed(n)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkReposSample_TopNByWatchEvents(b *testing.B) {
	archivePath := archivePath

	gzFile, err := samples.FS.Open(archivePath)
	if err != nil {
		b.Fatal(err)
	}

	defer func() {
		_ = gzFile.Close()
	}()

	var events []github.EventCSV
	if err := csvtargz.DecodeFromFile(gzFile, github.EventsCSVFilename, &events); err != nil {
		b.Fatal(err)
	}

	_ = gzFile.Close()

	gzFile, err = samples.FS.Open(archivePath)
	if err != nil {
		b.Fatal(err)
	}

	var commits []github.CommitCSV
	if err := csvtargz.DecodeFromFile(gzFile, github.CommitsCSVFilename, &commits); err != nil {
		b.Fatal(err)
	}

	_ = gzFile.Close()

	gzFile, err = samples.FS.Open(archivePath)
	if err != nil {
		b.Fatal(err)
	}

	var repoCSVs []github.RepoCSV
	if err := csvtargz.DecodeFromFile(gzFile, github.ReposCSVFilename, &repoCSVs); err != nil {
		b.Fatal(err)
	}

	_ = gzFile.Close()

	n := 10

	for i := 0; i < b.N; i++ {
		repos := github.NewReposSample(events, repoCSVs, commits)

		_, err := repos.TopNByWatchEvents(n)
		if err != nil {
			b.Fatal(err)
		}
	}
}
