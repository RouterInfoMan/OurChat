<script lang="ts">
	import { goto } from '$app/navigation';
	import { preventDefault } from 'svelte/legacy';

	let email = '';
	let errorMessage = '';
	let successMessage = '';
	let isLoading = false;

	async function handleSubmit() {
		try {
			isLoading = true;
			errorMessage = '';
			successMessage = '';

			// Basic validation
			if (!email) {
				errorMessage = 'Adresa de email este obligatorie';
				isLoading = false;
				return;
			}

			const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
			if (!emailRegex.test(email)) {
				errorMessage = 'Adresa de email nu este validă';
				isLoading = false;
				return;
			}

			// Send request to API
			const response = await fetch('/forgot-password', {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({ email })
			});

			const data = await response.json();

			if (!response.ok) {
				throw new Error(data.message || 'Eroare la procesarea cererii');
			}

			successMessage = 'Un email pentru resetarea parolei a fost trimis la adresa indicată.';

			// Redirect after 1.5 seconds to /newpassword
			setTimeout(() => {
				goto('/newpassword');
			}, 1500);
		} catch (error: any) {
			errorMessage = error.message || 'A apărut o eroare la procesarea cererii';
			console.error('Eroare resetare parolă:', error);
		} finally {
			isLoading = false;
		}
	}
</script>

<div class="container">
	<h1>Resetare Parolă</h1>

	<form on:submit|preventDefault={handleSubmit}>
		<input type="email" placeholder="Emailul tău" bind:value={email} required />

		<button type="submit" disabled={isLoading}>
			{isLoading ? 'Se trimite...' : 'Resetează parola'}
		</button>

		{#if errorMessage}
			<div class="error">{errorMessage}</div>
		{/if}

		{#if successMessage}
			<div class="success">{successMessage}</div>
		{/if}
	</form>
</div>

<style>
	.container {
		max-width: 400px;
		margin: 80px auto;
		padding: 20px;
		background: white;
		border-radius: 10px;
		box-shadow: 0 5px 20px rgba(0, 0, 0, 0.2);
		display: flex;
		flex-direction: column;
		gap: 20px;
	}

	h1 {
		text-align: center;
		color: #333;
	}

	form {
		display: flex;
		flex-direction: column;
		gap: 16px;
	}

	input {
		padding: 14px;
		border-radius: 8px;
		border: 1px solid #ccc;
		font-size: 16px;
	}

	button {
		padding: 12px;
		border-radius: 8px;
		background-color: #5977ff;
		color: white;
		font-size: 16px;
		border: none;
		cursor: pointer;
	}

	.error {
		color: #ff4444;
		text-align: center;
	}

	.success {
		color: #44aa44;
		text-align: center;
	}
</style>

