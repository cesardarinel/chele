## 1. Model & Migration

- [x] 1.1 Add `direction` field to `Schedule` model (`CharField` with choices `income`/`expense`, default `expense`)
- [x] 1.2 Run `python manage.py makemigrations` and `python manage.py migrate`

## 2. Backend Logic

- [x] 2.1 Update `process_due_schedules` to check `s.direction`: if `income` use `+amount` and add to balance, if `expense` use `-amount` and subtract
- [x] 2.2 Update `schedule_create` view to save `direction` from POST data
- [x] 2.3 Update `schedule_edit` view to save `direction` from POST data

## 3. Frontend

- [x] 3.1 Add direction toggle (radio buttons or select) to `schedule_form.html` between cuenta and monto fields
- [x] 3.2 Update `schedule_list.html` to show income in green with `+` and expense in red with `-`

## 4. Tests

- [x] 4.1 Add test for creating an income schedule
- [x] 4.2 Add test for income schedule execution (transaction amount positive, balance increases)
- [x] 4.3 Verify all existing tests still pass

## 5. Finalize

- [x] 5.1 Run full test suite
- [x] 5.2 Commit and push all changes
