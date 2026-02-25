#!/bin/bash
set -e

BASE="http://localhost:8080"

echo "=== REGISTER ==="
TOKENS=$(curl -s -X POST $BASE/auth/register -H 'Content-Type: application/json' \
  -d '{"username":"demo","password":"demo123"}')
echo "$TOKENS"

ACCESS=$(echo "$TOKENS" | python3 -c 'import sys,json; print(json.load(sys.stdin)["accessToken"])')
AUTH="Authorization: Bearer $ACCESS"
CT="Content-Type: application/json"

echo ""
echo "=== CREATE EXERCISES ==="

create_ex() {
  local data="$1"
  local result=$(curl -s -X POST $BASE/api/exercises -H "$AUTH" -H "$CT" -d "$data")
  local id=$(echo "$result" | python3 -c 'import sys,json; print(json.load(sys.stdin)["id"])')
  echo "$id"
}

EX1=$(create_ex '{"name":"Бицепс блок","muscles":["Бицепс","Предплечье"],"category":"Блок","description":"Сгибание руки на блочном тренажёре. Локоть фиксирован, тянем рукоять к плечу.","youtubeLinks":["https://youtube.com/watch?v=ex1"]}')
echo "Бицепс блок: $EX1"

EX2=$(create_ex '{"name":"Пулемёт блок","muscles":["Спина","Бицепс"],"category":"Блок","description":"Тяга двумя руками на нижнем блоке со скручиванием в конце."}')
echo "Пулемёт: $EX2"

EX3=$(create_ex '{"name":"Вертикальный блок тяга","muscles":["Спина","Бицепс"],"category":"Блок","description":"Тянем к груди, лопатки сводим."}')
echo "Вертикальный блок: $EX3"

EX4=$(create_ex '{"name":"Пронатор обратный + резина","muscles":["Предплечье","Пронатор"],"category":"Резина","description":"Обратная пронация + добивка резиной в отказ."}')
echo "Пронатор: $EX4"

EX5=$(create_ex '{"name":"Кисть пояс петля","muscles":["Кисть","Предплечье"],"category":"Свободный вес","description":"Работа на кисть с поясом и петлёй."}')
echo "Кисть: $EX5"

EX6=$(create_ex '{"name":"Гиперэкстензия","muscles":["Поясница","Ягодицы"],"category":"Свободный вес","description":"Разгибание спины на тренажёре."}')
echo "Гиперэкстензия: $EX6"

EX7=$(create_ex '{"name":"Бицепс узкий хват изогнутый гриф","muscles":["Бицепс"],"category":"Свободный вес","description":"Гриф 10кг + блины. Многоповторный режим."}')
echo "Бицепс узкий хват: $EX7"

EX8=$(create_ex '{"name":"Боковое давление","muscles":["Предплечье","Плечо"],"category":"Армрестлинг","description":"Боковое давление на столе или с резиной."}')
echo "Боковое давление: $EX8"

EX9=$(create_ex '{"name":"Крюк на столе","muscles":["Бицепс","Кисть"],"category":"Армрестлинг","description":"Кисть заворачиваем внутрь, тянем на себя.","youtubeLinks":["https://youtube.com/watch?v=ex2"]}')
echo "Крюк на столе: $EX9"

# Comments
curl -s -X POST "$BASE/api/exercises/$EX1/comments" -H "$AUTH" -H "$CT" -d '{"text":"Лучше с канатной рукояткой"}' > /dev/null
curl -s -X POST "$BASE/api/exercises/$EX2/comments" -H "$AUTH" -H "$CT" -d '{"text":"Контролировать негативную фазу"}' > /dev/null
curl -s -X POST "$BASE/api/exercises/$EX9/comments" -H "$AUTH" -H "$CT" -d '{"text":"Локоть строго на подушке"}' > /dev/null
echo "Comments added"

echo ""
echo "=== CREATE WORKOUTS ==="

add_ex_to_wk() {
  local wk_id="$1"
  local ex_id="$2"
  local name="$3"
  local comment="$4"
  local result=$(curl -s -X POST "$BASE/api/workouts/$wk_id/exercises" -H "$AUTH" -H "$CT" \
    -d "{\"exerciseId\":\"$ex_id\",\"name\":\"$name\",\"comment\":\"$comment\"}")
  echo "$result" | python3 -c 'import sys,json; exs=json.load(sys.stdin)["exercises"]; print(exs[-1]["id"])'
}

add_set() {
  local wk_id="$1"
  local we_id="$2"
  local w="$3"
  local r="$4"
  local f="$5"
  curl -s -X POST "$BASE/api/workouts/$wk_id/exercises/$we_id/sets" -H "$AUTH" -H "$CT" \
    -d "{\"weight\":$w,\"reps\":$r,\"toFailure\":$f}" > /dev/null
}

# Workout 1: Jan 28
WK1=$(curl -s -X POST "$BASE/api/workouts" -H "$AUTH" -H "$CT" \
  -d '{"date":"2025-01-28","weekday":"Среда","comment":"Хорошая тренировка"}')
WK1_ID=$(echo "$WK1" | python3 -c 'import sys,json; print(json.load(sys.stdin)["id"])')
echo "Тренировка 28 янв: $WK1_ID"

WE=$(add_ex_to_wk "$WK1_ID" "$EX1" "Бицепс блок" "Правая нога вперёд")
add_set "$WK1_ID" "$WE" 15 15 false
add_set "$WK1_ID" "$WE" 20 15 false
add_set "$WK1_ID" "$WE" 25 12 true
echo "  Бицепс блок: 3 подхода"

