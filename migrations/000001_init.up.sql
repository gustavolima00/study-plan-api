CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE
    subjects (
        id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
        user_id UUID,
        name VARCHAR(100) NOT NULL,
        description TEXT,
        created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        is_custom BOOLEAN NOT NULL DEFAULT TRUE,
        CONSTRAINT unique_subject_name_per_user UNIQUE (user_id, name)
    );

CREATE TABLE
    subject_relations (
        parent_id UUID NOT NULL REFERENCES subjects (id) ON DELETE CASCADE,
        child_id UUID NOT NULL REFERENCES subjects (id) ON DELETE CASCADE,
        PRIMARY KEY (parent_id, child_id),
        CHECK (parent_id != child_id)
    );

CREATE TABLE
    study_sessions (
        id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
        user_id UUID NOT NULL,
        title VARCHAR(100),
        notes TEXT,
        date DATE NOT NULL,
        created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        session_state SMALLINT NOT NULL DEFAULT 0
    );

CREATE TABLE
    session_events (
        id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
        session_id UUID NOT NULL REFERENCES study_sessions (id) ON DELETE CASCADE,
        event_type SMALLINT NOT NULL DEFAULT 0,
        event_time TIMESTAMP NOT NULL,
        device_info VARCHAR(100)
    );

CREATE TABLE
    session_subjects (
        session_id UUID NOT NULL REFERENCES study_sessions (id) ON DELETE CASCADE,
        subject_id UUID NOT NULL REFERENCES subjects (id) ON DELETE CASCADE,
        PRIMARY KEY (session_id, subject_id)
    );