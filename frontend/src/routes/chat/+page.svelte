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

	// Gestionare selectare conversații
	//conversationItems.forEach((item) => {
	//	item.addEventListener('click', function (this: Element) {
	//		if (!messagesContainer) return;

	//		// Elimină clasa active de la toate conversațiile
	//		conversationItems.forEach((conv) => {
	//			conv.classList.remove('active');
	//		});

	//		// Adaugă clasa active la conversația selectată
	//		this.classList.add('active');

	//		// Actualizează informațiile din header
	//		const nameElement = this.querySelector('.conv-details h3');
	//		const headerName = document.querySelector('.current-conversation h2');

	//		if (nameElement && headerName) {
	//			const name = nameElement.textContent || '';
	//			headerName.textContent = name;

	//			// Actualizează avatarul din header
	//			const avatarImg = this.querySelector('.avatar-wrapper img') as HTMLImageElement | null;
	//			const headerAvatar = document.querySelector(
	//				'.current-conversation .avatar-wrapper img'
	//			) as HTMLImageElement | null;

	//			if (avatarImg && headerAvatar && avatarImg.getAttribute('src')) {
	//				headerAvatar.setAttribute('src', avatarImg.getAttribute('src') || '');
	//			}

	//			// Curăță mesageria și adaugă statusul
	//			messagesContainer.innerHTML = '';

	//			const statusMessage = document.createElement('div');
	//			statusMessage.className = 'status-message';
	//			const statusParagraph = document.createElement('p');
	//			statusParagraph.textContent = `Ai selectat conversația: ${name}`;
	//			statusMessage.appendChild(statusParagraph);
	//			messagesContainer.appendChild(statusMessage);
	//		}
	//	});
	//});

	// Funcție pentru integrarea cu backend-ul de read receipts
	// Aceasta va fi implementată când vei conecta cu backend-ul real
	//	function updateMessageStatus(messageId: string, status: 'delivered' | 'seen'): void {
	//		// Găsește mesajul după ID
	//		const messageElement = document.querySelector(`[data-message-id="${messageId}"]`);
	//		if (!messageElement) return;

	//		// Găsește read receipt-ul
	//		const readReceipt = messageElement.querySelector('.read-receipt');
	//		if (!readReceipt) return;

	//		// Actualizează statusul
	//		if (status === 'delivered') {
	//			readReceipt.className = 'read-receipt delivered';
	//			readReceipt.textContent = 'Livrat';
	//		} else if (status === 'seen') {
	//			readReceipt.className = 'read-receipt seen';
	//			readReceipt.textContent = 'Văzut';
	//		}
	//	}
	//});

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
			<a href="#" class="icon-btn">
				<svg viewBox="0 0 24 24" width="24" height="24"
					><path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"></path><circle cx="9" cy="7" r="4"
					></circle><path d="M23 21v-2a4 4 0 0 0-3-3.87"></path><path d="M16 3.13a4 4 0 0 1 0 7.75"
					></path></svg
				>
			</a>
		</div>
		<div class="bottom-icons">
			<a href="#" class="icon-btn settings-btn">
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
			>
				<img src="/default-avatar.png" alt="Profil" />
			</div>
		</div>
	</div>

	{#if loading}
		Se încarcă...
	{:else if chats !== null}
		<!-- Sidebar-ul cu lista de conversații -->
		<div class="sidebar-conversations">
			<div class="conversations-header">
				<h2>Chats</h2>
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
							<p>Mesaj recent</p>
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
						<div class="status-message">
							<p>Ai selectat conversația: {selected_chat_data.name}</p>
						</div>

						<!-- Exemplu de mesaj cu receipt
						<div class="message-sent">
							<div class="message-bubble">Salut! Cum merge proiectul?</div>
							<div class="message-info">
								<span class="message-time">14:25</span>
								<span class="read-receipt seen">Văzut</span>
							</div>
						</div>
						-->
						{#if typeof chat_messages === 'string'}
							Eroare la încărcarea mesajelor: {chat_messages}
						{:else if chat_messages === null}
							Se încarcă mesajele...
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
						<button class="attachment-btn">
							<svg viewBox="0 0 24 24" width="24" height="24"
								><path
									d="M21.44 11.05l-9.19 9.19a6 6 0 0 1-8.49-8.49l9.19-9.19a4 4 0 0 1 5.66 5.66l-9.2 9.19a2 2 0 0 1-2.83-2.83l8.49-8.48"
								></path></svg
							>
						</button>
						<textarea
							bind:value={currentMessageText}
							placeholder="Tastează un mesaj..."
							class="message-input"
							rows="1"
							onkeypress={(e) => {
								if (e.key === 'Enter' && !e.shiftKey) {
									e.preventDefault();
									sendMessage();
								}
							}}
						></textarea>
						<button class="send-btn" onclick={sendMessage}>
							<svg viewBox="0 0 24 24" width="24" height="24"
								><polygon points="22 2 15 22 11 13 2 9 22 2"></polygon></svg
							>
						</button>
					</div>
				</div>
			{:else if selected_chat_data === true}
				Se încarcă conversația...
			{:else}
				eroare
			{/if}
		{:else}
			Alege un chat din stânga.
		{/if}
	{:else}
		eroare
	{/if}
</div>

{#if show_new_chat_popover}
	<div class="popover-overlay">
		<div class="popover-content">
			<label>
				Cu cine veri să vorbești? (format: [1,2,3] cu id-uri de utilizatori):
				<input type="text" bind:value={new_chat_entites} />
			</label>
			<div style="margin-top:10px;">
				<button onclick={createChat}>Adaugă</button>
				<button
					onclick={() => {
						show_new_chat_popover = false;
					}}>Închide</button
				>
			</div>
			<div style="color: red; margin-top: 10px;">
				{#if typeof new_chat_creating === 'string'}
					Eroare: {new_chat_creating}
				{:else if new_chat_creating === true}
					Se creează chat...
				{/if}
			</div>
		</div>
	</div>
{/if}

<style>
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
	}

	/* Sidebar-ul cu iconițe (albastru) */
	.sidebar-icons {
		width: 60px;
		background-color: #1900fb;
		display: flex;
		flex-direction: column;
		justify-content: space-between;
		padding: 15px 0;
		z-index: 10;
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
		color: rgb(255, 255, 255);
		transition: background-color 0.2s;
	}

	.icon-btn:hover {
		background-color: rgba(255, 255, 255, 0.1);
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
	}

	.profile-img img {
		width: 100%;
		height: 100%;
		object-fit: cover;
	}

	/* Sidebar-ul cu conversații */
	.sidebar-conversations {
		width: 260px;
		background-color: rgb(0, 145, 255);
		border-right: 1px solid #000000;
		display: flex;
		flex-direction: column;
	}

	/* Header pentru secțiunea de conversații */
	.conversations-header {
		padding: 10px 20px; /* Reduce padding-ul de sus și jos la 10px, la fel ca în chat-header */
		border-bottom: 1px solid #000000;
		height: 60px; /* Setează o înălțime fixă */
		display: flex;
		align-items: center; /* Centrează vertical conținutul */
	}

	.conversations-header h2 {
		margin: 0;
		font-size: 22px;
		font-weight: 600;
		color: #050505;
	}

	.conversation-list {
		flex: 1;
		overflow-y: auto;
	}

	.conversation-item {
		display: flex;
		align-items: center;
		padding: 12px 15px;
		cursor: pointer;
		transition: background-color 0.2s;
		border-bottom: 1px solid #000000;
	}

	.conversation-item.active {
		background-color: rgb(47, 64, 255);
	}

	.conversation-item:hover {
		background-color: #f5f7fa;
	}

	.avatar-wrapper {
		width: 50px;
		height: 50px;
		border-radius: 50%;
		overflow: hidden;
		flex-shrink: 0;
		background-color: #000000;
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
		font-size: 14px;
		font-weight: 500;
		color: #050505;
		white-space: nowrap;
		overflow: hidden;
		text-overflow: ellipsis;
	}

	.conv-details p {
		margin: 4px 0 0;
		font-size: 12px;
		color: #ffffff;
		white-space: nowrap;
		overflow: hidden;
		text-overflow: ellipsis;
	}

	/* Zona principală de chat */
	.chat-area {
		flex: 1;
		display: flex;
		flex-direction: column;
		background: linear-gradient(to left, #4cb1ff, #88a8ff);
	}

	.chat-header {
		height: 60px; /* Înălțime fixă, 60px (40px conținut + 20px total padding) */
		padding: 10px 20px;
		border-bottom: 1px solid #000000;
		background-color: rgb(0, 145, 255);
		display: flex;
		align-items: center;
	}

	.current-conversation {
		display: flex;
		align-items: center;
	}

	.current-conversation .avatar-wrapper {
		width: 40px;
		height: 40px;
	}

	.current-conversation h2 {
		margin: 0 0 0 15px;
		font-size: 16px;
		font-weight: 500;
		color: #050505;
	}

	.messages-container {
		flex: 1;
		padding: 20px;
		overflow-y: auto;
		display: flex;
		flex-direction: column;
		gap: 10px;
	}

	.message-bubble {
		padding: 10px 15px;
		border-radius: 18px;
		max-width: 100%;
		word-wrap: break-word;
		overflow-wrap: break-word;
		word-break: break-word;
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
		background-color: #1877f2;
		color: white;
		border-bottom-right-radius: 4px;
	}

	.message-received .message-bubble {
		background-color: white;
		color: #050505;
		border-bottom-left-radius: 4px;
		box-shadow: 0 1px 2px rgba(0, 0, 0, 0.1);
	}

	/* Stiluri pentru read receipt și informații mesaj */
	.message-info {
		display: flex;
		align-items: center;
		font-size: 11px;
		color: #65676b;
		margin-top: 2px;
	}

	.message-sent .message-info {
		justify-content: flex-end;
		padding-right: 8px;
	}

	.message-received .message-info {
		padding-left: 8px;
	}

	.message-time {
		margin-right: 4px;
	}

	.read-receipt {
		display: flex;
		align-items: center;
	}

	.read-receipt:before {
		content: '';
		display: inline-block;
		width: 14px;
		height: 14px;
		margin-right: 3px;
		background-position: center;
		background-repeat: no-repeat;
		background-size: contain;
	}

	.read-receipt.delivered:before {
		background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='16' height='16' viewBox='0 0 24 24' fill='none' stroke='%2365676B' stroke-width='2' stroke-linecap='round' stroke-linejoin='round'%3E%3Cpolyline points='20 6 9 17 4 12'%3E%3C/polyline%3E%3C/svg%3E");
	}

	.read-receipt.seen:before {
		background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='16' height='16' viewBox='0 0 24 24' fill='none' stroke='%231877F2' stroke-width='2' stroke-linecap='round' stroke-linejoin='round'%3E%3Cpolyline points='20 6 9 17 4 12'%3E%3C/polyline%3E%3C/svg%3E");
		color: #1877f2;
	}

	.status-message {
		align-self: center;
		background-color: #e4e6eb;
		padding: 8px 16px;
		border-radius: 18px;
		color: #65676b;
		font-size: 13px;
		max-width: 70%;
		text-align: center;
	}

	.status-message p {
		margin: 0;
	}

	.message-input-area {
		padding: 10px 15px;
		background-color: rgb(126, 199, 255);
		border-top: 1px solid #000000;
		display: flex;
		align-items: center;
		gap: 10px;
	}

	.attachment-btn,
	.send-btn {
		width: 36px;
		height: 36px;
		border-radius: 50%;
		border: none;
		background: none;
		display: flex;
		align-items: center;
		justify-content: center;
		cursor: pointer;
		flex-shrink: 0;
	}

	.attachment-btn {
		color: #65676b;
	}

	.send-btn {
		color: #1877f2;
	}

	.attachment-btn:hover,
	.send-btn:hover {
		background-color: #0037b8;
	}

	.attachment-btn svg,
	.send-btn svg {
		width: 20px;
		height: 20px;
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
		background-color: #a9d7ff;
		border-radius: 20px;
		padding: 9px 12px;
		resize: none;
		font-size: 15px;
		max-height: 120px;
		min-height: 20px;
		outline: none;
	}

	.message-input:focus {
		outline: none;
	}

	/* Media queries pentru responsivitate */
	@media (max-width: 768px) {
		.sidebar-conversations {
			width: 200px;
		}
	}

	@media (max-width: 576px) {
		.chat-layout {
			flex-direction: column;
		}

		.sidebar-icons {
			width: 100%;
			height: 50px;
			flex-direction: row;
			padding: 0 10px;
		}

		.top-icons,
		.bottom-icons {
			flex-direction: row;
		}

		.sidebar-conversations {
			width: 100%;
			height: 80px;
			overflow-x: auto;
			overflow-y: hidden;
			display: flex;
			flex-direction: row;
		}

		.conversations-header {
			width: 80px;
			height: 100%;
			display: flex;
			align-items: center;
			justify-content: center;
			border-right: 1px solid #e4e6eb;
			border-bottom: none;
		}

		.conversation-list {
			display: flex;
			flex-direction: row;
		}

		.conversation-item {
			min-width: 200px;
			border-right: 1px solid #e4e6eb;
			border-bottom: none;
		}

		.chat-area {
			height: calc(100vh - 130px);
		}
	}
	.popover-overlay {
		position: fixed;
		top: 0;
		left: 0;
		right: 0;
		bottom: 0;
		background: rgba(0, 0, 0, 0.2);
		display: flex;
		align-items: center;
		justify-content: center;
		z-index: 1000;
	}
	.popover-content {
		background: #fff;
		padding: 24px 20px 16px 20px;
		border-radius: 12px;
		box-shadow: 0 2px 16px rgba(0, 0, 0, 0.18);
		min-width: 300px;
		max-width: 90vw;
		display: flex;
		flex-direction: column;
		gap: 12px;
	}
</style>
