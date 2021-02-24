package database

const schema = `
CREATE TABLE IF NOT EXISTS errors (
	hash TEXT NOT NULL,
	error_id TEXT NOT NULL,
	message TEXT NOT NULL,
	timestamp TIMESTAMP NOT NULL,
	data TEXT NOT NULL
) engine=MergeTree PRIMARY KEY (error_id) ORDER BY error_id;
`
