## table structure for storing todos

CREATE TABLE IF NOT EXISTS todos (
    id UUID,
    title String,
    desc String,
    status String,
    created_at DateTime,
    effort_hours Int32
) ENGINE = MergeTree()
ORDER BY (id);
