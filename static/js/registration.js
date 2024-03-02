$(document).ready(function () {
    $('#signupForm').submit(function (event) {
        event.preventDefault();
    
        var username = $('#regname').val();
        var password = $('#regpass').val();
        var address = $('#address').val();
        var email = $('#email').val();
    
        var formData = {
            name: username,
            password: password,
            address: address,
            email: email
        };
    
        fetch('http://localhost:8080/auth/signup', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(formData),
        })
        .then(response => {
            if (!response.ok) {
                throw new Error('Failed to submit data.');
            }
            return response.json();
        })
        .then(() => {
            $('#signupModal').modal('hide');
            $('#emailConfirmModal').modal('show');
        })
        .catch(error => {
            console.error('Error:', error.message);
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
        
        var email = $('#logmail').val();
        var password = $('#logpass').val();
    
        var formData = {
            email: email,
            password: password,
        };
    
        fetch('http://localhost:8080/auth/login', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(formData),
        })
        .then(response => {
            if (!response.ok) {
                throw new Error('Failed to log in.');
            }
            return response.json();
        })
        .then(data => {
            console.log('Server response:', data);
    
            if (data && data.token) {
                setCookie('jwtToken', data.token, 30);
                location.reload();
            } else {
                alert('Token not found in server response.');
            }
        })
        .catch(error => {
            console.error('Error:', error.message);
            alert('Failed to log in.');
        });
    });
});