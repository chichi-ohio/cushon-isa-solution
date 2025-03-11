document.addEventListener('DOMContentLoaded', function() {
    const form = document.getElementById('investmentForm');
    const successMessage = document.getElementById('successMessage');
    const errorMessage = document.getElementById('errorMessage');

    form.addEventListener('submit', async function(e) {
        e.preventDefault();

        // Hide any existing messages
        successMessage.classList.add('hidden');
        errorMessage.classList.add('hidden');

        // Get form data
        const formData = {
            name: form.name.value,
            email: form.email.value,
            fund: form.fund.value,
            amount: parseFloat(form.amount.value)
        };

        try {
            const response = await fetch('/api/invest', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify(formData)
            });

            const data = await response.json();

            if (response.ok) {
                // Show success message
                successMessage.classList.remove('hidden');
                successMessage.classList.add('fade-in');
                form.reset();
            } else {
                // Show error message
                errorMessage.textContent = data.error || 'An error occurred while processing your investment.';
                errorMessage.classList.remove('hidden');
                errorMessage.classList.add('fade-in');
            }
        } catch (error) {
            console.error('Error:', error);
            errorMessage.textContent = 'An error occurred while submitting your investment. Please try again.';
            errorMessage.classList.remove('hidden');
            errorMessage.classList.add('fade-in');
        }
    });

    // Add input validation
    const amountInput = document.getElementById('amount');
    amountInput.addEventListener('input', function() {
        const value = parseFloat(this.value);
        if (value < 1) {
            this.setCustomValidity('Minimum investment amount is £1');
        } else {
            this.setCustomValidity('');
        }
    });
}); 