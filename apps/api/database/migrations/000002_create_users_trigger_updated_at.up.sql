CREATE OR REPLACE FUNCTION fn_set_updated_at_to_now()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = NOW();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER trig_b_u_updated_at_to_now
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE PROCEDURE fn_set_updated_at_to_now();