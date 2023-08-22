CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE pessoa (
	id uuid NOT NULL,
	nome varchar(100) NOT NULL,
	apelido varchar(32) NOT NULL,
	nascimento varchar(10) NOT NULL,
	stack _text NULL,
	CONSTRAINT pessoa_pkey PRIMARY KEY (id)
);
