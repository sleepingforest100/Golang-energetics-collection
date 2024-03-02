document.addEventListener('DOMContentLoaded', function () {
    if (auth('admin')) {
        fetchUserData();
        sendMailInit();
    } else {
        window.location.href = '../index-go.html';
    }
});

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

    fetch('http://localhost:8080/users', requestOptions)
        .then(response => response.json())
        .then(data => {
            const tableBody = document.getElementById('tableBody');
            data.forEach(item => {
                const row = document.createElement('tr');
                row.innerHTML = `
            <td>${item.ID}</td>
            <td>${item.name}</td>
            <td>${item.email}</td>
            <td>${item.address}</td>
            <td>
                <select class="roleSelect">
                    <option value="user" ${item.role === "user" ? "selected" : ""}>USER</option>
                    <option value="admin" ${item.role === "admin" ? "selected" : ""}>ADMIN</option>
                </select>
            </td>
                `;
                const selectBox = row.querySelector('.roleSelect');
                selectBox.addEventListener('change', function() {
                    const selectedRole = this.value;
                    const userId = item.ID;
                    fetch('/your-server-endpoint', {
                        method: 'POST',
                        headers: {
                            'Authorization': `Bearer ${token}`,
                            'Content-Type': 'application/json'
                        },
                        body: JSON.stringify({
                            userId: userId,
                            role: selectedRole
                        })
                    })
                        .then(response => {
                            if (response.ok) {
                                alert('Role updated successfully');
                                console.log('Role updated successfully');
                            } else {
                                console.error('Failed to update role');
                            }
                        })
                        .catch(error => {
                            console.error('Error:', error);
                        });
                });
                tableBody.appendChild(row);
            });
        })
        .catch(error => console.error('Error fetching data:', error));
}

function sendMailInit () {
    const token = getCookie('jwtToken');
    if (!token) {
        console.error('Token not found.');
        return;
    }

    const form = document.getElementById('mailForm');

    form.addEventListener('submit', function(event) {
        event.preventDefault();

        const mailHeader = document.getElementById('mailHeader').value;
        const mailBody = document.getElementById('mailBody').value;

        const mailData = {
            subject: mailHeader,
            body: mailBody
        };

        fetch('#', {
            method: 'POST',
            headers: {
                'Authorization': `Bearer ${token}`,
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(mailData)
        })
            .then(response => {
                if (!response.ok) {
                    throw new Error('Failed to send mail');
                }
                return response.json();
            })
            .then(data => {
                console.log('Mail sent successfully:', data);
            })
            .catch(error => {
                console.error('Error sending mail:', error.message);
            });
    });
}