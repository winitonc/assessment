CREATE TABLE IF NOT EXISTS expenses (
	id SERIAL PRIMARY KEY,
	title TEXT,
	amount FLOAT,
	note TEXT,
	tags TEXT []
);

INSERT INTO
	expenses (title, amount, note, tags)
values
	(
		'Title IT',
		10.99,
		'Note IT',
		'{tagsIT1,tagsIT2}'
	)