# Pastebin Clone - Node.js Monolithic Application

A monolithic pastebin application built with Node.js, Express, and MongoDB.

## Features

- User authentication (register/login)
- Create, read, update, and delete pastes
- Public, private, and unlisted paste visibility
- Paste expiration
- View count tracking
- User-specific paste management

## Prerequisites

- Node.js (v14 or higher)
- MongoDB (v4.4 or higher)
- npm or yarn package manager

## Installation

1. Clone the repository:
```bash
git clone https://github.com/yourusername/pastebin-mono-nodejs.git
cd pastebin-mono-nodejs
```

2. Install dependencies:
```bash
npm install
```

3. Create a `.env` file in the root directory and configure your environment variables:
```env
NODE_ENV=development
PORT=3000
MONGODB_URI=mongodb://localhost:27017/pastebin
JWT_SECRET=your_jwt_secret_key_here
```

4. Start the development server:
```bash
npm run dev
```

## API Endpoints

### Authentication
- POST `/api/auth/register` - Register a new user
- POST `/api/auth/login` - Login user

### Pastes
- POST `/api/pastes` - Create a new paste
- GET `/api/pastes/public` - Get all public pastes
- GET `/api/pastes/:id` - Get a single paste by ID
- PUT `/api/pastes/:id` - Update a paste
- DELETE `/api/pastes/:id` - Delete a paste
- GET `/api/pastes/user/me` - Get user's pastes

## Security

- Password hashing using bcrypt
- JWT authentication
- Input validation and sanitization
- Protected routes
- Role-based access control

## Development

To run the application in development mode with hot reload:

```bash
npm run dev
```

## Testing

To run tests:

```bash
npm test
```

## Production

For production deployment:

1. Update environment variables in `.env`
2. Build and start the application:
```bash
npm start
```

## License

MIT

## Contributing

1. Fork the repository
2. Create your feature branch
3. Commit your changes
4. Push to the branch
5. Create a new Pull Request 