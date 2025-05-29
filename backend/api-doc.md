# OurChat API Documentation

## Table of Contents
- [Authentication](#authentication)
  - [Register](#register)
  - [Login](#login)
  - [Logout](#logout)
  - [Request Password Reset](#request-password-reset)
  - [Reset Password](#reset-password)
- [User](#user)
  - [Get Profile](#get-profile)
  - [Update Profile](#update-profile)
  - [Upload Profile Picture](#upload-profile-picture)
  - [Get Users by IDs](#get-users-by-ids)
  - [Search Users](#search-users)
- [Media](#media)
  - [Upload Media File](#upload-media-file)
  - [Get Media File](#get-media-file)
- [Chats](#chats)
  - [Get Chats](#get-chats)
  - [Create Chat](#create-chat)
  - [Get Chat](#get-chat)
  - [Get Chat Members](#get-chat-members)
- [Messages](#messages)
  - [Get Messages](#get-messages)
  - [Send Text Message](#send-text-message)
  - [Send Media Message](#send-media-message)
  - [Mark Messages as Read](#mark-messages-as-read)
  - [Search Messages](#search-messages)

## Authentication

All protected endpoints require a JWT token in the Authorization header:
```
Authorization: Bearer <token>
```

### Register

Register a new user account.

**URL**: `/api/register`
**Method**: `POST`
**Auth required**: No

**Request Body**:
```json
{
  "username": "testuser",
  "email": "test@example.com",
  "password": "password123"
}
```

**Success Response**:
- **Code**: 201 Created
- **Content**:
```json
{
  "user_id": 1,
  "message": "User registered successfully"
}
```

**Error Responses**:
- **Code**: 400 Bad Request (Missing required fields, username too short)
- **Code**: 409 Conflict (Username or email already exists)
- **Code**: 500 Internal Server Error

### Login

Login to get an authentication token.

**URL**: `/api/login`
**Method**: `POST`
**Auth required**: No

**Request Body**:
```json
{
  "username": "testuser",
  "password": "password123"
}
```

**Success Response**:
- **Code**: 200 OK
- **Content**:
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user_id": 1,
  "message": "Login successful"
}
```

**Error Responses**:
- **Code**: 400 Bad Request (Invalid request)
- **Code**: 401 Unauthorized (Invalid username or password)
- **Code**: 500 Internal Server Error

### Logout

Logout and invalidate the current session.

**URL**: `/api/logout`
**Method**: `POST`
**Auth required**: Yes

**Success Response**:
- **Code**: 200 OK
- **Content**:
```json
{
  "message": "Logout successful"
}
```

**Error Responses**:
- **Code**: 401 Unauthorized (Invalid or missing token)
- **Code**: 500 Internal Server Error

### Request Password Reset

Initiates the password reset process for a user.

**URL**: `/api/request-password-reset`
**Method**: `POST`
**Auth required**: No

**Request Body**:
```json
{
  "email": "user@example.com"
}
```

**Success Response**:
- **Code**: 200 OK
- **Content**:
```json
{
  "message": "If your email is registered, you will receive a password reset link shortly",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Error Responses**:
- **Code**: 400 Bad Request (Invalid request or missing email)
- **Code**: 500 Internal Server Error

### Reset Password

Resets a user's password using a valid reset token.

**URL**: `/api/reset-password`
**Method**: `POST`
**Auth required**: No

**Request Body**:
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "new_password": "newSecurePassword123"
}
```

**Success Response**:
- **Code**: 200 OK
- **Content**:
```json
{
  "message": "Password has been reset successfully. Please login with your new password."
}
```

**Error Responses**:
- **Code**: 400 Bad Request (Invalid request, missing fields, or invalid/expired token)
- **Code**: 500 Internal Server Error

## User

### Get Profile

Get the current user's profile information.

**URL**: `/api/profile`
**Method**: `GET`
**Auth required**: Yes

**Success Response**:
- **Code**: 200 OK
- **Content**:
```json
{
  "id": 1,
  "username": "testuser",
  "email": "test@example.com",
  "profile_picture_url": "/api/media/profiles/abc123def456.jpg",
  "status": "online",
  "created_at": "2025-05-15T10:30:45Z",
  "last_login": "2025-05-15T15:20:30Z"
}
```

**Error Responses**:
- **Code**: 401 Unauthorized (Invalid or missing token)
- **Code**: 500 Internal Server Error

### Update Profile

Update the current user's profile information.

**URL**: `/api/profile`
**Method**: `PUT`
**Auth required**: Yes

**Request Body**:
```json
{
  "email": "newemail@example.com",
  "status": "away"
}
```

**Success Response**:
- **Code**: 200 OK
- **Content**: Updated profile information
```json
{
  "id": 1,
  "username": "testuser",
  "email": "newemail@example.com",
  "profile_picture_url": "/api/media/profiles/abc123def456.jpg",
  "status": "away",
  "created_at": "2025-05-15T10:30:45Z",
  "last_login": "2025-05-15T15:20:30Z"
}
```

**Error Responses**:
- **Code**: 400 Bad Request (Invalid request, invalid status, no valid fields)
- **Code**: 401 Unauthorized (Invalid or missing token)
- **Code**: 409 Conflict (Email already in use)
- **Code**: 500 Internal Server Error

### Upload Profile Picture

Upload or update the current user's profile picture.

**URL**: `/api/profile/picture`
**Method**: `POST`
**Auth required**: Yes
**Content-Type**: `multipart/form-data`

**Request Body** (Form Data):
- `profile_picture`: Image file (JPEG, PNG, or GIF, max 5MB)

**Success Response**:
- **Code**: 200 OK
- **Content**:
```json
{
  "message": "Profile picture updated successfully",
  "url": "/api/media/profiles/abc123def456.jpg"
}
```

**Error Responses**:
- **Code**: 400 Bad Request (No file provided, invalid file type, file too large)
- **Code**: 401 Unauthorized (Invalid or missing token)
- **Code**: 500 Internal Server Error

### Get Users by IDs

Get basic information for a list of users by their IDs.

**URL**: `/api/users`
**Method**: `GET` or `POST`
**Auth required**: Yes

**GET Method**:
```
GET /api/users?ids=1,2,3
```

**POST Method - Request Body**:
```json
{
  "user_ids": [1, 2, 3]
}
```

**Success Response**:
- **Code**: 200 OK
- **Content**:
```json
{
  "1": {
    "id": 1,
    "username": "testuser1",
    "status": "online",
    "profile_picture_url": "/api/media/profiles/user1_profile.jpg"
  },
  "2": {
    "id": 2,
    "username": "testuser2",
    "status": "offline",
    "profile_picture_url": null
  }
}
```

**Error Responses**:
- **Code**: 400 Bad Request (Invalid request or missing IDs)
- **Code**: 401 Unauthorized (Invalid or missing token)
- **Code**: 500 Internal Server Error

### Search Users

Search for users by partial username match.

**URL**: `/api/users/search`
**Method**: `GET`
**Auth required**: Yes

**Query Parameters**:
- `q`: Search term (required, minimum 4 characters)
- `limit`: Maximum number of results (optional, default: 20, max: 50)

**Success Response**:
- **Code**: 200 OK
- **Content**:
```json
{
  "users": [
    {
      "id": 1,
      "username": "john_doe",
      "status": "online",
      "profile_picture_url": "/api/media/profiles/john_profile.jpg"
    },
    {
      "id": 5,
      "username": "johnny123",
      "status": "offline",
      "profile_picture_url": null
    }
  ],
  "count": 2,
  "query": "john"
}
```

**Error Responses**:
- **Code**: 400 Bad Request (Missing search query, query too short)
- **Code**: 401 Unauthorized (Invalid or missing token)
- **Code**: 500 Internal Server Error

## Media

### Upload Media File

Upload a media file for use in messages.

**URL**: `/api/media/upload`
**Method**: `POST`
**Auth required**: Yes
**Content-Type**: `multipart/form-data`

**Request Body** (Form Data):
- `media`: Media file (images, videos, audio, PDFs, max 50MB)

**Supported File Types**:
- Images: JPEG, JPG, PNG, GIF, WebP
- Videos: MP4, WebM, AVI, MOV
- Audio: MP3, WAV, OGG, MPEG
- Documents: PDF, TXT

**Success Response**:
- **Code**: 201 Created
- **Content**:
```json
{
  "id": 123,
  "filename": "abc123def456.jpg",
  "original_filename": "vacation_photo.jpg",
  "file_size": 1024567,
  "mime_type": "image/jpeg",
  "uploaded_by": 1,
  "uploaded_at": "2025-05-28T15:30:45Z",
  "url": "/api/media/files/abc123def456.jpg"
}
```

**Error Responses**:
- **Code**: 400 Bad Request (No file provided, invalid file type, file too large)
- **Code**: 401 Unauthorized (Invalid or missing token)
- **Code**: 500 Internal Server Error

### Get Media File

Retrieve a media file or profile picture.

**URL**: `/api/media/{type}/{filename}`
**Method**: `GET`
**Auth required**: Yes

**URL Parameters**:
- `type`: File type (`files` for media files, `profiles` for profile pictures)
- `filename`: The filename of the media file

**Success Response**:
- **Code**: 200 OK
- **Content**: Raw file data with appropriate Content-Type header

**Error Responses**:
- **Code**: 400 Bad Request (Invalid request)
- **Code**: 401 Unauthorized (Invalid or missing token)
- **Code**: 403 Forbidden (Access denied - not authorized to view this file)
- **Code**: 404 Not Found (File not found)
- **Code**: 500 Internal Server Error

**Access Control**:
- Users can access their own uploaded files
- Users can access media files shared in chats they're members of
- Users can access profile pictures of people they share chats with

## Chats

### Get Chats

Get all chats for the current user.

**URL**: `/api/chats`
**Method**: `GET`
**Auth required**: Yes

**Success Response**:
- **Code**: 200 OK
- **Content**:
```json
[
  {
    "id": 1,
    "type": "direct",
    "name": "jane_doe",
    "created_at": "2025-05-15T10:20:30Z",
    "updated_at": "2025-05-15T16:45:20Z",
    "is_active": true
  },
  {
    "id": 2,
    "type": "group",
    "name": "Project Team",
    "created_at": "2025-05-15T11:20:30Z",
    "updated_at": "2025-05-15T15:30:45Z",
    "is_active": true
  }
]
```

**Note**: For direct chats, the `name` field contains the other user's username.

**Error Responses**:
- **Code**: 401 Unauthorized (Invalid or missing token)
- **Code**: 500 Internal Server Error

### Create Chat

Create a new direct or group chat.

**URL**: `/api/chats`
**Method**: `POST`
**Auth required**: Yes

**Request Body - Direct Chat**:
```json
{
  "type": "direct",
  "users": [2]
}
```

**Request Body - Group Chat**:
```json
{
  "type": "group",
  "name": "Project Team",
  "users": [2, 3, 4]
}
```

**Success Response - Direct Chat**:
- **Code**: 200 OK (if chat already exists)
- **Content**: Existing chat object

**Success Response - Group Chat**:
- **Code**: 201 Created
- **Content**: New chat object
```json
{
  "id": 2,
  "type": "group",
  "name": "Project Team",
  "created_at": "2025-05-15T11:20:30Z",
  "updated_at": "2025-05-15T11:20:30Z",
  "is_active": true
}
```

**Error Responses**:
- **Code**: 400 Bad Request (Invalid chat type, missing required fields, invalid user count for direct chat)
- **Code**: 401 Unauthorized (Invalid or missing token)
- **Code**: 500 Internal Server Error

### Get Chat

Get details about a specific chat.

**URL**: `/api/chats/{chatID}`
**Method**: `GET`
**Auth required**: Yes

**URL Parameters**:
- `chatID`: ID of the chat to get

**Success Response**:
- **Code**: 200 OK
- **Content**:
```json
{
  "id": 1,
  "type": "direct",
  "name": "jane_doe",
  "created_at": "2025-05-15T10:20:30Z",
  "updated_at": "2025-05-15T16:45:20Z",
  "is_active": true
}
```

**Error Responses**:
- **Code**: 400 Bad Request (Invalid chat ID)
- **Code**: 401 Unauthorized (Invalid or missing token)
- **Code**: 403 Forbidden (Not a member of this chat)
- **Code**: 500 Internal Server Error

### Get Chat Members

Get all members of a specific chat.

**URL**: `/api/chats/{chatID}/members`
**Method**: `GET`
**Auth required**: Yes

**URL Parameters**:
- `chatID`: ID of the chat to get members for

**Success Response**:
- **Code**: 200 OK
- **Content**:
```json
[
  {
    "id": 1,
    "user_id": 1,
    "chat_id": 1,
    "role": "admin",
    "joined_at": "2025-05-15T10:20:30Z",
    "last_read_at": "2025-05-15T12:30:45Z",
    "username": "testuser1",
    "status": "online",
    "profile_picture_url": null
  },
  {
    "id": 2,
    "user_id": 2,
    "chat_id": 1,
    "role": "member",
    "joined_at": "2025-05-15T10:20:30Z",
    "last_read_at": "2025-05-15T11:45:20Z",
    "username": "testuser2",
    "status": "offline",
    "profile_picture_url": null
  }
]
```

**Error Responses**:
- **Code**: 400 Bad Request (Invalid chat ID)
- **Code**: 401 Unauthorized (Invalid or missing token)
- **Code**: 403 Forbidden (Not a member of this chat)
- **Code**: 500 Internal Server Error

## Messages

### Get Messages

Get messages from a specific chat with pagination.

**URL**: `/api/chats/{chatID}/messages`
**Method**: `GET`
**Auth required**: Yes

**URL Parameters**:
- `chatID`: ID of the chat to get messages from

**Query Parameters**:
- `limit`: Maximum number of messages to retrieve (default: 50)
- `offset`: Offset for pagination (default: 0)

**Success Response**:
- **Code**: 200 OK
- **Content**:
```json
[
  {
    "id": 6,
    "sender_id": 2,
    "chat_id": 1,
    "content": "Check out this photo!",
    "message_type": "media",
    "media_file_id": 123,
    "media_file": {
      "id": 123,
      "filename": "abc123def456.jpg",
      "original_filename": "sunset.jpg",
      "file_size": 1024567,
      "mime_type": "image/jpeg",
      "uploaded_by": 2,
      "uploaded_at": "2025-05-15T11:15:20Z",
      "url": "/api/media/files/abc123def456.jpg"
    },
    "created_at": "2025-05-15T11:20:30Z",
    "is_read": true
  },
  {
    "id": 5,
    "sender_id": 1,
    "chat_id": 1,
    "content": "Hello, how are you?",
    "message_type": "text",
    "media_file_id": null,
    "media_file": null,
    "created_at": "2025-05-15T11:15:20Z",
    "is_read": true
  }
]
```

**Error Responses**:
- **Code**: 400 Bad Request (Invalid chat ID)
- **Code**: 401 Unauthorized (Invalid or missing token)
- **Code**: 403 Forbidden (Not a member of this chat)
- **Code**: 500 Internal Server Error

### Send Text Message

Send a text message to a chat.

**URL**: `/api/chats/{chatID}/messages`
**Method**: `POST`
**Auth required**: Yes

**URL Parameters**:
- `chatID`: ID of the chat to send a message to

**Request Body**:
```json
{
  "message_type": "text",
  "content": "Hello, this is a test message!"
}
```

**Success Response**:
- **Code**: 201 Created
- **Content**:
```json
{
  "id": 7,
  "sender_id": 1,
  "chat_id": 1,
  "content": "Hello, this is a test message!",
  "message_type": "text",
  "media_file_id": null,
  "media_file": null,
  "created_at": "2025-05-15T12:30:45Z",
  "is_read": false
}
```

**Error Responses**:
- **Code**: 400 Bad Request (Invalid chat ID, empty message content)
- **Code**: 401 Unauthorized (Invalid or missing token)
- **Code**: 403 Forbidden (Not a member of this chat)
- **Code**: 500 Internal Server Error

### Send Message with Media

Send a message with a previously uploaded media file.

**URL**: `/api/chats/{chatID}/messages`
**Method**: `POST`
**Auth required**: Yes

**URL Parameters**:
- `chatID`: ID of the chat to send a message to

**Request Body**:
```json
{
  "message_type": "media",
  "media_file_id": 123,
  "content": "Check out this awesome photo!"
}
```

**Success Response**:
- **Code**: 201 Created
- **Content**:
```json
{
  "id": 8,
  "sender_id": 1,
  "chat_id": 1,
  "content": "Check out this awesome photo!",
  "message_type": "media",
  "media_file_id": 123,
  "media_file": {
    "id": 123,
    "filename": "abc123def456.jpg",
    "original_filename": "vacation.jpg",
    "file_size": 1024567,
    "mime_type": "image/jpeg",
    "uploaded_by": 1,
    "uploaded_at": "2025-05-15T12:25:30Z",
    "url": "/api/media/files/abc123def456.jpg"
  },
  "created_at": "2025-05-15T12:30:45Z",
  "is_read": false
}
```

**Error Responses**:
- **Code**: 400 Bad Request (Invalid chat ID, invalid or missing media_file_id)
- **Code**: 401 Unauthorized (Invalid or missing token)
- **Code**: 403 Forbidden (Not a member of this chat, or media file doesn't belong to user)
- **Code**: 500 Internal Server Error

### Send Media Message

Upload a media file and send it as a message in one request.

**URL**: `/api/chats/{chatID}/messages/media`
**Method**: `POST`
**Auth required**: Yes
**Content-Type**: `multipart/form-data`

**URL Parameters**:
- `chatID`: ID of the chat to send a media message to

**Request Body** (Form Data):
- `media`: Media file (images, videos, audio, PDFs, max 50MB)
- `caption`: Optional text caption for the media (optional)

**Success Response**:
- **Code**: 201 Created
- **Content**:
```json
{
  "id": 9,
  "sender_id": 1,
  "chat_id": 1,
  "content": "Beautiful sunset!",
  "message_type": "media",
  "media_file_id": 124,
  "media_file": {
    "id": 124,
    "filename": "def789ghi012.jpg",
    "original_filename": "sunset.jpg",
    "file_size": 2048567,
    "mime_type": "image/jpeg",
    "uploaded_by": 1,
    "uploaded_at": "2025-05-15T12:35:30Z",
    "url": "/api/media/files/def789ghi012.jpg"
  },
  "created_at": "2025-05-15T12:35:45Z",
  "is_read": false
}
```

**Error Responses**:
- **Code**: 400 Bad Request (Invalid chat ID, no media file, invalid file type, file too large)
- **Code**: 401 Unauthorized (Invalid or missing token)
- **Code**: 403 Forbidden (Not a member of this chat)
- **Code**: 500 Internal Server Error

### Mark Messages as Read

Mark all messages in a chat as read for the current user.

**URL**: `/api/chats/{chatID}/messages/read`
**Method**: `POST`
**Auth required**: Yes

**URL Parameters**:
- `chatID`: ID of the chat to mark messages as read

**Success Response**:
- **Code**: 200 OK
- **Content**:
```json
{
  "message": "Messages marked as read"
}
```

**Error Responses**:
- **Code**: 400 Bad Request (Invalid chat ID)
- **Code**: 401 Unauthorized (Invalid or missing token)
- **Code**: 500 Internal Server Error

### Search Messages

Search for messages in a chat containing specific text.

**URL**: `/api/chats/{chatID}/messages/search`
**Method**: `GET`
**Auth required**: Yes

**URL Parameters**:
- `chatID`: ID of the chat to search messages in

**Query Parameters**:
- `q`: Search query text (required)

**Success Response**:
- **Code**: 200 OK
- **Content**:
```json
[
  {
    "id": 7,
    "sender_id": 1,
    "chat_id": 1,
    "content": "Hello, this is a test message!",
    "message_type": "text",
    "media_file_id": null,
    "media_file": null,
    "created_at": "2025-05-15T12:30:45Z",
    "is_read": true
  },
  {
    "id": 3,
    "sender_id": 2,
    "chat_id": 1,
    "content": "Let's do a test call tomorrow",
    "message_type": "text",
    "media_file_id": null,
    "media_file": null,
    "created_at": "2025-05-15T10:45:30Z",
    "is_read": true
  }
]
```

**Error Responses**:
- **Code**: 400 Bad Request (Invalid chat ID, missing search query)
- **Code**: 401 Unauthorized (Invalid or missing token)
- **Code**: 403 Forbidden (Not a member of this chat)
- **Code**: 500 Internal Server Error

## Common HTTP Status Codes

- **200 OK**: Request successful
- **201 Created**: Resource created successfully
- **400 Bad Request**: Invalid request data
- **401 Unauthorized**: Authentication required or failed
- **403 Forbidden**: Access denied (insufficient permissions)
- **404 Not Found**: Resource not found
- **409 Conflict**: Resource already exists (duplicate)
- **500 Internal Server Error**: Server error

## File Upload Limits

- **Profile Pictures**: 5MB maximum, JPEG/PNG/GIF only
- **Media Files**: 50MB maximum, supports images, videos, audio, and documents

## Authentication Notes

- JWT tokens expire after 24 hours
- Include the token in the Authorization header: `Authorization: Bearer <token>`
- Tokens are invalidated on password reset and can be invalidated on logout (depending on implementation)