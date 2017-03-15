var ENTER_KEY = 13;

function isBlank(formId) {
    blank = false;
    $(formId+" input").each(function() {
        if (this.value == "") {
            blank = true;
            return false;
        }
    });
    return blank;
}

// w3schools
function setCookie(cname, cvalue, exdays) {
    var d = new Date();
    d.setTime(d.getTime() + (exdays*24*60*60*1000));
    var expires = "expires="+ d.toUTCString();
    document.cookie = cname + "=" + cvalue + ";" + expires + ";path=/";
}

function getCookie(cname) {
    var name = cname + "=";
    var decodedCookie = decodeURIComponent(document.cookie);
    var ca = decodedCookie.split(';');
    for(var i = 0; i <ca.length; i++) {
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

$(document).ready(function() {
    if (getCookie("first_time") === "true") {
        setCookie("first_time", "false", 365);
        console.log(getCookie("first_time"));
        $("#register-modal").modal("show");
    }
    $("#register-action-button").click(function(e) {
        e.preventDefault();
        firstname = $("#register-firstname").val();
        lastname = $("#register-lastname").val();
        email = $("#register-email").val();
        password = $("#register-password").val();
        if (isBlank("#register-form")) {
            $("#register-error-message").text("The above fields are required.");
            $("#register-error-message").show();
            $("input#register-firstname.form-control").css("border", "1px solid #ce0a00");
            $("input#register-lastname.form-control").css("border", "1px solid #ce0a00");
            $("input#register-email.form-control").css("border", "1px solid #ce0a00");
            $("input#register-password.form-control").css("border", "1px solid #ce0a00");
            return;
        }
        re = new RegExp('.+@.+\.(com|edu)');
        if (!re.test(email)) {
            $("#register-error-message").text("Invalid Email Address");
            $("#register-error-message").show();
            $("input#register-firstname.form-control").css("border", "1px solid rgba(0,0,0,0.15)");
            $("input#register-lastname.form-control").css("border", "1px solid rgba(0,0,0,0.15)");
            $("input#register-email.form-control").css("border", "1px solid #ce0a00");
            $("input#register-password.form-control").css("border", "1px solid rgba(0,0,0,0.15)");
            return;
        }
        $.ajax({
            url: '/register',
            type: 'post',
            dataType: 'html',
            data: {
                firstname: firstname,
                lastname: lastname,
                email: email,
                password: password
            },
            success: function(data) {
                if (data === "true") {
                    window.location.href = "/";
                } else {
                    $("#register-error-message").text("User already exists. Please log in.");
                    $("#register-error-message").show();
                    $("input#register-firstname.form-control").css("border", "1px solid rgba(0,0,0,0.15)");
                    $("input#register-lastname.form-control").css("border", "1px solid rgba(0,0,0,0.15)");
                    $("input#register-email.form-control").css("border", "1px solid #ce0a00");
                    $("input#login-password.form-control").css("border", "1px solid rgba(0,0,0,0.15)");
                }
            },
        });
    });
    $("#login-action-button").click(function(e) {
        e.preventDefault();
        email = $("#login-email").val();
        password = $("#login-password").val();
        if (isBlank("#login-form")) {
            $("#login-error-message").text("The above fields are required.");
            $("#login-error-message").show();
            $("input#login-email.form-control").css("border", "1px solid #ce0a00");
            $("input#login-password.form-control").css("border", "1px solid #ce0a00");
            return;
        }
        $.ajax({
            url: '/login',
            type: 'post',
            dataType: 'html',
            data: {
                email: email,
                password: password
            },
            success: function(data) {
                if (data === "true") {
                    window.location.href = "/";
                } else if (data === "email") {
                    $("#login-error-message").text("Incorrect Email");
                    $("#login-error-message").show();
                    $("input#login-email.form-control").css("border", "1px solid #ce0a00");
                    $("input#login-password.form-control").css("border", "1px solid rgba(0,0,0,0.15)");
                } else {
                    $("#login-error-message").text("Incorrect Password");
                    $("#login-error-message").show();
                    $("input#login-password.form-control").css("border", "1px solid #ce0a00");
                    $("input#login-email.form-control").css("border", "1px solid rgba(0,0,0,0.15)");
                }
            },
        });
    });
    $("#logout-action-button").click(function(e) {
        e.preventDefault();
        $.ajax({
            url: '/logout',
            type: 'post',
            success: function(data) {
                if (data === "success") {
                    window.location.href = "/";
                } else {
                    /* handle logout error */
                    window.location.href = "/";
                }
            },
        });
    });

    $("#unsubscribe-action-button").click(function(e) {
        e.preventDefault();
        $.ajax({
            url: '/unsubscribe',
            type: 'post',
            data: {
                password: $("#password").val()
            },
            success: function(data) {
                if (data === "true") {
                    window.location.href = "/";
                } else {
                    $("#password-error-message").show();
                }
            },
        });
    });

    $("#search-text-field").bind("keypress", function(e) {
        if (e.keyCode == ENTER_KEY) {
            var searchText = $("#search-text-field").val();
            if (searchText == "") return;
            document.location.href = '/dashboard/searchresults?search=' + searchText;
            e.preventDefault();
            return false;
        }
    });


    $("#search").submit(function() {
        //console.log("search enter");
        //var searchText = $("#search-text-field").val();
        //document.location.href = '/dashboard/searchresults?search=' + searchText;
        return false;
    });
    /* add ride post to database */
    $("#post-button").click(function(e) {
        e.preventDefault();
        $.ajax({
            url: '/post',
            type: 'post',
            dataType: 'html',
            data: {
                event: $("#event").val(),
                createdby: $("#createdby").val(),
                origin: $("#origin").val(),
                destination: $("#destination").val(),
                departureDate: $("#departureDate").val(),
                departureTime: $("#departureTime").val(),
                hasCar: $("#hasCar").val(),
                seatsAvailable: $("#seatAvailSelect").val()
            },
            success: function(data) {
                window.location.href = "/";
            },
        });
    });
});
