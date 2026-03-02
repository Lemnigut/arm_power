-- exercises: add is_single_hand
ALTER TABLE exercises ADD COLUMN is_single_hand BOOLEAN NOT NULL DEFAULT false;

-- workouts: add workout_type (gym/home)
ALTER TABLE workouts ADD COLUMN workout_type TEXT NOT NULL DEFAULT 'gym';

-- workout_exercises: add snapshot fields
ALTER TABLE workout_exercises ADD COLUMN is_single_hand BOOLEAN NOT NULL DEFAULT false;
ALTER TABLE workout_exercises ADD COLUMN weight_unit TEXT NOT NULL DEFAULT 'kg';

-- workout_sets: add hand for single-hand exercises
ALTER TABLE workout_sets ADD COLUMN hand TEXT NOT NULL DEFAULT 'right';
