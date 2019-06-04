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

document.addEventListener('DOMContentLoaded', function(event) {
    setTimeout(function() {
        if (isLoggedIn()) {
            populatePlaylistSnapshots();
        }
    }, 500);
});
