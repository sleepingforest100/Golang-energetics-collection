$(document).ready(function () {
    checkUserRole();
});

function checkUserRole() {
    const navLogin = document.getElementById('nav-login');
    const navProfile = document.getElementById('nav-profile');
    const navAdmin = document.getElementById('nav-admin');
    const navLogout = document.getElementById('nav-logout');

    navLogout.addEventListener('click', function(event) {
        event.preventDefault();
        document.cookie = 'jwtToken=;expires=Thu, 01 Jan 1970 00:00:00 UTC;path=/;';
        window.location.reload();
    });

    if (auth('admin')) {
        navProfile.classList.remove('hidden');
        navAdmin.classList.remove('hidden');
        navLogout.classList.remove('hidden');
        navLogin.classList.add('hidden');
    } else if (auth('user')) {
        navProfile.classList.remove('hidden');
        navLogout.classList.remove('hidden');
        navAdmin.classList.add('hidden');
        navLogin.classList.add('hidden');
    } else {
        navLogin.classList.remove('hidden');
        navProfile.classList.add('hidden');
        navAdmin.classList.add('hidden');
        navLogout.classList.add('hidden');
    }
}
