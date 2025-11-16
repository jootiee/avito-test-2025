-- Create pull_requests table
CREATE TABLE IF NOT EXISTS pull_requests (
    pull_request_id VARCHAR(255) PRIMARY KEY,
    pull_request_name VARCHAR(500) NOT NULL,
    author_id VARCHAR(255) NOT NULL REFERENCES users(user_id),
    status VARCHAR(20) NOT NULL CHECK (status IN ('OPEN', 'MERGED')),
    assigned_reviewers JSONB NOT NULL DEFAULT '[]'::jsonb,
    created_at TIMESTAMP DEFAULT NOW(),
    merged_at TIMESTAMP NULL
);

CREATE INDEX idx_prs_author ON pull_requests(author_id);
CREATE INDEX idx_prs_reviewers ON pull_requests USING GIN (assigned_reviewers);
