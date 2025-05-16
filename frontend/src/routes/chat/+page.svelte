<script lang="ts">
	import { onMount } from 'svelte';

	let loading = $state(true);
	let chats = $state(null);
	let error: Error | null = $state(null);

	onMount(async () => {
		// Selectează elementele DOM cu type assertions
		const messageInput = document.querySelector('.message-input') as HTMLTextAreaElement | null;
		const sendButton = document.querySelector('.send-btn') as HTMLButtonElement | null;
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

		// Funcția pentru a adăuga un mesaj nou cu receipt
		function addMessage(text: string, type: 'sent' | 'received'): void {
			if (!messagesContainer) return;

			const messageDiv = document.createElement('div');
			messageDiv.className = 'message-' + type;

			const messageBubble = document.createElement('div');
			messageBubble.className = 'message-bubble';
			messageBubble.textContent = text;

			messageDiv.appendChild(messageBubble);

			// Adaugă informații despre mesaj (ora și read receipt) doar pentru mesajele trimise
			if (type === 'sent') {
				const messageInfo = document.createElement('div');
				messageInfo.className = 'message-info';

				// Adaugă ora curentă
				const currentTime = new Date();
				const timeStr =
					currentTime.getHours() +
					':' +
					(currentTime.getMinutes() < 10 ? '0' : '') +
					currentTime.getMinutes();

				const messageTime = document.createElement('span');
				messageTime.className = 'message-time';
				messageTime.textContent = timeStr;

				// Adaugă read receipt inițial ca "delivered"
				const readReceipt = document.createElement('span');
				readReceipt.className = 'read-receipt delivered';
				readReceipt.textContent = 'Livrat';

				messageInfo.appendChild(messageTime);
				messageInfo.appendChild(readReceipt);
				messageDiv.appendChild(messageInfo);

				// Simulează schimbarea statusului la "seen" după un timp
				setTimeout(() => {
					readReceipt.className = 'read-receipt seen';
					readReceipt.textContent = 'Văzut';
				}, 3000);
			}

			messagesContainer.appendChild(messageDiv);

			// Scroll la ultimul mesaj
			messagesContainer.scrollTop = messagesContainer.scrollHeight;
		}

		// Funcția pentru a trimite un mesaj
		function sendMessage(): void {
			if (!messageInput || !messagesContainer) return;

			const text = messageInput.value.trim();

			if (text) {
				// Adaugă mesajul utilizatorului
				addMessage(text, 'sent');

				// Golește input-ul și resetează înălțimea
				messageInput.value = '';
				messageInput.style.height = 'auto';

				// Simulare răspuns automat
				setTimeout(() => {
					addMessage('Această funcționalitate va fi implementată complet în curând.', 'received');
				}, 1000);
			}
		}

		// Event listeners
		if (sendButton) {
			sendButton.addEventListener('click', sendMessage);
		}

		if (messageInput) {
			messageInput.addEventListener('keypress', function (e: KeyboardEvent) {
				if (e.key === 'Enter' && !e.shiftKey) {
					sendMessage();
					e.preventDefault();
				}
			});
		}

		// Gestionare selectare conversații
		conversationItems.forEach((item) => {
			item.addEventListener('click', function (this: Element) {
				if (!messagesContainer) return;

				// Elimină clasa active de la toate conversațiile
				conversationItems.forEach((conv) => {
					conv.classList.remove('active');
				});

				// Adaugă clasa active la conversația selectată
				this.classList.add('active');

				// Actualizează informațiile din header
				const nameElement = this.querySelector('.conv-details h3');
				const headerName = document.querySelector('.current-conversation h2');

				if (nameElement && headerName) {
					const name = nameElement.textContent || '';
					headerName.textContent = name;

					// Actualizează avatarul din header
					const avatarImg = this.querySelector('.avatar-wrapper img') as HTMLImageElement | null;
					const headerAvatar = document.querySelector(
						'.current-conversation .avatar-wrapper img'
					) as HTMLImageElement | null;

					if (avatarImg && headerAvatar && avatarImg.getAttribute('src')) {
						headerAvatar.setAttribute('src', avatarImg.getAttribute('src') || '');
					}

					// Curăță mesageria și adaugă statusul
					messagesContainer.innerHTML = '';

					const statusMessage = document.createElement('div');
					statusMessage.className = 'status-message';
					const statusParagraph = document.createElement('p');
					statusParagraph.textContent = `Ai selectat conversația: ${name}`;
					statusMessage.appendChild(statusParagraph);
					messagesContainer.appendChild(statusMessage);
				}
			});
		});

		// Funcție pentru integrarea cu backend-ul de read receipts
		// Aceasta va fi implementată când vei conecta cu backend-ul real
		function updateMessageStatus(messageId: string, status: 'delivered' | 'seen'): void {
			// Găsește mesajul după ID
			const messageElement = document.querySelector(`[data-message-id="${messageId}"]`);
			if (!messageElement) return;

			// Găsește read receipt-ul
			const readReceipt = messageElement.querySelector('.read-receipt');
			if (!readReceipt) return;

			// Actualizează statusul
			if (status === 'delivered') {
				readReceipt.className = 'read-receipt delivered';
				readReceipt.textContent = 'Livrat';
			} else if (status === 'seen') {
				readReceipt.className = 'read-receipt seen';
				readReceipt.textContent = 'Văzut';
			}
		}

		try {
			let req = await fetch('/api/chats', {
				method: 'GET',
				headers: {
					'Authorization': `Bearer ${localStorage.getItem('jwt_token')}`
				}
			});
			if (!req.ok) {
				throw new Error('Eroare la obținerea conversațiilor');
			}
			chats = await req.json();
		}
		catch (err) {
			error = err as Error;
		} finally {
			loading = false;
		}
	});
