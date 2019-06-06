package models

func SpPlaylist2dtoPlaylist(spPlaylist SpPlaylist, tracks []SpPlaylistTrack) DTOPlaylist {
	dtoPlaylist := DTOPlaylist{
		ID:         spPlaylist.ID,
		Name:       spPlaylist.Name,
		URI:        spPlaylist.URI,
		TracksHref: spPlaylist.Tracks.Href,
		// Tracks:     tracks,
	}

	//TODO: tracks

	return dtoPlaylist
}

func SpAddedTrack2dtoAddedTrack(spAddedTrack SpAddedTrack) DTOAddedTrack {
	return DTOAddedTrack{
		AddedAt: spAddedTrack.AddedAt.Unix(),
		Track:   SpTrack2dtoTrack(spAddedTrack.Track),
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
	dtoArtists := []DTOArtist{}
	for _, spA := range spArtists {
		dtoArtists = append(dtoArtists, SpArtist2dtoArtist(spA))
	}
	return dtoArtists
}
