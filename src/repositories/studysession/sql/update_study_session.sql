UPDATE study_sessions
SET title = :title,
    notes = :notes,
    session_state = :session_state,
    updated_at = CURRENT_TIMESTAMP
WHERE id = :id
RETURNING *