WE=$(add_ex_to_wk "$WK1_ID" "$EX2" "Пулемёт блок" "Колени подать, руки на них")
add_set "$WK1_ID" "$WE" 65 15 false
add_set "$WK1_ID" "$WE" 75 15 false
add_set "$WK1_ID" "$WE" 85 17 true
echo "  Пулемёт: 3 подхода"

WE=$(add_ex_to_wk "$WK1_ID" "$EX3" "Вертикальный блок тяга" "")
add_set "$WK1_ID" "$WE" 60 10 false
add_set "$WK1_ID" "$WE" 75 8 false
add_set "$WK1_ID" "$WE" 80 6 true
echo "  Вертикальный блок: 3 подхода"

# Workout 2: Jan 31
WK2=$(curl -s -X POST "$BASE/api/workouts" -H "$AUTH" -H "$CT" \
  -d '{"date":"2025-01-31","weekday":"Суббота"}')
WK2_ID=$(echo "$WK2" | python3 -c 'import sys,json; print(json.load(sys.stdin)["id"])')
echo ""
echo "Тренировка 31 янв: $WK2_ID"

WE=$(add_ex_to_wk "$WK2_ID" "$EX4" "Пронатор обратный + резина" "")
add_set "$WK2_ID" "$WE" 10 20 false
add_set "$WK2_ID" "$WE" 15 15 false
add_set "$WK2_ID" "$WE" 20 10 true
echo "  Пронатор: 3 подхода"

WE=$(add_ex_to_wk "$WK2_ID" "$EX5" "Кисть пояс петля" "На красной лавке")
add_set "$WK2_ID" "$WE" 10 10 false
add_set "$WK2_ID" "$WE" 15 8 false
add_set "$WK2_ID" "$WE" 10 25 true
echo "  Кисть: 3 подхода"

WE=$(add_ex_to_wk "$WK2_ID" "$EX6" "Гиперэкстензия" "")
add_set "$WK2_ID" "$WE" 0 10 false
add_set "$WK2_ID" "$WE" 35 10 false
add_set "$WK2_ID" "$WE" 45 10 false
echo "  Гиперэкстензия: 3 подхода"

# Workout 3: Feb 4
WK3=$(curl -s -X POST "$BASE/api/workouts" -H "$AUTH" -H "$CT" \
  -d '{"date":"2025-02-04","weekday":"Вторник","comment":"Лёгкая после соревнований"}')
WK3_ID=$(echo "$WK3" | python3 -c 'import sys,json; print(json.load(sys.stdin)["id"])')
echo ""
echo "Тренировка 4 фев: $WK3_ID"

WE=$(add_ex_to_wk "$WK3_ID" "$EX1" "Бицепс блок" "")
add_set "$WK3_ID" "$WE" 15 15 false
add_set "$WK3_ID" "$WE" 22 12 false
add_set "$WK3_ID" "$WE" 27 10 true
echo "  Бицепс блок: 3 подхода"

WE=$(add_ex_to_wk "$WK3_ID" "$EX7" "Бицепс узкий хват изогнутый гриф" "Гриф 10кг")
add_set "$WK3_ID" "$WE" 30 20 false
add_set "$WK3_ID" "$WE" 30 15 false
add_set "$WK3_ID" "$WE" 30 17 true
echo "  Бицепс узкий хват: 3 подхода"

# Workout 4: Feb 7
WK4=$(curl -s -X POST "$BASE/api/workouts" -H "$AUTH" -H "$CT" \
  -d '{"date":"2025-02-07","weekday":"Пятница"}')
WK4_ID=$(echo "$WK4" | python3 -c 'import sys,json; print(json.load(sys.stdin)["id"])')
echo ""
echo "Тренировка 7 фев: $WK4_ID"

WE=$(add_ex_to_wk "$WK4_ID" "$EX2" "Пулемёт блок" "")
add_set "$WK4_ID" "$WE" 70 15 false
add_set "$WK4_ID" "$WE" 80 12 false
add_set "$WK4_ID" "$WE" 90 14 true
echo "  Пулемёт: 3 подхода"

WE=$(add_ex_to_wk "$WK4_ID" "$EX3" "Вертикальный блок тяга" "")
add_set "$WK4_ID" "$WE" 65 10 false
add_set "$WK4_ID" "$WE" 80 8 false
add_set "$WK4_ID" "$WE" 85 5 true
echo "  Вертикальный блок: 3 подхода"

# Test copy
echo ""
echo "=== COPY WORKOUT ==="
COPY=$(curl -s -X POST "$BASE/api/workouts/$WK1_ID/copy" -H "$AUTH")
COPY_ID=$(echo "$COPY" | python3 -c 'import sys,json; print(json.load(sys.stdin)["id"])')
echo "Копия тренировки 28 янв: $COPY_ID"

echo ""
echo "=== VERIFY: GET ALL WORKOUTS ==="
curl -s "$BASE/api/workouts" -H "$AUTH" | python3 -c '
import sys, json
wks = json.load(sys.stdin)
print(f"Всего тренировок: {len(wks)}")
for w in wks:
    exs = w.get("exercises", [])
    total_sets = sum(len(e.get("sets",[])) for e in exs)
    print(f"  {w[\"date\"][:10]} {w[\"weekday\"]} — {len(exs)} упр., {total_sets} подх. | {w[\"comment\"]}")
'

echo ""
echo "=== VERIFY: GET ALL EXERCISES ==="
curl -s "$BASE/api/exercises" -H "$AUTH" | python3 -c '
import sys, json
exs = json.load(sys.stdin)
print(f"Всего упражнений: {len(exs)}")
for e in exs:
    print(f"  {e[\"name\"]} [{e[\"category\"]}] — {e[\"muscles\"]}")
'

echo ""
echo "=== ALL DONE ==="
echo "Credentials: demo / demo123"
echo "Access token (first 50 chars): ${ACCESS:0:50}..."
