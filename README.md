Below is a proposed directory structure for both the backend (Go) and frontend (Angular), followed by a step-by-step plan outlining how to build out each component and tie everything together.

---

## 1. Directory Structures

### 1.1 Backend (Go)

```
backend/
├── cmd/
│   ├── api/
│   │   ├── main.go
│   │   └── Dockerfile                   # Multi-stage Dockerfile for the API server
│   └── worker/
│       ├── main.go
│       └── Dockerfile                   # Multi-stage Dockerfile for the background worker
├── internal/
│   ├── auth/
│   │   ├── handlers.go                  # HTTP endpoints: login, signup
│   │   ├── middleware.go                # JWT validation, request context
│   │   └── jwt.go                       # Token generation + parsing
│   ├── bucket/
│   │   ├── handlers.go                  # CRUD endpoints for buckets
│   │   ├── service.go                   # Business logic: create bucket, rename via AI
│   │   └── model.go                     # GORM model definition for Bucket
│   ├── file/
│   │   ├── handlers.go                  # Upload endpoint, status checks, delete
│   │   ├── service.go                   # File storage + enqueue processing
│   │   └── model.go                     # GORM model for File, status enum, metadata
│   ├── quiz/
│   │   ├── handlers.go                  # Endpoints: start quiz, fetch quiz, submit answers
│   │   ├── service.go                   # Business logic: queue quiz job, calculate scores
│   │   └── model.go                     # GORM models for Quiz, Question, Answer, Attempt
│   ├── ai/
│   │   └── openai.go                    # Wrappers for OpenAI Embedding & ChatGPT calls
│   ├── db/
│   │   └── postgres.go                  # DB initialization (GORM + migrations)
│   ├── queue/
│   │   └── redis.go                     # Redis client + job enqueue helpers (e.g. Asynq)
│   └── utils/
│       ├── pdf.go                       # PDF text extraction (using pdfcpu or similar)
│       └── chunk.go                     # Text-chunking logic (split into ~2000 chars)
├── go.mod
├── go.sum
└── README.md                            # Overview + instructions for running locally
```

* **`cmd/api/`**

  * `main.go`: boots the HTTP server, registers routes, attaches middleware, reads config from environment.
  * `Dockerfile`: multi-stage build that compiles `main.go` and packages the binary into a minimal image.

* **`cmd/worker/`**

  * `main.go`: initializes a Redis client (e.g., Asynq), registers task handlers (file processing, quiz generation), and starts the worker loop.
  * `Dockerfile`: similar multi-stage build that compiles the worker binary.

