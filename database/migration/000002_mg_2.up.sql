CREATE TABLE compositions
(
	composition_id serial NOT NULL PRIMARY KEY,
    energetics_id integer,
    caffeine integer,
    taurine integer,
	FOREIGN KEY (energetics_id) REFERENCES energetics (energetics_id)
);

ALTER TABLE IF EXISTS energetics
    ADD COLUMN composition_id bigint
	REFERENCES compositions(composition_id);