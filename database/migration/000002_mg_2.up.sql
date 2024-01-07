CREATE TABLE composition
(
	composition_id serial NOT NULL PRIMARY KEY,
    energetics_id integer NOT NULL,
    caffeine integer,
    taurine integer,
	FOREIGN KEY (energetics_id) REFERENCES energetics (id)
);

ALTER TABLE IF EXISTS energetics
    ADD COLUMN composition_id bigint
	REFERENCES composition(composition_id);