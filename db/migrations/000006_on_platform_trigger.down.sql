DROP TRIGGER IF EXISTS t_before_insert_check_time_on_platform ON session.on_platform;
DROP FUNCTION IF EXISTS session.f_before_insert_check_time_on_platform();

DROP TRIGGER IF EXISTS t_before_update_check_time_on_platform ON session.on_platform;
DROP FUNCTION IF EXISTS session.f_before_update_check_time_on_platform();