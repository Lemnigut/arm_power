CREATE TABLE exercise_media (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    exercise_id UUID NOT NULL REFERENCES exercises(id),
    user_id UUID NOT NULL REFERENCES users(id),
    media_type TEXT NOT NULL,
    s3_key TEXT NOT NULL,
    thumbnail_s3_key TEXT,
    original_name TEXT,
    size_bytes BIGINT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_exercise_media_exercise_id ON exercise_media(exercise_id);
