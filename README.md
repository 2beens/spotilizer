# :notes: spotilizer
[![Build Status](https://travis-ci.com/2beens/spotilizer.svg?branch=master)](https://travis-ci.com/2beens/spotilizer)
[![GitHub license](https://img.shields.io/github/license/Naereen/StrapDown.js.svg)](https://github.com/Naereen/StrapDown.js/blob/master/LICENSE)
![golangci-lint](https://github.com/golangci/golangci-web/blob/master/src/assets/images/badge_a_plus_flat.svg "Using GolangCI Lint")
[![Maintenance](https://img.shields.io/badge/Maintained%3F-no-red.svg)](https://GitHub.com/Naereen/StrapDown.js/graphs/commit-activity)
[![made-for-VSCode](https://img.shields.io/badge/Made%20for-VSCode-1f425f.svg)](https://code.visualstudio.com/)

Study **golang** through using **Spotify API** and **Redis**.

A small project which aims to help managing personal Spotify lists and tracks, but mostly to do what I love to do and study **Go** in the process.

Came to my mind that I could find this useful, having a history/snapshots of my Spotify Data, because a lot of times it hapens that I screw up my lists (by removing a song from fav songs list), without ability to undo the operation, and also not remembering the name of the removed song. Sadly, Spotify Desktop/Browser Client App does not have undo for that operation. Also, in the future, it can be made to transfer playlists from other Music Services like Deezer, YouTube, etc.

## :wrench: Install & Run

### :one: Set Spotify App
Make sure you have a **Spotify App**, through which `spotilizer` interacts with Spotify API. That goes out of this scope, so not gonna explain that here. It's all nicelly explained at this location:

https://developer.spotify.com/documentation/web-api/

After Spotify App is created, now we need to 2 set `env. variables` like so:
```sh
export SPOTIFY_CLIENT_ID=<your_spotify_app_client_id_here>
export SPOTIFY_CLIENT_SECRET=<your_spotify_app_client_secret_here>
```

### :two: Redis - install and start service
Make sure `Redis` is installed, up and running. See: https://redis.io/

Maybe more convenient, you can also install and run it using `brew`:
``` sh
brew install redis
```
Then just run the service :
``` sh
brew services start redis
```

Spotilizer uses default Redis setup, so just starting the service is enough.

### :three: Get Spotilizer
Make sure `golang` is properly installed and set: https://golang.org/doc/install

``` sh
go get github.com/2beens/spotilizer`
cd $GOPATH/src/github.com/2beens/spotilizer`
go get ./...
```

Next, we need to get frontend dependencies via PMS `yarn`, which should be installed on your machine: https://yarnpkg.com/en/docs/install

We are still in the spotilizer project dir, now run:
``` sh
yarn install
```

### :four: Run Spotilizer

``` sh
go install
spotilizer
```

By default, logger output is terminal (can be changed to file. see source code `main.go` for more info).

### :five: Web Client
:point_right: Open browser (Chrome, ofc) and go to: http://localhost:8080

:point_right: Click at just about anything you see there :collision:

:point_right: Observe terminal output for what happens in the server

## :mag_right: Testing, Static code analysis ...
#### Linter
`golangci-lint` is used. Make sure it's installed: [GolangCI Lint Installation](https://github.com/golangci/golangci-lint#install). 

Run it like so:
``` sh
golangci-lint run
```

#### Unit Tests
Still work in progress, thus not many unit tests exists. Nevertheless:
``` sh
go test -v -cover ./...
```
