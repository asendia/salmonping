INSERT INTO online_listing (name, platform, url)
VALUES 
    ('Gojek: Salmon Fit Kebon Jeruk', 'gofood', 'https://gofood.co.id/jakarta/restaurant/salmon-fit-apartemen-menara-kebon-jeruk-06f0dcc6-14f4-4092-810f-2bcc81214d23')
    , ('Grab: Salmon Fit Kebon Jeruk', 'grabfood', 'https://food.grab.com/id/id/restaurant/salmon-fit-apartemen-menara-kebun-jeruk-delivery/6-C2XUWAX3PEU1JT')
    , ('Gojek: Salmon Fit Sudirman', 'gofood', 'https://gofood.co.id/jakarta/restaurant/salmon-fit-sudirman-815b2b33-584e-46d6-b12e-2d6da2f46f96')
    , ('Grab: Salmon Fit Sudirman', 'grabfood', 'https://food.grab.com/id/id/restaurant/salmon-fit-sudirman-delivery/6-C36EKGLYHB42DA')
    , ('Gojek: Salmon Fit Haji Nawi', 'gofood', 'https://gofood.co.id/jakarta/restaurant/salmon-fit-haji-nawi-9d68471b-5d49-468c-a162-091b1ea9b468')
    , ('Grab: Salmon Fit Haji Nawi', 'grabfood', 'https://food.grab.com/id/id/restaurant/salmon-fit-haji-nawi-delivery/6-C4LJLRKEKAEVME')
ON CONFLICT DO NOTHING;