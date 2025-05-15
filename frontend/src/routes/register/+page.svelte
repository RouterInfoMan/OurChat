<script lang="ts">
	import { goto } from '$app/navigation';
	import { preventDefault } from 'svelte/legacy';

	// Variables to store form data
	let username = $state('');
	let email = $state('');
	let password = $state('');
	let confirmPassword = $state('');
	let errorMessage = $state('');
	let isLoading = $state(false);

	// Function to handle form submission
	async function handleSubmit() {
		try {
			isLoading = true;
			errorMessage = '';

			// Basic form validation
			if (!username || !email || !password || !confirmPassword) {
				errorMessage = 'Toate câmpurile sunt obligatorii';
				isLoading = false;
				return;
			}

			// Email validation
			const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
			if (!emailRegex.test(email)) {
				errorMessage = 'Adresa de email nu este validă';
				isLoading = false;
				return;
			}

			// Password match validation
			if (password !== confirmPassword) {
				errorMessage = 'Parolele nu coincid';
				isLoading = false;
				return;
			}

			// Make the API call to register the user
			const response = await fetch('api/register', {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({ username, email, password })
			});

			const data = await response.json();

			if (!response.ok) {
				throw new Error(data.message || 'Eroare la înregistrare');
			}

			// If successful, redirect to login page
			goto('login');
		} catch (error: any) {
			errorMessage = error.message || 'A apărut o eroare la înregistrare';
			console.error('Eroare înregistrare:', error);
		} finally {
			isLoading = false;
		}
	}

	// Navigation to the login page
	function goToLogin() {
		goto('login');
	}

	// Navigation to the dashboard (not used in the registration page)
	function goToDashboard() {
		goto('dashboard');
	}
</script>

<div class="container">
	<!-- Logo Section -->
	<div class="logo-container" onclick={goToDashboard}>
		<img src="/ourchat_logo.png" alt="OurChat Logo" class="logo-image" />
		<div class="logo-text">OurChat</div>
	</div>

	<!-- Registration Form -->
	<h1>Create an account</h1>

	<form onsubmit={preventDefault(handleSubmit)}>
		<input type="text" placeholder="username" bind:value={username} required />

		<input type="email" placeholder="e-mail" bind:value={email} required />

		<input type="password" placeholder="password" bind:value={password} required />

		<input type="password" placeholder="confirm password" bind:value={confirmPassword} required />

		<button type="submit" disabled={isLoading}>
			{isLoading ? 'Loading...' : 'CONFIRM'}
		</button>

		{#if errorMessage}
			<div class="error">{errorMessage}</div>
		{/if}
	</form>

	<span class="login-link" onclick={goToLogin}>Already have an account? Log in</span>
</div>

<style>
	:global(body) {
		margin: 0;
		padding: 0;
		font-family: Arial, sans-serif;
		height: 100vh;
		background: linear-gradient(135deg, #6a5af9 0%, #4a91ff 100%);
		display: flex;
		justify-content: center;
		align-items: center;
	}

	.container {
		display: flex;
		flex-direction: column;
		align-items: center;
		width: 100%;
		max-width: 450px;
		padding: 20px;
	}

	.logo-container {
		display: flex;
		align-items: center;
		margin-bottom: 40px;
		cursor: pointer;
	}

	.logo-image {
		width: 120px;
		height: 120px;
		object-fit: contain;
	}

	.logo-text {
		font-size: 42px;
		font-weight: bold;
		color: #222;
		margin-left: 20px;
	}

	h1 {
		color: white;
		font-size: 36px;
		font-weight: normal;
		margin-bottom: 30px;
		text-align: center;
	}

	form {
		width: 100%;
		display: flex;
		flex-direction: column;
		gap: 20px;
	}

	input {
		width: 100%;
		padding: 16px;
		border-radius: 8px;
		border: none;
		background-color: rgba(255, 255, 255, 0.2);
		color: white;
		font-size: 18px;
		box-sizing: border-box;
		text-align: center;
		outline: none;
		transition: all 0.3s ease;
	}

	input::placeholder {
		color: rgba(255, 255, 255, 0.8);
	}

	input:focus {
		background-color: rgba(255, 255, 255, 0.3);
	}

	button {
		width: 50%;
		margin: 20px auto 0;
		padding: 14px;
		border-radius: 30px;
		border: none;
		background-color: #5977ff;
		color: white;
		font-size: 20px;
		font-weight: bold;
		cursor: pointer;
		transition: all 0.3s ease;
		box-shadow: 0 4px 8px rgba(0, 0, 0, 0.2);
	}

	button:hover {
		background-color: #4a68eb;
		transform: translateY(-2px);
		box-shadow: 0 6px 12px rgba(0, 0, 0, 0.2);
	}

	button:active {
		transform: translateY(0);
		box-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
	}

	.error {
		color: #ff4444;
		text-align: center;
		margin-top: 10px;
		font-size: 14px;
		background-color: rgba(255, 255, 255, 0.7);
		padding: 8px;
		border-radius: 4px;
	}

	.login-link {
		margin-top: 20px;
		color: white;
		text-decoration: none;
		cursor: pointer;
	}

	.login-link:hover {
		text-decoration: underline;
	}
</style>
