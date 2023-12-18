function login() {
    var email = document.getElementById('loginUsername').value;
    var password = document.getElementById('loginPassword').value;
    // Добавьте код для отправки данных на сервер для проверки логина и пароля
    // Используйте Fetch API для отправки запроса к серверу
    const credentials = {
        email: email,
        password: password,
        app_id: "1"
    };
    fetch('/auth/login', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(credentials)
    })
        .then(response => response.json())
        .then(data => {
            // Обработка полученных данных
            var resp = data.message

            // Object.keys(data.details).forEach(prop => {
            //    resp += " " + data.details[prop]
            // })
            for ( item in data.details ) {
                resp += " " + item;
            }
            document.getElementById('response').innerText = resp;aaaBB
        })
        .catch(error => console.error('Error:', error));
    console.log('Login:', username, 'Password:', password);
}

function register() {
    var username = document.getElementById('registerUsername').value;
    var password = document.getElementById('registerPassword').value;
    var confirmPassword = document.getElementById('confirmPassword').value;

    if (password !== confirmPassword) {
        alert('Пароли не совпадают');
        return;
    }

    // Добавьте код для отправки данных на сервер для регистрации нового пользователя
    console.log('Register:', username, 'Password:', password);
}

function resetPassword() {
    //todo
    var userEmail = document.getElementById('userEmail').value;
    console.log('email:',userEmail)

}

function openRegisterForm() {
    document.getElementById('forgotPassword').style.display = 'none';
    document.getElementById('returnToLogin').style.display = 'flex';
    document.getElementById('registerForm').style.display = 'flex';
    document.getElementById('loginForm').style.display = 'none';
    document.getElementById('regButton').style.display = 'none';

    document.getElementById('restorePassword').style.display = 'none';
}

function returnToLogin() {
    document.getElementById('forgotPassword').style.display = 'block';
    document.getElementById('returnToLogin').style.display = 'none';
    document.getElementById('registerForm').style.display = 'none';
    document.getElementById('loginForm').style.display = 'flex';
    document.getElementById('regButton').style.display = 'block';

    document.getElementById('restorePassword').style.display = 'none';
}

function forgotPassword() {
    // Добавьте код для восстановления пароля
    document.getElementById('restorePassword').style.display = 'flex';
    document.getElementById('loginForm').style.display = 'none';
    document.getElementById('regButton').style.display = 'none';

    console.log('Forgot Password');
}

document.getElementById('loginPassword').addEventListener('input', function() {
    var forgotPassword = document.getElementById('forgotPassword');
    forgotPassword.style.display = this.value ? 'none' : 'block';
});

document.getElementById('registerPassword').addEventListener('input', function() {
    var returnToLoginReg = document.getElementById('returnToLoginReg');
    returnToLoginReg.style.display = this.value ? 'block' : 'none';
});