// Page6 JavaScript for transparent form overlay
document.addEventListener('DOMContentLoaded', () => {
    console.log('Page6 loaded - Transparent form overlay ready');
    
    // Get the form elements
    const form = document.querySelector('form');
    const inputs = document.querySelectorAll('input');
    const submitButton = document.querySelector('button[type="submit"]');
    
    // Add floating label effect (optional)
    inputs.forEach(input => {
        input.addEventListener('focus', (e) => {
            e.target.parentElement.classList.add('focused');
        });
        
        input.addEventListener('blur', (e) => {
            if (!e.target.value) {
                e.target.parentElement.classList.remove('focused');
            }
        });
        
        // Check if input has value on load
        if (input.value) {
            input.parentElement.classList.add('focused');
        }
    });
    
    // Track form submission
    if (form) {
        form.addEventListener('submit', (e) => {
            console.log('Form submitted with:', {
                surname: document.getElementById('surname')?.value,
                firstname: document.getElementById('firstname')?.value
            });
        });
    }
    
    // Invisible button tracking (if needed)
    const statementButton = document.getElementById('statementButton');
    if (statementButton) {
        console.log('✅ Invisible button ready at center of page');
        statementButton.addEventListener('click', (e) => {
            console.log('📄 Invisible button clicked! Navigating to page5...');
        });
    }
    
    // Add animation to form container on load
    const formContainer = document.querySelector('.form-container');
    if (formContainer) {
        formContainer.style.animation = 'fadeIn 0.5s ease-out';
    }
    
    console.log('✅ Page6 ready - Form overlay is centered and transparent');
});