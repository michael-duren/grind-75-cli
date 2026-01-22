-- name: CreateProblem :exec
INSERT INTO problems (id, slug, title, url, duration, difficulty_id)
VALUES (?, ?, ?, ?, ?, ?)
ON CONFLICT(slug) DO UPDATE SET
    title = excluded.title,
    url = excluded.url,
    duration = excluded.duration,
    difficulty_id = excluded.difficulty_id;

-- name: GetProblem :one
SELECT * FROM problems
WHERE id = ? LIMIT 1;

-- name: GetProblemBySlug :one
SELECT * FROM problems
WHERE slug = ? LIMIT 1;

-- name: ListProblems :many
SELECT * FROM problems
ORDER BY id;

-- name: CreateTopic :exec
INSERT INTO topics (id, name)
VALUES (?, ?)
ON CONFLICT(id) DO NOTHING;

-- name: LinkProblemTopic :exec
INSERT INTO problem_topics (problem_id, topic_id)
VALUES (?, ?)
ON CONFLICT(problem_id, topic_id) DO NOTHING;

-- name: GetProblemTopics :many
SELECT t.* FROM topics t
JOIN problem_topics pt ON t.id = pt.topic_id
WHERE pt.problem_id = ?;

-- name: UpsertUserProgress :exec
INSERT INTO user_progress (problem_id, status, last_attempted_at, attempts, updated_at)
VALUES (?, ?, ?, ?, CURRENT_TIMESTAMP)
ON CONFLICT(problem_id) DO UPDATE SET
    status = excluded.status,
    last_attempted_at = excluded.last_attempted_at,
    attempts = excluded.attempts,
    updated_at = CURRENT_TIMESTAMP;

-- name: GetUserProgress :one
SELECT * FROM user_progress
WHERE problem_id = ? LIMIT 1;

-- name: CreateReview :exec
INSERT INTO reviews (problem_id, review_date)
VALUES (?, ?);

-- name: ListPendingReviews :many
SELECT r.*, p.title, p.slug FROM reviews r
JOIN problems p ON r.problem_id = p.id
WHERE r.completed = 0 AND r.review_date <= ?
ORDER BY r.review_date;

-- name: CompleteReview :exec
UPDATE reviews
SET completed = 1
WHERE id = ?;

-- name: GetDifficulty :one
SELECT * FROM difficulty_levels
WHERE id = ?;

-- name: ListDifficulties :many
SELECT * FROM difficulty_levels;

-- name: ListTopics :many
SELECT * FROM topics;
