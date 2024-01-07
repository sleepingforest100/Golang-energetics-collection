CREATE TABLE energetics
(
    energetics_id serial NOT NULL,
    name varchar(60) NOT NULL,
    taste varchar(30),
    description varchar(128),
    manufacturer_name varchar(35),
    manufacture_country varchar(35),
    pictureURL text,
    PRIMARY KEY (energetics_id)
);
