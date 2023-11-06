package release

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/Masterminds/semver/v3"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

type PreReleaseFormat string

const (
	// PreReleaseFormatSemVer is the default format for prerelease versions.
	PreReleaseFormatSemVer PreReleaseFormat = "semver"
	// PreReleaseFormatDate is the format for prerelease versions that use the
	// current date.
	PreReleaseFormatDate PreReleaseFormat = "date"
	// PreReleaseFormatDateTime is the format for prerelease versions that
	// use the current date and time.
	PreReleaseFormatDateTime PreReleaseFormat = "datetime"
)

func (p PreReleaseFormat) String() string {
	return string(p)
}

type SemVerBumpType string

const (
	// SemVerBumpTypeMajor is the bump type for major versions.
	SemVerBumpTypeMajor SemVerBumpType = "major"
	// SemVerBumpTypeMinor is the bump type for minor versions.
	SemVerBumpTypeMinor SemVerBumpType = "minor"
	// SemVerBumpTypePatch is the bump type for patch versions.
	SemVerBumpTypePatch SemVerBumpType = "patch"
)

func (s SemVerBumpType) String() string {
	return string(s)
}

// bumpPreRelease bumps the prerelease version of a given semver.Version.
// The bumping is done according to the given PreReleaseFormat.
func bumpPreRelease(smvFormat PreReleaseFormat, version semver.Version, preReleasePrefix string) string {
	tagPrefix := ""
	if strings.HasPrefix(version.Original(), "v") {
		tagPrefix = "v"
	}
	tagPrefix += fmt.Sprintf("%d.%d.%d", version.Major(), version.Minor(), version.Patch())

	switch smvFormat {
	case PreReleaseFormatSemVer:
		if version.Prerelease() == "" {
			return fmt.Sprintf("%s-%s.%s", tagPrefix, preReleasePrefix, "1")
		}
		intVers, err := strconv.Atoi(strings.TrimPrefix(version.Prerelease(), preReleasePrefix+"."))
		if err != nil {
			return fmt.Sprintf("%s.%d", version.String(), 1)
		}
		return fmt.Sprintf("%s-%s.%d", tagPrefix, preReleasePrefix, intVers+1)
	case PreReleaseFormatDate:
		return fmt.Sprintf("%s-%s.%s", tagPrefix, preReleasePrefix, time.Now().Format("20060102"))
	case PreReleaseFormatDateTime:
		return fmt.Sprintf("%s-%s.%s", tagPrefix, preReleasePrefix, time.Now().Format("200601021504"))
	}
	return fmt.Sprintf("%s-%s.%d", tagPrefix, preReleasePrefix, 1)
}

// GetLatestSemVerTagFromRepo returns the latest semver tag from a given git repository.
// If no semver tag is found, it returns a semver.Version with the value v0.0.0.
func GetLatestSemVerTagFromRepo(repo *git.Repository) (*semver.Version, error) {
	// get tags from repository
	tags, err := repo.Tags()
	if err != nil {
		return nil, err
	}
	// get tags from repository
	vs := []*semver.Version{}
	tags.ForEach(func(t *plumbing.Reference) error {
		gitTag := strings.TrimPrefix(t.Name().String(), "refs/tags/")
		// check if tag matches semver format
		isSemver := regexp.MustCompile("^[v]{0,1}[0-9]{1,}.[0-9]{1,}.[0-9]{1,}(-[a-zA-Z0-9.-]+){0,1}$").MatchString(gitTag)
		if !isSemver {
			return nil
		}
		//
		smv, err := semver.NewVersion(gitTag)
		if err != nil {
			return err
		}
		vs = append(vs, smv)
		return nil
	})
	// sort tags
	sort.Sort(semver.Collection(vs))

	// get latest version
	var latest *semver.Version
	if len(vs) > 0 {
		if vs[len(vs)-1].Prerelease() != "" {
			latest = vs[len(vs)-1]
		} else {
			for i := len(vs) - 1; i >= 0; i-- {
				if vs[i].Prerelease() == "" {
					latest = vs[i]
					break
				}
			}
		}
	}

	if latest == nil {
		return semver.MustParse("v0.0.0"), nil
	}

	return latest, nil
}

// BumpTag takes a semver.Version and a semVerBumpType and returns the
// bumped version as a string.
func BumpTag(latest *semver.Version, semVerType SemVerBumpType, preReleaseFormat PreReleaseFormat, preReleasePrefix string, isPreRelease bool) string {
	newTag := semver.Version{}
	switch semVerType {
	case SemVerBumpTypeMajor:
		newTag = latest.IncMajor()
	case SemVerBumpTypeMinor:
		newTag = latest.IncMinor()
	case SemVerBumpTypePatch:
		newTag = latest.IncPatch()
	}
	// if pre-release is enabled, bump the pre-release version
	if isPreRelease {
		vrs := bumpPreRelease(preReleaseFormat, newTag, preReleasePrefix)
		newTag = *semver.MustParse(vrs)
	}
	return newTag.String()
}
