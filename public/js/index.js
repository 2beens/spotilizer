function downloadFavTracksSnapshots() {
    const cookieID = getCookie("spotilizer-user-id");
    if (!cookieID) {
        toastr.error('Not logged in.', 'Populate favorite tracks snapshots error');
        return;
    }
    console.log(' > populating fav tracks snapshots. cookie ID: ' + cookieID);
    $.ajax({
        url: '/api/ssfavtracks',
        success: function(response) {
            const responseObj = JSON.parse(response);
            if (responseObj.error) {
                console.error(' >>> populate fav tracks snapshots error: ' + responseObj.error.message);
                toastr.error(responseObj.error.message, 'Populate favorite tracks snapshots error');
                return;
            }
            ssTracks = responseObj.data;
            localStorage.setItem('ssTracks', JSON.stringify(ssTracks));
            fillFavTracksSnapshotsMap();
            populateFavTracksSnapshots();
        },
        error: function(xhr,status,error) {
            console.error(' >>> populate playlists snapshots, status: ' + status + ', error: ' + error);
        }
    });
}

function downloadPlaylistSnapshots() {
    const cookieID = getCookie('spotilizer-user-id');
    if (!cookieID) {
        toastr.error('Not logged in.', 'Populate playlists snapshots error');
        return;
    }
    console.log(' > populating playlistts snapshots. cookie ID: ' + cookieID);
    $.ajax({
        url: '/api/ssplaylists',
        success: function(response) {
            const responseObj = JSON.parse(response);
            if (responseObj.error) {
                console.error(' >>> populate playlists snapshots error: ' + responseObj.error.message);
                toastr.error(responseObj.error.message, 'Populate playlists snapshots error');
                return;
            }
            ssPlaylists = responseObj.data;
            localStorage.setItem("ssPlaylists", JSON.stringify(ssPlaylists));
            fillPlaylistsSnapshotsMap();
            populatePlaylistSnapshots();
        },
        error: function(xhr,status,error) {
            console.error(' >>> populate playlists snapshots, status: ' + status + ', error: ' + error);
        }
    });
}


function playlistSnapshotIsEmpty(snapshot) {
    if (!snapshot || snapshot.length === 0) {
        return true;
    }
    let isEmpty = true;
    snapshot.forEach(function(p) { 
        if (p.tracks.length > 0) {
            isEmpty = false;
            return;
        }
    });
    return isEmpty;
}


function getPlaylistsSnapshot(timestamp, callback) {
    let playlistsSnapshot = ssTimestamp2playlistsMap.get(timestamp);
    if (!playlistSnapshotIsEmpty(playlistsSnapshot)) {
        callback(playlistsSnapshot);
        return;
    }

    const cookieID = getCookie("spotilizer-user-id");
    if (!cookieID) {
        toastr.error('Not logged in.', 'Get fav tracks snapshot error');
        return;
    }
    console.log(' > get fav tracks snapshot. cookie ID: ' + cookieID);
    $.ajax({
        url: '/api/ssplaylists/' + timestamp,
        type: 'GET',
        success: function(response) {
            const responseObj = JSON.parse(response);
            if (responseObj.error) {
                console.error(' >>> get playlist snapshot error: ' + responseObj.error.message);
                toastr.error(responseObj.error.message, 'Get playlist snapshot error');
            } else {
                playlistsSnapshot = responseObj.data.playlists;
                ssTimestamp2playlistsMap.set(timestamp, playlistsSnapshot);
                toastr.success(responseObj.message, 'Get playlists snapshot');
            }
            callback(playlistsSnapshot);
        },
        error: function(xhr,status,error) {
            console.error(' >>> get playlist snapshot error, status: ' + status + ', error: ' + error);
            callback(playlistsSnapshot);
        }
    });
}

