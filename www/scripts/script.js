function login() {
    clearErrors()
    var email = document.getElementById('loginUsername').value;
    var password = document.getElementById('loginPassword').value;
    // Добавьте код для отправки данных на сервер для проверки логина и пароля
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
            if (data.code !== 0) {
                showError('response',data);
                return false;
            }
            document.getElementById('response').innerText = data.message;
        })
        .catch(error => console.error('Error:', error));
    console.log('Login:', username, 'Password:', password);
}

function register() {
    clearErrors()
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
    clearErrors()
    document.getElementById('forgotPassword').style.display = 'none';
    document.getElementById('returnToLogin').style.display = 'flex';
    document.getElementById('registerForm').style.display = 'flex';
    document.getElementById('loginForm').style.display = 'none';
    document.getElementById('regButton').style.display = 'none';

    document.getElementById('restorePassword').style.display = 'none';
}

function returnToLogin() {
    clearErrors()
    document.getElementById('forgotPassword').style.display = 'block';
    document.getElementById('returnToLogin').style.display = 'none';
    document.getElementById('registerForm').style.display = 'none';
    document.getElementById('loginForm').style.display = 'flex';
    document.getElementById('regButton').style.display = 'block';

    document.getElementById('restorePassword').style.display = 'none';
}

function forgotPassword() {
    clearErrors();
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

// функция сообщения об ошибке
function showError(field, data) {
    clearErrors();
    const errData = JSON.stringify(data.details[0]);
    const parsedData = JSON.parse(errData);
    var errorSpan = document.createElement("span");
    var errorMessage = document.createTextNode(parsedData.field + " " + parsedData.description);

    errorSpan.appendChild(errorMessage);
    errorSpan.className = "errorMsg";
    //
    // var fieldLabel = document.getElementById('container');
    // while (fieldLabel.nodeName.toLowerCase() !== "label") {
    //     fieldLabel = fieldLabel.previousSibling;
    // }
    // fieldLabel.appendChild(errorSpan);
    document.getElementById(field).appendChild(errorSpan)
}

function clearErrors() {
    var errorContainer = document.getElementById('response');
    errorContainer.innerHTML = ''; // или errorContainer.textContent = '';
}