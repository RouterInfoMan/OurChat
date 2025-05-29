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
	let selected_chat_data:
		| {
				id: number;
				name: string;
		  }
		| boolean = $state(false);
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

	// User search and chat creation states
	let show_new_chat_popover = $state(false);
	let new_chat_creating: boolean | string = $state(false);
	let user_search_query = $state('');
	let search_results: any[] = $state([]);
	let selected_users: any[] = $state([]);
	let searching_users = $state(false);
	let search_error = $state('');
	let chat_type: 'direct' | 'group' = $state('direct');
	let group_name = $state('');

	// Profile management states
	let show_profile_modal = $state(false);
	let current_user: any = $state(null);
	let profile_picture_file: File | null = $state(null);
	let uploading_picture = $state(false);
	let profile_update_message = $state('');

	// Image loading with authorization
    let profileImageCache = {};
    let loadingImages = {};

	let messageInput: HTMLTextAreaElement | null;
	let currentMessageText = $state('');
	let isSending = $state(false);

    // Funcție pentru a încărca imaginea cu autorizare
    async function loadProfileImageWithAuth(imageUrl: string): Promise<string> {
        console.log('loadProfileImageWithAuth called with:', imageUrl); // DEBUG

        if (!imageUrl || !imageUrl.startsWith('/api/')) {
            console.log('Not an API URL, returning as-is'); // DEBUG
            return imageUrl || '/default-avatar.png';
        }

        // Verifică cache-ul
        if (profileImageCache[imageUrl]) {
            return profileImageCache[imageUrl];
        }

        try {
            console.log('Loading image with auth:', imageUrl); // DEBUG

            const response = await fetch(imageUrl, {
                headers: {
                    'Authorization': `Bearer ${localStorage.getItem('jwt_token')}`
                }
            });

            console.log('Image response status:', response.status); // DEBUG

            if (response.ok) {
                const blob = await response.blob();
                const blobUrl = URL.createObjectURL(blob);
                profileImageCache[imageUrl] = blobUrl;
                console.log('Image loaded successfully, blob URL:', blobUrl); // DEBUG
                return blobUrl;
            } else {
                console.error('Failed to load image:', response.status, response.statusText);
            }
        } catch (error) {
            console.error('Error loading profile image:', error);
        }

        profileImageCache[imageUrl] = '/default-avatar.png';
        return '/default-avatar.png';
    }

	// Funcție helper pentru componente
    // Funcție helper pentru componente
    async function getProfileImageUrl(profileUrl: string | null, uniqueId: string): Promise<string> {
        console.log('getProfileImageUrl called with:', profileUrl, uniqueId);

        if (!profileUrl) {
            console.log('No profile URL, returning default');
            return '/default-avatar.png';
        }

        const cacheKey = `${uniqueId}_${profileUrl}_${imageRefreshCounter}`; // Include counter
        console.log('Cache key:', cacheKey);

        if (profileImageCache[cacheKey]) {
            console.log('Found in cache:', profileImageCache[cacheKey]); // DEBUG
            return profileImageCache[cacheKey];
        }

        if (loadingImages[cacheKey]) {
            console.log('Already loading, returning default'); // DEBUG
            return '/default-avatar.png';
        }

        console.log('Starting to load image...'); // DEBUG
        loadingImages[cacheKey] = true;
        const imageUrl = await loadProfileImageWithAuth(profileUrl);
        loadingImages[cacheKey] = false;
        profileImageCache[cacheKey] = imageUrl;

        console.log('Final image URL:', imageUrl); // DEBUG
        return imageUrl;
    }
	onMount(async () => {
		messageInput = document.querySelector('.message-input') as HTMLTextAreaElement | null;

		if (messageInput) {
			messageInput.addEventListener('input', function () {
				this.style.height = 'auto';
				const newHeight = Math.min(this.scrollHeight, 120);
				this.style.height = newHeight + 'px';
			});
		}

		await loadEverything();
		await loadCurrentUser();
	});

	// Load current user profile
	async function loadCurrentUser() {
		try {
			const response = await fetch('/api/profile', {
				method: 'GET',
				headers: {
					Authorization: `Bearer ${localStorage.getItem('jwt_token')}`
				}
			});

			if (response.ok) {
				current_user = await response.json();
				console.log('Current user loaded:', current_user);
                console.log('Profile picture URL:', current_user.profile_picture_url);
			}
		} catch (error) {
			console.error('Error loading user profile:', error);
		}
	}

	// Search users function
	async function searchUsers() {
		if (user_search_query.trim().length < 3) {
			search_results = [];
			search_error = '';
			return;
		}

		try {
			searching_users = true;
			search_error = '';

			const response = await fetch(`/api/users/search?q=${encodeURIComponent(user_search_query.trim())}&limit=10`, {
				method: 'GET',
				headers: {
					Authorization: `Bearer ${localStorage.getItem('jwt_token')}`
				}
			});

			if (!response.ok) {
				throw new Error('Eroare la căutarea utilizatorilor');
			}

			const data = await response.json();
			search_results = data.users || [];
		} catch (error: any) {
			console.error('Error searching users:', error);
			search_error = error.message || 'Eroare la căutarea utilizatorilor';
			search_results = [];
		} finally {
			searching_users = false;
		}
	}

	// Add user to selected list
	function addUserToSelection(user: any) {
		if (chat_type === 'direct' && selected_users.length >= 1) {
			selected_users = [user];
		} else if (!selected_users.find(u => u.id === user.id)) {
			selected_users = [...selected_users, user];
		}
		user_search_query = '';
		search_results = [];
	}

	// Remove user from selected list
	function removeUserFromSelection(userId: number) {
		selected_users = selected_users.filter(u => u.id !== userId);
	}

	// Handle chat type change
	function handleChatTypeChange(type: 'direct' | 'group') {
		chat_type = type;
		if (type === 'direct' && selected_users.length > 1) {
			selected_users = selected_users.slice(0, 1);
		}
		if (type === 'direct') {
			group_name = '';
		}
	}

	// Create chat with selected users
	async function createChat() {
		try {
			new_chat_creating = true;

			if (selected_users.length === 0) {
				throw new Error('Selectează cel puțin un utilizator');
			}

			if (chat_type === 'group' && !group_name.trim()) {
				throw new Error('Introdu un nume pentru grup');
			}

			const user_ids = selected_users.map(u => u.id);

			const response = await fetch('/api/chats', {
				method: 'POST',
				headers: {
					Authorization: `Bearer ${localStorage.getItem('jwt_token')}`,
					'Content-Type': 'application/json'
				},
				body: JSON.stringify({
					type: chat_type,
					users: user_ids,
					name: chat_type === 'group' ? group_name.trim() : undefined
				})
			});

			if (!response.ok) {
				const errorText = await response.text();
				throw new Error(errorText || 'Eroare la crearea chat-ului');
			}

			new_chat_creating = false;
			show_new_chat_popover = false;
			resetChatCreation();
			await loadEverything();
		} catch (error: any) {
			console.error('Error:', error);
			new_chat_creating = error.message || 'A apărut o eroare la crearea chat-ului.';
		}
	}

	// Reset chat creation form
	function resetChatCreation() {
		selected_users = [];
		user_search_query = '';
		search_results = [];
		chat_type = 'direct';
		group_name = '';
		search_error = '';
		new_chat_creating = false;
	}
    let imageRefreshCounter = $state(0);

    // Upload profile picture
    async function uploadProfilePicture() {
        if (!profile_picture_file) return;

        try {
            uploading_picture = true;
            profile_update_message = '';

            const formData = new FormData();
            formData.append('profile_picture', profile_picture_file);

            const response = await fetch('/api/profile/picture', {
                method: 'POST',
                headers: {
                    Authorization: `Bearer ${localStorage.getItem('jwt_token')}`
                },
                body: formData
            });

            if (!response.ok) {
                const errorData = await response.json();
                throw new Error(errorData.message || 'Eroare la încărcarea pozei');
            }

            const data = await response.json();
            profile_update_message = 'Poza de profil a fost actualizată!';

            Object.keys(profileImageCache).forEach(key => {
                if (profileImageCache[key] && profileImageCache[key].startsWith('blob:')) {
                    URL.revokeObjectURL(profileImageCache[key]);
                }
                delete profileImageCache[key];
            });

            Object.keys(loadingImages).forEach(key => {
                delete loadingImages[key];
            });
            // Reload current user
            await loadCurrentUser();

            // Reset file input
            profile_picture_file = null;
            const fileInput = document.getElementById('profile-picture-input') as HTMLInputElement;
            if (fileInput) fileInput.value = '';
            imageRefreshCounter++;

        } catch (error: any) {
            console.error('Error uploading profile picture:', error);
            profile_update_message = error.message || 'Eroare la încărcarea pozei';
        } finally {
            uploading_picture = false;
        }
    }
	// Handle file selection
	function handleFileSelect(event: Event) {
		const target = event.target as HTMLInputElement;
		const file = target.files?.[0];

		if (file) {
			// Validate file type
			if (!file.type.startsWith('image/')) {
				profile_update_message = 'Te rog selectează un fișier imagine';
				return;
			}

			// Validate file size (5MB)
			if (file.size > 5 * 1024 * 1024) {
				profile_update_message = 'Imaginea este prea mare. Mărimea maximă este 5MB';
				return;
			}

			profile_picture_file = file;
			profile_update_message = '';
		}
	}

	// Send message function
	async function sendMessage(): Promise<void> {
		const text = currentMessageText.trim();
		if (!text || isSending) return;

		try {
			isSending = true;

			const tempMessage = {
				id: Date.now(),
				sender_id: current_user?.id || -1,
				chat_id: selected_chat,
				content: text,
				created_at: new Date().toISOString(),
				is_read: 'false'
			};

			if (Array.isArray(chat_messages)) {
				chat_messages = [...chat_messages, tempMessage];
			} else {
				chat_messages = [tempMessage];
			}

			currentMessageText = '';
			if (messageInput) {
				messageInput.value = '';
				messageInput.style.height = 'auto';
			}

			setTimeout(() => {
				const messagesContainer = document.querySelector('.messages-container');
				if (messagesContainer) {
					messagesContainer.scrollTop = messagesContainer.scrollHeight;
				}
			}, 50);

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

			await loadChatMessages();
		} catch (error) {
			console.error('Error sending message:', error);
		} finally {
			isSending = false;
		}
	}

    async function loadEverything() {
        loading = true;
        chats = null;
        error = null;
        selected_chat = null;
        show_new_chat_popover = false;

        try {
            const response = await fetch('/api/chats', {
                method: 'GET',
                headers: {
                    Authorization: `Bearer ${localStorage.getItem('jwt_token')}`
                }
            });
            if (!response.ok) {
                throw new Error('Eroare la obținerea conversațiilor');
            }
            const chatsData = await response.json();

            // Pentru fiecare chat direct, încearcă să obții poza celuilalt utilizator
            for (const chat of chatsData) {
                if (chat.type === 'direct') {
                    try {
                        const membersResponse = await fetch(`/api/chats/${chat.id}/members`, {
                            method: 'GET',
                            headers: {
                                Authorization: `Bearer ${localStorage.getItem('jwt_token')}`
                            }
                        });

                        if (membersResponse.ok) {
                            const members = await membersResponse.json();
                            const otherMember = members.find(member => member.user_id !== current_user?.id);
                            if (otherMember) {
                                chat.other_user_avatar = otherMember.profile_picture_url;
                            }
                        }
                    } catch (error) {
                        console.log('Could not load members for chat', chat.id);
                    }
                }
            }

            chats = chatsData;
        } catch (err) {
            error = err as Error;
        } finally {
            loading = false;
        }
    }

    async function loadChat() {
        selected_chat_data = true;
        loadChatMessages();

        try {
            // Încarcă detaliile chat-ului
            const chatResponse = await fetch(`/api/chats/${selected_chat}`, {
                method: 'GET',
                headers: {
                    Authorization: `Bearer ${localStorage.getItem('jwt_token')}`
                }
            });

            if (!chatResponse.ok) {
                throw new Error('Eroare la obținerea detaliilor conversației');
            }

            const chatData = await chatResponse.json();
            console.log('Chat data loaded:', chatData); // DEBUG

            // Pentru chat-uri direct, încearcă să obții informații despre celălalt utilizator
            if (chatData.type === 'direct') {
                try {
                    const membersResponse = await fetch(`/api/chats/${selected_chat}/members`, {
                        method: 'GET',
                        headers: {
                            Authorization: `Bearer ${localStorage.getItem('jwt_token')}`
                        }
                    });

                    if (membersResponse.ok) {
                        const members = await membersResponse.json();
                        console.log('Chat members:', members); // DEBUG

                        // Găsește celălalt membru (nu pe utilizatorul curent)
                        const otherMember = members.find(member => member.user_id !== current_user?.id);
                        if (otherMember) {
                            chatData.other_user_avatar = otherMember.profile_picture_url;
                            // Folosește username-ul din members în loc de name din chat
                            chatData.display_name = otherMember.username;
                        }
                    }
                } catch (error) {
                    console.log('Could not load chat members:', error);
                }
            } else {
                // Pentru group chats, folosește numele din chat
                chatData.display_name = chatData.name;
            }

            console.log('Final chat data:', chatData); // DEBUG
            selected_chat_data = chatData;
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
	<!-- Sidebar-ul cu iconițe -->
	<div class="sidebar-icons">
		<div class="top-icons">
			<a
				href="#"
				class="icon-btn"
				onclick={() => {
					show_new_chat_popover = true;
					resetChatCreation();
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
					show_profile_modal = true;
				}}
				title="Profile Settings"
			>
				{#await getProfileImageUrl(current_user?.profile_picture_url, `profile_${current_user?.id}`)}
					<img src="/default-avatar.png" alt="Loading..." class="loading" />
				{:then imageUrl}
					<img src={imageUrl} alt="Profil" />
				{/await}
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
                            {#if chat.type === 'direct' && chat.other_user_avatar}
                                {#await getProfileImageUrl(chat.other_user_avatar, `chat_list_${chat.id}`)}
                                    <img src="/default-avatar.png" alt="Loading..." class="loading" />
                                {:then imageUrl}
                                    <img src={imageUrl} alt={chat.name} />
                                {/await}
                            {:else}
                                <img src="/default-avatar.png" alt="Avatar" />
                            {/if}
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
				<div class="chat-area">
                    <div class="chat-header">
                        <div class="current-conversation">
                            <div class="avatar-wrapper">
                                {#if typeof selected_chat_data === 'object' && selected_chat_data.type === 'direct' && selected_chat_data.other_user_avatar}
                                    {#await getProfileImageUrl(selected_chat_data.other_user_avatar, `chat_header_${selected_chat}`)}
                                        <img src="/default-avatar.png" alt="Loading..." class="loading" />
                                    {:then imageUrl}
                                        <img src={imageUrl} alt="Chat Avatar" />
                                    {/await}
                                {:else}
                                    <img src="/default-avatar.png" alt="Avatar" />
                                {/if}
                            </div>
                            {#if typeof selected_chat_data === 'object'}
                                <h2>{selected_chat_data.display_name || selected_chat_data.name || 'Chat'}</h2>
                            {:else}
                                <h2>Loading...</h2>
                            {/if}
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
								<div class={message.sender_id === current_user?.id ? 'message-sent' : 'message-received'}>
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

<!-- New Chat Popover -->
{#if show_new_chat_popover}
	<div class="popover-overlay">
		<div class="popover-content large">
			<h3>Create New Chat</h3>

			<!-- Chat Type Selection -->
			<div class="chat-type-selector">
				<button
					class="type-btn {chat_type === 'direct' ? 'active' : ''}"
					onclick={() => handleChatTypeChange('direct')}
				>
					<svg viewBox="0 0 24 24" width="20" height="20">
						<path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"></path>
						<circle cx="12" cy="7" r="4"></circle>
					</svg>
					Direct Chat
				</button>
				<button
					class="type-btn {chat_type === 'group' ? 'active' : ''}"
					onclick={() => handleChatTypeChange('group')}
				>
					<svg viewBox="0 0 24 24" width="20" height="20">
						<path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"></path>
						<circle cx="9" cy="7" r="4"></circle>
						<path d="M23 21v-2a4 4 0 0 0-3-3.87"></path>
						<path d="M16 3.13a4 4 0 0 1 0 7.75"></path>
					</svg>
					Group Chat
				</button>
			</div>

			<!-- Group Name Input (only for group chats) -->
			{#if chat_type === 'group'}
				<div class="form-group">
					<label>Group Name:</label>
					<input
						type="text"
						bind:value={group_name}
						placeholder="Enter group name..."
						class="group-name-input"
					/>
				</div>
			{/if}

			<!-- User Search -->
			<div class="form-group">
				<label>Search Users {chat_type === 'direct' ? '(select 1)' : '(select multiple)'}:</label>
				<div class="search-container">
					<input
						type="text"
						bind:value={user_search_query}
						placeholder="Type username to search..."
						oninput={searchUsers}
						class="user-search-input"
					/>
					{#if searching_users}
						<div class="search-loading">
							<div class="loading-spinner tiny"></div>
						</div>
					{/if}
				</div>

				{#if search_error}
					<div class="search-error">{search_error}</div>
				{/if}

				<!-- Search Results -->
				{#if search_results.length > 0}
					<div class="search-results">
						{#each search_results as user}
							<div class="user-result" onclick={() => addUserToSelection(user)}>
								<div class="user-avatar">
									{#await getProfileImageUrl(user.profile_picture_url, `search_${user.id}`)}
										<img src="/default-avatar.png" alt="Loading..." class="loading" />
									{:then imageUrl}
										<img src={imageUrl} alt={user.username} />
									{/await}
								</div>
								<div class="user-info">
									<span class="username">{user.username}</span>
									<span class="user-status status-{user.status}">{user.status}</span>
								</div>
								<button class="add-user-btn">
									<svg viewBox="0 0 24 24" width="16" height="16">
										<circle cx="12" cy="12" r="10"></circle>
										<line x1="12" y1="8" x2="12" y2="16"></line>
										<line x1="8" y1="12" x2="16" y2="12"></line>
									</svg>
								</button>
							</div>
						{/each}
					</div>
				{/if}
			</div>

			<!-- Selected Users -->
			{#if selected_users.length > 0}
				<div class="form-group">
					<label>Selected Users:</label>
					<div class="selected-users">
						{#each selected_users as user}
							<div class="selected-user">
								<div class="user-avatar small">
									{#await getProfileImageUrl(user.profile_picture_url, `selected_${user.id}`)}
										<img src="/default-avatar.png" alt="Loading..." class="loading" />
									{:then imageUrl}
										<img src={imageUrl} alt={user.username} />
									{/await}
								</div>
								<span class="username">{user.username}</span>
								<button class="remove-user-btn" onclick={() => removeUserFromSelection(user.id)}>
									<svg viewBox="0 0 24 24" width="14" height="14">
										<line x1="18" y1="6" x2="6" y2="18"></line>
										<line x1="6" y1="6" x2="18" y2="18"></line>
									</svg>
								</button>
							</div>
						{/each}
					</div>
				</div>
			{/if}

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
						resetChatCreation();
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

<!-- Profile Modal -->
{#if show_profile_modal}
	<div class="popover-overlay">
		<div class="popover-content">
			<h3>Profile Settings</h3>

			{#if current_user}
				<div class="profile-info">
					<div class="current-profile-picture">
						{#await getProfileImageUrl(current_user.profile_picture_url, `modal_${current_user.id}`)}
							<img src="/default-avatar.png" alt="Loading..." class="loading" />
						{:then imageUrl}
							<img src={imageUrl} alt="Current Profile" />
						{/await}
					</div>
					<div class="user-details">
						<h4>{current_user.username}</h4>
						<p>{current_user.email}</p>
						<span class="user-status status-{current_user.status}">{current_user.status}</span>
					</div>
				</div>

				<div class="form-group">
					<label for="profile-picture-input">Change Profile Picture:</label>
					<input
						type="file"
						id="profile-picture-input"
						accept="image/*"
						onchange={handleFileSelect}
						class="file-input"
					/>
					{#if profile_picture_file}
						<div class="file-selected">
							<span>Selected: {profile_picture_file.name}</span>
							<button class="upload-btn" onclick={uploadProfilePicture} disabled={uploading_picture}>
								{#if uploading_picture}
									<div class="loading-spinner tiny"></div>
									Uploading...
								{:else}
									Upload
								{/if}
							</button>
						</div>
					{/if}
				</div>

				{#if profile_update_message}
					<div class="profile-message {profile_update_message.includes('Eroare') ? 'error' : 'success'}">
						{profile_update_message}
					</div>
				{/if}
			{/if}

			<div class="popover-actions">
				<button class="secondary-btn" onclick={() => goto('/')}>Go to Dashboard</button>
				<button
					class="secondary-btn"
					onclick={() => {
						show_profile_modal = false;
						profile_update_message = '';
						profile_picture_file = null;
					}}>Close</button
				>
			</div>
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

	/* Sidebar-ul cu iconițe */
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

	.profile-img img.loading,
	.user-avatar img.loading,
	.current-profile-picture img.loading {
		opacity: 0.5;
		filter: blur(1px);
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

	.avatar-wrapper.small {
		width: 30px;
		height: 30px;
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

	/* Welcome screen */
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

	/* Chat area */
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
		max-height: 90vh;
		overflow-y: auto;
		color: white;
		border: 1px solid rgba(255, 255, 255, 0.2);
	}

	.popover-content.large {
		min-width: 500px;
	}

	.popover-content h3 {
		margin: 0 0 25px 0;
		font-size: 24px;
		font-weight: 600;
		text-align: center;
	}

	/* Chat type selector */
	.chat-type-selector {
		display: flex;
		gap: 10px;
		margin-bottom: 20px;
	}

	.type-btn {
		flex: 1;
		padding: 12px;
		border: 2px solid rgba(255, 255, 255, 0.3);
		border-radius: 12px;
		background: transparent;
		color: white;
		cursor: pointer;
		transition: all 0.3s ease;
		display: flex;
		align-items: center;
		justify-content: center;
		gap: 8px;
		font-size: 14px;
		font-weight: 500;
	}

	.type-btn:hover {
		background: rgba(255, 255, 255, 0.1);
	}

	.type-btn.active {
		background: rgba(255, 255, 255, 0.2);
		border-color: rgba(255, 255, 255, 0.5);
	}

	.type-btn svg {
		stroke: currentColor;
		fill: none;
		stroke-width: 2;
	}

	/* Form groups */
	.form-group {
		margin-bottom: 20px;
	}

	.form-group label {
		display: block;
		margin-bottom: 8px;
		font-size: 16px;
		font-weight: 500;
	}

	.form-group input {
		width: 100%;
		padding: 12px 16px;
		border-radius: 10px;
		border: none;
		background: rgba(255, 255, 255, 0.2);
		color: white;
		font-size: 16px;
		outline: none;
		transition: all 0.3s ease;
	}

	.form-group input::placeholder {
		color: rgba(255, 255, 255, 0.7);
	}

	.form-group input:focus {
		background: rgba(255, 255, 255, 0.3);
		box-shadow: 0 0 0 2px rgba(255, 255, 255, 0.3);
	}

	/* Search container */
	.search-container {
		position: relative;
	}

	.search-loading {
		position: absolute;
		right: 12px;
		top: 50%;
		transform: translateY(-50%);
	}

	.search-error {
		color: #ff6b6b;
		font-size: 14px;
		margin-top: 8px;
	}

	/* Search results */
	.search-results {
		max-height: 200px;
		overflow-y: auto;
		background: rgba(255, 255, 255, 0.1);
		border-radius: 10px;
		margin-top: 10px;
	}

	.user-result {
		display: flex;
		align-items: center;
		padding: 12px;
		cursor: pointer;
		transition: all 0.3s ease;
		border-bottom: 1px solid rgba(255, 255, 255, 0.1);
	}

	.user-result:hover {
		background: rgba(255, 255, 255, 0.1);
	}

	.user-result:last-child {
		border-bottom: none;
	}

    .user-avatar {
        width: 40px;
        height: 40px;
        border-radius: 50%;
        overflow: hidden;
        flex-shrink: 0;
        border: 1px solid rgba(255, 255, 255, 0.3);
    }

    .user-avatar img {
        width: 100%;
        height: 100%;
        object-fit: cover; /* Asigură-te că imaginea se scalează corect */
        object-position: center; /* Centrează imaginea */
    }

    /* Pentru avatarele mici din selected users */
    .user-avatar.small {
        width: 32px; /* Mărește de la 30px la 32px */
        height: 32px;
        border-width: 1px;
    }

	.user-info {
		flex: 1;
		margin-left: 12px;
		display: flex;
		flex-direction: column;
		gap: 2px;
	}

	.username {
		font-weight: 500;
		font-size: 14px;
	}

	.user-status {
		font-size: 12px;
		padding: 2px 8px;
		border-radius: 12px;
		text-transform: capitalize;
	}

	.status-online {
		background: rgba(76, 175, 80, 0.3);
		color: #4caf50;
	}

	.status-offline {
		background: rgba(158, 158, 158, 0.3);
		color: #9e9e9e;
	}

	.status-away {
		background: rgba(255, 193, 7, 0.3);
		color: #ffc107;
	}

	.status-busy {
		background: rgba(244, 67, 54, 0.3);
		color: #f44336;
	}

	.add-user-btn {
		width: 32px;
		height: 32px;
		border-radius: 50%;
		border: none;
		background: rgba(255, 255, 255, 0.2);
		color: white;
		cursor: pointer;
		display: flex;
		align-items: center;
		justify-content: center;
		transition: all 0.3s ease;
	}

	.add-user-btn:hover {
		background: rgba(255, 255, 255, 0.3);
		transform: scale(1.1);
	}

	.add-user-btn svg {
		stroke: currentColor;
		fill: none;
		stroke-width: 2;
	}

	/* Selected users */
	.selected-users {
		display: flex;
		flex-wrap: wrap;
		gap: 8px;
		padding: 12px;
		background: rgba(255, 255, 255, 0.1);
		border-radius: 10px;
		min-height: 50px;
	}

	.selected-user {
		display: flex;
		align-items: center;
		gap: 8px;
		background: rgba(255, 255, 255, 0.2);
		padding: 6px 12px;
		border-radius: 20px;
		font-size: 14px;
	}

	.remove-user-btn {
		width: 20px;
		height: 20px;
		border-radius: 50%;
		border: none;
		background: rgba(255, 255, 255, 0.3);
		color: white;
		cursor: pointer;
		display: flex;
		align-items: center;
		justify-content: center;
		transition: all 0.3s ease;
	}

	.remove-user-btn:hover {
		background: rgba(255, 255, 255, 0.5);
	}

	.remove-user-btn svg {
		stroke: currentColor;
		fill: none;
		stroke-width: 2;
	}

	/* Profile modal styles */
	.profile-info {
		display: flex;
		align-items: center;
		gap: 20px;
		margin-bottom: 25px;
		padding: 20px;
		background: rgba(255, 255, 255, 0.1);
		border-radius: 15px;
	}

	.current-profile-picture {
		width: 80px;
		height: 80px;
		border-radius: 50%;
		overflow: hidden;
		border: 3px solid rgba(255, 255, 255, 0.3);
	}

	.current-profile-picture img {
		width: 100%;
		height: 100%;
		object-fit: cover;
	}

	.user-details h4 {
		margin: 0 0 8px 0;
		font-size: 20px;
		font-weight: 600;
	}

	.user-details p {
		margin: 0 0 8px 0;
		opacity: 0.8;
		font-size: 14px;
	}

	.file-input {
		padding: 8px 0;
		font-size: 14px;
	}

	.file-selected {
		display: flex;
		align-items: center;
		justify-content: space-between;
		gap: 12px;
		margin-top: 10px;
		padding: 8px 12px;
		background: rgba(255, 255, 255, 0.1);
		border-radius: 8px;
		font-size: 14px;
	}

	.upload-btn {
		padding: 6px 12px;
		border: none;
		border-radius: 15px;
		background: rgba(255, 255, 255, 0.2);
		color: white;
		cursor: pointer;
		font-size: 12px;
		transition: all 0.3s ease;
		display: flex;
		align-items: center;
		gap: 6px;
	}

	.upload-btn:hover:not(:disabled) {
		background: rgba(255, 255, 255, 0.3);
	}

	.upload-btn:disabled {
		opacity: 0.5;
		cursor: not-allowed;
	}

	.profile-message {
		padding: 10px 15px;
		border-radius: 8px;
		margin-top: 15px;
		font-size: 14px;
		text-align: center;
	}

	.profile-message.error {
		background: rgba(255, 107, 107, 0.2);
		border: 1px solid rgba(255, 107, 107, 0.3);
		color: #ff6b6b;
	}

	.profile-message.success {
		background: rgba(76, 175, 80, 0.2);
		border: 1px solid rgba(76, 175, 80, 0.3);
		color: #4caf50;
	}

	/* Button styles */
	.popover-actions {
		display: flex;
		gap: 12px;
		justify-content: center;
		margin-top: 25px;
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
	.messages-container::-webkit-scrollbar,
	.search-results::-webkit-scrollbar,
	.popover-content::-webkit-scrollbar {
		width: 6px;
	}

	.conversation-list::-webkit-scrollbar-track,
	.messages-container::-webkit-scrollbar-track,
	.search-results::-webkit-scrollbar-track,
	.popover-content::-webkit-scrollbar-track {
		background: rgba(255, 255, 255, 0.1);
	}

	.conversation-list::-webkit-scrollbar-thumb,
	.messages-container::-webkit-scrollbar-thumb,
	.search-results::-webkit-scrollbar-thumb,
	.popover-content::-webkit-scrollbar-thumb {
		background: rgba(255, 255, 255, 0.3);
		border-radius: 3px;
	}

	.conversation-list::-webkit-scrollbar-thumb:hover,
	.messages-container::-webkit-scrollbar-thumb:hover,
	.search-results::-webkit-scrollbar-thumb:hover,
	.popover-content::-webkit-scrollbar-thumb:hover {
		background: rgba(255, 255, 255, 0.5);
	}

	/* Media queries pentru responsivitate */
	@media (max-width: 768px) {
		.sidebar-conversations {
			width: 240px;
		}

		.popover-content {
			min-width: 350px;
			padding: 25px;
		}

		.popover-content.large {
			min-width: 400px;
		}

		.chat-type-selector {
			flex-direction: column;
		}

		.profile-info {
			flex-direction: column;
			text-align: center;
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

		.popover-content {
			min-width: 320px;
			margin: 15px;
		}

		.popover-content.large {
			min-width: 320px;
		}

		.popover-actions {
			flex-direction: column;
			gap: 8px;
		}

		.selected-users {
			flex-direction: column;
			align-items: stretch;
		}

		.selected-user {
			justify-content: space-between;
		}

		.file-selected {
			flex-direction: column;
			align-items: stretch;
			text-align: center;
		}
	}
</style>