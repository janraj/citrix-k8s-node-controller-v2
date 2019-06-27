#!/bin/sh
version=1.0.0
update_version() {
  ver=$(cat ../version/VERSION)
  version=$ver
  major=$(echo $ver | cut -d. -f1)
  minor=$(echo $ver | cut -d. -f2)
  patch=$(echo $ver | cut -d. -f3)
  echo "Current Version $ver"
  if [[ ${TRAVIS_COMMIT_MESSAGE} =~ "[PATCH]" ]]; then
    let "patch=patch+1"
    echo "Major $major, Minor $minor, Patch $patch"
    version=$major.$minor.$patch
  elif [[ ${TRAVIS_COMMIT_MESSAGE} =~ "[MINOR]" ]]; then 
    let "patch=0"
    let "minor=minor+1"
    echo "Major $major, Minor $minor, Patch $patch"
    version=$major.$minor.$patch
  elif [[ ${TRAVIS_COMMIT_MESSAGE} =~ "[MAJOR]" ]]; then 
    let "patch=0"
    let "minor=0"
    let "major=major+1"
    echo "Major $major, Minor $minor, Patch $patch"
    version=$major.$minor.$patch
  fi
  echo "$version"
  echo "$version" > '../version/VERSION'
  export cnc_version=$version
}

git_setup() {
  git config --global user.email "travis@travis-ci.org"
  git config --global user.name "Travis CI"
}

git_commit() {
  git checkout master
  git add ../version/VERSION
  git commit -m "[skip ci] Travis update: $dateAndMonth (Build $TRAVIS_BUILD_NUMBER)"
}

git_push() {
  git remote rm origin
  git remote add origin https://${GH_TOKEN}@https://github.com/janraj/citrix-k8s-node-controller.git > /dev/null 2>&1
  git push origin master --quiet
}

update_version
echo "New Version is $version"
git_setup
git_commit
if [ $? -eq 0 ]; then
  echo "Commit is success, pushing new version to GitHub"
  git_push
else
  echo "Some issue with Git Commit"
fi
echo "$version"
