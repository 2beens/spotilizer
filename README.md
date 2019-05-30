# spotilizer
Study **golang** through using **Spotify API**.

A small project which aims to help managing personal Spotify lists and tracks, but mostly to do what I love to do and study **Go** in the process.

Came to my mind that I could find this useful, because a lot of times it hapens that I screw up my lists (by removing a song from fav songs list), without ability to undo the operation. Sadly, Spotify Desktop Client App does not have undo for that operation.
Also, in the future, it can be made to transfer playlists from other Music Services like Deezer, YouTube, etc.

## Install & Run

Make sure `golang` is properly installed and set.

``` sh
go get github.com/2beens/spotilizer`
cd $GOPATH/src/github.com/2beens/spotilizer`
go get ./...
go install
spotilizer
```
