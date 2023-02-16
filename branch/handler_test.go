package branch

import (
	"reflect"
	"testing"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/leonsteinhaeuser/git-tag-identifier/release"
)

func Test_branchName(t *testing.T) {
	type args struct {
		repo *git.Repository
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "master",
			args: args{
				repo: func() *git.Repository {
					repo, err := git.PlainOpenWithOptions(".", &git.PlainOpenOptions{
						DetectDotGit: true,
					})
					if err != nil {
						t.Fatal(err)
					}
					return repo
				}(),
			},
			want: func() string {
				repo, err := git.PlainOpenWithOptions(".", &git.PlainOpenOptions{
					DetectDotGit: true,
				})
				if err != nil {
					t.Fatal(err)
				}
				hd, err := repo.Head()
				if err != nil {
					t.Fatal(err)
				}
				return hd.Name().Short()
			}(),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := branchName(tt.args.repo)
			if (err != nil) != tt.wantErr {
				t.Errorf("branchName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("branchName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_identifyBranch(t *testing.T) {
	type args struct {
		cfg    *Config
		branch string
	}
	tests := []struct {
		name    string
		args    args
		want    release.SemVerBumpType
		wantErr bool
	}{
		{
			name: "major",
			args: args{
				cfg: &Config{
					Major: Identifier{
						Branch: BranchIdentifier{
							Name: RegExIdentifier{
								RegEx: "^[a-z]+!/",
							},
						},
					},
					Minor: Identifier{
						Branch: BranchIdentifier{
							Name: RegExIdentifier{
								RegEx: "^(feat|feature)/",
							},
						},
					},
					Patch: Identifier{
						Branch: BranchIdentifier{
							Name: RegExIdentifier{
								RegEx: `^(enh|enhanc|enhancement|fix|bugfix|chore\([a-z0-9-]+\))/`,
							},
						},
					},
				},
				branch: "feat!/feature",
			},
			want:    "major",
			wantErr: false,
		},
		{
			name: "minor",
			args: args{
				cfg: &Config{
					Major: Identifier{
						Branch: BranchIdentifier{
							Name: RegExIdentifier{
								RegEx: "^[a-z]+!/",
							},
						},
					},
					Minor: Identifier{
						Branch: BranchIdentifier{
							Name: RegExIdentifier{
								RegEx: "^(feat|feature)/",
							},
						},
					},
					Patch: Identifier{
						Branch: BranchIdentifier{
							Name: RegExIdentifier{
								RegEx: `^(enh|enhanc|enhancement|fix|bugfix|chore\([a-z0-9-]+\))/`,
							},
						},
					},
				},
				branch: "feat/feature",
			},
			want:    "minor",
			wantErr: false,
		},
		{
			name: "patch",
			args: args{
				cfg: &Config{
					Major: Identifier{
						Branch: BranchIdentifier{
							Name: RegExIdentifier{
								RegEx: "^[a-z]+!/",
							},
						},
					},
					Minor: Identifier{
						Branch: BranchIdentifier{
							Name: RegExIdentifier{
								RegEx: "^(feat|feature)/",
							},
						},
					},
					Patch: Identifier{
						Branch: BranchIdentifier{
							Name: RegExIdentifier{
								RegEx: `^(enh|enhanc|enhancement|fix|bugfix|chore\([a-z0-9-]+\))/`,
							},
						},
					},
				},
				branch: "fix/feature",
			},
			want:    "patch",
			wantErr: false,
		},
		{
			name: "branch name format mismatch",
			args: args{
				cfg: &Config{
					Major: Identifier{
						Branch: BranchIdentifier{
							Name: RegExIdentifier{
								RegEx: `^(feat|feature|enh|enhanc|enhancement|fix|bugfix|chore\([a-z0-9-]+\))!/`,
							},
						},
					},
					Minor: Identifier{
						Branch: BranchIdentifier{
							Name: RegExIdentifier{
								RegEx: "^(feat|feature)/",
							},
						},
					},
					Patch: Identifier{
						Branch: BranchIdentifier{
							Name: RegExIdentifier{
								RegEx: `^(enh|enhanc|enhancement|fix|bugfix|chore\([a-z0-9-]+\))/`,
							},
						},
					},
				},
				branch: "failure/branchname",
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := identifyBranch(tt.args.cfg, tt.args.branch)
			if (err != nil) != tt.wantErr {
				t.Errorf("identifyBranch() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("identifyBranch() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIdentify(t *testing.T) {
	type args struct {
		cfg  *Config
		repo *git.Repository
	}
	tests := []struct {
		name    string
		args    args
		want    release.SemVerBumpType
		wantErr bool
	}{
		{
			name: "major",
			args: args{
				cfg: &Config{
					Major: Identifier{
						Branch: BranchIdentifier{
							Name: RegExIdentifier{
								RegEx: "^[a-z0-9)(]+/",
							},
						},
					},
					Minor: Identifier{
						Branch: BranchIdentifier{
							Name: RegExIdentifier{
								RegEx: "^[a-z0-9)(]!+/",
							},
						},
					},
					Patch: Identifier{
						Branch: BranchIdentifier{
							Name: RegExIdentifier{
								RegEx: "^[a-z0-9)(]!+/",
							},
						},
					},
				},
				repo: func() *git.Repository {
					repo, err := git.PlainOpenWithOptions(".", &git.PlainOpenOptions{
						DetectDotGit: true,
					})
					if err != nil {
						t.Fatal(err)
					}
					return repo
				}(),
			},
			want:    release.SemVerBumpTypeMajor,
			wantErr: false,
		},
		{
			name: "minor",
			args: args{
				cfg: &Config{
					Major: Identifier{
						Branch: BranchIdentifier{
							Name: RegExIdentifier{
								RegEx: "^[a-z0-9)(]!+/",
							},
						},
					},
					Minor: Identifier{
						Branch: BranchIdentifier{
							Name: RegExIdentifier{
								RegEx: "^[a-z0-9)(]+/",
							},
						},
					},
					Patch: Identifier{
						Branch: BranchIdentifier{
							Name: RegExIdentifier{
								RegEx: "^[a-z0-9)(]!+/",
							},
						},
					},
				},
				repo: func() *git.Repository {
					repo, err := git.PlainOpenWithOptions(".", &git.PlainOpenOptions{
						DetectDotGit: true,
					})
					if err != nil {
						t.Fatal(err)
					}
					return repo
				}(),
			},
			want:    release.SemVerBumpTypeMinor,
			wantErr: false,
		},
		{
			name: "patch",
			args: args{
				cfg: &Config{
					Major: Identifier{
						Branch: BranchIdentifier{
							Name: RegExIdentifier{
								RegEx: "^[a-z0-9)(]!+/",
							},
						},
					},
					Minor: Identifier{
						Branch: BranchIdentifier{
							Name: RegExIdentifier{
								RegEx: "^[a-z0-9)(]!+/",
							},
						},
					},
					Patch: Identifier{
						Branch: BranchIdentifier{
							Name: RegExIdentifier{
								RegEx: "^[a-z0-9)(]+/",
							},
						},
					},
				},
				repo: func() *git.Repository {
					repo, err := git.PlainOpenWithOptions(".", &git.PlainOpenOptions{
						DetectDotGit: true,
					})
					if err != nil {
						t.Fatal(err)
					}
					return repo
				}(),
			},
			want:    release.SemVerBumpTypePatch,
			wantErr: false,
		},
		{
			name: "error",
			args: args{
				cfg: &Config{},
				repo: &git.Repository{
					Storer: memory.NewStorage(),
				},
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "not implemented",
			args: args{
				cfg: &Config{
					Major: Identifier{
						Branch: BranchIdentifier{
							Name: RegExIdentifier{
								RegEx: "^sajdnksjnd",
							},
						},
					},
					Minor: Identifier{
						Branch: BranchIdentifier{
							Name: RegExIdentifier{
								RegEx: "^sajdnksjnd",
							},
						},
					},
					Patch: Identifier{
						Branch: BranchIdentifier{
							Name: RegExIdentifier{
								RegEx: "^sajdnksjnd",
							},
						},
					},
				},
				repo: func() *git.Repository {
					repo, err := git.PlainOpenWithOptions(".", &git.PlainOpenOptions{
						DetectDotGit: true,
					})
					if err != nil {
						t.Fatal(err)
					}
					return repo
				}(),
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Identify(tt.args.cfg, tt.args.repo)
			if (err != nil) != tt.wantErr {
				t.Errorf("Identify() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Identify() = %v, want %v", got, tt.want)
			}
		})
	}
}
