# git-tag-bump

[![testing](https://github.com/leonsteinhaeuser/git-tag-bump/actions/workflows/tests.yaml/badge.svg)](https://github.com/leonsteinhaeuser/git-tag-bump/actions/workflows/tests.yaml)

A simple tool to bump git tags (semver). It can be used to bump the patch, minor or major version of a tag. It can also be used to automatically determine the next version based on the last tag and the current branch name.

## Arguments

| Flag            | Type     | Required | Default | Description |
|-----------------|----------|----------|---------|-------------|
| `--auto-bump`   | `bool`   | false    | `false` | Automatically determine the next version based on the last tag and the branch name passed to it. |
| `--bump`        | `string` | false    | `patch` | The part of the version to bump. Can be `patch`, `minor` or `major`. |
| `--config`      | `string` | false    | `` | The path to the config file. If not defined, the default config will be used. |
| `--pre-release` | `bool`   | false    | `false` | Whether to create a pre-release tag. |
| `--pre-release-format` | `string` | false    | `semver` | The format of the pre-release tag. Can be `semver`, `date` or `datetime` |
| `--pre-release-prefix` | `string` | false    | `rc` | The prefix of the pre-release tag. Example: When defining the following tag `v1.0.0-rc.1`, `rc` would be the prefix and the number after it the format ***semver***. |
| `--repo-path`   | `string` | false    | `.` | The path to the git repository. If not defined, the current working directory will be used. |
| `--create` | `bool` | false | `false` | Whether to create and push the tag if it does not exist. Requires `--actor-name`, `--actor-mail` and the env variable `GITHUB_TOKEN` to be set. |
| `--actor-name` | `string` | false | `` | The name of the actor used to create the tag. Only used if `--create` is set. |
| `--actor-email` | `string` | false | `` | The mail of the actor used to create the tag. Only used if `--create` is set. |
| `--branch-name` | `string` | false | `` | The name of the branch to use. |
| `--v-prefix`   | `bool` | false    | `true` | Whether to prefix the tag with `v`. Example: `v1.0.0` instead of `1.0.0`. |

## Environment Variables

| Variable | Description |
|----------|-------------|
| `GITHUB_TOKEN` | The GitHub token used to authenticate with ***git*** in order to push the tag. Only necessary if `--create` is set. |

## Config

The config file is a simple YAML file. It can be used to define the branch name rules to determine the next version. The following example shows the default config:

```yaml
major:
  branch:
    name:
      regex: '^(feat|feature|enh|enhanc|enhancement|fix|bugfix|chore)(\([a-z0-9-]+\)){0,1}!\/'

minor:
  branch:
    name:
      regex: '^(feat|feature)(\([a-z0-9-]+\)){0,1}\/'

patch:
  branch:
    name:
      regex: '(enh|enhanc|enhancement|fix|bugfix|chore)(\([a-z0-9-]+\)){0,1}\/'
```

The config file can be passed to the tool using the `--config` flag. If no config file is passed, the default config will be used.

The config file consists of three parts: `major`, `minor` and `patch`. Each part defines the rules for the corresponding version part. The rules are defined using regular expressions. The first rule that matches the branch name will be used to determine the next version.

What this means is that the following branch names will result in the following versions. In the following example, we assume that the last tag is `v1.0.0`.

| Branch name             | Release Type    | Version  |
| ----------------------- | ------- | -------- |
| `feat!/abc`             | `major` | `v2.0.0` |
| `feat(ctx)!/abc`        | `major` | `v2.0.0` |
| `feature!/abc`          | `major` | `v2.0.0` |
| `feature(ctx)!/abc`     | `major` | `v2.0.0` |
| `enh!/abc`              | `major` | `v2.0.0` |
| `enh(ctx)!/abc`         | `major` | `v2.0.0` |
| `enhanc!/abc`           | `major` | `v2.0.0` |
| `enhanc(ctx)!/abc`      | `major` | `v2.0.0` |
| `enhancement!/abc`      | `major` | `v2.0.0` |
| `enhancement(ctx)!/abc` | `major` | `v2.0.0` |
| `fix!/abc`              | `major` | `v2.0.0` |
| `fix(ctx)!/abc`         | `major` | `v2.0.0` |
| `bugfix!/abc`           | `major` | `v2.0.0` |
| `bugfix(ctx)!/abc`      | `major` | `v2.0.0` |
| `chore!/abc`            | `major` | `v2.0.0` |
| `chore(ctx)!/abc`       | `major` | `v2.0.0` |
|                         |         |          |
| `feat/abc`              | `minor` | `v1.1.0` |
| `feat(ctx)/abc`         | `minor` | `v1.1.0` |
| `feature/abc`           | `minor` | `v1.1.0` |
| `feature(ctx)/abc`      | `minor` | `v1.1.0` |
|                         |         |          |
| `enh`                   | `patch` | `v1.0.1` |
| `enh(ctx)`              | `patch` | `v1.0.1` |
| `enhanc`                | `patch` | `v1.0.1` |
| `enhanc(ctx)`           | `patch` | `v1.0.1` |
| `enhancement`           | `patch` | `v1.0.1` |
| `enhancement(ctx)`      | `patch` | `v1.0.1` |
| `fix`                   | `patch` | `v1.0.1` |
| `fix(ctx)`              | `patch` | `v1.0.1` |
| `bugfix`                | `patch` | `v1.0.1` |
| `bugfix(ctx)`           | `patch` | `v1.0.1` |
| `chore`                 | `patch` | `v1.0.1` |
| `chore(ctx)`            | `patch` | `v1.0.1` |

## Using the tool in a CI/CD pipeline

The tool can be used in a CI/CD pipeline to automatically determine the next version and create a tag for it. The following example shows how to use the tool in a GitHub CI/CD pipeline:

Manually create a tag based on the branch name:

```yaml
name: Release

on:
  pull_request:
    types: [closed]

jobs:
  tag-branch:
    if: github.event.pull_request.merged == true && (github.base_ref == 'main' || github.base_ref == 'staging')
    runs-on: ubuntu-latest
    env:
      FROM_BRANCH: ${{ github.head_ref }}
      SHORT_BRANCH_NAME: ${{ github.base_ref }}
      GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
    steps:
      - name: Checkout
        uses: actions/checkout@v3
        with:
          fetch-depth: 0
          branch: ${{ env.SHORT_BRANCH_NAME }}

      - name: Extract branch name prefix
        run: |
          echo "BRANCH_NAME_PREFIX=$(echo ${{ env.FROM_BRANCH }} | cut -d'/' -f1)" >> $GITHUB_ENV

      - name: Is feature branch and does not contain a breaking change and is staging
        if: |
          (contains(env.BRANCH_NAME_PREFIX, 'feat') ||
          contains(env.BRANCH_NAME_PREFIX, 'feature') ||
          contains(env.BRANCH_NAME_PREFIX, 'chore') ||
          contains(env.BRANCH_NAME_PREFIX, 'refactor') ||
          contains(env.BRANCH_NAME_PREFIX, 'revert')) &&
          !endsWith(env.BRANCH_NAME_PREFIX, '!') && github.base_ref == 'staging'
        uses: leonsteinhaeuser/git-tag-bump@v1.0.2
        with:
          args: "--bump patch --v-prefix --pre-release --create"
      - name: Is feature branch and does not contain a breaking change and is not staging
        if: |
          (contains(env.BRANCH_NAME_PREFIX, 'feat') ||
          contains(env.BRANCH_NAME_PREFIX, 'feature') ||
          contains(env.BRANCH_NAME_PREFIX, 'chore') ||
          contains(env.BRANCH_NAME_PREFIX, 'refactor') ||
          contains(env.BRANCH_NAME_PREFIX, 'revert')) &&
          !endsWith(env.BRANCH_NAME_PREFIX, '!') && github.base_ref == 'main'
        uses: leonsteinhaeuser/git-tag-bump@v1.0.2
        with:
          args: "--bump patch --v-prefix --create"

      - name: Is patch branch and is staging
        if: |
          (contains(env.BRANCH_NAME_PREFIX, 'fix') ||
          contains(env.BRANCH_NAME_PREFIX, 'bugfix') ||
          contains(env.BRANCH_NAME_PREFIX, 'hotfix')) &&
          !endsWith(env.BRANCH_NAME_PREFIX, '!') && github.base_ref == 'staging'
        uses: leonsteinhaeuser/git-tag-bump@v1.0.2
        with:
          args: "--bump minor --v-prefix --pre-release --create"
      - name: Is patch branch and is not staging
        if: |
          (contains(env.BRANCH_NAME_PREFIX, 'fix') ||
          contains(env.BRANCH_NAME_PREFIX, 'bugfix') ||
          contains(env.BRANCH_NAME_PREFIX, 'hotfix')) &&
          !endsWith(env.BRANCH_NAME_PREFIX, '!') && github.base_ref == 'main'
        uses: leonsteinhaeuser/git-tag-bump@v1.0.2
        with:
          args: "--bump minor --v-prefix --create"

      - name: Is a breaking change and is staging
        if: |
          (contains(env.BRANCH_NAME_PREFIX, 'feat') ||
          contains(env.BRANCH_NAME_PREFIX, 'feature') ||
          contains(env.BRANCH_NAME_PREFIX, 'chore') ||
          contains(env.BRANCH_NAME_PREFIX, 'refactor') ||
          contains(env.BRANCH_NAME_PREFIX, 'revert') ||
          contains(env.BRANCH_NAME_PREFIX, 'fix') ||
          contains(env.BRANCH_NAME_PREFIX, 'bugfix') ||
          contains(env.BRANCH_NAME_PREFIX, 'hotfix')) &&
          endsWith(env.BRANCH_NAME_PREFIX, '!') && github.base_ref == 'staging'
        uses: leonsteinhaeuser/git-tag-bump@v1.0.2
        with:
          args: "--bump major --v-prefix --pre-release --create"
      - name: Is a breaking change and is not staging
        if: |
          (contains(env.BRANCH_NAME_PREFIX, 'feat') ||
          contains(env.BRANCH_NAME_PREFIX, 'feature') ||
          contains(env.BRANCH_NAME_PREFIX, 'chore') ||
          contains(env.BRANCH_NAME_PREFIX, 'refactor') ||
          contains(env.BRANCH_NAME_PREFIX, 'revert') ||
          contains(env.BRANCH_NAME_PREFIX, 'fix') ||
          contains(env.BRANCH_NAME_PREFIX, 'bugfix') ||
          contains(env.BRANCH_NAME_PREFIX, 'hotfix')) &&
          endsWith(env.BRANCH_NAME_PREFIX, '!') && github.base_ref == 'main'
        uses: leonsteinhaeuser/git-tag-bump@v1.0.2
        with:
          args: "--bump major --v-prefix --create"
```

Automatically determine the next version and create a tag for it:

```yaml
name: Release

on:
  pull_request:
    types: [closed]

jobs:
  tag-branch:
    if: github.event.pull_request.merged == true && (github.base_ref == 'main' || github.base_ref == 'staging')
    runs-on: ubuntu-latest
    env:
      FROM_BRANCH: ${{ github.head_ref }}
      SHORT_BRANCH_NAME: ${{ github.base_ref }}
      GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
    steps:
      - name: Checkout
        uses: actions/checkout@v3
        with:
          fetch-depth: 0
          branch: ${{ env.SHORT_BRANCH_NAME }}

      # feature --> staging
      - name: Pre-release
        if: github.base_ref == 'staging'
        uses: leonsteinhaeuser/git-tag-bump@v1.0.2
        with:
          args: >-
            --v-prefix
            --pre-release
            --pre-release-format 'datetime'
            --auto-bump
            --branch-name "${{ env.FROM_BRANCH }}"
            --create

      # staging --> main
      - name: Release
        if: github.head_ref == 'staging'
        uses: leonsteinhaeuser/git-tag-bump@v1.0.2
        with:
          args: >-
            --v-prefix
            --pre-release
            --pre-release-format 'datetime'
            --auto-bump
            --branch-name "${{ env.FROM_BRANCH }}"
            --create
```
