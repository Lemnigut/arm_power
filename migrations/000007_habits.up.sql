CREATE TABLE habits (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id),
    title TEXT NOT NULL,
    color TEXT NOT NULL DEFAULT '#E0A018',
    repeat_type TEXT NOT NULL DEFAULT 'daily',
    repeat_mode TEXT NOT NULL DEFAULT '',
    repeat_weekday INT,
    is_deleted BOOLEAN DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_habits_user_id ON habits(user_id);

CREATE TABLE habit_completions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    habit_id UUID NOT NULL REFERENCES habits(id),
    date DATE NOT NULL,
    is_deleted BOOLEAN DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_habit_completions_habit_id ON habit_completions(habit_id);
CREATE UNIQUE INDEX idx_habit_completions_habit_date ON habit_completions(habit_id, date);
