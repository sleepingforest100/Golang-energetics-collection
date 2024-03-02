document.addEventListener('DOMContentLoaded', function () {
    if (auth('admin')) {
        fetchUserData();
    } else {
        window.location.href = '../index-go.html';
    }
});

function fetchUserData() {
    fetch('your-api-endpoint')
        .then(response => response.json())
        .then(data => {
            const tableBody = document.getElementById('tableBody');
            data.forEach(item => {
                const row = document.createElement('tr');
                row.innerHTML = `
            <td>${item.id}</td>
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
                    const userId = item.id;
                    fetch('/your-server-endpoint', {
                        method: 'POST',
                        headers: {
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