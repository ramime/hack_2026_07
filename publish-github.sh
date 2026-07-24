#!/usr/bin/env bash
# publish-github.sh — Squash merge current release to public-release branch & push to public GitHub
set -e

VERSION_TAG="$1"

if [ -z "$VERSION_TAG" ]; then
    echo "Usage: ./publish-github.sh <version_tag> (e.g. ./publish-github.sh v0.0.1)"
    exit 1
fi

echo "==> Preparing public release for version: ${VERSION_TAG}..."

# Ensure working tree is clean
if [ -n "$(git status --porcelain)" ]; then
    echo "Error: Working tree is not clean. Commit all changes to private branch first."
    exit 1
fi

CURRENT_BRANCH=$(git branch --show-current)

# Create public-release orphan branch if it does not exist
if ! git rev-parse --verify public-release >/dev/null 2>&1; then
    echo "==> Creating orphan branch 'public-release'..."
    git checkout --orphan public-release
    git commit --allow-empty -m "Initial public release branch"
    git checkout "$CURRENT_BRANCH"
fi

echo "==> Merging '${CURRENT_BRANCH}' onto 'public-release' with --squash..."
git checkout public-release
git merge "$CURRENT_BRANCH" --squash --allow-unrelated-histories -m "Release ${VERSION_TAG}: AgencyPulse MVP"

echo "==> Setting tag ${VERSION_TAG}..."
git tag -a -f "${VERSION_TAG}" -m "Public Release ${VERSION_TAG}"

echo "==> Pushing to public GitHub remote 'github'..."
git push -u github public-release:main --force
git push github "${VERSION_TAG}" --force

echo "==> Returning to branch '${CURRENT_BRANCH}'..."
git checkout "$CURRENT_BRANCH"

echo "==> Public GitHub release ${VERSION_TAG} published successfully!"
