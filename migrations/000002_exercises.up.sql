CREATE TABLE exercises (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id),
    name TEXT NOT NULL,
    muscles TEXT[] DEFAULT '{}',
    category TEXT DEFAULT '',
    description TEXT DEFAULT '',
    youtube_links TEXT[] DEFAULT '{}',
    is_deleted BOOLEAN DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_exercises_user_id ON exercises(user_id);

CREATE TABLE exercise_comments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    exercise_id UUID NOT NULL REFERENCES exercises(id),
    user_id UUID NOT NULL REFERENCES users(id),
    text TEXT NOT NULL,
    is_deleted BOOLEAN DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_exercise_comments_exercise_id ON exercise_comments(exercise_id);
