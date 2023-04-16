package branch

import (
	"fmt"

	"github.com/go-git/go-git/v5"
	"github.com/leonsteinhaeuser/git-tag-bump/release"
)

var (
	ErrBranchNameFormat = fmt.Errorf("branch name format is invalid")
)

// branchName returns the name of the current branch
func branchName(repo *git.Repository) (string, error) {
	branch, err := repo.Head()
	if err != nil {
		return "", err
	}
	return branch.Name().Short(), nil
}

// IdentifyBranch identifies the bump type of a branch
// if the branch does not match any of the configured identifiers, an error is returned
func IdentifyBranch(cfg *Config, branch string) (release.SemVerBumpType, error) {
	// check if branch name matches any of the configured identifiers
	if cfg.Major.match(branch) {
		return release.SemVerBumpTypeMajor, nil
	}
	if cfg.Minor.match(branch) {
		return release.SemVerBumpTypeMinor, nil
	}
	if cfg.Patch.match(branch) {
		return release.SemVerBumpTypePatch, nil
	}
	return "", fmt.Errorf("%w: for branch %q", ErrBranchNameFormat, branch)
}

// Identifier identifies the bump type of a branch or pull request.
func Identify(cfg *Config, repo *git.Repository) (release.SemVerBumpType, error) {
	bn, err := branchName(repo)
	if err != nil {
		return "", err
	}
	bumpType, err := IdentifyBranch(cfg, bn)
	if err == nil {
		// we found a match for the branch name
		return bumpType, nil
	}
	// we did not find a match for the branch name
	// TODO: check if pull request was merged and has the correct labels
	return "", fmt.Errorf("not implemented for branch %s", bn)
}
