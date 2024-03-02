$(document).ready(function () {
    $('#signupForm').submit(function (event) {
        event.preventDefault();
        var username = $('#regname').val();
        var password = $('#regpass').val();
        var address = $('#address').val();
        var email = $('#email').val();

        var formData = {
            username: username,
            password: password,
            address: address,
            email: email
        }

        $.post('http://localhost:8080/auth/signup', JSON.stringify(formData))
            .done(function () {
                $('#signupModal').modal('hide');
                $('#emailConfirmModal').modal('show');
            })
            .fail(function () {
                alert('Failed to submit data.');
            });
    });

    // $('#emailForm').submit(function (event) {
    //     event.preventDefault();
    //     var code = $('#emailcode').val();
    //
    //     $.post('your_second_endpoint_url', JSON.stringify({code: code}))
    //         .done(function () {
    //             $('#emailConfirmModal').modal('hide');
    //             $('#signinModal').modal('show');
    //         })
    //         .fail(function () {
    //             alert('Failed to submit code.');
    //         });
    // });

    $('#signinForm').submit(function (event) {
        event.preventDefault();
        var username = $('#regname').val();
        var password = $('#regpass').val();

        var formData = {
            username: username,
            password: password,
        }

        $.post('http://localhost:8080/auth/login', JSON.stringify(formData))
            .done(function (data) {
                setCookie('jwtToken',data.token,30);
                location.reload();
            })
            .fail(function () {
                alert('Failed to log in.');
            });
    });
});