function getCookie(name) {
    const value = `; ${document.cookie}`;
    const parts = value.split(`; ${name}=`);
    if (parts.length === 2) return parts.pop().split(';').shift();
}

function decodeToken(token) {
    const tokenPayload = token.split('.')[1];
    const decodedPayload = JSON.parse(atob(tokenPayload));
    return decodedPayload;
}

function setCookie(name, value, days) {
    const expires = new Date();
    expires.setDate(expires.getDate() + days);
    document.cookie = `${name}=${value};expires=${expires.toUTCString()};path=/`;
}

function auth(role){
    const token = getCookie('jwtToken');
    if (token) {
        const decodedToken = decodeToken(token);
        return  decodedToken.role === role;
    }else {
        console.log('Token not found.');
        return false;
    }
}