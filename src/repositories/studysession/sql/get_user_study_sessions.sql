SELECT s.*,
    COALESCE(
        (
            SELECT JSON_AGG(
                    JSON_BUILD_OBJECT(
                        'name',
                        sub.name,
                        'description',
                        sub.description
                    )
                )
            FROM subjects sub
                JOIN session_subjects ss ON sub.id = ss.subject_id
            WHERE ss.session_id = s.id
        ),
        '[]'::json
    ) AS subjects,
    COALESCE(
        (
            SELECT JSON_AGG(
                    JSON_BUILD_OBJECT(
                        'event_type',
                        e.event_type,
                        'event_time',
                        to_char(
                            e.event_time::timestamp AT TIME ZONE 'UTC',
                            'YYYY-MM-DD"T"HH24:MI:SS"Z"'
                        )
                    )
                    ORDER BY e.event_time
                )
            FROM session_events e
            WHERE e.session_id = s.id
        ),
        '[]'::json
    ) AS events
FROM study_sessions s
WHERE s.user_id = $1