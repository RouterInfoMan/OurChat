<script lang="ts">
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';

	let loading = $state(true);
	let chats:
		| [
				{
					id: number;
					name: string;
				}
		  ]
		| null = $state(null);
	let error: Error | null = $state(null);

	let selected_chat: number | null = $state(null);
	// false = nu se încarcă, string = mesaj de eroare
	// true = se încarcă
	let selected_chat_data:
		| {
				id: number;
				name: string;
		  }
		| boolean = $state(false);
	// string = mesaj de eroare, null = se incarcă
	let chat_messages:
		| {
				id: number;
				sender_id: number;
				chat_id: number;
				content: string;
				created_at: string;
				is_read: string;
		  }[]
		| string
		| null = $state(null);

	let show_new_chat_popover = $state(false);
	let new_chat_entites = $state('');
	// false = nu se creează chat, string = mesaj de eroare
	// true = se creează chat
	let new_chat_creating: boolean | string = $state(false);

	let messageInput: HTMLTextAreaElement | null;
	let currentMessageText = $state('');
	let isSending = $state(false);

	onMount(async () => {
		// Selectează elementele DOM cu type assertions
		messageInput = document.querySelector('.message-input') as HTMLTextAreaElement | null;
		const messagesContainer = document.querySelector(
			'.messages-container'
		) as HTMLDivElement | null;
		const conversationItems = document.querySelectorAll('.conversation-item');

		// Auto-resize pentru textarea
		if (messageInput) {
			messageInput.addEventListener('input', function () {
				this.style.height = 'auto';
				const newHeight = Math.min(this.scrollHeight, 120);
				this.style.height = newHeight + 'px';
			});
		}

		loadEverything();
	});

	// Funcția pentru a trimite un mesaj
	async function sendMessage(): Promise<void> {
		//if (!messageInput || !selected_chat) return;

		const text = currentMessageText.trim();
		if (!text || isSending) return;

		try {
			isSending = true;

			// Crează un obiect temporar de mesaj pentru afișarea imediată în UI
			const tempMessage = {
				id: Date.now(), // ID temporar
				sender_id: -1, // Va fi înlocuit de backend
				chat_id: selected_chat,
				content: text,
				created_at: new Date().toISOString(),
				is_read: 'false'
			};

			// Adaugă mesajul temporar la lista de mesaje pentru feedback instant
			if (Array.isArray(chat_messages)) {
				chat_messages = [...chat_messages, tempMessage];
			} else {
				chat_messages = [tempMessage];
			}

			// Golește input-ul
			currentMessageText = '';
			if (messageInput) {
				messageInput.value = '';
				messageInput.style.height = 'auto';
			}

			// Scroll la ultima poziție după ce DOM-ul s-a actualizat
			setTimeout(() => {
				const messagesContainer = document.querySelector('.messages-container');
				if (messagesContainer) {
					messagesContainer.scrollTop = messagesContainer.scrollHeight;
				}
			}, 50);

			// Trimite mesajul la server
			const response = await fetch(`/api/chats/${selected_chat}/messages`, {
				method: 'POST',
				headers: {
					Authorization: `Bearer ${localStorage.getItem('jwt_token')}`,
					'Content-Type': 'application/json'
				},
				body: JSON.stringify({
					content: text
				})
			});

			if (!response.ok) {
				throw new Error('Eroare la trimiterea mesajului');
			}

			// Reîncarcă mesajele pentru a obține mesajul cu ID-ul real
			await loadChatMessages();
		} catch (error) {
			console.error('Error sending message:', error);
		} finally {
			isSending = false;
		}
	}

	// (Re)setează starea aplicației făcând cereri la backend
	async function loadEverything() {
		loading = true;
		chats = null;
		error = null;
		selected_chat = null;
		show_new_chat_popover = false;

		try {
			let req = await fetch('/api/chats', {
				method: 'GET',
				headers: {
					Authorization: `Bearer ${localStorage.getItem('jwt_token')}`
				}
			});
			if (!req.ok) {
				throw new Error('Eroare la obținerea conversațiilor');
			}
			chats = await req.json();
		} catch (err) {
			error = err as Error;
		} finally {
			loading = false;
		}
	}

	async function createChat() {
		try {
			new_chat_creating = true;
			let new_chat_users = JSON.parse(new_chat_entites) as number[];

			const response = await fetch('api/chats', {
				method: 'POST',
				headers: {
					Authorization: `Bearer ${localStorage.getItem('jwt_token')}`,
					'Content-Type': 'application/json'
				},
				body: JSON.stringify({
					type: new_chat_users.length > 1 ? 'group' : 'direct',
					users: new_chat_users,
					name: new_chat_users.length > 1 ? 'Group Chat ' + Math.random() : undefined
				})
			});

			if (!response.ok) {
				throw new Error(await response.text());
			}

			new_chat_creating = false;
			await loadEverything();
		} catch (error: any) {
			console.error('Error:', error);
			new_chat_creating = error.message || 'A apărut o eroare la crearea chat-ului.';
		}
	}

	async function loadChat() {
		selected_chat_data = true;
		loadChatMessages();
		try {
			const response = await fetch(`/api/chats/${selected_chat}`, {
				method: 'GET',
				headers: {
					Authorization: `Bearer ${localStorage.getItem('jwt_token')}`
				}
			});

			if (!response.ok) {
				throw new Error('Eroare la obținerea detaliilor conversației');
			}

			selected_chat_data = await response.json();
		} catch (error) {
			console.error('Error:', error);
			selected_chat_data = false;
		}
	}

	async function loadChatMessages() {
		try {
			const response = await fetch(`/api/chats/${selected_chat}/messages`, {
				method: 'GET',
				headers: {
					Authorization: `Bearer ${localStorage.getItem('jwt_token')}`
				}
			});

			if (!response.ok) {
				throw new Error('Eroare la obținerea mesajelor conversației');
			}

			chat_messages = await response.json();
		} catch (error) {
			console.error('Error:', error);
			chat_messages = (error as Error).message;
		}
	}
