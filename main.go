package main

import (
	_ "embed"
	"flag"
	"fmt"

	"github.com/go-git/go-git/v5"
	"github.com/leonsteinhaeuser/git-tag-identifier/branch"
	"github.com/leonsteinhaeuser/git-tag-identifier/release"
	"gopkg.in/yaml.v3"
)

var (
	preReleaseFormat = flag.String("prerelease-format", release.PreReleaseFormatSemVer.String(), "Prerelease format. Can be 'semver', 'date' or 'datetime'")
	preReleasePrefix = flag.String("prerelease-prefix", "rc", "Prerelease prefix")
	bumpType         = flag.String("bump", release.SemVerBumpTypePatch.String(), "Bump type (major, minor, patch")
	isPreRelease     = flag.Bool("pre-release", false, "Whether to create a pre-release")
	repoTarget       = flag.String("repo-path", ".", "Path to the repository")
	configPath       = flag.String("config", "", "Path to the config file")
	autoBump         = flag.Bool("auto-bump", false, "Whether to automatically bump the version based on the rules in the config file")

	//go:embed config.yaml
	configBts []byte
	// embed default config during build
	config *branch.Config
)

func init() {
	flag.Parse()
	// if config flag is set, read config from file
	if *configPath != "" {
		cfg, err := branch.ReadConfig(*configPath)
		if err != nil {
			panic(err)
		}
		config = cfg
	} else {
		err := yaml.Unmarshal(configBts, &config)
		if err != nil {
			panic(err)
		}
	}
}

func main() {
	repo, err := git.PlainOpen(*repoTarget)
	if err != nil {
		panic(err)
	}

	latest, err := release.GetLatestSemVerTagFromRepo(repo)
	if err != nil {
		panic(err)
	}

	bt := release.SemVerBumpType(*bumpType)
	if *autoBump {
		identifier, err := branch.Identify(config, repo)
		if err != nil {
			panic(err)
		}
		bt = identifier
	}

	newTag := release.BumpTag(
		latest,
		bt,
		release.PreReleaseFormat(*preReleaseFormat),
		*preReleasePrefix,
		*isPreRelease,
	)

	fmt.Println("Latest version:", newTag)
}
