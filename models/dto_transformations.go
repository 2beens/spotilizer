package models

func SpPlaylist2dtoPlaylist(spPlaylist SpPlaylist, tracks []SpPlaylistTrack) DTOPlaylist {
	dtoPlaylist := DTOPlaylist{
		ID:         spPlaylist.ID,
		Name:       spPlaylist.Name,
		URI:        spPlaylist.URI,
		TracksHref: spPlaylist.Tracks.Href,
		Tracks:     []DTOTrack{},
	}

	for _, t := range tracks {
		dtoPlaylist.Tracks = append(dtoPlaylist.Tracks, SpPlaylistTrack2dtoPlaylistTrack(t))
	}

	return dtoPlaylist
}

func SpPlaylistTrack2dtoPlaylistTrack(spPlTrack SpPlaylistTrack) DTOTrack {
	dtoTrack := DTOTrack{
		AddedAt:     spPlTrack.AddedAt.Unix(),
		AddedBy:     spPlTrack.AddedBy.ID,
		ID:          spPlTrack.Track.ID,
		DurationMs:  spPlTrack.Track.DurationMs,
		Name:        spPlTrack.Track.Name,
		TrackNumber: spPlTrack.Track.TrackNumber,
		URI:         spPlTrack.Track.URI,
		Artists:     SpArtists2dtoArtists(spPlTrack.Track.Artists),
	}
	return dtoTrack
}

func SpAddedTrack2dtoTrack(spAddedTrack SpAddedTrack) DTOTrack {
	return DTOTrack{
		AddedAt:     spAddedTrack.AddedAt.Unix(),
		ID:          spAddedTrack.Track.ID,
		DurationMs:  spAddedTrack.Track.DurationMs,
		Name:        spAddedTrack.Track.Name,
		TrackNumber: spAddedTrack.Track.TrackNumber,
		URI:         spAddedTrack.Track.URI,
		Artists:     SpArtists2dtoArtists(spAddedTrack.Track.Artists),
	}
}

func SpTrack2dtoTrack(spTrack SpTrack) DTOTrack {
	return DTOTrack{
		ID:          spTrack.ID,
		DurationMs:  spTrack.DurationMs,
		Name:        spTrack.Name,
		TrackNumber: spTrack.TrackNumber,
		URI:         spTrack.URI,
		Artists:     SpArtists2dtoArtists(spTrack.Artists),
	}
}

func SpArtist2dtoArtist(spArtist SpArtist) DTOArtist {
	return DTOArtist{
		Href: spArtist.Href,
		Name: spArtist.Name,
		Type: spArtist.Type,
	}
}

func SpArtists2dtoArtists(spArtists []SpArtist) []DTOArtist {
	var dtoArtists []DTOArtist
	for _, spA := range spArtists {
		dtoArtists = append(dtoArtists, SpArtist2dtoArtist(spA))
	}
	return dtoArtists
}
