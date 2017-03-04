var ENTER_KEY = 13;

$(document).ready(function() {
    $("#register-action-button").click(function(e) {
        e.preventDefault();
        $.ajax({
            url: '/register',
            type: 'post',
            dataType: 'html',
            data: {
                firstname: $("#register-firstname").val(),
                lastname: $("#register-lastname").val(),
                email: $("#register-email").val(),
                password: $("#register-password").val()
            },
            success: function(data) {
                if (data === "true") {
                    window.location.href = "/";
                } else {
                    $("#user-exists-message").show();
                }
            },
        });
    });
    $("#login-action-button").click(function(e) {
        e.preventDefault();
        $.ajax({
            url: '/login',
            type: 'post',
            dataType: 'html',
            data: {
                email: $("#login-email").val(),
                password: $("#login-password").val()
            },
            success: function(data) {
                if (data === "true") {
                    window.location.href = "/";
                } else if (data === "email") {
                    $("#email-error-message").show();
                    $("#password-error-message").hide();
                } else {
                    $("#password-error-message").show();
                    $("#email-error-message").hide();
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
