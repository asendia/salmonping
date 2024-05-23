INSERT INTO online_listing (name, platform, url)
VALUES 
    ('Kebon Jeruk', 'gofood', 'https://gofood.co.id/jakarta/restaurant/salmon-fit-apartemen-menara-kebon-jeruk-06f0dcc6-14f4-4092-810f-2bcc81214d23')
    , ('Kebon Jeruk', 'grabfood', 'https://food.grab.com/id/id/restaurant/salmon-fit-apartemen-menara-kebun-jeruk-delivery/6-C2XUWAX3PEU1JT')
    , ('Sudirman', 'gofood', 'https://gofood.co.id/jakarta/restaurant/salmon-fit-sudirman-815b2b33-584e-46d6-b12e-2d6da2f46f96')
    , ('Sudirman', 'grabfood', 'https://food.grab.com/id/id/restaurant/salmon-fit-sudirman-delivery/6-C36EKGLYHB42DA')
    , ('Haji Nawi', 'gofood', 'https://gofood.co.id/jakarta/restaurant/salmon-fit-haji-nawi-9d68471b-5d49-468c-a162-091b1ea9b468')
    , ('Haji Nawi', 'grabfood', 'https://food.grab.com/id/id/restaurant/salmon-fit-haji-nawi-delivery/6-C4LJLRKEKAEVME')
    , ('Tanjung Duren', 'gofood', 'https://gofood.co.id/jakarta/restaurant/salmon-fit-apartemen-menara-kebon-jeruk-06f0dcc6-14f4-4092-810f-2bcc81214d23')
    , ('Tanjung Duren', 'grabfood', 'https://food.grab.com/id/id/restaurant/salmon-fit-apartemen-menara-kebun-jeruk-delivery/6-C2XUWAX3PEU1JT')
ON CONFLICT DO NOTHING;

DO $$ 
DECLARE
    online_listing_id_value uuid;
BEGIN
    FOR online_listing_id_value IN (SELECT id FROM online_listing)
    LOOP
        -- Monday to Saturday
        FOR day_in_week IN 1..6
        LOOP
            -- 10AM-8PM, adjusted for Asia/Jakarta timezone (UTC+7), assuming server is in UTC
            INSERT INTO schedule (online_listing_id, day_of_week, opening_time, closing_time)
            VALUES (online_listing_id_value, day_in_week, '03:00:00', '13:00:00')
            ON CONFLICT DO NOTHING;
        END LOOP;
    END LOOP;
END $$;
