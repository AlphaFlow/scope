# Deploying

So you want to add some code.  This file will outline the instructions for doing so in the context of this repository.
Follow these steps in order.

>WARNING: We will describe how to create a new version below.  Any other projects that integrate this package will
also need their `go.mod` file updated so that they can access the new code.

### Create a pull request
First, create a pull request for your feature branch.  Have this PR reviewed.

### Merge branches (feature -> main)

Merge your approved PR on GitHub into dev by squashing commits. Then delete the merged feature branch. 

### Tag the release

Switch to the `main` branch locally. Then pull in updates:

```bash
git checkout main
git pull origin main
```

We use an [annotated tag](https://git-scm.com/book/en/v2/Git-Basics-Tagging) to tag the release.
Tag the release, following the format below:

```bash
git tag -a vX.y.z -m "Release X.y.z. Add filtering."
```

Things to note:
- This project uses [SemVer](https://semver.org/) version numbers.
- the annotated version tag __starts with a `v`__
- the message is short, starts with "Release X.y.z." and a concise message of the change.

Push tags:

```bash
git push origin main --tags
```
