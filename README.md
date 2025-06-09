# Quizgenie Monorepo

This monorepo contains both the **backend** (Go REST API + worker) and **frontend** (Angular SPA) for Quizgenie—a platform where users upload documents, have AI generate quizzes, and track their quiz history.

---

## Table of Contents

1. [Getting Started](#getting-started)
2. [Architecture Overview](#architecture-overview)
3. [Backend](#backend)

   * [Directory Structure](#backend-directory-structure)
   * [Key Components](#backend-key-components)
   * [Running Locally](#backend-running-locally)
4. [Frontend](#frontend)

   * [Directory Structure](#frontend-directory-structure)
   * [Key Components](#frontend-key-components)
   * [Building & Running](#frontend-building--running)
5. [Docker & Docker Compose](#docker--docker-compose)
6. [Environment Variables](#environment-variables)
7. [Component Workflows](#component-workflows)
8. [Database Schema](#database-schema)
9. [UI Development Roadmap](#ui-development-roadmap)
10. [Future Work](#future-work)

---

## Getting Started

Clone the repo and copy `.env.template` to `.env`:

```bash
cp .env.template .env
```

Populate `.env` with your settings (JWT secret, OpenAI key, etc.). See [Environment Variables](#environment-variables).

## Architecture Overview

* **API Server**: Go HTTP server with JWT auth, exposes REST endpoints.
* **Worker**: Go background process using Redis (Asynq) for file processing and quiz generation.
* **PostgreSQL**: stores users, buckets, files, quizzes, questions, attempts.
* **Redis**: job queue for background tasks.
* **Angular Frontend**: SPA served by Nginx, handles login, bucket/file UI, and quiz workflows.

---

## Backend

### Backend Directory Structure

```
backend/
├── cmd/
│   ├── api/          # API server entrypoint + Dockerfile
│   └── worker/       # Worker entrypoint + Dockerfile
├── internal/
│   ├── auth/         # JWT and auth handlers/middleware
│   ├── db/           # GORM Postgres init
│   ├── ai/           # OpenAI wrapper (embeddings & chat)
│   ├── bucket/       # Bucket CRUD and AI renaming
│   ├── file/         # File upload, storage, queue enqueue
│   ├── quiz/         # Quiz endpoints and service logic
│   └── utils/        # PDF extraction + text chunking
├── go.mod
├── go.sum
└── README.md         # (This file)
```

#### Key Components

* `cmd/api/main.go`: sets up routes, middleware, DB migrations and multi-stage Dockerfile.
* `cmd/worker/main.go`: registers Asynq task handlers for file & quiz jobs and its Dockerfile.
* `internal/auth`: JWT init, signup & login handlers, middleware.
* `internal/db`: GORM connection and `AutoMigrate()`.
* `internal/ai`: `GetEmbedding`, `GenerateBucketName`, `GenerateQuestions` using OpenAI.
* `internal/bucket`: create/list buckets; bucket renaming by AI.
* `internal/file`: multipart upload handler, stores file, enqueues `ProcessFile` task.
* `internal/quiz`: endpoints for quiz lifecycle; `GenerateQuiz` service enqueues & writes Q\&A.
* `internal/utils`: PDF text extraction and chunking logic.

### Backend Running Locally

With Docker Compose (see [Docker & Docker Compose](#docker--docker-compose)):

```bash
docker-compose up --build
```

API will be available at `http://localhost:8080`.

Direct local run (Go & Redis & Postgres installed):

```bash
cd backend
go run cmd/api/main.go   # in one shell
# in another shell
go run cmd/worker/main.go
```

---

## Frontend

### Frontend Directory Structure

```
frontend/
├── Dockerfile          # multi-stage Angular build → Nginx
├── angular.json
├── package.json
└── src/
    ├── app/
    │   ├── app.routes.ts
    │   ├── app.component.ts/html/css
    │   ├── app.config.ts  # providers (HttpClient, Router)
    │   ├── components/
    │   │   ├── login/
    │   │   ├── signup/
    │   │   ├── nav-bar/
    │   │   ├── bucket-list/
    │   │   └── bucket-detail/
    │   ├── guards/auth.guard.ts
    │   └── services/
    │       ├── auth.service.ts
    │       ├── env.service.ts
    │       └── token.interceptor.ts
    ├── index.html        # includes env.js loader
    ├── main.ts           # bootstrapApplication
    └── styles.css
```

### Key Frontend Components

* `NavBar`: header with “Home”, “Buckets” toggle, and Logout button.
* `LoginComponent` & `SignupComponent`: authentication forms.
* `AuthGuard`: protects routes when not logged in.
* `BucketListComponent`: renders bucket list in a drawer/panel toggled by NavBar.
* Services:

  * `AuthService`: login/signup, token storage & login state observable.
  * `EnvService`: reads `ALLOW_SIGNUP` flag.
  * `TokenInterceptor`: attaches JWT to HTTP requests.

### Building & Running

From the monorepo root (using Docker Compose) the frontend is served by Nginx at `http://localhost:9000` (configurable via `UI_HOST_PORT`).

For local Angular dev server:

```bash
cd frontend
npm install
npm start  # serves at http://localhost:4200
```

---

## Docker & Docker Compose

The root `docker-compose.yaml` brings up all services:

* **postgres** (with pgvector)
* **redis**
* **api** (Go HTTP server)
* **worker** (Go background jobs)
* **frontend** (Nginx + Angular)

```bash
docker-compose up --build
```

Ports:

* API: `localhost:8080`
* Frontend: `localhost:9000`

---

## Environment Variables

Copy `.env.template` → `.env` and fill in:

```dotenv
# PostgreSQL
db_host=postgres
db_port=5432
db_user=quizgenie
db_password=quizgenie
db_name=quizgenie

# Redis
redis_addr=redis:6379

# JWT & OpenAI
jwt_secret=<your_jwt_secret>
openai_api_key=<your_openai_key>

# File storage path
file_storage_path=/data/uploads

# Frontend
allow_signup=true
```

---

## Component Workflows

1. **Signup / Login**: user obtains JWT, stored in `localStorage`.
2. **Create Bucket**: upload first file via `POST /buckets`, placeholder name.
3. **ProcessFile**: worker extracts text, chunks, embeddings, renames bucket via AI, marks file complete.
4. **Bucket List**: drawer polls `GET /buckets` and shows AI-generated names.
5. **File Status**: detail view polls `GET /buckets/{id}/files` every 5s.
6. **Take Quiz**: settings → `POST /buckets/{id}/quizzes` → poll `/quizzes/{quizId}` until ready.
7. **Quiz**: fetch questions → take quiz (timed/practice) → submit answers → view report.
8. **History**: list attempts via `GET /buckets/{id}/attempts`.

---

## UI Development Roadmap

1. **Application Shell & Navigation**

   * Implement a **top-level header** (`NavBar`) and remove the permanent sidebar.
   * Add a **“Buckets” menu** in the header that toggles a drawer (or dropdown panel) containing the bucket list.

2. **NavBar Component**

   * Inject `AuthService` and bind `isLoggedIn()` to the template using `*ngIf` to show/hide the “Logout” button.
   * Highlight the active route link (use `routerLinkActive`).
   * **Add a “Buckets” toggle** that opens/closes the bucket list drawer or dropdown.
   * Ensure responsive behavior: on mobile, the drawer can cover full screen; on desktop, slides in from left.

3. **Authentication Forms**

   * Build `LoginComponent` and `SignupComponent` with Angular template-driven forms.
   * Add live validation feedback (required, email format, password match).
   * Display errors inline and consider a global toast service for server errors.
   * Use `EnvService.allowSignup` to conditionally render the signup link.

4. **Routing & Guards**

   * Define public (`/login`, `/signup`) and protected routes (`/buckets`, `/quizzes`, `/attempts`).
   * Apply `AuthGuard` to all protected routes in `app.routes.ts`.
   * On navigation start, redirect unauthenticated users to `/login` and preserve return URL.

5. **Bucket List Drawer/Dropdown**

   * Render `<app-bucket-list>` inside a hidden panel toggled by the NavBar’s “Buckets” button.
   * Ensure keyboard accessibility (focus trap, Esc to close).
   * Support both slide‑in (desktop) and overlay/full‑screen (mobile) modes.
   * Poll `GET /buckets` to refresh names; show loading indicator if fetching.

6. **Bucket Detail & File Status**

   * Develop `BucketDetailComponent` showing file list with status icons (pending/processing/completed/failed).
   * Poll `GET /buckets/{id}/files` every 5 seconds; unsubscribe on component destroy.
   * Enable file deletion and re-upload: attach click handlers to delete buttons and upload-more button.
   * Display messages if no files or if all files are processing.

7. **FileUpload Component**

   * Support both `<input type="file">` and drag-and-drop zones.
   * Show upload progress indicator; disable UI during upload.
   * Validate file extensions and size limits before upload.
   * On success, emit an event to parent to refresh file list.

8. **Quiz Settings & Launch**

   * Implement `QuizSettingsComponent` as a modal dialog with checkboxes for timed and practice modes.
   * Validate that at least one file is `completed` before enabling the “Start Quiz” button.
   * On submit, call `QuizService.createQuiz` and navigate to `QuizStatusComponent`.

9. **Quiz Status & Polling**

   * Create `QuizStatusComponent` to poll `QuizService.getQuizStatus(quizId)` every 3 seconds.
   * Show a full-screen loading card or overlay indicating “Preparing your quiz…”.
   * On status "ready", automatically navigate to `QuizTakingComponent`.

10. **Quiz Taking UI**

    * Build `QuizTakingComponent` to render one question at a time or all at once (configurable).
    * For timed mode: display a countdown timer per question and auto-advance or auto-submit.
    * For practice mode: reveal correctness immediately upon answer selection and display explanations.
    * Manage user selections in a form model; upon completion or timeout, collect answers and call `QuizService.submitAnswers`.

11. **Quiz Report & History**

    * Implement `QuizReportComponent` to fetch and display attempt details, highlighting correct/incorrect answers with colors and explanations.
    * Build `ReportHistoryComponent` (or integrate into NavBar menu) to list past attempts with date and score. Clicking an entry opens `QuizReportComponent`.

12. **Styling & Theming**

    * Choose a CSS approach (e.g. Tailwind or custom SASS).
    * Define a consistent color palette, typography, and spacing scale.
    * Create shared UI components: `Button`, `Card`, `Drawer`, `Modal`, and `Spinner`.
    * Apply responsive breakpoints for mobile-friendly design.

13. **Accessibility & Performance**

    * Ensure all interactive elements have keyboard focus states and ARIA labels.
    * Lazy-load feature modules (e.g. quiz, history) to reduce initial bundle size.
    * Use `OnPush` change detection where appropriate for performance.

14. **Component Testing**

    * Write unit tests for each component using Angular’s TestBed.
    * Mock service dependencies and verify UI logic (e.g. show/hide elements).
    * Add end-to-end (E2E) tests with Cypress or Protractor to cover critical workflows.

---

## Database Schema

See `backend/internal/db` models and `AutoMigrate()` in `cmd/api/main.go`. Tables include:

* `users`, `buckets`, `files`, `file_chunks`, `quizzes`, `questions`, `answers`, `attempts`, `attempt_answers`.

---

## Future Work

* Support DOCX/PPTX file types.
* More advanced AI-driven quiz generation.
* Role-based access control.
* UI enhancements (drag & drop, progress bars).
* CI/CD pipelines and integration tests.

Happy quizzing! 🎉
