
-- Keys Table
CREATE TABLE keys (
    id VARCHAR(50) PRIMARY KEY NOT NULL,
    value text NOT NULL,

    created_at timestamp default current_timestamp,
    updated_at timestamp default null,
    expires_at timestamp default null
);
