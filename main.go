package main

import (
	_ "embed"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/Masterminds/semver/v3"
	"github.com/go-git/go-git/v5"
	gconfig "github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/leonsteinhaeuser/git-tag-bump/branch"
	"github.com/leonsteinhaeuser/git-tag-bump/release"
	"gopkg.in/yaml.v3"
)

var (
	preReleaseFormat     = flag.String("pre-release-format", release.PreReleaseFormatSemVer.String(), "Prerelease format. Can be 'semver', 'date' or 'datetime'")
	preReleasePrefix     = flag.String("pre-release-prefix", "rc", "Prerelease prefix")
	bumpType             = flag.String("bump", release.SemVerBumpTypePatch.String(), "Bump type (major, minor, patch, none)")
	isPreRelease         = flag.Bool("pre-release", false, "Whether to create a pre-release")
	repoTarget           = flag.String("repo-path", ".", "Path to the repository")
	configPath           = flag.String("config", "", "Path to the config file")
	autoBump             = flag.Bool("auto-bump", false, "Whether to automatically bump the version based on the rules in the config file")
	createTag            = flag.Bool("create", false, "Whether to create a tag in the repository and push it to the remote")
	createTagLightweight = flag.Bool("lightweight", false, "Whether any tag created should be a lightweight tag")
	branchName           = flag.String("branch-name", "", "Name of the branch to check")
	vPrefix              = flag.Bool("v-prefix", true, "Whether to prefix the tag with a 'v'. E.g. v1.0.0 instead of 1.0.0")
	gitBaseTagOverride   = flag.String("git-base-tag", "", "Override the base tag to use for the bump. If not set, the latest tag will be used.")

	actorName = flag.String("actor-name", "", "The name of the actor used to create the tag. Only used if --create is set.")
	actorMail = flag.String("actor-mail", "", "The mail of the actor used to create the tag. Only used if --create is set.")

	actor *object.Signature = &object.Signature{}

	githubToken = os.Getenv("GITHUB_TOKEN")

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

	if *createTag && !*createTagLightweight && (*actorName == "" || *actorMail == "" || githubToken == "") {
		panic("Either --lightweight, or both --actor-name and --actor-mail must be set when --create is set")
	}

	if *createTag && *actorName != "" && *actorMail != "" {
		actor.Email = *actorMail
		actor.Name = *actorName
		actor.When = time.Now()
	}
}

func main() {
	repo, err := git.PlainOpen(*repoTarget)
	if err != nil {
		panic(err)
	}

	latest, err := release.GetLatestSemVerTagFromRepo(repo, *isPreRelease)
	if err != nil {
		panic(err)
	}

	// override the current latest identified tag with the one from the flag
	if *gitBaseTagOverride != "" {
		overrideTag := semver.MustParse(*gitBaseTagOverride)
		if overrideTag.Major() != latest.Major() || overrideTag.Minor() != latest.Minor() || overrideTag.Patch() != latest.Patch() {
			// if the major, minor or patch version of the override tag does not match the latest tag, use the override tag
			latest = overrideTag
		}
	}

	bt := release.SemVerBumpType(*bumpType)
	if *autoBump && *branchName == "" {
		identifier, err := branch.Identify(config, repo)
		if err != nil {
			panic(err)
		}
		bt = identifier
	}
	if *branchName != "" {
		smvTag, err := branch.IdentifyBranch(config, *branchName)
		if err != nil {
			panic(err)
		}
		bt = smvTag
	}

	newTag := release.BumpTag(
		latest,
		bt,
		release.PreReleaseFormat(*preReleaseFormat),
		*preReleasePrefix,
		*isPreRelease,
	)

	// add v prefix if enabled
	if *vPrefix {
		newTag = fmt.Sprintf("v%s", newTag)
	}

	if *createTag {
		rfc, err := repo.Head()
		if err != nil {
			panic(err)
		}

		var options *git.CreateTagOptions
		if !*createTagLightweight {
			options = &git.CreateTagOptions{
				//Tagger:  nil,
				Message: newTag,
				Tagger: func() *object.Signature {
					if actor.Name == "" || actor.Email == "" {
						return nil
					}
					return actor
				}(),
			}
		}

		// create the tag in the repository for the current commit hash
		pmbrfc, err := repo.CreateTag(newTag, rfc.Hash(), options)
		if err != nil {
			panic(fmt.Sprintf("Could not create tag: %q, exited with error: %s", newTag, err))
		}

		// push the tag to the remote
		refTag := pmbrfc.Name().String()
		err = repo.Push(&git.PushOptions{
			FollowTags: true,
			RefSpecs: []gconfig.RefSpec{
				gconfig.RefSpec(fmt.Sprintf("%s:%s", refTag, refTag)),
			},
			Progress: os.Stdout,
			Auth:     &http.BasicAuth{Username: "bot", Password: githubToken},
		})
		if err != nil {
			panic(err)
		}
	}

	fmt.Println(newTag)
}
