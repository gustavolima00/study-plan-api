create
or replace VIEW study_session_view as (
    select
        s.*,
        (
            select
                JSON_AGG (
                    JSON_BUILD_OBJECT ('name', sub.name, 'description', sub.description)
                )
            from
                subjects sub
                join session_subjects ss on sub.id = ss.subject_id
            where
                ss.session_id = s.id
        ) as subjects,
        (
            select
                JSON_AGG (
                    JSON_BUILD_OBJECT (
                        'event_type',
                        e.event_type,
                        'event_time',
                        e.event_time,
                        'device_info',
                        e.device_info
                    )
                    order by
                        e.event_time
                )
            from
                session_events e
            where
                e.session_id = s.id
        ) as events
    from
        study_sessions s
);