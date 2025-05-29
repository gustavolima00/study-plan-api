select s.*,
    (
        select JSON_AGG (
                JSON_BUILD_OBJECT (
                    'name',
                    sub.name,
                    'description',
                    sub.description
                )
            )
        from subjects sub
            join session_subjects ss on sub.id = ss.subject_id
        where ss.session_id = s.id
    ) as subjects,
    (
        select JSON_AGG (
                JSON_BUILD_OBJECT (
                    'event_type',
                    e.event_type,
                    'event_time',
                    to_char (
                        e.event_time::timestamp at time zone 'UTC',
                        'YYYY-MM-DD"T"HH24:MI:SS"Z"'
                    )
                )
                order by e.event_time
            )
        from session_events e
        where e.session_id = s.id
    ) as events
from study_sessions s
where s.user_id = $1
    and s.session_state = 'active'