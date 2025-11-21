# Achievement System Frontend

Frontend application for the Achievement Management System built with SvelteKit 5, TypeScript, and Tailwind CSS.

## Features

### For Students (Mahasiswa)
- Create and manage personal achievements
- Submit achievements for verification
- Track achievement status (draft, submitted, verified, rejected)
- View statistics and performance metrics

### For Lecturers (Dosen Wali)
- Review submitted achievements
- Verify or reject achievements with notes
- View all achievements under supervision

### For Admins
- Manage all achievements
- Manage users (activate/deactivate)
- View comprehensive statistics
- Monitor system activity

## Tech Stack

- **Framework**: SvelteKit 5 with Runes
- **Language**: TypeScript
- **Styling**: Tailwind CSS
- **State Management**: Svelte Stores
- **API Client**: Custom fetch wrapper with JWT authentication

## Project Structure

```
frontend/
├── src/
│   ├── lib/
│   │   ├── api.ts              # API client
│   │   ├── stores/
│   │   │   └── auth.ts         # Authentication store
│   │   └── assets/
│   ├── routes/
│   │   ├── +layout.svelte      # Main layout with navbar
│   │   ├── +page.svelte        # Landing/redirect page
│   │   ├── dashboard/          # Dashboard page
│   │   ├── login/              # Login & Register page
│   │   ├── achievements/       # Achievement management
│   │   ├── verification/       # Verification page (lecturers)
│   │   ├── users/              # User management (admin)
│   │   └── statistics/         # Statistics & reports
│   ├── app.css                 # Global styles
│   └── app.html                # HTML template
├── package.json
├── svelte.config.js
├── tailwind.config.js
├── tsconfig.json
└── vite.config.ts
```

## Getting Started

### Prerequisites
- Node.js 18+ 
- npm or pnpm

### Installation

1. Install dependencies:
```bash
npm install
```

2. Start the development server:
```bash
npm run dev
```

The application will be available at `http://localhost:5173`

### Building for Production

```bash
npm run build
```

Preview the production build:
```bash
npm run preview
```

## API Integration

The frontend communicates with the backend API at `/api/v1`. All API calls include JWT authentication tokens when available.

### API Endpoints Used

- **Auth**: `/auth/login`, `/auth/register`
- **Achievements**: `/achievements`, `/achievements/:id`
- **Verification**: `/verification/pending`, `/verification/:id/verify`, `/verification/:id/reject`
- **Users**: `/users`, `/users/:id/toggle-status`
- **Reports**: `/reports/statistics`

## Authentication Flow

1. User logs in via `/login`
2. JWT token is stored in localStorage
3. Token is included in all API requests via Authorization header
4. User data is stored in Svelte store
5. Protected routes redirect to login if not authenticated

## Available Routes

- `/` - Landing page (redirects to dashboard or login)
- `/login` - Login & Registration
- `/dashboard` - Main dashboard with statistics
- `/achievements` - Achievement management
- `/verification` - Achievement verification (lecturers & admin)
- `/users` - User management (admin only)
- `/statistics` - Detailed statistics and reports

## Styling

Custom utility classes defined in `app.css`:

- `.btn`, `.btn-primary`, `.btn-secondary`, `.btn-success`, `.btn-danger`
- `.btn-sm` - Small button variant
- `.card` - Card container
- `.input`, `.form-input` - Form inputs
- `.form-label` - Form labels

## Development

### Adding a New Page

1. Create a new folder in `src/routes/`
2. Add `+page.svelte` file
3. Update navigation in `+layout.svelte` if needed

### Adding a New API Endpoint

1. Add the endpoint to `src/lib/api.ts`
2. Use the `fetchApi` helper for automatic auth headers

### State Management

Authentication state is managed via Svelte stores in `src/lib/stores/auth.ts`. Access it in components with:

```typescript
import { authStore } from '$lib/stores/auth';

$: user = $authStore.user;
$: isAuthenticated = $authStore.isAuthenticated;
```

## License

MIT
