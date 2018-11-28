# How to release kubor

## Rules

1. Everything which should be released needs to be on a specific tag that matches the `^v\d+.\d+.\d+|snapshot-.+$` regexp pattern.
2. Real releases should follow the `^v\d+.\d+.\d+$` regexp pattern.
3. Releases should be only created from the [`master`](https://github.com/levertonai/kubor/tree/master) branch.
4. The master branch has to be always stable.
5. Snapshot releases (for testing new features) should be following the `^snapshot-.+$` regexp pattern.
6. Snapshot releases should be cleaned up after successful testing.

## Trigger the release

Create a tag/release with that matches the `^v\d+.\d+.\d+|snapshot-.+$` regexp pattern by going to the [Draft release page](https://github.com/levertonai/kubor/releases/new). Select the branch to create the release from.

This will create following the [.travis.yml](.travis.yml) on the release page with the name of this tag the corresponding files. See [build status page](https://travis-ci.org/levertonai/kubor) to follow the process.
