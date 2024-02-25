-- SQLite
SELECT * FROM books;
SELECT * FROM shelves;

INSERT INTO shelves (id, theme, size, creat_at, update_at) VALUES (NULL, 'drama', 5, '2024-02-21 10:30:19.795642+08:00', '2024-02-21 10:30:19.795642+08:00');

INSERT INTO books (id, author, title, shelf_id, creat_at, update_at) VALUES (NULL, 'John', 'The Book 0',1, '', '');
INSERT INTO books (id, author, title, shelf_id, creat_at, update_at) VALUES (NULL, "Jane Doe", "The Book 2", 1, '2019-01-01 00:00:00', '2019-01-01 00:00:00');
INSERT INTO books (id, author, title, shelf_id, creat_at, update_at) VALUES (NULL, "John Doe", "The Book 3", 1, '2019-01-01 00:00:00', '2019-01-01 00:00:00');
INSERT INTO books (id, author, title, shelf_id, creat_at, update_at) VALUES (NULL, "Jane Doe", "The Book 4", 1, '2019-01-01 00:00:00', '2019-01-01 00:00:00');
INSERT INTO books (id, author, title, shelf_id, creat_at, update_at) VALUES (NULL, "John Doe", "The Book 5", 2, '2019-01-01 00:00:00', '2019-01-01 00:00:00');