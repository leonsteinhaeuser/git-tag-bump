package release

import (
	"fmt"
	"testing"
	"time"

	"github.com/Masterminds/semver/v3"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/storage/memory"
)

func TestPreReleaseFormat_String(t *testing.T) {
	tests := []struct {
		name string
		p    PreReleaseFormat
		want string
	}{
		{
			name: "PreReleaseFormatSemVer",
			p:    PreReleaseFormatSemVer,
			want: "semver",
		},
		{
			name: "PreReleaseFormatDate",
			p:    PreReleaseFormatDate,
			want: "date",
		},
		{
			name: "PreReleaseFormatDateAndTime",
			p:    PreReleaseFormatDateTime,
			want: "datetime",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.String(); got != tt.want {
				t.Errorf("PreReleaseFormat.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_semVerBumpType_String(t *testing.T) {
	tests := []struct {
		name string
		s    SemVerBumpType
		want string
	}{
		{
			name: "SemVerBumpTypeMajor",
			s:    SemVerBumpTypeMajor,
			want: "major",
		},
		{
			name: "SemVerBumpTypeMinor",
			s:    SemVerBumpTypeMinor,
			want: "minor",
		},
		{
			name: "SemVerBumpTypePatch",
			s:    SemVerBumpTypePatch,
			want: "patch",
		},
		{
			name: "SemVerBumpTypePatch",
			s:    "",
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.String(); got != tt.want {
				t.Errorf("semVerBumpType.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_bumpPreRelease(t *testing.T) {
	type args struct {
		smvFormat        PreReleaseFormat
		version          semver.Version
		preReleasePrefix string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "PreReleaseFormatSemVer with alpha prerelease",
			args: args{
				smvFormat:        PreReleaseFormatSemVer,
				version:          *semver.MustParse("1.0.0-alpha.1"),
				preReleasePrefix: "alpha",
			},
			want: "1.0.0-alpha.2",
		},
		{
			name: "PreReleaseFormatSemVer",
			args: args{
				smvFormat:        PreReleaseFormatSemVer,
				version:          *semver.MustParse("1.0.0"),
				preReleasePrefix: "alpha",
			},
			want: "1.0.0-alpha.1",
		},
		{
			name: "PreReleaseFormatDate",
			args: args{
				smvFormat:        PreReleaseFormatDate,
				version:          *semver.MustParse("1.0.0"),
				preReleasePrefix: "alpha",
			},
			want: fmt.Sprintf("1.0.0-alpha.%s", time.Now().Format("20060102")),
		},
		{
			name: "PreReleaseFormatDateAndTime",
			args: args{
				smvFormat:        PreReleaseFormatDateTime,
				version:          *semver.MustParse("1.0.0"),
				preReleasePrefix: "alpha",
			},
			want: fmt.Sprintf("1.0.0-alpha.%s", time.Now().Format("200601021504")),
		},
		//
		//
		//
		{
			name: "PreReleaseFormatSemVer with v prefix",
			args: args{
				smvFormat:        PreReleaseFormatSemVer,
				version:          *semver.MustParse("v1.0.0"),
				preReleasePrefix: "alpha",
			},
			want: "v1.0.0-alpha.1",
		},
		{
			name: "PreReleaseFormatDate with v prefix",
			args: args{
				smvFormat:        PreReleaseFormatDate,
				version:          *semver.MustParse("v1.0.0"),
				preReleasePrefix: "alpha",
			},
			want: fmt.Sprintf("v1.0.0-alpha.%s", time.Now().Format("20060102")),
		},
		{
			name: "PreReleaseFormatDateAndTime with v prefix",
			args: args{
				smvFormat:        PreReleaseFormatDateTime,
				version:          *semver.MustParse("v1.0.0"),
				preReleasePrefix: "alpha",
			},
			want: fmt.Sprintf("v1.0.0-alpha.%s", time.Now().Format("200601021504")),
		},

		{
			name: "empty",
			args: args{
				smvFormat:        "",
				version:          *semver.MustParse("v1.0.0"),
				preReleasePrefix: "alpha",
			},
			want: "v1.0.0-alpha.1",
		},

		{
			name: "PreReleaseFormatSemVer with alpha prerelease and not a number",
			args: args{
				smvFormat:        PreReleaseFormatSemVer,
				version:          *semver.MustParse("1.0.0-alpha.abc"),
				preReleasePrefix: "alpha",
			},
			want: "1.0.0-alpha.abc.1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := bumpPreRelease(tt.args.smvFormat, tt.args.version, tt.args.preReleasePrefix); got != tt.want {
				t.Errorf("bumpPreRelease() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_GetLatestSemVerTagFromRepo(t *testing.T) {
	type args struct {
		repo         *git.Repository
		isPreRelease bool
	}
	tests := []struct {
		name    string
		args    args
		want    *semver.Version
		wantErr bool
	}{
		{
			name: "no tags",
			args: args{
				repo: func() *git.Repository {
					repo, err := git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
						URL:  "https://github.com/leonsteinhaeuser/observer",
						Tags: git.AllTags,
					})
					if err != nil {
						t.Fatal(err)
					}
					return repo
				}(),
			},
			want:    semver.MustParse("2.0.1"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetLatestSemVerTagFromRepo(tt.args.repo, tt.args.isPreRelease)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetLatestSemVerTagFromRepo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got.String() != tt.want.String() {
				t.Errorf("GetLatestSemVerTagFromRepo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_BumpTag(t *testing.T) {
	type args struct {
		latest           *semver.Version
		semVerType       SemVerBumpType
		preReleaseFormat PreReleaseFormat
		preReleasePrefix string
		isPreRelease     bool
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "patch",
			args: args{
				latest:           semver.MustParse("1.0.0"),
				semVerType:       SemVerBumpTypePatch,
				preReleaseFormat: PreReleaseFormatSemVer,
				preReleasePrefix: "alpha",
				isPreRelease:     false,
			},
			want: "1.0.1",
		},
		{
			name: "patch pre-release",
			args: args{
				latest:           semver.MustParse("1.0.0"),
				semVerType:       SemVerBumpTypePatch,
				preReleaseFormat: PreReleaseFormatSemVer,
				preReleasePrefix: "alpha",
				isPreRelease:     true,
			},
			want: "1.0.1-alpha.1",
		},
		{
			name: "patch pre-release with existing pre-release",
			args: args{
				latest:           semver.MustParse("1.0.0-alpha.1"),
				semVerType:       SemVerBumpTypePatch,
				preReleaseFormat: PreReleaseFormatSemVer,
				preReleasePrefix: "alpha",
				isPreRelease:     true,
			},
			want: "1.0.1-alpha.1",
		},
		{
			name: "minor",
			args: args{
				latest:           semver.MustParse("1.0.0"),
				semVerType:       SemVerBumpTypeMinor,
				preReleaseFormat: PreReleaseFormatSemVer,
				preReleasePrefix: "alpha",
				isPreRelease:     false,
			},
			want: "1.1.0",
		},
		{
			name: "minor pre-release",
			args: args{
				latest:           semver.MustParse("1.0.0"),
				semVerType:       SemVerBumpTypeMinor,
				preReleaseFormat: PreReleaseFormatSemVer,
				preReleasePrefix: "alpha",
				isPreRelease:     true,
			},
			want: "1.1.0-alpha.1",
		},
		{
			name: "minor pre-release with existing pre-release",
			args: args{
				latest:           semver.MustParse("1.1.0-alpha.1"),
				semVerType:       SemVerBumpTypeMinor,
				preReleaseFormat: PreReleaseFormatSemVer,
				preReleasePrefix: "alpha",
				isPreRelease:     true,
			},
			want: "1.2.0-alpha.1",
		},
		{
			name: "major",
			args: args{
				latest:           semver.MustParse("1.0.0"),
				semVerType:       SemVerBumpTypeMajor,
				preReleaseFormat: PreReleaseFormatSemVer,
				preReleasePrefix: "alpha",
				isPreRelease:     false,
			},
			want: "2.0.0",
		},
		{
			name: "major",
			args: args{
				latest:           semver.MustParse("1.0.0"),
				semVerType:       SemVerBumpTypeMajor,
				preReleaseFormat: PreReleaseFormatSemVer,
				preReleasePrefix: "alpha",
				isPreRelease:     true,
			},
			want: "2.0.0-alpha.1",
		},
		{
			name: "major pre-release with existing pre-release",
			args: args{
				latest:           semver.MustParse("1.0.0-alpha.1"),
				semVerType:       SemVerBumpTypeMajor,
				preReleaseFormat: PreReleaseFormatSemVer,
				preReleasePrefix: "alpha",
				isPreRelease:     true,
			},
			want: "2.0.0-alpha.1",
		},
		// tests added during bug fixes
		{
			name: "pre-release with existing pre-release and bump type none",
			args: args{
				latest:           semver.MustParse("v1.1.0"),
				semVerType:       SemVerBumpTypeNone,
				preReleaseFormat: PreReleaseFormatSemVer,
				preReleasePrefix: "rc",
				isPreRelease:     true,
			},
			want: "1.1.0-rc.1",
		},
		{
			name: "release with existing pre-release and bump type none",
			args: args{
				latest:           semver.MustParse("v1.1.0"),
				semVerType:       SemVerBumpTypeNone,
				preReleaseFormat: PreReleaseFormatSemVer,
				preReleasePrefix: "rc",
				isPreRelease:     false,
			},
			want: "1.1.0",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BumpTag(tt.args.latest, tt.args.semVerType, tt.args.preReleaseFormat, tt.args.preReleasePrefix, tt.args.isPreRelease); got != tt.want {
				t.Errorf("BumpTag() = %v, want %v", got, tt.want)
			}
		})
	}
}
