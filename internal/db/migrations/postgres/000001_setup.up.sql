-- generate uuid v7 as per RFC 9562.
--
-- https://datatracker.ietf.org/doc/html/rfc9562


create or replace function uuid_generate_v7() returns uuid as $$ begin -- use random v4 uuid as starting point (which has the same variant we need)
  -- then overlay timestamp
  -- then set version 7 by flipping the 2 and 1 bit in the version 4 string
  return encode(
    set_bit(
      set_bit(
        overlay(
          uuid_send(gen_random_uuid()) placing substring(
            int8send(
              floor(
                extract(
                  epoch
                  from clock_timestamp()
                ) * 1000
              )::bigint
            )
            from 3
          )
          from 1 for 6
        ),
        52,
        1
      ),
      53,
      1
    ),
    'hex'
  )::uuid;
end $$ language plpgsql volatile;

-- create a function that updates the updated_at column with the current timestamp

create or replace function update_updated_at()
returns trigger as $$
begin
  new.updated_at = now();
  return new;
end;
$$ language plpgsql;
