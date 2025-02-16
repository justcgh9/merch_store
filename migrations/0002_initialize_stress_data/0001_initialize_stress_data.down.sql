DELETE FROM inventory
WHERE username LIKE 'user%';

DELETE FROM balance
WHERE username LIKE 'user%';

DELETE FROM users
WHERE username LIKE 'user%';
