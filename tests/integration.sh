#!/usr/bin/env bash
set -euo pipefail

# -----------------------------------------------------------------------------
# This script will:
#  1) Sign up a new user
#  2) Log in and extract the JWT token
#  3) Create a bucket
#  4) List all buckets to verify
#  5) Upload a dummy file into the bucket
#  6) List all files in the bucket to verify
#  7) Create a quiz for that bucket
#  8) Poll the quiz status until it changes (it will stay "pending" until you implement quiz‐generation logic)
#  9) List attempts for that bucket (should be empty at first)
# 10) (Optional) Once quiz is "ready", fetch questions and submit an (empty) attempt
# 11) List attempts again
#
# Requirements: `curl` and `jq` must be installed on your PATH.
# -----------------------------------------------------------------------------

API="http://localhost:8080"

echo "🔹 1) Signing up a new user..."
signup_resp=$(curl -s -X POST "$API/signup" \
  -H "Content-Type: application/json" \
  -d '{
    "username":"alice",
    "password":"password123",
    "email":"alice@example.com"
  }')
echo "   → signup response: $signup_resp"

echo
echo "🔹 2) Logging in to get JWT..."
login_resp=$(curl -s -X POST "$API/login" \
  -H "Content-Type: application/json" \
  -d '{
    "username":"alice",
    "password":"password123"
  }')
echo "   → login response: $login_resp"
TOKEN=$(echo "$login_resp" | jq -r '.token')
echo "   → extracted token: $TOKEN"

echo
echo "🔹 3) Creating a new bucket..."
create_bucket_resp=$(curl -s -X POST "$API/buckets" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{}')
echo "   → create bucket response: $create_bucket_resp"
BUCKET_ID=$(echo "$create_bucket_resp" | jq -r '.id')
echo "   → new BUCKET_ID = $BUCKET_ID"

echo
echo "🔹 4) Listing all buckets for user..."
list_buckets_resp=$(curl -s -X GET "$API/buckets" \
  -H "Authorization: Bearer $TOKEN")
echo "   → list buckets response: $list_buckets_resp"

echo
echo "🔹 5) Uploading a dummy file into bucket $BUCKET_ID..."
# create a small dummy file
echo "Hello, QuizGenie!" > dummy.txt
upload_file_resp=$(curl -s -X POST "$API/buckets/$BUCKET_ID/files" \
  -H "Authorization: Bearer $TOKEN" \
  -F "file=@dummy.txt" \
  -F "filename=dummy.txt")
echo "   → upload file response: $upload_file_resp"
FILE_ID=$(echo "$upload_file_resp" | jq -r '.fileId')
echo "   → new FILE_ID = $FILE_ID"

echo
echo "🔹 6) Listing files in bucket $BUCKET_ID..."
list_files_resp=$(curl -s -X GET "$API/buckets/$BUCKET_ID/files" \
  -H "Authorization: Bearer $TOKEN")
echo "   → list files response: $list_files_resp"

echo
echo "🔹 7) Creating a quiz for bucket $BUCKET_ID..."
create_quiz_resp=$(curl -s -X POST "$API/buckets/$BUCKET_ID/quizzes" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "timedMode":false,
    "practiceMode":true
  }')
echo "   → create quiz response: $create_quiz_resp"
QUIZ_ID=$(echo "$create_quiz_resp" | jq -r '.quizId')
echo "   → new QUIZ_ID = $QUIZ_ID"

echo
echo "🔹 8) Polling quiz status for quiz $QUIZ_ID until 'ready' (or until timeout)..."
attempts=0
while true; do
  status_resp=$(curl -s -X GET "$API/quizzes/$QUIZ_ID" \
    -H "Authorization: Bearer $TOKEN")
  STATUS=$(echo "$status_resp" | jq -r '.status')
  echo "   → status check #$((++attempts)): $STATUS"
  if [[ "$STATUS" == "ready" ]]; then
    echo "   → quiz is ready!"
    break
  fi
  if (( attempts >= 10 )); then
    echo "   → timed out waiting for quiz to become 'ready'. Continuing anyway."
    break
  fi
  sleep 2
done

echo
echo "🔹 9) Listing attempts for bucket $BUCKET_ID (should be empty initially)..."
list_attempts_resp=$(curl -s -X GET "$API/buckets/$BUCKET_ID/attempts" \
  -H "Authorization: Bearer $TOKEN")
echo "   → list attempts response: $list_attempts_resp"

echo
echo "🔹 10) (Optional) If quiz is 'ready', fetch questions and submit a dummy attempt..."
if [[ "$STATUS" == "ready" ]]; then
  echo "   → Fetching questions for quiz $QUIZ_ID..."
  questions_resp=$(curl -s -X GET "$API/quizzes/$QUIZ_ID/questions" \
    -H "Authorization: Bearer $TOKEN")
  echo "   → questions response: $questions_resp"

  # Suppose we pick the first question & first answer to submit.
  QID=$(echo "$questions_resp" | jq -r '.[0].questionId')
  AID=$(echo "$questions_resp" | jq -r '.[0].answers[0].id')

  echo "   → Submitting one answer (QID=$QID, AID=$AID)..."
  submit_resp=$(curl -s -X POST "$API/quizzes/$QUIZ_ID/attempts" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $TOKEN" \
    -d "{
      \"answers\":[ { \"questionId\": $QID, \"answerId\": $AID } ]
    }")
  echo "   → submit attempt response: $submit_resp"
  ATTEMPT_ID=$(echo "$submit_resp" | jq -r '.attemptId')
  echo "   → new ATTEMPT_ID = $ATTEMPT_ID"

  echo
  echo "   → 11) Listing attempts again..."
  list_attempts_again=$(curl -s -X GET "$API/buckets/$BUCKET_ID/attempts" \
    -H "Authorization: Bearer $TOKEN")
  echo "      → list attempts now: $list_attempts_again"

  echo
  echo "   → 12) Fetching attempt details for attempt $ATTEMPT_ID..."
  details_resp=$(curl -s -X GET "$API/attempts/$ATTEMPT_ID" \
    -H "Authorization: Bearer $TOKEN")
  echo "      → attempt details: $details_resp"
else
  echo "   → Quiz never reached 'ready'; skipping questions/submit steps."
fi

echo
echo "✅ All done!"
