CREATE TABLE users (
	id INTEGER PRIMARY KEY,
	name TEXT NOT NULL UNIQUE,
	password_hash TEXT NOT NULL,
	balance INTEGER NOT NULL,
	invites INTEGER NOT NULL
);

CREATE TABLE services (
	id INTEGER PRIMARY KEY,
	name TEXT NOT NULL,
	type TEXT NOT NULL,
	created_at INTEGER NOT NULL,
	expires_at INTEGER NOT NULL,
	prolong INTEGER NOT NULL,
	prolong_price INTEGER NOT NULL,
	user_id INTEGER NOT NULL,
	location_id INTEGER NOT NULL,
	FOREIGN KEY (user_id)
	REFERENCES users (id),
	FOREIGN KEY (location_id)
	REFERENCES service_locations (id)
);

CREATE TABLE service_locations (
	id INTEGER PRIMARY KEY,
	name TEXT NOT NULL,
	address TEXT NOT NULL,
	services TEXT NOT NULL
);

CREATE TABLE service_prices (
	amount INTEGER NOT NULL,
	type TEXT NOT NULL UNIQUE
);

CREATE TABLE transactions (
	id INTEGER PRIMARY KEY,
	payment_id TEXT NOT NULL,
	amount INTEGER NOT NULL,
	status TEXT NOT NULL,
	timestamp INTEGER NOT NULL,
	url TEXT NOT NULL,
	user_id INTEGER NOT NULL,
	FOREIGN KEY (user_id)
	REFERENCES users (id)
);

CREATE TABLE invites (
	id INTEGER PRIMARY KEY,
	code TEXT NOT NULL UNIQUE,
	used INTEGER NOT NULL,
	user_id INTEGER NOT NULL,
	FOREIGN KEY (user_id)
	REFERENCES users (id)
);
