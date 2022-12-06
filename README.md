# Semantic Versioning Action

This action calculates the next version relying on semantic versioning.

## Strategies

If `auto` bump, it will try to extract the closest tag and calculate the next semantic version. If not, it will bump respecting the value passed.

### Branch Names

These are the prefixes we expect when `auto` bump:

- `^bugfix/.+` - `patch`
- `^feature/.+` - `minor`
- `^major/.+` - `major`
- `^misc/.+` - `build`
- `^docs?/.+` - `build`

### Scenarios

In case of `force_prelease` is `true`, it will always create a pre-release version. Otherwise, it will create a final version.

#### Auto Bump

- Not a valid source branch prefix - Increments prerelease version.

    ```text
        v0.1.0 results in v0.1.0-pre.1
        v1.5.3-pre.2 results in v1.5.3-pre.3
    ```

- Source branch is prefixed with `misc/` or `docs/` and dest branch is `main` - Increments build version.

    ```text
    v1.5.3-pre.2 results in v1.5.3-pre.3
    ```

- Source branch is prefixed with `bugfix/` and dest branch is `main` - Increments patch version.

    ```text
    v0.1.0 results in v0.1.1-pre.1
    v1.5.3-pre.2 results in v1.5.4-pre.1
    ```

- Source branch is prefixed with `feature/` and dest branch is `main` - Increments minor version.

    ```text
    v0.1.0 results in v0.2.0-pre.1
    v1.5.3-pre.2 results in v1.6.0-pre.1
    ```

- Source branch is prefixed with `major/` and dest branch is `main` - Increments major version.

    ```text
    v0.1.0 results in v1.0.0-pre.1
    v1.5.3-pre.2 results in v2.0.0-pre.1
    ```

## Github Environment Variables

Here are the environment variables we take from Github Actions so far:

- `GITHUB_SHA`

## Example usage

### Basic

Uses `auto` bump strategy to calculate the next semantic version.

```yaml
- id: semver-tag
  uses: snapfi/semver-action
- name: "Created tag"
  run: echo "tag ${{ steps.semver-tag.outputs.semver_tag }}"
```

### Custom

```yaml
- id: semver-tag
  uses: snapfi/semver-action
  with:
    prefix: ""
    prerelease_id: "alpha"
    branch_name: "master"
    debug: "true"
- name: "Created tag"
  run: echo "tag ${{ steps.semver-tag.outputs.semver_tag }}"
```

## Inputs

| parameter | required | description | default |
| --- | --- | --- | --- |
| bump | false | Bump strategy for semantic versioning. Can be `auto`, `major`, `minor`, `patch`. | auto |
| base_version | false | Version to use as base for the generation, skips version bumps. | |
| prefix | false | Prefix used to prepend the final version.| v |
| prerelease_id | false | Text representing the prerelease identifier. | pre |
| force_prerelease | false | Force the generation of a prerelease version. Usually used in trunk-based development where the main branch is always a prerelease version. | false |
| branch_name | false | The branch name. | main  |
| repo_dir | false | The repository path. | current dir |
| debug | false | Enables debug mode. | false |

## Outpus

| parameter     | description                                      |
| ---           | ---                                              |
| semver_tag    | The calculated semantic version.                 |
| is_prerelease | True if calculated tag is prerelease.            |
| previous_tag  | The tag used to calculate next semantic version. |
| ancestor_tag  | The ancestor tag based on specific pattern.      |
