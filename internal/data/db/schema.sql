-- Enable foreign keys
PRAGMA foreign_keys = ON;

CREATE TABLE difficulty_levels (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL UNIQUE
);

INSERT INTO difficulty_levels (id, name) VALUES 
('Easy', 'Easy'),
('Medium', 'Medium'),
('Hard', 'Hard');

CREATE TABLE topics (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL UNIQUE
);

CREATE TABLE problems (
    id INTEGER PRIMARY KEY,
    slug TEXT NOT NULL UNIQUE,
    title TEXT NOT NULL,
    url TEXT NOT NULL,
    duration INTEGER NOT NULL,
    difficulty_id TEXT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (difficulty_id) REFERENCES difficulty_levels(id)
);

CREATE TABLE problem_topics (
    problem_id INTEGER NOT NULL,
    topic_id TEXT NOT NULL,
    PRIMARY KEY (problem_id, topic_id),
    FOREIGN KEY (problem_id) REFERENCES problems(id) ON DELETE CASCADE,
    FOREIGN KEY (topic_id) REFERENCES topics(id) ON DELETE CASCADE
);

CREATE TABLE user_progress (
    problem_id INTEGER NOT NULL PRIMARY KEY,
    status TEXT NOT NULL CHECK(status IN ('New', 'Completed', 'Struggling')),
    last_attempted_at DATETIME,
    attempts INTEGER NOT NULL DEFAULT 0,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (problem_id) REFERENCES problems(id) ON DELETE CASCADE
);

CREATE TABLE reviews (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    problem_id INTEGER NOT NULL,
    review_date DATETIME NOT NULL,
    completed BOOLEAN NOT NULL DEFAULT 0,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (problem_id) REFERENCES problems(id) ON DELETE CASCADE
);

CREATE TABLE reminders (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    last_sent_at DATETIME NOT NULL,
    type TEXT NOT NULL
);

-- Indexes
CREATE INDEX idx_problems_difficulty ON problems(difficulty_id);
CREATE INDEX idx_user_progress_status ON user_progress(status);
CREATE INDEX idx_reviews_review_date ON reviews(review_date);
CREATE INDEX idx_reviews_completed ON reviews(completed);
