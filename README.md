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
| `--create` | `bool` | false | `false` | Whether to create and push the tag if it does not exist. |
| `--branch-name` | `string` | false | `` | The name of the branch to use. |
| `--v-prefix`   | `bool` | false    | `true` | Whether to prefix the tag with `v`. Example: `v1.0.0` instead of `1.0.0`. |

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
