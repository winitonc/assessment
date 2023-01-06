CREATE TABLE IF NOT EXISTS expenses (
	id SERIAL PRIMARY KEY,
	title TEXT,
	amount FLOAT,
	note TEXT,
	tags TEXT []
);

INSERT INTO
	expenses (id, title, amount, note, tags)
values
	(
		2,
		'Title IT',
		10.99,
		'Note IT',
		ARRAY ['tagsIT1', 'tagsIT2']
	)