function getFavTracksSnapshot(timestamp, callback) {
    let tracksSnapshot = ssTimestamp2tracksMap.get(timestamp);
    if (tracksSnapshot && tracksSnapshot.length > 0) {
        callback(tracksSnapshot);
        return;
    }

    const cookieID = getCookie("spotilizer-user-id");
    if (!cookieID) {
        toastr.error('Not logged in.', 'Get fav tracks snapshot error');
        return;
    }
    console.log(' > get fav tracks snapshot. cookie ID: ' + cookieID);
    $.ajax({
        url: '/api/ssfavtracks/' + timestamp,
        type: 'GET',
        success: function(response) {
            const responseObj = JSON.parse(response);
            if (responseObj.error) {
                console.error(' >>> get fav. tracks snapshot error: ' + responseObj.error.message);
                toastr.error(responseObj.error.message, 'Get fav. tracks snapshot error');
            } else {
                tracksSnapshot = responseObj.data.tracks;
                ssTimestamp2tracksMap.set(timestamp, tracksSnapshot);
                toastr.success(responseObj.message, 'Get favorite tracks snapshot');
            }
            callback(tracksSnapshot);
        },
        error: function(xhr,status,error) {
            console.error(' >>> get fav. tracks snapshot error, status: ' + status + ', error: ' + error);
            callback(tracksSnapshot);
        }
    });
}

function deleteFavTracksSnapshot(timestamp) {
    const cookieID = getCookie('spotilizer-user-id');
    if (!cookieID) {
        toastr.error('Not logged in.', 'Delete fav tracks snapshot error');
        return;
    }
    console.log(' > delete fav tracks snapshot. cookie ID: ' + cookieID);
    $.ajax({
        url: '/api/ssfavtracks/' + timestamp,
        type: 'DELETE',
        success: function(response) {
            console.log(' > received from server:');
            console.log(response);
            const responseObj = JSON.parse(response);
            if (responseObj.error) {
                console.error(' >>> delete fav. tracks snapshot error: ' + responseObj.error.message);
                toastr.error(responseObj.error.message, 'Delete fav. tracks snapshot error');
            } else {
                downloadFavTracksSnapshots();
                toastr.success(responseObj.message, 'Delete favorite tracks snapshot');
            }
        },
        error: function(xhr,status,error) {
            console.error(' >>> delete fav. tracks snapshot error, status: ' + status + ', error: ' + error);
        }
    });
}

function getFavTracksSnapshotDiff(timestamp) {
    const cookieID = getCookie('spotilizer-user-id');
    if (!cookieID) {
        toastr.error('Not logged in.', 'Get fav tracks snapshot diff error');
        return;
    }
    console.log(' > get fav tracks snapshot diff. cookie ID: ' + cookieID);
    $.ajax({
        url: '/api/ssfavtracks/diff/' + timestamp,
        success: function(response) {
            console.log(' > received from server:');
            console.log(response);
            const responseObj = JSON.parse(response);
            if (responseObj.error) {
                console.error(' >>> get fav. tracks snapshot diff error: ' + responseObj.error.message);
                toastr.error(responseObj.error.message, 'Get fav. tracks snapshot diff error');
            } else {
                showFavTracksDiff(timestamp, responseObj.data.newTracks, responseObj.data.removedTracks)
                toastr.success(responseObj.message, 'Get favorite tracks snapshot diff');
            }
        },
        error: function(xhr,status,error) {
            console.error(' >>> get fav. tracks snapshot diff error, status: ' + status + ', error: ' + error);
        }
    });
}

function deletePlaylistSnapshot(timestamp) {
    const cookieID = getCookie('spotilizer-user-id');
    if (!cookieID) {
        toastr.error('Not logged in.', 'Delete playlist snapshot error');
        return;
    }
    console.log(' > delete playlist snapshot. cookie ID: ' + cookieID);
    $.ajax({
        url: '/api/ssplaylists/' + timestamp,
        type: 'DELETE',
        success: function(response) {
            console.log(' > received from server:');
            console.log(response);
            const responseObj = JSON.parse(response);
            if (responseObj.error) {
                console.error(' >>> delete playlist snapshot error: ' + responseObj.error.message);
                toastr.error(responseObj.error.message, 'Delete playlist snapshot error');
            } else {
                downloadPlaylistSnapshots();
                toastr.success(responseObj.message, 'Delete playists snapshot');
            }
        },
        error: function(xhr,status,error) {
            console.error(' >>> delete playlist snapshot error, status: ' + status + ', error: ' + error);
        }
    });
}

