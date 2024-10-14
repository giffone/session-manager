DROP TRIGGER IF EXISTS t_before_insert_check_time_on_campus ON session.on_campus;
DROP FUNCTION IF EXISTS session.f_before_insert_check_time_on_campus();

DROP TRIGGER IF EXISTS t_before_update_check_time_on_campus ON session.on_campus;
DROP FUNCTION IF EXISTS session.f_before_update_check_time_on_campus();