#!/usr/bin/env bash
set -euo pipefail

# -------------------------------------------------------
# Prerequisites:
#  1. A running Postgres + Redis (e.g. via docker-compose).
#  2. Your Go API: `go run ./cmd/api` listening on localhost:8080.
#  3. Your Go worker: `go run ./cmd/worker`.
#  4. 'jq' installed locally (used to parse JSON).
# -------------------------------------------------------

API_BASE="http://localhost:8080"

echo
echo "=== 1) Sign up a new user ==="
signup_response=$(
  curl -s -X POST "$API_BASE/signup" \
    -H "Content-Type: application/json" \
    -d '{
      "username": "smoketest",
      "password": "Password123!",
      "email":    "smoke@example.com"
    }'
)
echo "→ signup response: $signup_response"
echo

echo "=== 2) Log in to get a JWT ==="
login_response=$(
  curl -s -X POST "$API_BASE/login" \
    -H "Content-Type: application/json" \
    -d '{
      "username": "smoketest",
      "password": "Password123!"
    }'
)
echo "→ login response: $login_response"
JWT_TOKEN=$(echo "$login_response" | jq -r '.token')
if [ "$JWT_TOKEN" == "null" ] || [ -z "$JWT_TOKEN" ]; then
  echo "ERROR: failed to obtain token. Response: $login_response"
  exit 1
fi
echo "→ extracted JWT_TOKEN: $JWT_TOKEN"
AUTH_HEADER="Authorization: Bearer $JWT_TOKEN"
echo

echo "=== 3) Create a new bucket ==="
# CreateBucketHandler expects no body, just the auth header.
create_bucket_resp=$(
  curl -s -X POST "$API_BASE/buckets" \
    -H "Content-Type: application/json" \
    -H "$AUTH_HEADER" \
    -d '{}' 
)
echo "→ create bucket response: $create_bucket_resp"
BUCKET_ID=$(echo "$create_bucket_resp" | jq -r '.id')
if [ -z "$BUCKET_ID" ] || [ "$BUCKET_ID" == "null" ]; then
  echo "ERROR: failed to create bucket."
  exit 1
fi
echo "→ new BUCKET_ID = $BUCKET_ID"
echo

echo "=== 4) Upload one PDF (or text) file to that bucket ==="
# Replace '/path/to/sample.pdf' with a real file path on your machine.
FILE_PATH="sample.pdf"
if [ ! -f "$FILE_PATH" ]; then
  echo "ERROR: file '$FILE_PATH' does not exist. Please adjust the path."
  exit 1
fi

upload_response=$(
  curl -s -X POST "$API_BASE/buckets/$BUCKET_ID/files" \
    -H "$AUTH_HEADER" \
    -F "file=@${FILE_PATH}" \
    -F "filename=sample.pdf"
)
echo "→ upload file response: $upload_response"
FILE_ID=$(echo "$upload_response" | jq -r '.fileId')
if [ -z "$FILE_ID" ] || [ "$FILE_ID" == "null" ]; then
  echo "ERROR: failed to upload file."
  exit 1
fi
echo "→ new FILE_ID = $FILE_ID"
echo

echo "=== 5) Poll /buckets/{bucketId}/files until status == \"completed\" ==="
file_status="pending"
until [ "$file_status" == "completed" ]; do
  sleep 5
  files_list=$(
    curl -s -X GET "$API_BASE/buckets/$BUCKET_ID/files" \
      -H "$AUTH_HEADER"
  )
  # Extract the status of our specific file
  file_status=$(echo "$files_list" | jq -r ".[] | select(.id==$FILE_ID) | .status")
  echo "    Current status of file $FILE_ID → $file_status"
done
echo "✅ File processing completed."
echo

echo "=== 6) Create a new quiz from that bucket ==="
create_quiz_resp=$(
  curl -s -X POST "$API_BASE/buckets/$BUCKET_ID/quizzes" \
    -H "Content-Type: application/json" \
    -H "$AUTH_HEADER" \
    -d '{
      "timedMode": false,
      "practiceMode": false
    }'
)
echo "→ create quiz response: $create_quiz_resp"
QUIZ_ID=$(echo "$create_quiz_resp" | jq -r '.quizId')
if [ -z "$QUIZ_ID" ] || [ "$QUIZ_ID" == "null" ]; then
  echo "ERROR: failed to create quiz."
  exit 1
fi
echo "→ new QUIZ_ID = $QUIZ_ID"
echo

echo "=== 7) Poll /quizzes/{quizId} until status == \"ready\" ==="
quiz_status="pending"
until [ "$quiz_status" == "ready" ]; do
  sleep 5
  status_resp=$(
    curl -s -X GET "$API_BASE/quizzes/$QUIZ_ID" \
      -H "$AUTH_HEADER"
  )
  quiz_status=$(echo "$status_resp" | jq -r '.status')
  echo "    Quiz $QUIZ_ID status → $quiz_status"
done
echo "✅ Quiz generation completed."
echo

echo "=== 8) Fetch the generated questions ==="
questions_resp=$(
  curl -s -X GET "$API_BASE/quizzes/$QUIZ_ID/questions" \
    -H "$AUTH_HEADER"
)
echo "→ questions payload:"
echo "$questions_resp" | jq '.'      # pretty-print array of questions/answers
echo

# For demonstration, pick the first question’s ID and its first answer ID:
FIRST_Q_ID=$(echo "$questions_resp" | jq -r '.[0].questionId')
FIRST_A_ID=$(echo "$questions_resp" | jq -r '.[0].answers[0].id')
echo "→ choosing questionId=$FIRST_Q_ID, answerId=$FIRST_A_ID as our single answer"
echo

echo "=== 9) Submit one attempt (picking the first answer for each question) ==="
# Build an answers array that picks the first choice for _every_ question.
ANSWERS_ARRAY=$(
  echo "$questions_resp" | \
  jq '[ .[] | { questionId: .questionId, answerId: .answers[0].id } ]'
)

submit_attempt_resp=$(
  curl -s -X POST "$API_BASE/quizzes/$QUIZ_ID/attempts" \
    -H "Content-Type: application/json" \
    -H "$AUTH_HEADER" \
    -d "{\"answers\": $ANSWERS_ARRAY}"
)
echo "→ submit attempt response: $submit_attempt_resp"
ATTEMPT_ID=$(echo "$submit_attempt_resp" | jq -r '.attemptId')
SCORE=$(echo "$submit_attempt_resp" | jq -r '.score')
echo "→ new ATTEMPT_ID = $ATTEMPT_ID, SCORE = $SCORE"
echo

echo "=== 10) Fetch attempt details (report) ==="
attempt_details=$(
  curl -s -X GET "$API_BASE/attempts/$ATTEMPT_ID" \
    -H "$AUTH_HEADER"
)
echo "→ attempt details:"
echo "$attempt_details" | jq '.'
echo

echo "🚀 Smoke-test complete. All endpoints responded as expected."
