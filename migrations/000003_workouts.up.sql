CREATE TABLE workouts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id),
    date DATE NOT NULL,
    weekday TEXT DEFAULT '',
    comment TEXT DEFAULT '',
    is_deleted BOOLEAN DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_workouts_user_id ON workouts(user_id);

CREATE TABLE workout_exercises (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    workout_id UUID NOT NULL REFERENCES workouts(id),
    exercise_id UUID,
    name TEXT NOT NULL,
    sort_order INT DEFAULT 0,
    comment TEXT DEFAULT '',
    is_deleted BOOLEAN DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_workout_exercises_workout_id ON workout_exercises(workout_id);

CREATE TABLE workout_sets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    workout_exercise_id UUID NOT NULL REFERENCES workout_exercises(id),
    set_number INT NOT NULL,
    weight NUMERIC(10,2) DEFAULT 0,
    reps INT DEFAULT 10,
    to_failure BOOLEAN DEFAULT false,
    is_deleted BOOLEAN DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_workout_sets_workout_exercise_id ON workout_sets(workout_exercise_id);
