package main

import (
	_ "embed"
	"flag"
	"fmt"

	"github.com/go-git/go-git/v5"
	gconfig "github.com/go-git/go-git/v5/config"
	"github.com/leonsteinhaeuser/git-tag-identifier/branch"
	"github.com/leonsteinhaeuser/git-tag-identifier/release"
	"gopkg.in/yaml.v3"
)

var (
	preReleaseFormat = flag.String("pre-release-format", release.PreReleaseFormatSemVer.String(), "Prerelease format. Can be 'semver', 'date' or 'datetime'")
	preReleasePrefix = flag.String("pre-release-prefix", "rc", "Prerelease prefix")
	bumpType         = flag.String("bump", release.SemVerBumpTypePatch.String(), "Bump type (major, minor, patch")
	isPreRelease     = flag.Bool("pre-release", false, "Whether to create a pre-release")
	repoTarget       = flag.String("repo-path", ".", "Path to the repository")
	configPath       = flag.String("config", "", "Path to the config file")
	autoBump         = flag.Bool("auto-bump", false, "Whether to automatically bump the version based on the rules in the config file")
	createTag        = flag.Bool("create", false, "Whether to create a tag in the repository and push it to the remote")

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

	if *createTag {
		rfc, err := repo.Head()
		if err != nil {
			panic(err)
		}
		// create the tag in the repository for the current commit hash
		pmbrfc, err := repo.CreateTag(newTag, rfc.Hash(), &git.CreateTagOptions{
			//Tagger:  nil,
			Message: newTag,
		})
		if err != nil {
			panic(err)
		}

		// push the tag to the remote
		refTag := pmbrfc.Name().String()
		err = repo.Push(&git.PushOptions{
			FollowTags: true,
			RefSpecs: []gconfig.RefSpec{
				gconfig.RefSpec(fmt.Sprintf("%s:%s", refTag, refTag)),
			},
		})
		if err != nil {
			panic(err)
		}
	}

	fmt.Println(newTag)
}
