document.addEventListener('DOMContentLoaded', function () {

    if (auth('admin')||auth('user')) {

        fetchUserData();

        function fetchUserData() {
            const token = getCookie('jwtToken');
            if (!token) {
                console.error('Token not found.');
                return;
            }
            const requestOptions = {
                method: 'GET',
                headers: {
                    'Authorization': `Bearer ${token}`,
                    'Content-Type': 'application/json'
                }
            };
            fetch('http://localhost:8080/user', requestOptions)
                .then(response => {
                    if (!response.ok) {
                        throw new Error('Failed to fetch user data');
                    }
                    return response.json();
                })
                .then(data => {
                    document.getElementById('profile-info').innerHTML = `
                <div>Username: ${data.name}</div>
                <div>Email: ${data.email}</div>
                <div>Address: ${data.address}</div>
            `;
                    document.getElementById('edname').value = data.name;
                    //document.getElementById('edemail').value = data.email;
                    document.getElementById('edaddress').value = data.address;
                })
                .catch(error => console.error('Error fetching user data:', error));
        }

        document.getElementById('profileEditForm').addEventListener('submit', function (event) {
            event.preventDefault();
            const token = getCookie('jwtToken');
            if (!token) {
                console.error('Token not found.');
                return;
            }


            const formData = {
                name: document.getElementById('edname').value,
                // email: document.getElementById('edemail').value,
                address: document.getElementById('edaddress').value,
            };


            fetch('http://localhost:8080/user', {
                method: 'PUT',
                headers: {
                    'Authorization': `Bearer ${token}`,
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify(formData)
            })
                .then(response => {
                    if (response.ok) {
                        console.log('User data updated successfully');
                        $('#profileEditModal').modal('hide');
                        fetchUserData();
                    } else {
                        throw new Error('Failed to update user data');
                    }
                })
                .catch(error => console.error('Error updating user data:', error));
        });
        document.getElementById('changePassForm').addEventListener('submit', function (event) {
            event.preventDefault();
            const token = getCookie('jwtToken');
            if (!token) {
                console.error('Token not found.');
                return;
            }
            const formData = {
                oldpassword: document.getElementById('oldpass').value,
                password: document.getElementById('newpass').value
            };

            fetch('http://localhost:8080/auth/reset', {
                method: 'PUT',
                headers: {
                    'Authorization': `Bearer ${token}`,
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify(formData)
            })
                .then(response => {
                    if (response.ok) {
                        console.log('User data updated successfully');
                        alert('You\'ve successfully changed your password!');
                        $('#changePassModal').modal('hide');
                        document.getElementById('oldpass').value ='';
                        document.getElementById('newpass').value ='';
                    } else {
                        throw new Error('Failed to update user data');
                    }
                })
                .catch(error => console.error('Error updating user data:', error));
        });
    } else {
        window.location.href = '../index-go.html'
    }
});