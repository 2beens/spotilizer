function downloadFavTracksSnapshots() {
    var cookieID = getCookie("spotilizer-user-id");
    console.log(' > populating fav tracks snapshots. cookie ID: ' + cookieID);
    $.ajax({
        url: '/api/ssfavtracks',
        headers: {
            'Authorization': 'Bearer ' + window.accessToken
        },
        success: function(response) {
            var responseObj = JSON.parse(response);
            if (responseObj.error) {
                console.error(' >>> populate fav tracks snapshots error: ' + responseObj.error.message);
                toastr.error(responseObj.error.message, 'Populate favorite tracks snapshots error');
                return;
            }
            ssTracks = responseObj.data;
            localStorage.setItem("ssTracks", JSON.stringify(ssTracks));
            populateFavTracksSnapshots();
        },
        error: function(xhr,status,error) {
            console.error(' >>> populate playlists snapshots, status: ' + status + ', error: ' + error);
        }
    });
}

function downloadPlaylistSnapshots() {
    var cookieID = getCookie("spotilizer-user-id");
    console.log(' > populating playlistts snapshots. cookie ID: ' + cookieID);

    $.ajax({
        url: '/api/ssplaylists',
        headers: {
            'Authorization': 'Bearer ' + window.accessToken
        },
        success: function(response) {
            var responseObj = JSON.parse(response);
            console.warn(' > received playlists:')
            console.warn(responseObj);
            if (responseObj.error) {
                console.error(' >>> populate playlists snapshots error: ' + responseObj.error.message);
                toastr.error(responseObj.error.message, 'Populate playlists snapshots error');
                return;
            }
            ssPlaylists = responseObj.data;
            localStorage.setItem("ssPlaylists", JSON.stringify(ssPlaylists));
            populatePlaylistSnapshots();
        },
        error: function(xhr,status,error) {
            console.error(' >>> populate playlists snapshots, status: ' + status + ', error: ' + error);
        }
    });
}

function populateFavTracksSnapshots() {
    var tracksSnapshotsUL = $('#tracks-snapshots');
    if (!tracksSnapshotsUL || !ssTracks) {
        return;
    }
    ssTracks.forEach(function(ts) {
        var timestamp = new Date(ts.timestamp * 1000);
        var timestampStr = timestamp.toISOString().slice(0, 19).replace('T', ' ');
        tracksSnapshotsUL.append(`<li onclick="showFavTracksSnapshot('${ts.timestamp}')">${timestampStr}: <span class="badge badge-info">${ts.tracks.length}</span> playlists</li>`);
    });
}

function populatePlaylistSnapshots() {
    var playlistSnapshotsUL = $('#playlist-snapshots');
    if (!playlistSnapshotsUL || !ssPlaylists) {
        return;
    }
    ssPlaylists.forEach(function(ps) {
        var timestamp = new Date(ps.timestamp * 1000);
        var timestampStr = timestamp.toISOString().slice(0, 19).replace('T', ' ');
        playlistSnapshotsUL.append(`<li onclick="showPlaylistSnapshot(${ps.timestamp})">${timestampStr}: <span class="badge badge-info">${ps.playlists.length}</span> playlists</li>`);
    });
}

function showPlaylistSnapshot(timestamp) {
    console.log(timestamp);
}

function showFavTracksSnapshot(timestamp) {
    console.log(timestamp);
}

function callDebug() {
    console.log(' > calling debug function on server ...');
    makeRequest('/debug', function(response) {
        console.log(' > received from server: ' + response);
    });
}

function saveCurrentTracks() {
    lastCalledFunc = saveCurrentTracks;
    makeRequest("/save_current_tracks", function(response) {
        console.log(' > received from server: ' + response);
        var respObj = JSON.parse(response);
        if(checkIfRefreshTokenNeeded(respObj)) {
            console.log(' > refresh token needed ...');
            refreshTokenFunc();
            return;
        }
        if (respObj.error) {
            toastr.error(respObj.error.message, 'Save fav tracks error');
        } else {
            toastr.success(respObj.message, 'Save fav tracks');
        }
    });
}

function saveCurrentPlaylists() {
    lastCalledFunc = saveCurrentPlaylists;
    makeRequest("/save_current_playlists", function(response) {
        console.log(' > received from server: ' + response);
        var respObj = JSON.parse(response);
        if(checkIfRefreshTokenNeeded(respObj)) {
            console.log(' > refresh token needed ...');
            refreshTokenFunc();
            return;
        }
        if (respObj.error) {
            toastr.error(respObj.error.message, 'Save current playlists error');
        } else {
            toastr.success(respObj.message, 'Save current playlists');
        }
    }, function(xhr, status, error) {
        toastr.error("Status: " + status + ", error: " + JSON.stringify(error), 'Save current playlists error');
    });
}

var ssPlaylists = null;
var ssTracks = null;
function getDataFromLocalStorage() {
    ssPlaylistsJson = localStorage.getItem("ssPlaylists");
    ssTracksJson = localStorage.getItem("ssTracks");
    countPlaylists = 0;
    countTracks = 0;
    if (ssPlaylistsJson) {
        ssPlaylists = JSON.parse(ssPlaylistsJson);
        countPlaylists = ssPlaylists.length;
    }
    if (ssTracksJson) {
        ssTracks = JSON.parse(ssTracksJson);
        countTracks = ssTracks.length;
    }
    console.info(` > gotten [${countPlaylists}] playlists and [${countTracks}] fav tracks from dataStorage`)
}

function refreshData() {
    // start both population functions in parallel
    setTimeout(function() {
        if (isLoggedIn()) {
            downloadFavTracksSnapshots();
        }
    }, 400);
    setTimeout(function() {
        if (isLoggedIn()) {
            downloadPlaylistSnapshots();
        }
    }, 650);
}

document.addEventListener('DOMContentLoaded', function(event) {
    if (isLoggedIn()) {
        $('#spotify-controls-div').removeClass('invisible-elem');
        $('#refresh-button-div').removeClass('invisible-elem');
        $('#snapshots-data').removeClass('invisible-elem');
        getDataFromLocalStorage();
        populateFavTracksSnapshots();
        populatePlaylistSnapshots();
    } else {
        $('#spotify-controls-div').addClass('invisible-elem');
        $('#refresh-button-div').addClass('invisible-elem');
        $('#snapshots-data').addClass('invisible-elem');
    }
});
