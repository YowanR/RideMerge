/* w3schools */
function setCookie(cname, cvalue, exdays) {
    var d = new Date();
    d.setTime(d.getTime() + (exdays*24*60*60*1000));
    var expires = "expires="+ d.toUTCString();
    document.cookie = cname + "=" + cvalue + ";" + expires + ";path=/";
}

$(function() {
    setTimeout(function() {
        setCookie("first_time", "true", 365);
        window.location.href = "/";
    }, 3000);
});
