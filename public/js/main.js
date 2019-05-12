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
    $.ajax({
        url: 'https://api.spotify.com/v1/me',
        headers: {
            'Authorization': 'Bearer ' + accessToken
        },
        success: function (response) {
            console.log('------------------------------- response from spotify ------------')
            console.log(response);
            console.log('------------------------------------------------------------------')
            // $('#login').hide();
            // $('#loggedin').show();
        }
    });
}

(function () {
    var accessToken = getCookie("accessToken");
    var refreshToken = getCookie("refreshToken");
    console.log(" > loaded cookie AT: " + accessToken);
    console.log(" > loaded cookie RT: " + refreshToken);
})()