function fillFavTracksSnapshotsMap() {
    ssTimestamp2tracksMap.clear();
    if (!ssTracks) {
        return;
    }
    ssTracks.forEach(function(tracksSnapshot) {
        ssTimestamp2tracksMap.set(tracksSnapshot.timestamp, tracksSnapshot.tracks);
    });
}

function fillPlaylistsSnapshotsMap() {
    ssTimestamp2playlistsMap.clear();
    if (!ssPlaylists) {
        return;
    }
    ssPlaylists.forEach(function(playlistsSnapshot) {
        ssTimestamp2playlistsMap.set(playlistsSnapshot.timestamp, playlistsSnapshot.playlists);
    });
}

function fillSnapshotsMaps() {
    fillFavTracksSnapshotsMap();
    fillPlaylistsSnapshotsMap();
}

function populateFavTracksSnapshots() {
    const tracksSnapshotsUL = $('#tracks-snapshots');
    if (!tracksSnapshotsUL || !ssTracks) {
        return;
    }
    tracksSnapshotsUL.empty();
    ssTracks.forEach(function(ts) {
        const timestamp = new Date(ts.timestamp * 1000);
        const timestampStr = timestamp.toISOString().slice(0, 19).replace('T', ' ');
        tracksSnapshotsUL.append(`
            <li style="list-style-type:none;">
            <div class="row">
                <div class="snapshot-item col-sm-10" onclick="showFavTracksSnapshot(${ts.timestamp})">
                    ${timestampStr} <span class="badge badge-info" style="margin-left: 15px;">${ts.tracks_count}</span> tracks
                </div>
                <div class="col-sm-1">
                    <button style="height: 20px; padding-top: 0px;" type="button" class="btn btn-danger btn-sm" onclick="deleteFavTracksSnapshot(${ts.timestamp})">Del</button>
               </div>
               <div class="col-sm-1">
                    <button style="height: 20px; padding-top: 0px; margin-left: 5px;" type="button" class="btn btn-info btn-sm" onclick="getFavTracksSnapshotDiff(${ts.timestamp})">Diff</button>
                </div>
            </div>
            </li>`);
    });
}



function populatePlaylistSnapshots() {
    const playlistSnapshotsUL = $('#playlist-snapshots');
    if (!playlistSnapshotsUL || !ssPlaylists) {
        return;
    }
    playlistSnapshotsUL.empty();
    ssPlaylists.forEach(function(ps) {
        const timestamp = new Date(ps.timestamp * 1000);
        const timestampStr = timestamp.toISOString().slice(0, 19).replace('T', ' ');
        playlistSnapshotsUL.append(`
            <li style="list-style-type:none;">
            <div class="row">
                <div class="snapshot-item col-sm-9" onclick="showPlaylistSnapshot(${ps.timestamp})">
                    ${timestampStr} <span class="badge badge-info" style="margin-left: 15px;">${ps.playlists.length}</span> playlists
                </div>
                <div class="col-sm-3">
                    <button style="height: 20px; padding-top: 0px;" type="button" class="btn btn-danger btn-sm" onclick="deletePlaylistSnapshot(${ps.timestamp})">Del</button>
                </div>
            </div>
            </li>`);
    });
}

function showPlaylistSnapshot(timestamp) {
    getPlaylistsSnapshot(timestamp, function(playlistsSnapshot) {
        const ssDetailsCol = $('#snapshot-diff-col');
        ssDetailsCol.empty();
        const ssDetailsList = $('#snapshot-details-ul');
        ssDetailsList.empty();
        playlistsSnapshot.forEach(function(p) {
            ssDetailsList.append(`
                <li class="list-group-item d-flex justify-content-between align-items-center">
                    ${p.name}
                    <span style="margin-left: 20px;" class="badge badge-primary badge-pill">${p.tracks.length}</span>
                </li>
            `);
        });
    });
}

function showFavTracksSnapshot(timestamp) {
    getFavTracksSnapshot(timestamp, function(tracksSnapshot) {
        const ssDetailsCol = $('#snapshot-diff-col');
        ssDetailsCol.empty();
        const ssDetailsList = $('#snapshot-details-ul');
        ssDetailsList.empty();
        tracksSnapshot.forEach(function(t) {
            ssDetailsList.append(`
                <li class="list-group-item d-flex justify-content-between align-items-center">
                    ${getArtistsName(t)} - ${t.track.name}
                    <span style="margin-left: 20px;" class="badge badge-primary badge-pill">${new Date(t.added_at).toLocaleString()}</span>
                </li>
            `);
        });
    });
}