</script>

<div class="chat-layout">
	<!-- Sidebar-ul albastru îngust -->
	<div class="sidebar-icons">
		<div class="top-icons">
			<a href="#" class="icon-btn">
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
			<div class="profile-img">
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
				<div class="conversation-item active">
					<div class="avatar-wrapper">
						<img src="/default-avatar.png" alt="Avatar" />
					</div>
					<div class="conv-details">
						<h3>Conversație</h3>
						<p>Mesaj recent</p>
					</div>
				</div>
				<div class="conversation-item">
					<div class="avatar-wrapper">
						<img src="/default-avatar.png" alt="Avatar" />
					</div>
					<div class="conv-details">
						<h3>Altă conversație</h3>
						<p>Mesaj recent</p>
					</div>
				</div>
			</div>
		</div>

		<!-- Zona principală de chat -->
		<div class="chat-area">
			<div class="chat-header">
				<div class="current-conversation">
					<div class="avatar-wrapper">
						<img src="/default-avatar.png" alt="Avatar" />
					</div>
					<h2>Conversație</h2>
				</div>
			</div>

			<div class="messages-container">
				<div class="status-message">
					<p>Ai selectat conversația: Conversație</p>
				</div>

				<!-- Exemplu de mesaj cu receipt -->
				<div class="message-sent">
					<div class="message-bubble">Salut! Cum merge proiectul?</div>
					<div class="message-info">
						<span class="message-time">14:25</span>
						<span class="read-receipt seen">Văzut</span>
					</div>
				</div>
			</div>

			<div class="message-input-area">
				<button class="attachment-btn">
					<svg viewBox="0 0 24 24" width="24" height="24"
						><path
							d="M21.44 11.05l-9.19 9.19a6 6 0 0 1-8.49-8.49l9.19-9.19a4 4 0 0 1 5.66 5.66l-9.2 9.19a2 2 0 0 1-2.83-2.83l8.49-8.48"
						></path></svg
					>
				</button>
				<textarea placeholder="Tastează un mesaj..." class="message-input" rows="1"></textarea>
				<button class="send-btn">
					<svg viewBox="0 0 24 24" width="24" height="24"
						><polygon points="22 2 15 22 11 13 2 9 22 2"></polygon></svg
					>
				</button>
			</div>
		</div>
	{:else}
		eroare
	{/if}
</div>

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
</style>
