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
			const response = await fetch('api/request-password-reset', {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({ email })
			});

			const data = await response.json();

			if (!response.ok) {
				throw new Error(data.message || 'Eroare la procesarea cererii');
			}

			successMessage = 'Un email pentru resetarea parolei a fost trimis la adresa indicată.';
			console.log('Email trimis cu succes:', data);
			// If token exists, redirect immediately to /newpassword with token
			if (data.token) {
				console.log(`Received token: ${data.token}`);
				
				setTimeout(() => {
					goto(`/newPassword?token=${encodeURIComponent(data.token)}`);
				}, 1500);
			} else {
				console.log(`A INTRAT AICI!!: ${data.token}`);
				// fallback if token not returned
				setTimeout(() => {
					goto('/newPassword');
				}, 1500);
			}
		} catch (error: any) {
			errorMessage = error.message || 'A apărut o eroare la procesarea cererii';
			console.error('Eroare resetare parolă:', error);
		} finally {
			isLoading = false;
		}
	}

	// Function for navigating back to login
	function goToLogin() {
		goto('login');
	}
</script>

<div class="container">
	<div class="logo-container" onclick={goToLogin}>
		<img src="/ourchat_logo.png" alt="OurChat Logo" class="logo-image" />
		<div class="logo-text">OurChat</div>
	</div>

	<h1>Reset Your Password</h1>

	<form onsubmit={preventDefault(handleSubmit)}>
		<input type="email" placeholder="email address" bind:value={email} required />

		<button type="submit" disabled={isLoading}>
			{isLoading ? 'Loading...' : 'CONFIRM'}
		</button>

		{#if errorMessage}
			<div class="error">{errorMessage}</div>
		{/if}

		{#if successMessage}
			<div class="success">{successMessage}</div>
		{/if}
	</form>

	<span class="login-link" onclick={goToLogin}>Back to login</span>
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
		text-align: center;
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

	button:disabled {
		opacity: 0.7;
		cursor: not-allowed;
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

	.success {
		color: #44aa44;
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