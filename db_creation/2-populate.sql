INSERT INTO notification_service.rate_limit_rules (notification_type, max_count, duration)
VALUES
    ('STATUS', 2, EXTRACT(EPOCH FROM INTERVAL '1 minute')),
    ('NEWS', 1, EXTRACT(EPOCH FROM INTERVAL '1 day')),
    ('MARKETING', 3, EXTRACT(EPOCH FROM INTERVAL '1 hour'));

