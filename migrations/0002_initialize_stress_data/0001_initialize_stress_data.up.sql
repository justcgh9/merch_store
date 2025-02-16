INSERT INTO users (username, password)
SELECT 'user' || i, '$2a$10$jCyIG8r5D1dW2c01cg/oo.3HYBe.VXc3x/x3/DbMKFkyaLbnVN2Ia'
FROM generate_series(1, 100000) AS s(i);

INSERT INTO Balance (username, balance)
SELECT 'user' || i, 100000000
FROM generate_series(1, 100000) AS s(i);

INSERT INTO inventory (username)
SELECT 'user' || i
FROM generate_series(1, 100000) AS s(i);


