
-- Keys Table
CREATE TABLE keys (
    id UUID PRIMARY KEY NOT NULL,
    value text NOT NULL,

    created_at timestamp default current_timestamp,
    updated_at timestamp default null,
    expires_at timestamp default null
);