* **`internal/`**

  * **`auth/`**: JWT-based authentication—login/signup handlers, middleware to protect routes, and token utilities.
  * **`bucket/`**: endpoints and logic to create buckets, rename via AI (after the first file is processed), list buckets, delete buckets.
  * **`file/`**: handles multipart file uploads, stores the raw file on disk (e.g. under a per-user folder), inserts a DB record with status = “Pending,” then enqueues a “ProcessFile” job. Exposes endpoints to list user’s files (with `status: Pending|Processing|Completed|Failed`) and to delete.
  * **`quiz/`**: once files are processed, user can request “Generate quiz.” That enqueues a “GenerateQuiz” job. The job collects embeddings from processed files, calls OpenAI to create Q\&As, writes into DB, then marks quiz as ready. Handlers let the frontend poll for quiz status, fetch the question set when ready, submit answers, and fetch historical attempts.
  * **`ai/`**: wrapper functions around OpenAI embedding and chat completion endpoints:

    * `GetEmbedding(text) → []float32`
    * `GenerateBucketName(textChunks) → string`
    * `GenerateQuestions(context, options) → []QuestionData`
  * **`db/`**: uses GORM (or another ORM) to connect to Postgres. Contains an `AutoMigrate()` call for all models at startup.
  * **`queue/`**: uses a Redis-based job library (e.g. [Asynq](https://github.com/hibiken/asynq)) to enqueue background tasks. Exposes `EnqueueFileProcessing(fileID)` and `EnqueueQuizGeneration(quizID)` helpers.
  * **`utils/`**: PDF text extraction (using something like `pdfcpu` or a Go PDF library), plus a simple chunk splitter to break large text into 2000-character chunks.

* **Root files (`go.mod`, `go.sum`, `README.md`)** describe dependencies and how to build/run.

---

### 1.2 Frontend (Angular)

```
frontend/
├── src/
│   ├── app/
│   │   ├── components/
│   │   │   ├── login/
│   │   │   │   ├── login.component.ts
│   │   │   │   ├── login.component.html
│   │   │   │   └── login.component.css
│   │   │   ├── nav-bar/
│   │   │   │   ├── nav-bar.component.ts
│   │   │   │   ├── nav-bar.component.html
│   │   │   │   └── nav-bar.component.css
│   │   │   ├── bucket-list/
│   │   │   │   ├── bucket-list.component.ts
│   │   │   │   ├── bucket-list.component.html
│   │   │   │   └── bucket-list.component.css
│   │   │   ├── file-upload/
│   │   │   │   ├── file-upload.component.ts
│   │   │   │   ├── file-upload.component.html
│   │   │   │   └── file-upload.component.css
│   │   │   ├── file-status/
│   │   │   │   ├── file-status.component.ts
│   │   │   │   ├── file-status.component.html
│   │   │   │   └── file-status.component.css
│   │   │   ├── quiz-settings/
│   │   │   │   ├── quiz-settings.component.ts
│   │   │   │   ├── quiz-settings.component.html
│   │   │   │   └── quiz-settings.component.css
│   │   │   ├── quiz-taking/
│   │   │   │   ├── quiz-taking.component.ts
│   │   │   │   ├── quiz-taking.component.html
│   │   │   │   └── quiz-taking.component.css
│   │   │   ├── quiz-report/
│   │   │   │   ├── quiz-report.component.ts
│   │   │   │   ├── quiz-report.component.html
│   │   │   │   └── quiz-report.component.css
│   │   │   └── report-history/
│   │   │       ├── report-history.component.ts
│   │   │       ├── report-history.component.html
│   │   │       └── report-history.component.css
│   │   ├── models/
│   │   │   ├── user.model.ts
│   │   │   ├── bucket.model.ts
│   │   │   ├── file.model.ts
│   │   │   ├── quiz.model.ts
│   │   │   └── attempt.model.ts
│   │   ├── services/
│   │   │   ├── auth.service.ts
│   │   │   ├── bucket.service.ts
│   │   │   ├── file.service.ts
│   │   │   └── quiz.service.ts
│   │   ├── guards/
│   │   │   └── auth.guard.ts                 # Prevent unauthenticated access
│   │   ├── app-routing.module.ts
│   │   └── app.module.ts
│   ├── assets/
│   │   └── logo.png
│   ├── environments/
│   │   ├── environment.ts
│   │   └── environment.prod.ts
│   ├── index.html
│   ├── main.ts
│   ├── polyfills.ts
│   └── styles.css
├── angular.json
├── package.json
├── package-lock.json
├── tsconfig.json
└── Dockerfile                              # Multi-stage: build → serve via nginx
```

* **`app/components`**

  * **`login`**: a simple form (username/password) that calls `AuthService.login()` to retrieve a JWT.
  * **`nav-bar`**: always visible across authenticated routes; includes a “+ New Bucket” button.
  * **`bucket-list`**: sits in a left sidebar; shows each bucket tab (with AI-generated name). A “+” icon at the top lets you create/upload and thereby create a new bucket.
  * **`file-upload`**: triggered when no buckets exist or when clicking “+ New Bucket.” Contains a drag-and-drop / file chooser control; calls `FileService.upload(file)` → backend.
  * **`file-status`**: next to each file in the currently selected bucket, showing a spinner/pending icon or a green check or red error icon. Automatically polls the backend every few seconds to refresh statuses.
  * **`quiz-settings`**: form fields—checkbox for “Timed (30s per question)”, checkbox for “Practice Mode.” Submit → `QuizService.createQuiz(bucketId, options)`.
  * **`quiz-taking`**: once a quiz is ready, this component renders each question (one at a time or all at once, depending on UX). In Practice Mode, each answer reveals correctness immediately; in “normal” mode, answers are hidden until the very end.
  * **`quiz-report`**: after submission, shows final score, question-by-question details (selected answer, correct answer, explanation).
  * **`report-history`**: lists past quiz attempts in that bucket; clicking one opens `quiz-report` in read-only mode.

* **`app/models`**

  * `user.model.ts`: stores token, username, etc.
  * `bucket.model.ts`: `id`, `name`, `createdAt`, `updatedAt`.
  * `file.model.ts`: `id`, `filename`, `status: 'pending'|'processing'|'done'|'failed'`, timestamps.
  * `quiz.model.ts`: `id`, `bucketId`, `status: 'pending'|'generating'|'ready'|'failed'`, creation timestamps.
  * `attempt.model.ts`: `id`, `quizId`, `score`, `answersSelected`, `completedAt`.

* **`app/services`**

  * `auth.service.ts`: handles `login()`, `signup()`, token persistence (e.g. localStorage). Sets `Authorization: Bearer <token>` header on all HTTP calls.
  * `bucket.service.ts`: `listBuckets()`, `getBucket(bucketId)`, `createBucket(initialFile)`, `renameBucket(bucketId)`.
  * `file.service.ts`: `uploadFile(bucketId, file)`, `getFiles(bucketId)`, `pollFileStatus(bucketId)`, `deleteFile(bucketId, fileId)`.
  * `quiz.service.ts`: `createQuiz(bucketId, options)`, `getQuizStatus(quizId)`, `fetchQuiz(quizId)`, `submitAnswers(quizId, answers)`, `listAttempts(bucketId)`.

* **`guards/auth.guard.ts`**

  * Blocks all routes except `/login` when no valid JWT is present. Redirects to `/login`.

* **`Dockerfile`** (multi-stage)

  1. **Stage 1**: Uses `node:16-alpine`, installs dependencies, runs `ng build --prod`.
  2. **Stage 2**: Uses `nginx:alpine`, copies `dist/` into `/usr/share/nginx/html`, and supplies a minimal `nginx/default.conf` to serve the static files.

---

## 2. Verbose Plan Document

### 2.1 Architecture Overview

* **Frontend (Angular SPA)**

  * Handles user authentication via JWT.
  * Renders a left sidebar (bucket tabs) plus a main pane.
  * Allows file upload, shows file statuses, and drives quiz workflows (settings, taking, reporting).
  * Polls backend to keep bucket/file/quiz statuses up to date.

* **Backend (Go REST API + Worker)**

  * **API Server** (stateless): exposes HTTP endpoints for authentication, bucket management, file uploads, status polling, quiz creation, quiz retrieval, answer submission, and report history. Uses JWT for auth. Persists everything to Postgres.
  * **Background Worker**: a separate process that connects to Redis via a job library (e.g. Asynq). Listens for two job types:

    1. **ProcessFileJob**: once a file is uploaded → extract text (PDF or plain text), chunk it, store chunks (embedding-ready) in the database, call OpenAI embeddings on each chunk (store embeddings), then generate or update the bucket name via a lightweight AI call. Finally, mark file as “Completed” (or “Failed” if anything goes wrong).
    2. **GenerateQuizJob**: once the user requests a quiz → collect embeddings from all processed chunks in that bucket, select a sample (using pool/sample multipliers), assemble `contextStr`, call OpenAI Chat Completion to get a JSON array of question objects, write questions and answer choices into DB tables, then mark quiz as “Ready.”
  * **PostgreSQL** (Docker-ized): stores all models—Users, Buckets, Files, FileChunks (with embeddings), Quizzes, Questions, Answers, Attempts, AttemptAnswers, etc.
  * **Redis** (Docker-ized): job queue for file processing and quiz generation, as well as ephemeral counters (if needed for polling).

* **Containerization**

  * Each component lives in its own Docker container:

    1. **postgres** (official `postgres:15` or a pgvector image).
    2. **redis** (official `redis:alpine`).
    3. **api** (Go binary built with `GOOS=linux`, packaged with only the binary + CA certs).
    4. **worker** (Go binary for background jobs).
    5. **frontend** (Angular static served by `nginx:alpine`).
  * Use a `docker-compose.yaml` at the monorepo root to wire them together, setting appropriate environment variables (DB\_HOST, REDIS\_URL, JWT\_SECRET, etc.).

---

### 2.2 Database Schema (Postgres)

#### 2.2.1 `users` table

| Column         | Type         | Notes                         |
| -------------- | ------------ | ----------------------------- |
| id             | SERIAL PK    | Auto-incrementing primary key |
| username       | VARCHAR(50)  | Unique, indexed               |
| email          | VARCHAR(100) | Unique, optional              |
| password\_hash | VARCHAR(255) | Hashed (e.g. bcrypt)          |
| created\_at    | TIMESTAMP    | Default NOW()                 |
| updated\_at    | TIMESTAMP    | On update current timestamp   |

#### 2.2.2 `buckets` table

| Column      | Type         | Notes                                                |
| ----------- | ------------ | ---------------------------------------------------- |
| id          | SERIAL PK    |                                                      |
| user\_id    | INT FK       | REFERENCES users(id) ON DELETE CASCADE               |
| name        | VARCHAR(255) | AI-generated name                                    |
| created\_at | TIMESTAMP    | Default NOW()                                        |
| updated\_at | TIMESTAMP    | On update current timestamp                          |
| description | TEXT         | Optional (if you want to store more bucket metadata) |

#### 2.2.3 `files` table

| Column        | Type         | Notes                                                        |
| ------------- | ------------ | ------------------------------------------------------------ |
| id            | SERIAL PK    |                                                              |
| bucket\_id    | INT FK       | REFERENCES buckets(id) ON DELETE CASCADE                     |
| filename      | VARCHAR(255) | Original filename                                            |
| storage\_path | VARCHAR(500) | Path on disk (e.g. `/data/user_<id>/bucket_<id>/file.pdf`)   |
| status        | VARCHAR(20)  | Enum: `'pending'`, `'processing'`, `'completed'`, `'failed'` |
| error\_msg    | TEXT         | Nullable—if processing failed, store the error               |
| created\_at   | TIMESTAMP    | Default NOW()                                                |
| updated\_at   | TIMESTAMP    | On update current timestamp                                  |

#### 2.2.4 `file_chunks` table

| Column       | Type         | Notes                                                           |
| ------------ | ------------ | --------------------------------------------------------------- |
| id           | SERIAL PK    |                                                                 |
| file\_id     | INT FK       | REFERENCES files(id) ON DELETE CASCADE                          |
| chunk\_index | INT          | 0,1,2,…—the sequence of each chunk                              |
| content      | TEXT         | Raw text chunk extracted from the PDF or text file              |
| embedding    | VECTOR(1536) | pgvector column—store output of OpenAI embedding for each chunk |
| created\_at  | TIMESTAMP    | Default NOW()                                                   |
| updated\_at  | TIMESTAMP    | On update current timestamp                                     |

#### 2.2.5 `quizzes` table

| Column         | Type        | Notes                                                                       |
| -------------- | ----------- | --------------------------------------------------------------------------- |
| id             | SERIAL PK   |                                                                             |
| bucket\_id     | INT FK      | REFERENCES buckets(id) ON DELETE CASCADE                                    |
| status         | VARCHAR(20) | Enum: `'pending'`, `'generating'`, `'ready'`, `'failed'`                    |
| timed\_mode    | BOOLEAN     | `true` if 30-second timing per question is enabled                          |
| practice\_mode | BOOLEAN     | `true` if solutions should appear immediately upon submitting each question |
| created\_at    | TIMESTAMP   | Default NOW()                                                               |
| updated\_at    | TIMESTAMP   | On update current timestamp                                                 |
| error\_msg     | TEXT        | Nullable—if quiz generation fails, store the error                          |

#### 2.2.6 `questions` table

| Column      | Type      | Notes                                              |
| ----------- | --------- | -------------------------------------------------- |
| id          | SERIAL PK |                                                    |
| quiz\_id    | INT FK    | REFERENCES quizzes(id) ON DELETE CASCADE           |
| text        | TEXT      | Question prompt                                    |
| explanation | TEXT      | General explanation about why the question matters |
| created\_at | TIMESTAMP | Default NOW()                                      |
| updated\_at | TIMESTAMP | On update current timestamp                        |

#### 2.2.7 `answers` table

| Column       | Type      | Notes                                      |
| ------------ | --------- | ------------------------------------------ |
| id           | SERIAL PK |                                            |
| question\_id | INT FK    | REFERENCES questions(id) ON DELETE CASCADE |
| text         | TEXT      | One answer choice                          |
| is\_correct  | BOOLEAN   | `true` if this choice is correct           |
| explanation  | TEXT      | Explanation why it’s correct or incorrect  |
| created\_at  | TIMESTAMP | Default NOW()                              |
| updated\_at  | TIMESTAMP | On update current timestamp                |

#### 2.2.8 `attempts` table

| Column      | Type      | Notes                                    |
| ----------- | --------- | ---------------------------------------- |
| id          | SERIAL PK |                                          |
| quiz\_id    | INT FK    | REFERENCES quizzes(id) ON DELETE CASCADE |
| user\_id    | INT FK    | REFERENCES users(id) ON DELETE CASCADE   |
| score       | FLOAT     | Percentage score (e.g. 85.0)             |
| created\_at | TIMESTAMP | Default NOW()                            |
| updated\_at | TIMESTAMP | On update current timestamp              |

#### 2.2.9 `attempt_answers` table

| Column       | Type      | Notes                                                                         |
| ------------ | --------- | ----------------------------------------------------------------------------- |
| id           | SERIAL PK |                                                                               |
| attempt\_id  | INT FK    | REFERENCES attempts(id) ON DELETE CASCADE                                     |
| question\_id | INT FK    | REFERENCES questions(id) ON DELETE CASCADE                                    |
| answer\_id   | INT FK    | REFERENCES answers(id) ON DELETE CASCADE                                      |
| is\_correct  | BOOLEAN   | Redundant but helps when querying—store whether the chosen answer was correct |
| created\_at  | TIMESTAMP | Default NOW()                                                                 |

---

### 2.3 Backend: Detailed Component Plan

#### 2.3.1 Environment & Configuration

1. **Environment Variables** (read from `.env` or Docker Compose):

   * `PORT` (default: 8080 for the API).
   * `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME` (for PostgreSQL).
   * `REDIS_ADDR` (e.g. `redis:6379`).
   * `JWT_SECRET` (for signing tokens).
   * `OPENAI_API_KEY` (for AI calls).
   * `FILE_STORAGE_PATH` (where uploaded files are saved on disk).
   * `BASE_URL` or `FRONTEND_URL` (for CORS config).

2. **Bootstrap (`cmd/api/main.go`)**:

   * Load env vars (e.g. via `github.com/joho/godotenv`).
   * Initialize GORM + AutoMigrate all models (`autoMigrate(&User{}, &Bucket{}, &File{}, &FileChunk{}, &Quiz{}, &Question{}, &Answer{}, &Attempt{}, &AttemptAnswer{})`).
   * Initialize Redis client + Asynq client.
   * Register route handlers with a router (e.g. [chi](https://github.com/go-chi/chi) or Gin).
   * Add JWT middleware to protect relevant endpoints (all except `/login`, `/signup`).
   * Start listening on `0.0.0.0:PORT`.

3. **Worker (`cmd/worker/main.go`)**:

   * Load same env vars.
   * Initialize GORM to talk to the same Postgres (so it can update DB records).
   * Initialize Redis + Asynq server.
   * Register two task handlers:

     * `ProcessFileTask`: receives `FileID` → fetch file record, update status to “processing,” invoke `ProcessFile(fileID)` service.
     * `GenerateQuizTask`: receives `QuizID` → fetch quiz record, update status to “generating,” invoke `GenerateQuiz(quizID)` service.
   * Start the Asynq worker loop (concurrency configurable via env var).

#### 2.3.2 Authentication Module

* **Password hashing**: use `golang.org/x/crypto/bcrypt`. On signup, hash password and store `password_hash`. On login, compare bcrypt.

* **JWT**: use a library like `github.com/golang-jwt/jwt`.

  * `auth/jwt.go` has `GenerateToken(userID, username)` → `tokenString`, and `ParseToken(tokenString)` → `userID`.
  * Middleware (`auth/middleware.go`) checks `Authorization: Bearer <token>` header, verifies signature, rejects if invalid/expired. On success, attaches `userID` to request context.

* **Endpoints (`auth/handlers.go`)**:

  1. **POST `/signup`**

     * Accept `{"username": "...", "password": "...", "email": "..."}`.
     * Check uniqueness, hash password, create user row, Return HTTP 201 + JWT in JSON.
  2. **POST `/login`**

     * Accept `{"username": "...", "password": "..."}`.
     * Lookup user, compare bcrypt, if valid → `GenerateToken(...)`, return `{"token":"...", "username":"..."}`. Otherwise 401.

#### 2.3.3 Bucket & File Management

* **When user first logs in**: frontend checks `bucket.service.listBuckets()`. If empty, show a big “Upload your first file to create a bucket” prompt (modal or inline upload).

* **`POST /api/buckets`** (in `bucket/handlers.go`):

  * Request must include a single file upload (multipart form): `file: <File>`. The handler:

    1. Create a new bucket row immediately with `name="(processing…)"` or placeholder.
    2. Save the file to disk under `FILE_STORAGE_PATH/user_<userID>/bucket_<bucketID>/original_filename.ext`.
    3. Create a new file record with `bucket_id`, `filename`, `storage_path`, `status='pending'`.
    4. Enqueue a `ProcessFileTask` (pass `fileID`) to Redis.
    5. Return HTTP 201 with `{"bucketId": X, "fileId": Y, "bucketName": "(processing…)"}`.

* **`GET /api/buckets`**:

  * Returns a list of buckets belonging to `userID` (from JWT), each with `{id, name, createdAt}`.

* **Inside `ProcessFileService.ProcessFile(fileID int)`**:

  1. Update `files` record: `status='processing'`.
  2. If PDF (detect by extension), call `utils/pdf.go:ExtractTextFromPDF(path)` → a large string. If plain text, just read it.
  3. Call `utils/chunk.go:ChunkText(fullText, 2000)` → `[]string{chunk0, chunk1, …}`.
  4. For each chunk:

     * Create a `file_chunks` row with `file_id`, `chunk_index`, `content=chunkText`, `embedding=nil`.
     * Immediately call `Ai.GetEmbedding(chunkText)` → `[]float32`.
     * Update that `file_chunks` row’s `embedding` column.
  5. (Only after uploading at least one chunk) assemble a combined summary prompt (e.g. first 500 characters across chunks) and call `Ai.GenerateBucketName(textChunks)`.

     * Update the `buckets` row’s `name` with the returned string.
  6. Update `files` record: `status='completed'`.
  7. If any step fails, set `files.status='failed'` and store `error_msg`.

* **`GET /api/buckets/{bucketId}/files`**:

  * Returns all files under that bucket with `{id, filename, status}` so the frontend can display an icon or color code.
  * Does **not** expose chunks or embeddings—only the original file metadata and status.

* **`DELETE /api/buckets/{bucketId}/{fileId}`**:

  * Marks the file record as deleted, physically deletes the file from disk, and optionally deletes all `file_chunks` rows for that file. Returns 200.

* **Bucket name updates**:

  * The bucket name starts as placeholder in the DB. As soon as a chunk’s embedding is written (after the first chunk), the worker calls `GenerateBucketName`. At that point, the bucket gets a human-friendly name, which the frontend will pick up on its next polling cycle.

#### 2.3.4 Quiz Generation Workflow

* **Front-end logic**:

  1. “Take Quiz” button is initially disabled until at least one file’s status is `'completed'`.
  2. When user clicks “Take Quiz”:

     * Frontend calls `GET /api/buckets/{bucketId}/quiz-status`.
     * If any file under that bucket has `status='pending'` or `status='processing'`, return `{canStart:false, message:"Some files are still processing. Proceed anyway?"}`.
     * Otherwise `{canStart:true}`.
     * If `canStart=false` and user chooses “Proceed anyway,” continue. If “Wait,” return to main view.

* **API Endpoint: `POST /api/buckets/{bucketId}/quizzes`** (in `quiz/handlers.go`):

  * Payload: `{ timedMode: bool, practiceMode: bool }`.
  * Handler:

    1. Create a new `quizzes` row:

       * `bucket_id = bucketId`
       * `status = 'pending'`
       * `timed_mode = payload.timedMode`
       * `practice_mode = payload.practiceMode`
    2. Enqueue `GenerateQuizTask` with `quizID`.
    3. Return HTTP 202 with `{quizId, status:"pending"}`.

* **Background: `GenerateQuizService.GenerateQuiz(quizID int)`**:

  1. Update `quizzes.status='generating'`.
  2. Fetch all processed `file_chunks` for that bucket (i.e. `JOIN files ON files.id=file_chunks.file_id WHERE files.bucket_id = <bucketId> AND files.status='completed'`). Collect their `embedding` vectors and text content.
  3. Use a similarity logic (cosine distance with pgvector or in Go) to pick a semantic pool—e.g. top `N = max(1, int(questionCount * semanticMultiplier))`. Then randomly sample `M = max(1, int(questionCount * randomMultiplier))`. For now, default `semanticMultiplier=3.0`, `randomMultiplier=1.0`, `questionCount=10`.
  4. Concatenate the chosen chunks’ text into one large `context` string.
  5. Call `Ai.GenerateQuestions(context, questionCount, numChoices, difficulty)` → returns a JSON array of question objects:

     ```json
     [
       {
         "question": "...",
         "explanation": "...",
         "answers": [
           {"text":"...","is_correct":true,"explanation":"..."},
           {"text":"...","is_correct":false,"explanation":"..."},
           …
         ]
       },
       …
     ]
     ```
  6. For each question object:

     * Insert a row into `questions(quiz_id, text, explanation)`.
     * For each answer choice, insert `answers(question_id, text, is_correct, explanation)`.
  7. Update `quizzes.status='ready'`.

* **Endpoint: `GET /api/quizzes/{quizId}`**

  * Returns `{id, bucketId, status}`. Frontend polls every few seconds until `status == 'ready'`.

* **Endpoint: `GET /api/quizzes/{quizId}/questions`**

  * Returns an array of `{ questionId, text, answers: [ {id, text, isCorrect (only for practice-mode?), explanation }, … ] }`. In “practiceMode” frontends, you may choose to expose `isCorrect` client-side to let the frontend reveal it immediately. Otherwise omit it, and do correctness checking on the backend.

* **Endpoint: `POST /api/quizzes/{quizId}/attempts`**

  * Payload:

    ```json
    {
      "answers": [
        { "questionId": 123, "answerId": 456 },
        …
      ]
    }
    ```
  * Handler:

    1. Create a new `attempts(quiz_id, user_id, score=null)` row.
    2. For each object in `answers`, insert into `attempt_answers(attempt_id, question_id, answer_id, is_correct)` by checking whether `answer.is_correct == true`.
    3. Compute `score = (correctCount / totalCount) * 100` and update `attempts.score = score`.
    4. Return `{attemptId, score}`.

* **Endpoint: `GET /api/buckets/{bucketId}/attempts`**

  * Returns a list of all quiz attempts in that bucket for the logged-in user: `[ { attemptId, quizId, score, createdAt }, … ]`. Frontend uses this to populate “Report History.”

* **Endpoint: `GET /api/attempts/{attemptId}`**

  * Returns detailed payload:

    ```json
    {
      "attemptId": 789,
      "quizId": 456,
      "score": 85.0,
      "details": [
        {
          "questionId": 123,
          "questionText": "...",
          "selectedAnswerId": 456,
          "selectedAnswerText": "...",
          "isCorrect": true,
          "correctAnswerText": "…",
          "explanation": "…"
        },
        …
      ]
    }
    ```

---

### 2.4 Frontend: Component & UX Flow

1. **App Initialization**

   * On startup, Angular checks for a valid JWT in `localStorage`.
   * If none or expired → redirect to `/login`.
   * If valid → fetch `bucket.service.listBuckets()`.

2. **Login Screen (`/login`)**

   * Simple form: username + password → `authService.login()`. On success, store JWT in `localStorage`, navigate to `/buckets`.

3. **Buckets View (`/buckets`)**

   * **Left Sidebar** (vertical list of tabs): each bucket name. The top has a “+ New Bucket” button.
   * **Main Pane**: if no bucket is selected, show a prompt “Select a bucket or + New Bucket.” If there is at least one bucket, the first one is auto-selected on load.

4. **Creating a New Bucket**

   * Click “+ New Bucket” → open a modal or overlay with the `FileUploadComponent`.
   * User picks one (or multiple) files (for now, just accept one file at a time). On drop or select → call `fileService.uploadFile(bucketId, file)` (but since bucket doesn’t exist yet, first call `bucketService.createBucket(file)`).
   * The API returns `{bucketId, fileId}`. Immediately close the modal and switch sidebar’s active tab to the new bucket. The new bucket appears with name “(processing…)” until the worker updates it.

5. **File Upload & Status View**

   * Once bucket is created, the bucket view automatically calls `fileService.getFiles(bucketId)`. Displays a table or list: each row = a file, showing its `filename` and an icon depending on `status`:

     * `'pending'` → gray clock icon
     * `'processing'` → spinner
     * `'completed'` → green check
     * `'failed'` → red exclamation
   * The frontend polls `fileService.getFiles(bucketId)` every 5 seconds to refresh statuses.
   * If `file.status==='completed'` for the very first file in that bucket, immediately enable the “Take Quiz” button in the top toolbar. Also, once that first chunk has been processed, the worker calls OpenAI to generate a bucket name—on the next poll, the sidebar’s bucket tab label updates to the AI-generated name.

6. **Take Quiz / Quiz Settings**

   * The top of the main pane has a “Take Quiz” button. Initially disabled. When at least one file has `status==='completed'`, enable it.
   * When user clicks “Take Quiz,” open the `QuizSettingsComponent` as a modal or inline form. It contains:

     * Checkbox: “Timed Mode (30 seconds per question).”
     * Checkbox: “Practice Mode (show correct/incorrect immediately).”
     * “Start Quiz” button.
   * On “Start Quiz”: call `quizService.createQuiz(bucketId, {timedMode, practiceMode})`. API returns `{quizId, status:'pending'}`. Immediately navigate to `/quizzes/{quizId}/status`.

7. **Quiz Status / Loading Screen**

   * `QuizStatusComponent` polls `quizService.getQuizStatus(quizId)` every 3 seconds.
   * Show a loading spinner + text: “Preparing your quiz…”
   * If the API ever returns `{status:'ready'}`, navigate to `/quizzes/{quizId}/take`. If `{status:'failed'}`, show “Failed to generate quiz. Try again.”

8. **Taking the Quiz (`/quizzes/{quizId}/take`)**

   * First, fetch `quizService.getQuizQuestions(quizId)`. Backend returns JSON array of question objects, including each answer’s `text` and `isCorrect` only if `practiceMode==true`. If `practiceMode==false`, omit `isCorrect` on the payload—frontend just displays answer options.
   * **UI Layout**:

     * If `timedMode==true`, display a countdown timer (30 seconds) in the corner; at the end of 30 seconds, automatically proceed to next question (or auto-submit current selection if user has selected one).
     * Otherwise, user can click “Next Question” manually.
     * For each question: show the prompt text, render each answer as a radio button.
     * If `practiceMode==true`, once user picks a radio button, immediately reveal correctness (change background to green/red and show `explanation`). Disable further changes for that question. Then show a “Next” or “Continue” button.
     * If `practiceMode==false`, do not reveal anything until they have answered all questions and click “Submit Quiz.”

9. **Submit Quiz & Show Report**

   * On final question (or after “Submit Quiz”), collect user’s selected answers into a payload:

     ```json
     {
       "answers": [
         { "questionId": 123, "answerId": 456 },
         …
       ]
     }
     ```
   * Send `quizService.submitAnswers(quizId, answers)`. Backend returns `{attemptId, score}`. Navigate to `/attempts/{attemptId}`.

10. **Report View (`/attempts/{attemptId}`)**

    * Fetch detailed attempt: `quizService.getAttemptDetails(attemptId)`. That returns:

      ```json
      {
        "attemptId": 789,
        "quizId": 456,
        "score": 85,
        "details": [
          {
            "questionText": "...",
            "selectedAnswerText": "...",
            "isCorrect": true,
            "correctAnswerText": "...",
            "explanation": "..."
          },
          …
        ]
      }
      ```
    * Render each question in an ordered list: show question text, then show all answer choices. Highlight the user’s chosen answer (green if correct, red if incorrect), and also highlight the actual correct answer (border if they didn’t pick it). Show the explanation under each answer that the user selected or the correct answer.

11. **Report History (`/buckets/{bucketId}/history`)**

    * A list view that calls `quizService.listAttempts(bucketId)` and displays each attempt with date/time and score. Each entry is a link to `/attempts/{attemptId}`.

12. **Adding/Removing Files from an Existing Bucket**

    * In the bucket’s main pane, next to the file list, include “Upload More Files” (reuse `FileUploadComponent`) and a delete “X” icon next to each file.
    * When a new file is uploaded to an **existing** bucket, the bucket name does not change (unless you want to re-AI-generate—could leave that as future work). The new file goes through `ProcessFileTask` as usual.
    * If a user deletes a file from an existing bucket, the `files.status` and `file_chunks` are removed. Future quiz generations will ignore that file’s chunks.

---

### 2.5 Background Workers & Job Queue

1. **Choose a Job Library**

   * Use [Asynq](https://github.com/hibiken/asynq) (Redis-backed) because it has a simple API and is well-documented.
   * In `internal/queue/redis.go`, initialize:

     ```go
     import "github.com/hibiken/asynq"

     var RedisClient *asynq.Client
     var RedisServer *asynq.Server

     func InitRedis(addr, password string) {
       RedisClient = asynq.NewClient(asynq.RedisClientOpt{Addr: addr, Password: password})
       RedisServer = asynq.NewServer(
         asynq.RedisClientOpt{Addr: addr, Password: password},
         asynq.Config{Concurrency: 10},
       )
     }
     ```
   * In `cmd/api/main.go`, call `InitRedis(os.Getenv("REDIS_ADDR"), os.Getenv("REDIS_PASSWORD"))`.
   * When you want to enqueue a file job:

     ```go
     taskPayload := map[string]interface{}{"file_id": fileID}
     task := asynq.NewTask("ProcessFile", taskPayload)
     RedisClient.Enqueue(task)
     ```
   * Worker sets up a mux:

     ```go
     mux := asynq.NewServeMux()
     mux.HandleFunc("ProcessFile", ProcessFileTaskHandler)
     mux.HandleFunc("GenerateQuiz", GenerateQuizTaskHandler)
     RedisServer.Run(mux)
     ```

2. **File Processing Handler**

   * Receives `{{"file_id": 123}}` from the payload.
   * Calls `ProcessFileService.ProcessFile(123)`. On any panic or error, log it and update `files.status='failed'`.

3. **Quiz Generation Handler**

   * Receives `{{"quiz_id": 456}}`.
   * Calls `GenerateQuizService.GenerateQuiz(456)`. On error → `quizzes.status='failed'`.

4. **Scaling**

   * If concurrency needs to scale, simply increase `Asynq.Config.Concurrency` or run multiple worker replicas behind the same Redis.

---

### 2.6 AI Integration (OpenAI)

1. **`internal/ai/openai.go`**

   * Initialize a client:

     ```go
     import "github.com/sashabaranov/go-openai"

     var OpenAIClient *goopenai.Client

     func InitOpenAI(apiKey string) {
       OpenAIClient = goopenai.NewClient(apiKey)
     }
     ```
   * **`GetEmbedding(text string) ([]float32, error)`**: call `OpenAIClient.CreateEmbeddings(...)` with model=`"text-embedding-3-small"`. Return the embedding vector.
   * **`GenerateBucketName(chunks []string) (string, error)`**:

     ```go
     prompt := "Based on the following text snippets, provide a concise bucket name:\n\n" + strings.Join(chunks[:min(len(chunks), 3)], "\n\n")
     resp, err := OpenAIClient.CreateChatCompletion(ctx, goopenai.ChatCompletionRequest{
       Model: "gpt-4",
       Messages: []goopenai.ChatCompletionMessage{
         {Role: "system", Content: "You are an AI that generates short, descriptive names."},
         {Role: "user", Content: prompt},
       },
       MaxTokens: 16,
     })
     return resp.Choices[0].Message.Content, nil
     ```
   * **`GenerateQuestions(context string, count int, choices int, difficulty string) ([]QuestionData, error)`**: format a chat prompt very similar to the old Django version, call ChatCompletion with `model=“gpt-4”`, parse JSON output into a slice of a local struct `QuestionData`.

2. **Error Handling & Retries**

   * If an OpenAI call fails, either retry up to 3 times with exponential backoff or mark the job as failed. Log the exact error in the `error_msg` column.

---

### 2.7 Dockerization

#### 2.7.1 `backend/cmd/api/Dockerfile`

```dockerfile
# Stage 1: Build Go binary
FROM golang:1.20-alpine AS builder
RUN apk add --no-cache git gcc musl-dev
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
WORKDIR /app/cmd/api
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /quizgenie_api

# Stage 2: Final image
FROM alpine:latest
RUN apk --no-cache add ca-certificates

# Create a group and user to run the app (non-root)
RUN addgroup -S appuser && adduser -S -G appuser appuser
USER appuser

WORKDIR /home/appuser
COPY --from=builder /quizgenie_api .
COPY --from=builder /app/.env .   # Or rely on docker-compose to mount env

EXPOSE 8080
ENTRYPOINT ["./quizgenie_api"]
```

#### 2.7.2 `backend/cmd/worker/Dockerfile`

```dockerfile
# Stage 1: Build worker binary
FROM golang:1.20-alpine AS builder
RUN apk add --no-cache git gcc musl-dev
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
WORKDIR /app/cmd/worker
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /quizgenie_worker

# Stage 2: Final image
FROM alpine:latest
RUN apk --no-cache add ca-certificates

RUN addgroup -S appuser && adduser -S -G appuser appuser
USER appuser

WORKDIR /home/appuser
COPY --from=builder /quizgenie_worker .
COPY --from=builder /app/.env .   # Or mount via docker-compose

ENTRYPOINT ["./quizgenie_worker"]
```

#### 2.7.3 `frontend/Dockerfile`

```dockerfile
# Stage 1: Build Angular app
FROM node:16-alpine AS build
WORKDIR /app
COPY package.json package-lock.json ./
RUN npm install
COPY . .
RUN npm run build --prod

# Stage 2: Serve with nginx
FROM nginx:alpine
COPY --from=build /app/dist/frontend /usr/share/nginx/html
# Optional: custom nginx config if you need fallback to index.html for Angular routes
COPY nginx/default.conf /etc/nginx/conf.d/default.conf

EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
```

* **`nginx/default.conf`** (in `frontend/nginx/default.conf`):

  ```nginx
  server {
      listen 80;
      server_name _;

      root /usr/share/nginx/html;
      index index.html;

      location / {
          try_files $uri $uri/ /index.html;
      }

      location /api/ {
          proxy_pass http://api:8080/;
          proxy_set_header Host $host;
          proxy_set_header X-Real-IP $remote_addr;
          proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
      }
  }
  ```

#### 2.7.4 `docker-compose.yaml` (Monorepo root)

```yaml
version: "3.8"
services:
  postgres:
    image: postgres:15
    environment:
      - POSTGRES_DB=${DB_NAME:-quizgenie}
      - POSTGRES_USER=${DB_USER:-quizgenie}
      - POSTGRES_PASSWORD=${DB_PASSWORD:-quizgenie}
    volumes:
      - pgdata:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U ${DB_USER:-quizgenie}"]
      interval: 5s
      timeout: 5s
      retries: 5

  redis:
    image: redis:alpine
    ports:
      - "6379:6379"
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 5s
      timeout: 5s
      retries: 5

  api:
    build:
      context: ./backend
      dockerfile: ./cmd/api/Dockerfile
    ports:
      - "8080:8080"
    environment:
      - DB_HOST=postgres
      - DB_PORT=5432
      - DB_USER=${DB_USER:-quizgenie}
      - DB_PASSWORD=${DB_PASSWORD:-quizgenie}
      - DB_NAME=${DB_NAME:-quizgenie}
      - REDIS_ADDR=redis:6379
      - JWT_SECRET=${JWT_SECRET:-supersecretjwt}
      - OPENAI_API_KEY=${OPENAI_API_KEY}
      - FILE_STORAGE_PATH=/data/uploads
    volumes:
      - ./backend/uploads:/data/uploads      # Persist user uploads
      - ./backend/.env:/app/.env             # If you use a .env file
    depends_on:
      - postgres
      - redis

  worker:
    build:
      context: ./backend
      dockerfile: ./cmd/worker/Dockerfile
    environment:
      - DB_HOST=postgres
      - DB_PORT=5432
      - DB_USER=${DB_USER:-quizgenie}
      - DB_PASSWORD=${DB_PASSWORD:-quizgenie}
      - DB_NAME=${DB_NAME:-quizgenie}
      - REDIS_ADDR=redis:6379
      - JWT_SECRET=${JWT_SECRET:-supersecretjwt}
      - OPENAI_API_KEY=${OPENAI_API_KEY}
      - FILE_STORAGE_PATH=/data/uploads
    volumes:
      - ./backend/uploads:/data/uploads
      - ./backend/.env:/app/.env
    depends_on:
      - postgres
      - redis

  frontend:
    build:
      context: ./frontend
      dockerfile: Dockerfile
    ports:
      - "80:80"
    depends_on:
      - api

volumes:
  pgdata:
```

* **Key points**:

  * `api` and `worker` share the same volume `./backend/uploads:/data/uploads` so the worker can read newly uploaded files.
  * Both read the same `.env`, or you can supply a separate `.env` in the compose file.
  * `frontend` is listening on port 80, proxies any `/api/` calls to `api:8080`.

---

### 2.8 Implementation Roadmap

Below is a suggested step-by-step order of implementation:

1. **Initial Setup & Boilerplate**

   * Create the monorepo with two subfolders: `backend/` and `frontend/`.
   * Initialize Go modules inside `backend/` (`go mod init github.com/yourusername/quizgenie-backend`).
   * Initialize Angular project inside `frontend/` (`ng new frontend --routing --style=css`).
   * Commit a minimal `README.md` in each folder outlining high-level structure.

2. **Database Models & GORM Migrations**

   * In `backend/internal/db/postgres.go`, set up GORM to connect to Postgres (read env vars).
   * Define all model structs (User, Bucket, File, FileChunk, Quiz, Question, Answer, Attempt, AttemptAnswer) with GORM tags.
   * In `cmd/api/main.go`, call `db.AutoMigrate()` for every model.
   * Test that, when running the API container, the database tables appear in Postgres.

3. **Authentication & JWT**

   * Implement `internal/auth/jwt.go` (token generation and parsing).
   * Implement `internal/auth/handlers.go`: `/signup`, `/login`.
   * Add `auth/middleware.go` and apply it to a sample protected route (e.g. `GET /api/ping` returns `“pong”` only if a valid token is supplied).
   * Manually test with `curl` or Postman: signup → login → get token → call `/api/ping`.

4. **Bucket Endpoints**

   * Create `internal/bucket/model.go` (GORM struct).
   * Create `internal/bucket/service.go`: `CreateBucket(userID) → bucketID`, `RenameBucket(bucketID, newName)`.
   * Create `internal/bucket/handlers.go`:

     * `POST /api/buckets` (initially just create a bucket with placeholder name).
     * `GET /api/buckets`.
   * Wire routes in `cmd/api/main.go`.
   * Test with Postman: `POST /api/buckets` (returns bucket ID), then `GET /api/buckets` should show it.

5. **File Upload & Status**

   * In `internal/file/handlers.go`, implement `POST /api/buckets/{bucketId}/files` to accept a single multipart file:

     * Save to `FILE_STORAGE_PATH/user_<userId>/bucket_<bucketId>/…`.
     * Insert a `files` row with `status='pending'`.
     * Enqueue `ProcessFileTask`.
     * Return IDs to the client.
   * Implement `GET /api/buckets/{bucketId}/files` that returns file metadata + status.
   * Create `internal/file/model.go` (GORM struct for File) and `internal/file/service.go` containing the business logic to store on disk and enqueue a job.
   * In `cmd/api/main.go`, register these routes under JWT middleware.
   * Write a minimal `ProcessFileTaskHandler` in `cmd/worker/main.go` that simply logs “processing fileID” (for now, do nothing else).
   * Manually test upload, check that the file lands on disk, DB record appears as “pending,” and the worker logs the job.
   * Update `ProcessFileTaskHandler` to change the file’s DB record to “completed” after (fake) a short sleep. Have it rename the bucket to “Sample Bucket Name.” Then confirm the API’s `GET /api/buckets` reflects that new name.

6. **Actual File Processing & AI**

   * In `internal/utils/pdf.go`, implement PDF → text extraction (e.g. using `github.com/pdfcpu/pdfcpu`).
   * In `internal/utils/chunk.go`, write a simple loop to break a large string into `[]string` of length \~2000.
   * In the `ProcessFileTaskHandler`, replace the “fake sleep” with calls to:

     1. `ExtractTextFromPDF(path)` (or `os.ReadFile` if it’s plain .txt).
     2. `ChunkText(text, 2000)`.
     3. For each chunk, create `file_chunks` record, call `ai.GetEmbedding(chunk)`, update that row.
     4. Once the first embedding is stored, call `ai.GenerateBucketName(chunks[:3])`. Update `buckets.name`.
     5. Update `files.status = 'completed'`.
   * Confirm that: after upload, the worker processes PDF, stores chunks and embeddings, and bucket name changes.

7. **Quiz Models & Generation**

   * Define `internal/quiz/model.go` with GORM structs for Quiz, Question, Answer, Attempt, AttemptAnswer.
   * Migrate.
   * Implement `internal/quiz/service.go: GenerateQuiz(quizID)`: exactly as outlined in 2.3.4. Use GORM to query all `file_chunks` → retrieve embeddings (vectors). Use pgvector’s cosine distance (via SQL) or do a simple in-memory selection in Go. Then call `ai.GenerateQuestions(...)`, parse JSON to Go structs, insert into DB.
   * In `cmd/worker/main.go`, add `GenerateQuizTaskHandler` that wraps `GenerateQuiz(quizID)`.
   * In API, create `POST /api/buckets/{bucketId}/quizzes` and `GET /api/quizzes/{quizId}` + `GET /api/quizzes/{quizId}/questions`. These handlers will call `EnqueueGenerateQuiz(quizID)` and return quiz status or question payload.
   * Manually test: upload a few files, wait until “completed,” call `POST /api/buckets/1/quizzes`, then check Redis/worker → after a few seconds, DB gets questions. `GET /api/quizzes/1/questions` returns the question list.

8. **Quiz Attempt Endpoints**

   * Implement `POST /api/quizzes/{quizId}/attempts` (insert Attempt + AttemptAnswer rows, compute score).
   * Implement `GET /api/buckets/{bucketId}/attempts` (list all attempts for that bucket by this user).
   * Implement `GET /api/attempts/{attemptId}` (detailed report).
   * Confirm with Postman that this works.

9. **Finalize API Routes & Middleware**

   * Add CORS middleware to allow requests from `http://localhost:4200` (Angular dev server) or the nginx domain.
   * Add request logging middleware (e.g. \[chi/middleware.Logger]).
   * Harden error handling (return proper HTTP codes, format responses as JSON with `{"error":"…"}`).

10. **Frontend Implementation**

    1. **Authentication**

       * Create `auth.service.ts` with methods `login()`, `signup()`, `logout()`. Store JWT in `localStorage`.
       * Create `auth.guard.ts` that checks for token existence/validity (`AuthService.isLoggedIn()`). Protect all routes except `/login`.
       * Build `login.component.ts/html/css`. Form binds to `username`/`password`, calls `AuthService.login()`, navigates to `/buckets` on success.

    2. **App Routing**

       * Define routes:

         ```
         /login            → LoginComponent
         /buckets          → BucketListComponent (default view with sidebar)
         /buckets/:id      → BucketDetailComponent (file status, “Take Quiz” button)
         /buckets/:id/history → ReportHistoryComponent
         /quizzes/:quizId/status → QuizStatusComponent (loading screen)
         /quizzes/:quizId/take   → QuizTakingComponent
         /attempts/:attemptId    → QuizReportComponent
         ```
       * In `app.module.ts`, import `RouterModule.forRoot(...)` with these routes, add `AuthGuard`.

    3. **Navbar & Layout**

       * Build `nav-bar.component.ts/html/css`: top bar with “Logout” button, maybe a small app logo.
       * Overall layout in `app.component.html`:

         ```html
         <nav-bar></nav-bar>
         <div class="main-container">
           <aside class="sidebar">
             <bucket-list></bucket-list>
           </aside>
           <section class="content">
             <router-outlet></router-outlet>
           </section>
         </div>
         ```
       * Add styles so the sidebar is fixed width (200px), and `.content` flexes to fill remaining space.

    4. **Bucket List & Detail**

       * `BucketListComponent` calls `bucketService.getBuckets()` and displays each as a clickable tab. The “+” icon at the top triggers a modal or navigation to `/buckets/new`.
       * `BucketDetailComponent` (activated when route is `/buckets/:id`):

         * On `ngOnInit()`, call `fileService.getFiles(bucketId)`.
         * Display a table: each row = `{filename, statusIcon, deleteButton}`.
         * If no files exist, show `file-upload` component inline with instructions.
         * If at least one file is `'completed'`, enable “Take Quiz” button.

    5. **File Upload**

       * `FileUploadComponent` includes `<input type="file">` or drag-and-drop area. On change, run `fileService.uploadFile(bucketId, file)`.
       * After a successful upload response, immediately call the parent’s `refreshFiles()` to show the new row with `'pending'`.
       * Parent begins polling every 5s: `fileService.getFiles(bucketId)` to update status icons.

    6. **Quiz Settings & Status**

       * When user clicks “Take Quiz,” open `QuizSettingsComponent` (modal or inline).
       * On “Start Quiz,” call `quizService.createQuiz(bucketId, options)`. That returns `{quizId, status}`.
       * Redirect to `/quizzes/{quizId}/status`.

    7. **Quiz Status**

       * `QuizStatusComponent` polls `quizService.getQuizStatus(quizId)` every 3 seconds.
       * Show a spinner + message. Once `status==='ready'`, navigate to `/quizzes/{quizId}/take`.

    8. **Quiz Taking**

       * `QuizTakingComponent`: on `ngOnInit`, call `quizService.getQuizQuestions(quizId)` → `questions: Question[]`.
       * Track user’s answers in a local array: `selectedAnswers: { [questionId:number]: answerId }`.
       * If `timedMode`, start a 30-second countdown per question (using `setInterval`). If time expires, auto-go to next question and store `selectedAnswers[questionId] = null` (or skip).
       * If `practiceMode`, on radio button change, immediately show correctness (green/red) by checking the answer’s `isCorrect` boolean (which the backend must include in the payload). Also reveal `explanation`. Then show a “Next Question” button.
       * If `practiceMode===false`, disable revealing until the end. Show “Next” after they select (but don’t reveal correctness until “Submit Quiz”).
       * At the end, gather `selectedAnswers` into payload, call `quizService.submitQuiz(quizId, selectedAnswers)` → returns `attemptId`. Navigate to `/attempts/{attemptId}`.

    9. **Quiz Report & History**

       * `QuizReportComponent`: fetch from `quizService.getAttempt(attemptId)`, render as described.
       * `ReportHistoryComponent`: fetch `quizService.getAttempts(bucketId)` and list with date/score. Clicking navigates to `/attempts/{attemptId}`.

    10. **Polish**

        * Add a loading spinner component to show whenever HTTP calls are in flight (optional).
        * Add proper error handling & toasts for failures (e.g. file upload fails, or quiz generation error).
        * Style the sidebar tabs so the active bucket is highlighted.
        * Implement route guards so if someone tries `/buckets/999` without permission, redirect to `/buckets`.

---

### 2.9 Putting It All Together

1. **Local Development**

   * Clone the repo.
   * Copy `.env.example` → `.env`, fill in values:

     ```
     DB_USER=quizgenie
     DB_PASSWORD=quizgenie
     DB_NAME=quizgenie
     DB_HOST=postgres
     DB_PORT=5432
     REDIS_ADDR=redis:6379
     JWT_SECRET=supersecretjwt
     OPENAI_API_KEY=<your_key>
     FILE_STORAGE_PATH=/data/uploads
     ```
   * Run `docker-compose up --build` from the monorepo root. This brings up:

     1. **postgres**
     2. **redis**
     3. **api** (on localhost:8080)
     4. **worker** (background jobs)
     5. **frontend** (Angular app on localhost:80)

2. **First-Time Initialization**

   * The API container’s entrypoint registers DB migrations (AutoMigrate).
   * Once API starts, GORM ensures all tables exist.
   * Worker auto-connects to Redis + DB.

3. **Walkthrough**

   1. Open `http://localhost/` → Angular’s login page.
   2. Signup with a new user. On success, Angular stores JWT and navigates to `/buckets`.
   3. The sidebar is empty—“You have no buckets.” A “+ New Bucket” button is prominent.
   4. Click “+ New Bucket” → `FileUploadComponent` opens. Choose a PDF.
   5. Angular calls `POST /api/buckets` with that file. API:

      * Creates bucket row with placeholder name.
      * Creates file row with `status='pending'`.
      * Saves file under `/data/uploads/user_<id>/bucket_<id>/original.pdf`.
      * Enqueues `ProcessFileTask`.
      * Returns `{bucketId, fileId}`.
   6. Angular switches to new bucket’s detail. Sidebar shows tab labeled “(processing…)”.
   7. Meanwhile, worker picks up the `ProcessFileTask`:

      * Changes file status to `processing`.
      * Extracts PDF → text → chunks → embeddings → store.
      * Calls `GenerateBucketName(...)` → e.g. “Introduction to Biology.” Updates bucket’s `name`.
      * Marks file status `='completed'`.
   8. Angular’s polling sees that file is now “completed,” updates the icon to a green check. Sidebar’s “(processing…)” changes to “Introduction to Biology.” “Take Quiz” button becomes enabled.
   9. Click “Take Quiz.” A “Quiz Settings” modal appears. Check “Timed Mode” or “Practice Mode,” click “Start Quiz.”
   10. Angular calls `POST /api/buckets/1/quizzes`, API:

       * Inserts quiz row with `status='pending'`.
       * Enqueues `GenerateQuizTask(quizID)`. Returns `{quizId:2,status:'pending'}`.
       * Frontend navigates to `/quizzes/2/status`. A loading spinner with “Preparing your quiz….”
   11. Worker picks up `GenerateQuizTask(2)`:

       * Updates `status='generating'`.
       * Fetches all processed `file_chunks` for bucket 1, selects a semantic pool + sample.
       * Builds `contextStr`. Calls `Ai.GenerateQuestions(contextStr, count=10, choices=4, difficulty='REGULAR')`.
       * Inserts `questions` + `answers` rows, updates `quizzes.status='ready'`.
   12. Frontend’s `/quizzes/2/status` polling sees `status='ready'`, navigates to `/quizzes/2/take`.
   13. Angular fetches `GET /api/quizzes/2/questions` → array of 10 questions.
   14. User goes through quiz. In practice mode, each answer reveals correctness instantly; in timed mode, each question has a 30s timer.
   15. At the end, Angular calls `POST /api/quizzes/2/attempts` with selected answer IDs. Backend writes attempt + answers, calculates `score=85`, returns `{attemptId:5,score:85}`.
   16. Frontend navigates to `/attempts/5`, fetches detailed report → displays question texts, user’s answer (colored), correct answer, explanations.
   17. User clicks “History” in the sidebar or a “Back to Report History” link, which navigates to `/buckets/1/history`. Angular calls `GET /api/buckets/1/attempts` → displays a list of past attempts.

4. **Ongoing Maintenance & Extensibility**

   * **Adding New Features**: you can iterate on `internal/quiz/service.go` to add new quiz types (e.g. true/false), or `internal/file/service.go` to support additional file formats (DOCX, PPTX).
   * **Scaling Workers**: increase concurrency or spawn multiple `worker` replicas in Docker Compose.
   * **Testing**: write unit tests for each service (e.g. `GenerateBucketName`, `ProcessFile`). Integration tests can spin up a test Postgres + Redis locally.
   * **CI/CD**: add GitHub Actions that build the backend (with `go test ./...`), build the frontend (`npm run build --prod`), then push Docker images to a container registry.
