CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE study_sessions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL,
    title VARCHAR(100),
    notes TEXT,
    date DATE NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    session_state VARCHAR(100) NOT NULL DEFAULT ''
);
CREATE TABLE session_events (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    session_id UUID NOT NULL REFERENCES study_sessions (id) ON DELETE CASCADE,
    event_type VARCHAR(100) NOT NULL DEFAULT '',
    event_time TIMESTAMP NOT NULL
);