INSERT INTO apps (id, name, secret)
VALUES (1, 'test', 'test-secret', 1)
ON CONFLICT DO NOTHING;