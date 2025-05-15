<script lang="ts">
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';
	import { preventDefault } from 'svelte/legacy';

	let newPassword = '';
	let confirmPassword = '';
	let token = '';
	let errorMessage = '';
	let successMessage = '';
	let isLoading = false;

	onMount(() => {
		const url = new URL(window.location.href);
		token = url.searchParams.get('token') || '';
	});

	async function handleSubmit() {
		errorMessage = '';
		successMessage = '';

		if (!newPassword || !confirmPassword) {
			errorMessage = 'Toate câmpurile sunt obligatorii.';
			return;
		}

		if (newPassword !== confirmPassword) {
			errorMessage = 'Parolele nu se potrivesc.';
			return;
		}

		if (!token) {
			errorMessage = 'Token-ul lipsește. Vă rugăm să accesați din nou linkul de resetare.';
			return;
		}

		isLoading = true;

		try {
			const response = await fetch('api/reset-password', {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({ token, new_password: newPassword })
			});

			const data = await response.json();

			if (!response.ok) {
				throw new Error(data.message || 'Eroare la resetarea parolei');
			}

			successMessage = 'Parola a fost resetată cu succes.';
			setTimeout(() => goto('/login'), 2000);
		} catch (error: any) {
			errorMessage = error.message || 'A apărut o eroare la resetarea parolei.';
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

	<h1>Set New Password</h1>

	<form onsubmit={preventDefault(handleSubmit)}>
		<input
			type="password"
			placeholder="new password"
			bind:value={newPassword}
			required
		/>
		<input
			type="password"
			placeholder="confirm password"
			bind:value={confirmPassword}
			required
		/>

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
		margin-bottom: 20px;
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