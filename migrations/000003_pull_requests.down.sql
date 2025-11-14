DROP TABLE IF EXISTS pull_requests PURGE;
    created_at TIMESTAMP DEFAULT NOW(),
    merged_at TIMESTAMP NULL
);

DROP INDEX IF EXISTS idx_prs_author;
DROP INDEX IF EXISTS idx_prs_status;
DROP INDEX IF EXISTS idx_prs_reviewers;