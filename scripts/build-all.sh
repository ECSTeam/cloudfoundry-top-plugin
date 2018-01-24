#!/bin/sh
#
# To use this script, install github-release via:
#   go get github.com/aktau/github-release
# Then verify GITHUB_TOKEN variable is set
#   export GITHUB_TOKEN=XXXXXXX (your personal token)
#
# Example:
# ./scripts/build-all.sh release v0.7.3
#
#
set -e
export GITHUB_USER=ecsteam
export GITHUB_REPO=cloudfoundry-top-plugin

if [[ "$1" = "release" ]] ; then
	TAG="$2"
	: ${TAG:?"Usage: build_all.sh [release] [TAG]"}


	if $GOPATH/bin/github-release info --tag $TAG > /dev/null 2>&1 ; then
		echo "$TAG exists, remove it or increment"
		exit 1
	else
		MAJOR=`echo $TAG | sed 's/^v//' | awk 'BEGIN {FS = "." } ; { printf $1;}'`
		MINOR=`echo $TAG | sed 's/^v//' | awk 'BEGIN {FS = "." } ; { printf $2;}'`
		BUILD=`echo $TAG | sed 's/^v//' | awk 'BEGIN {FS = "." } ; { printf $3;}'`
		VERSION=`echo $TAG | sed 's/^v//'`

		`sed -i "" -e "1,/Major:.*/s/Major:.*/Major: $MAJOR,/" \
			-e "1,/Minor:.*/s/Minor:.*/Minor: $MINOR,/" \
			-e "1,/Build:.*/s/Build:.*/Build: $BUILD,/" main.go`
	fi
fi

LINUX32_FILENAME="top-plugin-linux32"
LINUX64_FILENAME="top-plugin-linux64"

MAC_FILENAME="top-plugin-darwin"

WIN32_FILENAME="top-plugin32.exe"
WIN64_FILENAME="top-plugin64.exe"

LINUX_ARM6_FILENAME="top-plugin-linux-arm6"

echo "Compile ${LINUX32_FILENAME}"
GOOS=linux GOARCH=386 go build -o "${LINUX32_FILENAME}"
LINUX32_SHA1=`cat $LINUX32_FILENAME | openssl sha1`
mkdir -p bin/linux32
mv $LINUX32_FILENAME bin/linux32

echo "Compile ${LINUX64_FILENAME}"
GOOS=linux GOARCH=amd64 go build -o "${LINUX64_FILENAME}"
LINUX64_SHA1=`cat $LINUX64_FILENAME | openssl sha1`
mkdir -p bin/linux64
mv $LINUX64_FILENAME bin/linux64

echo "Compile ${MAC_FILENAME}"
GOOS=darwin GOARCH=amd64 go build -o ${MAC_FILENAME}
OSX_SHA1=`cat $MAC_FILENAME | openssl sha1`
mkdir -p bin/osx
mv $MAC_FILENAME bin/osx

echo "Compile ${WIN32_FILENAME}"
GOOS=windows GOARCH=386 go build -o ${WIN32_FILENAME}
WIN32_SHA1=`cat $WIN32_FILENAME | openssl sha1`
mkdir -p bin/win32
mv $WIN32_FILENAME bin/win32

echo "Compile ${WIN64_FILENAME}"
GOOS=windows GOARCH=amd64 go build -o ${WIN64_FILENAME}
WIN64_SHA1=`cat $WIN64_FILENAME | openssl sha1`
mkdir -p bin/win64
mv $WIN64_FILENAME bin/win64

echo "Compile ${LINUX_ARM6_FILENAME}"
GOOS=linux GOARCH=arm GOARM=6 go build -o $LINUX_ARM6_FILENAME
LINUXARM6_SHA1=`cat $LINUX_ARM6_FILENAME | openssl sha1`
mkdir -p bin/linux_arm6
mv $LINUX_ARM6_FILENAME bin/linux_arm6

echo "Compile complete"
NOW=`TZ=UC date +'%Y-%m-%dT%TZ'`

cat repo-index.yml |
sed "s/__osx-sha1__/$OSX_SHA1/" |
sed "s/__win32-sha1__/$WIN32_SHA1/" |
sed "s/__win64-sha1__/$WIN64_SHA1/" |
sed "s/__linux32-sha1__/$LINUX32_SHA1/" |
sed "s/__linux64-sha1__/$LINUX64_SHA1/" |
sed "s/__linuxarm6-sha1__/$LINUXARM6_SHA1/" |
sed "s/__TAG__/$TAG/" |
sed "s/__VERSION__/$VERSION/" |
sed "s/__TODAY__/$NOW/" |
cat

if [[ "$1" = "release" ]] ; then

	echo "Git commit and push"

	git commit -am "Build version $TAG"
	git push

	echo "Create git release $TAG"

	$GOPATH/bin/github-release release \
    --tag $TAG \
    --name "Cloud Foundry top plugin $TAG" \
    --description "$TAG release"

	echo "Upload top-plugin-darwin"
	$GOPATH/bin/github-release upload \
    --tag $TAG \
    --name "top-plugin-darwin" \
    --file bin/osx/top-plugin-darwin

	echo "Upload top-plugin-linux32"
	$GOPATH/bin/github-release upload \
    --tag $TAG \
    --name "top-plugin-linux32" \
    --file bin/linux64/top-plugin-linux32

	echo "Upload top-plugin-linux64"
	$GOPATH/bin/github-release upload \
    --tag $TAG \
    --name "top-plugin-linux64" \
    --file bin/linux64/top-plugin-linux64

	echo "Upload top-plugin32.exe"
	$GOPATH/bin/github-release upload \
    --tag $TAG \
    --name "top-plugin32.exe" \
    --file bin/win64/top-plugin32.exe

	echo "Upload top-plugin64.exe"
	$GOPATH/bin/github-release upload \
    --tag $TAG \
    --name "top-plugin64.exe" \
    --file bin/win64/top-plugin64.exe

	echo "Upload top-plugin-linux-arm6"
	$GOPATH/bin/github-release upload \
    --tag $TAG \
    --name "top-plugin-linux-arm6" \
    --file bin/linux_arm6/top-plugin-linux-arm6

	#git commit -am "Build version $TAG"
	#git tag -a $TAG -m "Top Plugin v$TAG"
	#echo "Tagged release, 'git push --follow-tags' to push it to github, upload the binaries to github"
	echo "copy the output above to the cli repo you plan to deploy in"
fi
