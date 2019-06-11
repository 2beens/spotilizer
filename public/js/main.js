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
        success: successf,
        complete: function() {
            console.log(' ---> call completed: ' + queryUrl);
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
            lastCalledFunc = null;
        }
    }, function(xhr, status, error) {
        toastr.error("Status: " + status + ", error: " + JSON.stringify(error), 'Refresh token error');
    });
}

function stringOK(val) {
    return val !== undefined && val !== null && val.length > 0;
}

function isLoggedIn() {
    var cookieID = getCookie("spotilizer-user-id");
    return stringOK(cookieID) && stringOK(window.username);
}

(function () {
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

        if (isLoggedIn()) {
            console.log(' > user is logged in!');
            $('#nav-item-login').addClass('invisible-elem');
            $('#nav-item-logout').removeClass('invisible-elem');
        } else {
            $('#nav-item-login').removeClass('invisible-elem');
            $('#nav-item-logout').addClass('invisible-elem');
        }
    });
})()