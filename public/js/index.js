function populatePlaylistSnapshots() {
    var playlistSnapshotsUL = $('#playlist-snapshots');
    if (!playlistSnapshotsUL) {
        return;
    }
    var cookieID = getCookie("spotilizer-user-id");
    console.log(' > populating playlistts snapshots. cookie ID: ' + cookieID);

    $.ajax({
        url: '/api/ssplaylists',
        headers: {
            'Authorization': 'Bearer ' + window.accessToken
        },
        success: function(response) {
            var responseObj = JSON.parse(response);
            if (responseObj.error) {
                console.error(' >>> populate playlists snapshots error: ' + responseObj.error.message);
                toastr.error(responseObj.error.message, 'Populate playlists snapshots error');
                return;
            }
            responseObj.data.forEach(function(ps) {
                var timestamp = new Date(ps.timestamp * 1000);
                var timestampStr = timestamp.toISOString().slice(0, 19).replace('T', ' ');
                playlistSnapshotsUL.append(`<li><a href="#">${timestampStr}: <span class="badge badge-info">${ps.playlists.length}</span> playlists</a></li>`);
            });
        },
        error: function(xhr,status,error) {
            console.error(' >>> populate playlists snapshots, status: ' + status + ', error: ' + error);
        }
    });
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

document.addEventListener('DOMContentLoaded', function(event) {
    setTimeout(function() {
        if (isLoggedIn()) {
            populatePlaylistSnapshots();
        }
    }, 500);

    if (isLoggedIn()) {
        $('#spotify-controls-div').removeClass('invisible-elem');
        $('#playlists-data').removeClass('invisible-elem');
    } else {
        $('#spotify-controls-div').addClass('invisible-elem');
        $('#playlists-data').addClass('invisible-elem');
    }
});
