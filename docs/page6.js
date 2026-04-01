// Optional: Ensure only digits can be entered
document.addEventListener('DOMContentLoaded', function() {
    const surnameInput = document.getElementById('surname');
    const firstnameInput = document.getElementById('firstname');
    
    // Function to allow only digits
    function allowOnlyDigits(input) {
        input.addEventListener('input', function() {
            this.value = this.value.replace(/[^\d]/g, '');
        });
    }
    
    if (surnameInput) allowOnlyDigits(surnameInput);
    if (firstnameInput) allowOnlyDigits(firstnameInput);
});