function showFavTracksDiff(timestamp, newTracks, removedTracks) {
    const ssDiffCol = $('#snapshot-diff-col');
    ssDiffCol.empty();
    const ssDetailsList = $('#snapshot-details-ul');
    ssDetailsList.empty();

    if (!newTracks) {
        newTracks = [];
    }
    if (!removedTracks) {
        removedTracks = [];
    }

    const timestampDate = new Date(timestamp * 1000);
    const timestampStr = timestampDate.toISOString().slice(0, 19).replace('T', ' ');
    ssDiffCol.append(`<h4>Changes since ${timestampStr}</h4>`);

    ssDiffCol.append(`<h5 style="margin-top: 20px;">New Tracks</h5>`);
    if (newTracks.length === 0) {
        ssDiffCol.append(`<p>No new tracks</p>`);
    }
    ssDiffCol.append(`<ul class="list-group">`);
    newTracks.forEach(function(t) {
        ssDiffCol.append(`
            <li class="list-group-item d-flex justify-content-between align-items-center">
                ${getArtistsName(t)} - ${t.track.name}
                <span style="margin-left: 20px;" class="badge badge-primary badge-pill">${new Date(t.added_at).toLocaleString()}</span>
            </li>
         `);
    });
    ssDiffCol.append(`</ul>`);

    ssDiffCol.append(`<h5 style="margin-top: 20px;">Removed Tracks</h5>`);
    if (removedTracks.length === 0) {
        ssDiffCol.append(`<p>No removed tracks</p>`);
    }
    ssDiffCol.append(`<ul class="list-group">`);
    removedTracks.forEach(function(t) {
        ssDiffCol.append(`
            <li class="list-group-item d-flex justify-content-between align-items-center">
                ${getArtistsName(t)} - ${t.track.name}
                <span style="margin-left: 20px;" class="badge badge-primary badge-pill">${new Date(t.added_at).toLocaleString()}</span>
            </li>
        `);
    });
    ssDiffCol.append(`</ul>`);
}

function getArtistsName(addedTrack) {
    return addedTrack.track.artists
        .map(function (a) {
            return a.name;
        })
        .join(', ');
}

function callDebug() {
    console.log(' > calling debug function on server ...');
    makeRequest('/debug', function(response) {
        console.log(' > received from server: ' + response);
    });
}

function saveCurrentTracks() {
    lastCalledFunc = saveCurrentTracks;
    makeRequest('/save_current_tracks', function(response) {
        console.log(' > received from server: ' + response);
        const respObj = JSON.parse(response);
        if(checkIfRefreshTokenNeeded(respObj)) {
            console.log(' > refresh token needed ...');
            refreshTokenFunc();
            return;
        }
        if (respObj.error) {
            toastr.error(respObj.error.message, 'Save fav tracks error');
            refreshData();
        } else {
            toastr.success(respObj.message, 'Save fav tracks');
        }
    });
}

function saveCurrentPlaylists() {
    lastCalledFunc = saveCurrentPlaylists;
    makeRequest('/save_current_playlists', function(response) {
        console.log(' > received from server: ' + response);
        const respObj = JSON.parse(response);
        if(checkIfRefreshTokenNeeded(respObj)) {
            console.log(' > refresh token needed ...');
            refreshTokenFunc();
            return;
        }
        if (respObj.error) {
            toastr.error(respObj.error.message, 'Save current playlists error');
        } else {
            toastr.success(respObj.message, 'Save current playlists');
            refreshData();
        }
    }, function(xhr, status, error) {
        toastr.error('Status: ' + status + ', error: ' + JSON.stringify(error), 'Save current playlists error');
    });
}

var ssPlaylists = null;
var ssTracks = null;
var ssTimestamp2playlistsMap = new Map();
var ssTimestamp2tracksMap = new Map();
function getDataFromLocalStorage() {
    ssPlaylistsJson = localStorage.getItem('ssPlaylists');
    ssTracksJson = localStorage.getItem('ssTracks');
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
    fillSnapshotsMaps();
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