</script>

<div class="chat-layout">
	<!-- Sidebar-ul albastru îngust -->
	<div class="sidebar-icons">
		<div class="top-icons">
			<a
				href="#"
				class="icon-btn"
				onclick={() => {
					show_new_chat_popover = true;
					new_chat_entites = '';
					new_chat_creating = false;
				}}
				title="New Chat"
			>
				<svg viewBox="0 0 24 24" width="24" height="24"
					><circle cx="12" cy="12" r="10"></circle><line x1="12" y1="8" x2="12" y2="16"></line><line
						x1="8"
						y1="12"
						x2="16"
						y2="12"
					></line></svg
				>
			</a>
			<a href="#" class="icon-btn" title="Find People">
				<svg viewBox="0 0 24 24" width="24" height="24"
					><path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"></path><circle cx="9" cy="7" r="4"
					></circle><path d="M23 21v-2a4 4 0 0 0-3-3.87"></path><path d="M16 3.13a4 4 0 0 1 0 7.75"
					></path></svg
				>
			</a>
		</div>
		<div class="bottom-icons">
			<a href="#" class="icon-btn settings-btn" title="Settings">
				<svg viewBox="0 0 24 24" width="24" height="24"
					><circle cx="12" cy="12" r="3"></circle><path
						d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 0 1 0 2.83 2 2 0 0 1-2.83 0l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-2 2 2 2 0 0 1-2-2v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 0 1-2.83 0 2 2 0 0 1 0-2.83l.06-.06a1.65 1.65 0 0 0 .33-1.82 1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1-2-2 2 2 0 0 1 2-2h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 0 1 0-2.83 2 2 0 0 1 2.83 0l.06.06a1.65 1.65 0 0 0 1.82.33H9a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 2-2 2 2 0 0 1 2 2v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 0 1 2.83 0 2 2 0 0 1 0 2.83l-.06.06a1.65 1.65 0 0 0-.33 1.82V9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 2 2 2 2 0 0 1-2 2h-.09a1.65 1.65 0 0 0-1.51 1z"
					></path></svg
				>
			</a>
			<div
				class="profile-img"
				onclick={() => {
					goto('/');
				}}
				title="Go to Dashboard"
			>
				<img src="/default-avatar.png" alt="Profil" />
			</div>
		</div>
	</div>

	{#if loading}
		<div class="loading-container">
			<div class="loading-spinner"></div>
			<p>Se încarcă...</p>
		</div>
	{:else if error}
		<div class="error-container">
			<div class="error-message">
				<h2>Oops! Ceva nu a mers bine</h2>
				<p>{error.message}</p>
				<button class="retry-btn" onclick={loadEverything}>Încearcă din nou</button>
			</div>
		</div>
	{:else if chats !== null}
		<!-- Sidebar-ul cu lista de conversații -->
		<div class="sidebar-conversations">
			<div class="conversations-header">
				<h2>Conversations</h2>
			</div>

			<div class="conversation-list">
				{#each chats as chat}
					<div
						class="conversation-item {chat.id === selected_chat ? 'active' : ''}"
						onclick={() => {
							selected_chat = chat.id;
							loadChat();
						}}
					>
						<div class="avatar-wrapper">
							<img src="/default-avatar.png" alt="Avatar" />
						</div>
						<div class="conv-details">
							<h3>{chat.name}</h3>
							<p>Click to open chat</p>
						</div>
					</div>
				{/each}
			</div>
		</div>

		<!-- Zona principală de chat -->
		{#if selected_chat !== null}
			{#if typeof selected_chat_data === 'object'}
				<!-- Zona de mesaje -->
				<div class="chat-area">
					<div class="chat-header">
						<div class="current-conversation">
							<div class="avatar-wrapper">
								<img src="/default-avatar.png" alt="Avatar" />
							</div>
							<h2>{selected_chat_data.name}</h2>
						</div>
					</div>

					<div class="messages-container">
						{#if typeof chat_messages === 'string'}
							<div class="error-in-chat">
								<p>Eroare la încărcarea mesajelor: {chat_messages}</p>
							</div>
						{:else if chat_messages === null}
							<div class="loading-messages">
								<div class="loading-spinner small"></div>
								<p>Se încarcă mesajele...</p>
							</div>
						{:else if chat_messages.length === 0}
							<div class="empty-chat">
								<p>Nu există mesaje încă. Începe conversația!</p>
							</div>
						{:else}
							{#each [...chat_messages].reverse() as message}
								<!-- todo - verificare message.sender_id === user_id (care e?) -->
								<div class={message.sender_id === 1 ? 'message-received' : 'message-sent'}>
									<div class="message-bubble">{message.content}</div>
									<div class="message-info">
										<span class="message-time"
											>{new Date(message.created_at).toLocaleTimeString([], {
												hour: '2-digit',
												minute: '2-digit'
											})}</span
										>
										<span class="read-receipt delivered">Livrat</span>
									</div>
								</div>
							{/each}
						{/if}
					</div>

					<div class="message-input-area">
						<button class="attachment-btn" title="Attach File">
							<svg viewBox="0 0 24 24" width="24" height="24"
								><path
									d="M21.44 11.05l-9.19 9.19a6 6 0 0 1-8.49-8.49l9.19-9.19a4 4 0 0 1 5.66 5.66l-9.2 9.19a2 2 0 0 1-2.83-2.83l8.49-8.48"
								></path></svg
							>
						</button>
						<textarea
							bind:value={currentMessageText}
							placeholder="Type a message..."
							class="message-input"
							rows="1"
							onkeypress={(e) => {
								if (e.key === 'Enter' && !e.shiftKey) {
									e.preventDefault();
									sendMessage();
								}
							}}
						></textarea>
						<button class="send-btn" onclick={sendMessage} disabled={isSending} title="Send Message">
							{#if isSending}
								<div class="loading-spinner tiny"></div>
							{:else}
								<svg viewBox="0 0 24 24" width="24" height="24"
									><polygon points="22 2 15 22 11 13 2 9 22 2"></polygon></svg
								>
							{/if}
						</button>
					</div>
				</div>
			{:else if selected_chat_data === true}
				<div class="loading-container">
					<div class="loading-spinner"></div>
					<p>Se încarcă conversația...</p>
				</div>
			{:else}
				<div class="error-container">
					<div class="error-message">
						<h2>Eroare la încărcare</h2>
						<p>Nu s-au putut încărca detaliile conversației</p>
					</div>
				</div>
			{/if}
		{:else}
			<div class="welcome-chat">
				<div class="welcome-content">
					<img src="/ourchat_logo.png" alt="OurChat Logo" class="welcome-logo" />
					<h2>Welcome to OurChat!</h2>
					<p>Select a conversation from the left sidebar to start chatting, or create a new chat using the + button.</p>
				</div>
			</div>
		{/if}
	{/if}
</div>

{#if show_new_chat_popover}
	<div class="popover-overlay">
		<div class="popover-content">
			<h3>Create New Chat</h3>
			<label>
				Enter user IDs (format: [1,2,3]):
				<input type="text" bind:value={new_chat_entites} placeholder="[1,2,3]" />
			</label>
			<div class="popover-actions">
				<button class="primary-btn" onclick={createChat} disabled={new_chat_creating === true}>
					{#if new_chat_creating === true}
						<div class="loading-spinner tiny"></div>
						Creating...
					{:else}
						Create Chat
					{/if}
				</button>
				<button
					class="secondary-btn"
					onclick={() => {
						show_new_chat_popover = false;
					}}>Cancel</button
				>
			</div>
			{#if typeof new_chat_creating === 'string'}
				<div class="popover-error">
					Error: {new_chat_creating}
				</div>
			{/if}
		</div>
	</div>
{/if}

<style>
	/* Global body styling */
	:global(body) {
		margin: 0;
		padding: 0;
		font-family: 'Arial', sans-serif;
		background: linear-gradient(135deg, #6a5af9 0%, #4a91ff 100%);
		overflow: hidden;
	}

	/* Layout de bază */
	.chat-layout {
		display: flex;
		position: fixed;
		top: 0;
		left: 0;
		right: 0;
		bottom: 0;
		width: 100%;
		height: 100vh;
		background: linear-gradient(135deg, #6a5af9 0%, #4a91ff 100%);
	}

	/* Loading and Error States */
	.loading-container,
	.error-container {
		flex: 1;
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		color: white;
		text-align: center;
		padding: 40px;
	}

	.loading-spinner {
		width: 40px;
		height: 40px;
		border: 4px solid rgba(255, 255, 255, 0.3);
		border-top: 4px solid white;
		border-radius: 50%;
		animation: spin 1s linear infinite;
		margin-bottom: 20px;
	}

	.loading-spinner.small {
		width: 30px;
		height: 30px;
		border-width: 3px;
		margin-bottom: 15px;
	}

	.loading-spinner.tiny {
		width: 20px;
		height: 20px;
		border-width: 2px;
		margin: 0;
	}

	@keyframes spin {
		0% { transform: rotate(0deg); }
		100% { transform: rotate(360deg); }
	}

	.error-message h2 {
		margin: 0 0 10px 0;
		font-size: 24px;
	}

	.error-message p {
		margin: 0 0 20px 0;
		opacity: 0.8;
	}

	.retry-btn {
		background-color: rgba(255, 255, 255, 0.2);
		color: white;
		border: none;
		padding: 12px 24px;
		border-radius: 20px;
		cursor: pointer;
		font-size: 16px;
		transition: all 0.3s ease;
	}

	.retry-btn:hover {
		background-color: rgba(255, 255, 255, 0.3);
		transform: translateY(-2px);
	}

	/* Sidebar-ul cu iconițe (albastru) */
	.sidebar-icons {
		width: 60px;
		background: linear-gradient(180deg, #5e4bff 0%, #4a68eb 100%);
		display: flex;
		flex-direction: column;
		justify-content: space-between;
		padding: 15px 0;
		z-index: 10;
		box-shadow: 2px 0 10px rgba(0, 0, 0, 0.1);
	}

	.top-icons,
	.bottom-icons {
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: 15px;
	}

	.icon-btn {
		width: 40px;
		height: 40px;
		border-radius: 50%;
		display: flex;
		align-items: center;
		justify-content: center;
		color: white;
		transition: all 0.3s ease;
		cursor: pointer;
	}

	.icon-btn:hover {
		background-color: rgba(255, 255, 255, 0.2);
		transform: scale(1.1);
	}

	.icon-btn svg {
		width: 24px;
		height: 24px;
		stroke: currentColor;
		stroke-width: 2;
		fill: none;
	}

	.profile-img {
		width: 40px;
		height: 40px;
		border-radius: 50%;
		overflow: hidden;
		margin-top: 10px;
		border: 2px solid rgba(255, 255, 255, 0.3);
		cursor: pointer;
		transition: all 0.3s ease;
	}

	.profile-img:hover {
		border-color: rgba(255, 255, 255, 0.6);
		transform: scale(1.1);
	}

	.profile-img img {
		width: 100%;
		height: 100%;
		object-fit: cover;
	}

	/* Sidebar-ul cu conversații */
	.sidebar-conversations {
		width: 280px;
		background: rgba(255, 255, 255, 0.1);
		backdrop-filter: blur(10px);
		border-right: 1px solid rgba(255, 255, 255, 0.2);
		display: flex;
		flex-direction: column;
		color: white;
	}

	/* Header pentru secțiunea de conversații */
	.conversations-header {
		padding: 20px;
		border-bottom: 1px solid rgba(255, 255, 255, 0.2);
		background: rgba(255, 255, 255, 0.05);
	}

	.conversations-header h2 {
		margin: 0;
		font-size: 22px;
		font-weight: 600;
		color: white;
	}

	.conversation-list {
		flex: 1;
		overflow-y: auto;
	}

	.conversation-item {
		display: flex;
		align-items: center;
		padding: 15px 20px;
		cursor: pointer;
		transition: all 0.3s ease;
		border-bottom: 1px solid rgba(255, 255, 255, 0.1);
	}

	.conversation-item:hover {
		background-color: rgba(255, 255, 255, 0.1);
	}

	.conversation-item.active {
		background-color: rgba(255, 255, 255, 0.2);
		border-left: 4px solid white;
	}

	.avatar-wrapper {
		width: 50px;
		height: 50px;
		border-radius: 50%;
		overflow: hidden;
		flex-shrink: 0;
		border: 2px solid rgba(255, 255, 255, 0.3);
	}

	.avatar-wrapper img {
		width: 100%;
		height: 100%;
		object-fit: cover;
	}

	.conv-details {
		margin-left: 15px;
		overflow: hidden;
	}

	.conv-details h3 {
		margin: 0;
		font-size: 16px;
		font-weight: 500;
		color: white;
		white-space: nowrap;
		overflow: hidden;
		text-overflow: ellipsis;
	}

	.conv-details p {
		margin: 4px 0 0;
		font-size: 13px;
		color: rgba(255, 255, 255, 0.7);
		white-space: nowrap;
		overflow: hidden;
		text-overflow: ellipsis;
	}

	/* Welcome screen when no chat selected */
	.welcome-chat {
		flex: 1;
		display: flex;
		align-items: center;
		justify-content: center;
		color: white;
		text-align: center;
	}

	.welcome-content {
		max-width: 400px;
		padding: 40px 20px;
	}

	.welcome-logo {
		width: 120px;
		height: 120px;
		margin-bottom: 30px;
		opacity: 0.9;
	}

	.welcome-content h2 {
		font-size: 28px;
		margin: 0 0 20px 0;
		font-weight: 600;
	}

	.welcome-content p {
		font-size: 16px;
		line-height: 1.5;
		opacity: 0.8;
		margin: 0;
	}

	/* Zona principală de chat */
	.chat-area {
		flex: 1;
		display: flex;
		flex-direction: column;
		background: rgba(255, 255, 255, 0.05);
		backdrop-filter: blur(10px);
	}

	.chat-header {
		padding: 20px;
		border-bottom: 1px solid rgba(255, 255, 255, 0.2);
		background: rgba(255, 255, 255, 0.1);
		color: white;
	}

	.current-conversation {
		display: flex;
		align-items: center;
	}

	.current-conversation .avatar-wrapper {
		width: 50px;
		height: 50px;
	}

	.current-conversation h2 {
		margin: 0 0 0 15px;
		font-size: 20px;
		font-weight: 500;
		color: white;
	}

	.messages-container {
		flex: 1;
		padding: 20px;
		overflow-y: auto;
		display: flex;
		flex-direction: column;
		gap: 12px;
	}

	/* Message states */
	.loading-messages,
	.error-in-chat,
	.empty-chat {
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		height: 100%;
		color: rgba(255, 255, 255, 0.8);
		text-align: center;
	}

	.error-in-chat p,
	.empty-chat p {
		background: rgba(255, 255, 255, 0.1);
		padding: 15px 20px;
		border-radius: 10px;
		margin: 0;
	}

	/* Message bubbles */
	.message-bubble {
		padding: 12px 16px;
		border-radius: 18px;
		max-width: 100%;
		word-wrap: break-word;
		overflow-wrap: break-word;
		word-break: break-word;
		box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
	}

	.message-sent {
		align-self: flex-end;
		max-width: 70%;
	}

	.message-received {
		align-self: flex-start;
		max-width: 70%;
	}

	.message-sent .message-bubble {
		background: linear-gradient(135deg, #5977ff 0%, #4a68eb 100%);
		color: white;
		border-bottom-right-radius: 6px;
	}

	.message-received .message-bubble {
		background: rgba(255, 255, 255, 0.9);
		color: #333;
		border-bottom-left-radius: 6px;
	}

	/* Message info */
	.message-info {
		display: flex;
		align-items: center;
		font-size: 11px;
		color: rgba(255, 255, 255, 0.7);
		margin-top: 4px;
	}

	.message-sent .message-info {
		justify-content: flex-end;
		padding-right: 8px;
	}

	.message-received .message-info {
		padding-left: 8px;
	}

	.message-time {
		margin-right: 6px;
	}

	.read-receipt {
		display: flex;
		align-items: center;
		font-size: 10px;
	}

	.read-receipt:before {
		content: '';
		display: inline-block;
		width: 12px;
		height: 12px;
		margin-right: 3px;
		background-position: center;
		background-repeat: no-repeat;
		background-size: contain;
	}

	.read-receipt.delivered:before {
		background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='16' height='16' viewBox='0 0 24 24' fill='none' stroke='rgba(255,255,255,0.7)' stroke-width='2' stroke-linecap='round' stroke-linejoin='round'%3E%3Cpolyline points='20 6 9 17 4 12'%3E%3C/polyline%3E%3C/svg%3E");
	}

	.read-receipt.seen:before {
		background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='16' height='16' viewBox='0 0 24 24' fill='none' stroke='%23ffffff' stroke-width='2' stroke-linecap='round' stroke-linejoin='round'%3E%3Cpolyline points='20 6 9 17 4 12'%3E%3C/polyline%3E%3C/svg%3E");
		color: #ffffff;
	}

	/* Message input area */
	.message-input-area {
		padding: 20px;
		background: rgba(255, 255, 255, 0.1);
		border-top: 1px solid rgba(255, 255, 255, 0.2);
		display: flex;
		align-items: center;
		gap: 12px;
	}

	.attachment-btn,
	.send-btn {
		width: 44px;
		height: 44px;
		border-radius: 50%;
		border: none;
		background: rgba(255, 255, 255, 0.2);
		display: flex;
		align-items: center;
		justify-content: center;
		cursor: pointer;
		flex-shrink: 0;
		transition: all 0.3s ease;
		color: white;
	}

	.attachment-btn:hover,
	.send-btn:hover {
		background: rgba(255, 255, 255, 0.3);
		transform: scale(1.1);
	}

	.send-btn:disabled {
		opacity: 0.5;
		cursor: not-allowed;
		transform: none;
	}

	.attachment-btn svg,
	.send-btn svg {
		width: 22px;
		height: 22px;
		stroke-width: 2;
		fill: none;
		stroke: currentColor;
	}

	.send-btn svg {
		fill: currentColor;
		stroke-width: 0;
	}

	.message-input {
		flex: 1;
		border: none;
		background: rgba(255, 255, 255, 0.2);
		border-radius: 22px;
		padding: 12px 18px;
		resize: none;
		font-size: 15px;
		max-height: 120px;
		min-height: 20px;
		outline: none;
		color: white;
		font-family: inherit;
		transition: all 0.3s ease;
	}

	.message-input::placeholder {
		color: rgba(255, 255, 255, 0.7);
	}

	.message-input:focus {
		background: rgba(255, 255, 255, 0.3);
		box-shadow: 0 0 0 2px rgba(255, 255, 255, 0.3);
	}

	/* Popover styles */
	.popover-overlay {
		position: fixed;
		top: 0;
		left: 0;
		right: 0;
		bottom: 0;
		background: rgba(0, 0, 0, 0.5);
		display: flex;
		align-items: center;
		justify-content: center;
		z-index: 1000;
		backdrop-filter: blur(5px);
	}

	.popover-content {
		background: linear-gradient(135deg, #6a5af9 0%, #4a91ff 100%);
		padding: 30px;
		border-radius: 20px;
		box-shadow: 0 10px 40px rgba(0, 0, 0, 0.3);
		min-width: 400px;
		max-width: 90vw;
		color: white;
		border: 1px solid rgba(255, 255, 255, 0.2);
	}

	.popover-content h3 {
		margin: 0 0 20px 0;
		font-size: 24px;
		font-weight: 600;
		text-align: center;
	}

	.popover-content label {
		display: block;
		margin-bottom: 20px;
		font-size: 16px;
		font-weight: 500;
	}

	.popover-content input {
		width: 100%;
		padding: 12px 16px;
		border-radius: 10px;
		border: none;
		background: rgba(255, 255, 255, 0.2);
		color: white;
		font-size: 16px;
		margin-top: 8px;
		outline: none;
		transition: all 0.3s ease;
	}

	.popover-content input::placeholder {
		color: rgba(255, 255, 255, 0.7);
	}

	.popover-content input:focus {
		background: rgba(255, 255, 255, 0.3);
		box-shadow: 0 0 0 2px rgba(255, 255, 255, 0.3);
	}

	.popover-actions {
		display: flex;
		gap: 12px;
		justify-content: center;
	}

	.primary-btn,
	.secondary-btn {
		padding: 12px 24px;
		border-radius: 25px;
		border: none;
		font-size: 16px;
		font-weight: 600;
		cursor: pointer;
		transition: all 0.3s ease;
		display: flex;
		align-items: center;
		gap: 8px;
	}

	.primary-btn {
		background: rgba(255, 255, 255, 0.2);
		color: white;
	}

	.primary-btn:hover:not(:disabled) {
		background: rgba(255, 255, 255, 0.3);
		transform: translateY(-2px);
	}

	.primary-btn:disabled {
		opacity: 0.5;
		cursor: not-allowed;
		transform: none;
	}

	.secondary-btn {
		background: transparent;
		color: white;
		border: 2px solid rgba(255, 255, 255, 0.3);
	}

	.secondary-btn:hover {
		background: rgba(255, 255, 255, 0.1);
		border-color: rgba(255, 255, 255, 0.5);
	}

	.popover-error {
		color: #ff6b6b;
		text-align: center;
		margin-top: 15px;
		font-size: 14px;
		background: rgba(255, 255, 255, 0.1);
		padding: 10px;
		border-radius: 8px;
		border: 1px solid rgba(255, 107, 107, 0.3);
	}

	/* Scrollbar styling */
	.conversation-list::-webkit-scrollbar,
	.messages-container::-webkit-scrollbar {
		width: 6px;
	}

	.conversation-list::-webkit-scrollbar-track,
	.messages-container::-webkit-scrollbar-track {
		background: rgba(255, 255, 255, 0.1);
	}

	.conversation-list::-webkit-scrollbar-thumb,
	.messages-container::-webkit-scrollbar-thumb {
		background: rgba(255, 255, 255, 0.3);
		border-radius: 3px;
	}

	.conversation-list::-webkit-scrollbar-thumb:hover,
	.messages-container::-webkit-scrollbar-thumb:hover {
		background: rgba(255, 255, 255, 0.5);
	}

	/* Media queries pentru responsivitate */
	@media (max-width: 768px) {
		.sidebar-conversations {
			width: 240px;
		}

		.welcome-content {
			padding: 20px;
		}

		.welcome-logo {
			width: 80px;
			height: 80px;
		}

		.welcome-content h2 {
			font-size: 24px;
		}

		.popover-content {
			min-width: 320px;
			padding: 20px;
		}
	}

	@media (max-width: 576px) {
		.chat-layout {
			flex-direction: column;
		}

		.sidebar-icons {
			width: 100%;
			height: 60px;
			flex-direction: row;
			padding: 0 15px;
		}

		.top-icons,
		.bottom-icons {
			flex-direction: row;
		}

		.sidebar-conversations {
			width: 100%;
			height: 120px;
			overflow-x: auto;
			overflow-y: hidden;
		}

		.conversations-header {
			writing-mode: vertical-rl;
			width: 80px;
			height: 100%;
			display: flex;
			align-items: center;
			justify-content: center;
			border-right: 1px solid rgba(255, 255, 255, 0.2);
			border-bottom: none;
		}

		.conversation-list {
			display: flex;
			flex-direction: row;
			height: 100%;
		}

		.conversation-item {
			min-width: 200px;
			height: 100%;
			border-right: 1px solid rgba(255, 255, 255, 0.1);
			border-bottom: none;
		}

		.chat-area {
			height: calc(100vh - 180px);
		}

		.welcome-content {
			padding: 15px;
		}

		.popover-content {
			min-width: 280px;
			margin: 20px;
		}

		.popover-actions {
			flex-direction: column;
			gap: 8px;
		}
	}
</style>