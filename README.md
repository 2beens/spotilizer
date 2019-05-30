# :notes: spotilizer
Study **golang** through using **Spotify API**.

A small project which aims to help managing personal Spotify lists and tracks, but mostly to do what I love to do and study **Go** in the process.

Came to my mind that I could find this useful, because a lot of times it hapens that I screw up my lists (by removing a song from fav songs list), without ability to undo the operation. Sadly, Spotify Desktop Client App does not have undo for that operation.
Also, in the future, it can be made to transfer playlists from other Music Services like Deezer, YouTube, etc.

## Install & Run

### :one: Set Spotify App
Make sure you have a **Spotify App**, through which `spotilizer` interacts with Spotify API. That goes out of this scope, so not gonna explain that here. It's all nicelly explained at this location:

https://developer.spotify.com/documentation/web-api/

After Spotify App is created, now we need to 2 set `env. variables` like so:
```sh
export SPOTIFY_CLIENT_ID=<your_spotify_app_client_id_here>
export SPOTIFY_CLIENT_SECRET=<your_spotify_app_client_secret_here>
```

### :two: Spotilizer get & run
Make sure `golang` is properly installed and set.

``` sh
go get github.com/2beens/spotilizer`
cd $GOPATH/src/github.com/2beens/spotilizer`
go get ./...
go install
spotilizer
```

By default, logger output is terminal (can be changed to file. see source code `main.go` for more info).

### :three: Web Client
:point_right: Open browser (Chrome, ofc) and go to: http://localhost:8080

:point_right: Click at just about anything you see there :collision:
