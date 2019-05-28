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

function eraseCookie(name) {   
    document.cookie = name+'=; Max-Age=-99999999;';  
}

function onLogout() {
    console.log(' > onLogout initiated ...');
    eraseCookie("spotilizer-user-id");
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

function saveCurrentTracks() {
    makeRequest("/save_current_tracks", function(response) {
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

function stringOK(val) {
    return val !== undefined && val !== null && val.length > 0;
}

(function () {
    var cookieID = getCookie("spotilizer-user-id");
    console.log(' > main script function: cookie ID: ' + cookieID);
    window.accessToken = getCookie("accessToken");
    window.refreshToken = getCookie("refreshToken");

    // remove unnecessary path info during login or logout
    if (window.location.pathname === '/callback' || window.location.pathname === '/logout') {
        history.pushState({}, null, "/");
    }
    
    // set navbar active button
    document.addEventListener("DOMContentLoaded", function(event) {
        $('.nav-bar-a').each(function(index) {
            $(this).removeClass('active');
        });
        if (window.location.pathname === '/about') {
            $('#nav-bar-about').addClass('active');
        } else if (window.location.pathname === '/contact') {
            $('#nav-bar-contact').addClass('active');
        } else {
            $('#nav-bar-home').addClass('active');
        }

        if (stringOK(cookieID) && stringOK(window.username)) {
            console.log(' > cookieID: ' + cookieID + ', username: ' + username);
            $('#nav-item-login').addClass('invisible-elem');
            $('#nav-item-logout').removeClass('invisible-elem');
        } else {
            $('#nav-item-login').removeClass('invisible-elem');
            $('#nav-item-logout').addClass('invisible-elem');
        }
    });

    console.log(' > main script function: finished');
})()