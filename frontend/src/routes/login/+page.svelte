<script lang="ts">
	import { goto } from '$app/navigation';
	import { preventDefault } from 'svelte/legacy';

	// Variables for form data
	let username = $state('');
	let password = $state('');
	let errorMessage = $state('');
	let isLoading = $state(false);

	// Function to handle form submission
	async function handleSubmit() {
		try {
			isLoading = true;
			errorMessage = '';

			// Basic validation
			if (!username || !password) {
				errorMessage = 'Toate câmpurile sunt obligatorii';
				isLoading = false;
				return;
			}

			// API call to authenticate the user with username
			const response = await fetch('api/login', {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({ username, password })
			});

			const data = await response.json();

			
			if (!response.ok) {
				if (response.headers.get('Content-Type')?.includes('application/json')) {
					throw new Error(data.message || 'Eroare la autentificare');
				} else {
					throw new Error('Eroare la autentificare. Verificați mesajul de eroare de la server');
				}
			}

			// On successful login, store the JWT token and redirect to the dashboard
			if (data.token) {
				localStorage.setItem('jwt_token', data.token);
				goto('dashboard');
			} else {
				throw new Error('Nu s-a primit token de autentificare');
			}
		} catch (error: any) {
			errorMessage = error.message || 'A apărut o eroare la autentificare';
			console.error('Eroare autentificare:', error);
		} finally {
			isLoading = false;
		}
	}

	// Function for navigating to the registration page
	function goToRegister() {
		goto('register');
	}

	// Function for navigating to the password reset page
	function goToForgotPassword() {
		goto('forgotPassword');
	}

	// Function for navigating to the dashboard (if not logged in yet, redirect to login page)
	function goToDashboard() {
		goto('dashboard');
	}
</script>

<div class="container">
	<div class="logo-container" onclick={goToDashboard}>
		<img src="/ourchat_logo.png" alt="OurChat Logo" class="logo-image" />
		<div class="logo-text">OurChat</div>
	</div>

	<h1>Log in to your account</h1>

	<form onsubmit={preventDefault(handleSubmit)}>
		<input type="text" placeholder="username" bind:value={username} required />

		<input type="password" placeholder="password" bind:value={password} required />

		<button type="submit" disabled={isLoading}>
			{isLoading ? 'Loading...' : 'CONFIRM'}
		</button>

		{#if errorMessage}
			<div class="error">{errorMessage}</div>
		{/if}
	</form>

	<div class="forgot-password" onclick={goToForgotPassword}>Forgot Password</div>

	<span class="register-link" onclick={goToRegister}>Don't have an account? Sign up</span>
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

	button,
	.forgot-password {
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
		text-align: center;
	}

	button:hover,
	.forgot-password:hover {
		background-color: #4a68eb;
		transform: translateY(-2px);
		box-shadow: 0 6px 12px rgba(0, 0, 0, 0.2);
	}

	button:active,
	.forgot-password:active {
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

	.register-link {
		margin-top: 20px;
		color: white;
		text-decoration: none;
		cursor: pointer;
	}

	.register-link:hover {
		text-decoration: underline;
	}

	.forgot-password {
		margin-top: 20px;
		background-color: rgba(255, 255, 255, 0.2);
	}

	.forgot-password:hover {
		background-color: rgba(255, 255, 255, 0.3);
	}
</style>
