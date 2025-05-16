<script lang="ts">
	import { goto } from '$app/navigation';

	// Verifică dacă utilizatorul este autentificat
	import { onMount } from 'svelte';

	let username = $state('');

	onMount(async () => {
		try {
			// Încercăm să obținem informații despre utilizator
			const response = await fetch('/api/profile', {
				method: 'GET',
				headers: {
					Authorization: `Bearer ${localStorage.getItem('jwt_token')}`
				}
			});

			if (!response.ok) {
				throw new Error('Sesiune expirată sau invalidă');
			}

			const data = await response.json();
			username = data.username || 'Utilizator';
		} catch (error) {
			console.error('Eroare la obținerea datelor utilizatorului:', error);
			// În caz de eroare, delogăm utilizatorul
			onLogout();
		}
	});

	// Funcție pentru navigare la pagina de login
	function goToLogin() {
		goto('login');
	}

	function onLogout() {
		localStorage.removeItem('jwt_token'); // Remove the JWT
		goto('login'); // Redirect to login page
	}

	// Funcție pentru navigare la pagina de register
	function goToRegister() {
		goto('register');
	}

	function gotoChat() {
		goto('chat');
	}
</script>

<div class="dashboard">
	<div class="header">
		<div class="logo-container">
			<img src="/ourchat_logo.png" alt="OurChat Logo" class="logo-image" />
			<div class="logo-text">OurChat</div>
		</div>

		{#if true}
			<div class="user-info">
				<span class="username">Bine ai venit, {username}</span>
				<button class="logout-btn" onclick={onLogout}>Logout</button>
			</div>
		{/if}
	</div>

	<div class="content">
		{#if true}
			<h1 class="welcome">Dashboard OurChat</h1>
			<p class="message">
				Aceasta este pagina principală a aplicației OurChat. Din acest dashboard vei putea accesa
				toate funcționalitățile aplicației, odată ce acestea vor fi implementate.
			</p>

			<div class="features">
				<div class="feature-card">
					<div class="feature-icon">💬</div>
					<div class="feature-title">Chat Individual</div>
					<div class="feature-desc">Comunică privat cu orice utilizator OurChat</div>
				</div>

				<div class="feature-card">
					<div class="feature-icon">👥</div>
					<div class="feature-title">Chat de Grup</div>
					<div class="feature-desc">
						Creează grupuri pentru a comunica cu mai mulți utilizatori simultan
					</div>
				</div>

				<div class="feature-card">
					<div class="feature-icon">🔒</div>
					<div class="feature-title">Criptare End-to-End</div>
					<div class="feature-desc">Mesajele tale sunt securizate prin criptare end-to-end</div>
				</div>

				<div class="feature-card">
					<div class="feature-icon">📁</div>
					<div class="feature-title">Partajare Fișiere</div>
					<div class="feature-desc">Trimite și primește fișiere multimedia în timp real</div>
				</div>

				<div class="feature-card" onclick={gotoChat}>
					<div class="feature-icon">!!!</div>
					<div class="feature-title">Accesează chat-ul</div>
					<div class="feature-desc">Dă clic aici</div>
				</div>
			</div>
		{:else}
			<h1 class="welcome">Bine ai venit la OurChat</h1>
			<p class="message">
				OurChat este o aplicație de mesagerie securizată cu criptare end-to-end. Pentru a începe, te
				rugăm să te autentifici sau să îți creezi un cont.
			</p>

			<div class="auth-buttons">
				<button class="auth-button" onclick={goToLogin}>Login</button>
				<button class="auth-button" onclick={goToRegister}>Register</button>
			</div>
		{/if}
	</div>
</div>

<style>
	.dashboard {
		min-height: 100vh;
		background: linear-gradient(135deg, #6a5af9 0%, #4a91ff 100%);
		color: white;
		display: flex;
		flex-direction: column;
	}

	.header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		padding: 20px;
		background-color: rgba(0, 0, 0, 0.2);
	}

	.logo-container {
		display: flex;
		align-items: center;
		cursor: pointer;
	}

	.logo-image {
		width: 50px;
		height: 50px;
		object-fit: contain;
	}

	.logo-text {
		font-size: 24px;
		font-weight: bold;
		color: white;
		margin-left: 10px;
	}

	.user-info {
		display: flex;
		align-items: center;
	}

	.username {
		margin-right: 15px;
		font-size: 16px;
	}

	.logout-btn {
		padding: 8px 15px;
		border-radius: 20px;
		border: none;
		background-color: rgba(255, 255, 255, 0.2);
		color: white;
		cursor: pointer;
		transition: all 0.3s ease;
	}

	.logout-btn:hover {
		background-color: rgba(255, 255, 255, 0.3);
	}

	.content {
		flex: 1;
		padding: 30px;
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
	}

	.welcome {
		font-size: 36px;
		margin-bottom: 20px;
		text-align: center;
	}

	.message {
		font-size: 18px;
		margin-bottom: 40px;
		text-align: center;
		max-width: 600px;
	}

	.features {
		display: grid;
		grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
		gap: 20px;
		width: 100%;
		max-width: 1000px;
	}

	.feature-card {
		background-color: rgba(255, 255, 255, 0.1);
		padding: 20px;
		border-radius: 10px;
		display: flex;
		flex-direction: column;
		align-items: center;
		text-align: center;
		transition: all 0.3s ease;
	}

	.feature-card:hover {
		background-color: rgba(255, 255, 255, 0.2);
		transform: translateY(-5px);
	}

	.feature-icon {
		font-size: 48px;
		margin-bottom: 15px;
	}

	.feature-title {
		font-size: 20px;
		font-weight: bold;
		margin-bottom: 10px;
	}

	.feature-desc {
		font-size: 14px;
	}

	.auth-buttons {
		display: flex;
		gap: 20px;
		margin-top: 40px;
	}

	.auth-button {
		padding: 10px 20px;
		border-radius: 20px;
		border: none;
		background-color: rgba(255, 255, 255, 0.2);
		color: white;
		cursor: pointer;
		transition: all 0.3s ease;
	}

	.auth-button:hover {
		background-color: rgba(255, 255, 255, 0.3);
	}

	@media (max-width: 768px) {
		.welcome {
			font-size: 28px;
		}

		.message {
			font-size: 16px;
		}

		.auth-buttons {
			flex-direction: column;
			gap: 10px;
		}
	}
</style>
