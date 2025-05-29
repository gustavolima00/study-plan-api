INSERT INTO study_sessions (user_id, title, notes, date, session_state)
VALUES (:user_id, :title, :notes, :date, :session_state)
RETURNING *