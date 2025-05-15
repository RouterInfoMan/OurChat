<script lang="ts">
	import { goto } from '$app/navigation';

	let newPassword = '';
	let confirmPassword = '';
	let errorMessage = '';
	let successMessage = '';
	let isLoading = false;

	async function handleNewPassword() {
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

		isLoading = true;

		try {
			const response = await fetch('/reset-password', {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({ password: newPassword })
			});

			const data = await response.json();

			if (!response.ok) {
				throw new Error(data.message || 'Eroare la resetarea parolei');
			}

			successMessage = 'Parola a fost resetată cu succes.';
			setTimeout(() => goto('/login'), 2000);
		} catch (error: any) {
			errorMessage = error.message;
		} finally {
			isLoading = false;
		}
	}
</script>

<div class="container">
	<h1>Setează o parolă nouă</h1>

	<input
		type="password"
		placeholder="Parola nouă"
		bind:value={newPassword}
		required
	/>
	<input
		type="password"
		placeholder="Confirmă parola"
		bind:value={confirmPassword}
		required
	/>

	<button on:click={handleNewPassword} disabled={isLoading}>
		{isLoading ? 'Se procesează...' : 'Salvează parola'}
	</button>

	{#if errorMessage}
		<div class="error">{errorMessage}</div>
	{/if}
	{#if successMessage}
		<div class="success">{successMessage}</div>
	{/if}
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

	button:disabled {
		opacity: 0.7;
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
