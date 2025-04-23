// Form validation
document.addEventListener('DOMContentLoaded', function() {
    // Password confirmation check for registration form
    const registerForm = document.querySelector('form[action="/register"]');
    if (registerForm) {
        registerForm.addEventListener('submit', function(event) {
            const password = document.getElementById('password').value;
            const confirmPassword = document.getElementById('confirm_password').value;

            if (password !== confirmPassword) {
                event.preventDefault();
                alert('Passwords do not match!');
            }
        });
    }

    // Auto-hide success messages after 5 seconds
    const successMessages = document.querySelectorAll('.message.success');
    if (successMessages.length > 0) {
        setTimeout(function() {
            successMessages.forEach(message => {
                message.style.opacity = '0';
                setTimeout(() => {
                    message.style.display = 'none';
                }, 500);
            });
        }, 5000);
    }
});