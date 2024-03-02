document.addEventListener('DOMContentLoaded', function () {

    fetchUserData();

    function fetchUserData() {
        fetch('your-api-endpoint')
            .then(response => response.json())
            .then(data => {
                document.getElementById('profile-info').innerHTML = `
            <div>Username: ${data.username}</div>
            <div>Email: ${data.email}</div>
            <div>Address: ${data.address}</div>
          `;
                document.getElementById('edname').value = data.username;
                document.getElementById('edemail').value = data.email;
                document.getElementById('edaddress').value = data.address;
            })
            .catch(error => console.error('Error fetching user data:', error));
    }

    document.getElementById('profileEditForm').addEventListener('submit', function (event) {
        event.preventDefault();

        const formData = {
            username: document.getElementById('edname').value,
            email: document.getElementById('edemail').value,
            address: document.getElementById('edaddress').value,
        };

        fetch('your-api-endpoint', {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(formData)
        })
            .then(response => {
                if (response.ok) {
                    console.log('User data updated successfully');
                    fetchUserData();
                } else {
                    throw new Error('Failed to update user data');
                }
            })
            .catch(error => console.error('Error updating user data:', error));
    });
    document.getElementById('changePassForm').addEventListener('submit', function (event) {
        event.preventDefault();

        const formData = {
            oldpassword: document.getElementById('oldpass').value,
            newpassword: document.getElementById('newpass').value
        };

        fetch('your-api-endpoint', {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(formData)
        })
            .then(response => {
                if (response.ok) {
                    console.log('User data updated successfully');
                    alert('You\'ve successfully changed your password!')
                } else {
                    throw new Error('Failed to update user data');
                }
            })
            .catch(error => console.error('Error updating user data:', error));
    });
});