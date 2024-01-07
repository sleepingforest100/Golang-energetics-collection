CREATE TABLE energetics
(
    id serial NOT NULL,
    name "char"[] NOT NULL,
    taste "char"[],
    manufacturer_name "char"[],
    manufacture_country "char"[],
    PRIMARY KEY (id)
);
