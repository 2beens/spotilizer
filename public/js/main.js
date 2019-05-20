function getCookie(cname) {
    var name = cname + "=";
    var decodedCookie = decodeURIComponent(document.cookie);
    var ca = decodedCookie.split(';');
    for (var i = 0; i < ca.length; i++) {
        var c = ca[i];
        while (c.charAt(0) == ' ') {
            c = c.substring(1);
        }
        if (c.indexOf(name) == 0) {
            return c.substring(name.length, c.length);
        }
    }
    return "";
}

function setCookie(cname, cvalue, daysValid) {
    console.log(` > setting cookie [${cname}] to new value`);
    var now = new Date();
    now.setTime(now.getTime() + (daysValid * 24 * 60 * 60 * 1000));
    var expires = "expires=" + now.toUTCString();
    document.cookie = cname + "=" + cvalue + ";" + expires + ";path=/";
}

function printSelfSpotifyInfo() {
    var spotilizerId = getCookie("spotilizer-user-id");
    console.log(' > spotilizer-user-id: ' + spotilizerId);
    makeRequest('https://api.spotify.com/v1/me', function (response) {
        console.log('------------------------------- response from spotify ------------')
        console.log(response);
        console.log('------------------------------------------------------------------')
    });
}

function makeRequest(queryUrl, callback) {
    console.log(' ---> making a request call to: ' + queryUrl);
    $.ajax({
        url: queryUrl,
        headers: {
            'Authorization': 'Bearer ' + window.accessToken
        },
        success: callback,
    });
}

function saveCurrentPlaylists() {
    makeRequest("/save_current_playlists", function(response) {
        console.log(' > received from server: ' + response);
        $('#query-result').val(response);
    });
}

function executeUrlQueryRequest() {
    var queryUrl = $('#query-text').val();
    if (!queryUrl) {
        console.log(' > query URL error ...');
        return;
    }
    console.log(' > URL: ' + queryUrl);
    makeRequest(queryUrl, function(response) {
        console.log(' > received from Spotify: ' + response);
        $('#query-result').val(response);
    });
}

function getFavPlaylist() {
    console.log(' > getting fav playlist ...');
    $.ajax({
        url: 'https://api.spotify.com/v1/me/tracks',
        headers: {
            'Accept': 'application/json',
            'Content-Type': 'application/json',
            'Authorization': 'Bearer ' + window.accessToken
        },
        success: function(response) {
            console.log(response);
        },
    });
}

(function () {
    window.accessToken = getCookie("accessToken");
    window.refreshToken = getCookie("refreshToken");
    console.log(" > loaded cookie AT: " + window.accessToken);
    console.log(" > loaded cookie RT: " + window.refreshToken);
    console.log(' > main script finished');
})()