# OurChat API Documentation

## Table of Contents
- [Authentication](#authentication)
  - [Register](#register)
  - [Login](#login)
  - [Logout](#logout)
- [User](#user)
  - [Get Profile](#get-profile)
  - [Update Profile](#update-profile)
  - [Get Users by IDs](#get-users-by-ids)
- [Chats](#chats)
  - [Get Chats](#get-chats)
  - [Create Chat](#create-chat)
  - [Get Chat](#get-chat)
  - [Get Chat Members](#get-chat-members)
- [Messages](#messages)
  - [Get Messages](#get-messages)
  - [Send Message](#send-message)
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
- **Code**: 400 Bad Request (Missing required fields)
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
    "status": "online"
  },
  "2": {
    "id": 2,
    "username": "testuser2",
    "status": "offline"
  },
  "3": {
    "id": 3,
    "username": "testuser3",
    "status": "away"
  }
}
```

**Error Responses**:
- **Code**: 400 Bad Request (Invalid request or missing IDs)
- **Code**: 401 Unauthorized (Invalid or missing token)
- **Code**: 500 Internal Server Error

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
    "name": "",
    "created_at": "2025-05-15T10:20:30Z",
    "updated_at": "2025-05-15T10:20:30Z",
    "is_active": true
  },
  {
    "id": 2,
    "type": "group",
    "name": "Project Team",
    "created_at": "2025-05-15T11:20:30Z",
    "updated_at": "2025-05-15T11:20:30Z",
    "is_active": true
  }
]
```

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
- **Code**: 200 OK
- **Content**: Chat object

**Success Response - Group Chat**:
- **Code**: 201 Created
- **Content**: Chat object
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
- **Code**: 400 Bad Request (Invalid chat type, missing required fields)
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
  "name": "",
  "created_at": "2025-05-15T10:20:30Z",
  "updated_at": "2025-05-15T10:20:30Z",
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
    "status": "online"
  },
  {
    "id": 2,
    "user_id": 2,
    "chat_id": 1,
    "role": "member",
    "joined_at": "2025-05-15T10:20:30Z",
    "last_read_at": "2025-05-15T11:45:20Z",
    "username": "testuser2",
    "status": "offline"
  }
]
```

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
    "id": 5,
    "sender_id": 2,
    "chat_id": 1,
    "content": "Hello, how are you?",
    "created_at": "2025-05-15T11:20:30Z",
    "is_read": true
  },
  {
    "id": 4,
    "sender_id": 1,
    "chat_id": 1,
    "content": "Hi there!",
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

### Send Message

Send a new message to a chat.

**URL**: `/api/chats/{chatID}/messages`
**Method**: `POST`
**Auth required**: Yes

**URL Parameters**:
- `chatID`: ID of the chat to send a message to

**Request Body**:
```json
{
  "content": "Hello, this is a test message!"
}
```

**Success Response**:
- **Code**: 201 Created
- **Content**:
```json
{
  "id": 6,
  "sender_id": 1,
  "chat_id": 1,
  "content": "Hello, this is a test message!",
  "created_at": "2025-05-15T12:30:45Z",
  "is_read": false
}
```

**Note**: *Currently lacks membership verification. User may be able to send messages to chats they aren't members of.*

**Error Responses**:
- **Code**: 400 Bad Request (Invalid chat ID, empty message)
- **Code**: 401 Unauthorized (Invalid or missing token)
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

**Note**: *Currently lacks membership verification. User may be able to mark messages as read in chats they aren't members of.*

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
    "id": 6,
    "sender_id": 1,
    "chat_id": 1,
    "content": "Hello, this is a test message!",
    "created_at": "2025-05-15T12:30:45Z",
    "is_read": true
  },
  {
    "id": 3,
    "sender_id": 2,
    "chat_id": 1,
    "content": "Let's do a test call tomorrow",
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
