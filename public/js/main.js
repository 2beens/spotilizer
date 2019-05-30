toastr.options = {
    "positionClass": "toast-top-left",
}

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

function makeRequest(queryUrl, successf, errorf) {
    console.log(' ---> making a request call to: ' + queryUrl);
    if (errorf === undefined || errorf === null) {
        errorf = function(xhr,status,error) {
            console.log(' ---> error occured, status: ' + status + ', error: ' + error);
        };
    }
    $.ajax({
        url: queryUrl,
        headers: {
            'Authorization': 'Bearer ' + window.accessToken
        },
        success: successf,
        complete: function() {
            console.log(' ---> call completed!');
        },
        error: errorf
    });
}

function checkIfRefreshTokenNeeded(response) {
    if (!response.error) {
        return false;
    }
    if (!response.error.status || !response.error.message) {
        return false;
    }
    if (response.error.status !== 401) {
        return false;
    }
    if(!response.error.message.includes('access token expired')) {
        return false;
    }
    return true;
}

var lastCalledFunc = null;

function refreshTokenFunc() {
    console.log(' > making a refresh token call...');
    makeRequest("/refresh_token", function(response) {
        console.log(' > refresh token, received from server: ' + response);
        var respObj = JSON.parse(response);
        if (respObj.error) {
            toastr.error(respObj.error.message, 'Refresh token error');
            return;
        }
        // call last failed function if needed
        // TODO: watch for inifinite looping here - not calling last called function over and over again
        //       when, e.g. internet connection is gone
        if (lastCalledFunc !== null) {
            lastCalledFunc();
        }
    }, function(xhr, status, error) {
        toastr.error("Status: " + status + ", error: " + JSON.stringify(error), 'Refresh token error');
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

function saveCurrentTracks() {
    lastCalledFunc = saveCurrentTracks;
    makeRequest("/save_current_tracks", function(response) {
        console.log(' > received from server: ' + response);
        var respObj = JSON.parse(response);
        if (respObj.error) {
            toastr.error(respObj.error.message, 'Save fav tracks error');
        } else {
            toastr.success(respObj.message, 'Save fav tracks');
        }
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
            toastr.info('Page lodaded ...', 'Spotilizer', {timeOut: 1000})
        }

        if (stringOK(cookieID) && stringOK(window.username)) {
            console.log(' > cookieID: ' + cookieID + ', username: ' + username);
            $('#nav-item-login').addClass('invisible-elem');
            $('#nav-item-logout').removeClass('invisible-elem');
            $('#spotify-controls-div').removeClass('invisible-elem');
        } else {
            $('#nav-item-login').removeClass('invisible-elem');
            $('#nav-item-logout').addClass('invisible-elem');
            $('#spotify-controls-div').addClass('invisible-elem');
        }
    });

    console.log(' > main script function: finished');
})()