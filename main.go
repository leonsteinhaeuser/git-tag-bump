package main

import (
	"flag"
	"fmt"

	"github.com/go-git/go-git/v5"
	"github.com/leonsteinhaeuser/git-tag-identifier/release"
)

var (
	preReleaseFormat = flag.String("prerelease-format", release.PreReleaseFormatSemVer.String(), "Prerelease format. Can be 'semver', 'date' or 'datetime'")
	preReleasePrefix = flag.String("prerelease-prefix", "rc", "Prerelease prefix")
	bumpType         = flag.String("bump", release.SemVerBumpTypePatch.String(), "Bump type (major, minor, patch")
	isPreRelease     = flag.Bool("pre-release", false, "Whether to create a pre-release")
	repoTarget       = flag.String("repo-path", ".", "Path to the repository")
)

func init() {
	flag.Parse()
}

func main() {
	repo, err := git.PlainOpen(*repoTarget)
	if err != nil {
		panic(err)
	}

	latest, err := release.GetLatestSemVerTagFromRepoPath(repo, release.SemVerBumpType(*bumpType))
	if err != nil {
		panic(err)
	}

	newTag := release.BumpVersion(
		latest,
		release.SemVerBumpType(*bumpType),
		release.PreReleaseFormat(*preReleaseFormat),
		*preReleasePrefix,
		*isPreRelease,
	)

	fmt.Println("Latest version:", newTag)
}
