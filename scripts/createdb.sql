CREATE SCHEMA IF NOT EXISTS %s;
SET search_path = %s;

CREATE TABLE IF NOT EXISTS featureofinterest
(
  id serial NOT NULL,
  description character varying(255),
  encodingtype integer,
  feature public.geometry,
  CONSTRAINT featureofinterest_pkey PRIMARY KEY (id)
)
WITH (
  OIDS=FALSE
);

CREATE TABLE IF NOT EXISTS thing
(
  id serial NOT NULL,
  description character varying(255),
  properties jsonb,
  CONSTRAINT thing_pkey PRIMARY KEY (id)
)
WITH (
  OIDS=FALSE
);

CREATE TABLE IF NOT EXISTS location
(
  id serial NOT NULL,
  description character varying(255),
  encodingtype integer,
  location public.geometry,
  CONSTRAINT location_pkey PRIMARY KEY (id)
)
WITH (
  OIDS=FALSE
);

CREATE TABLE IF NOT EXISTS thing_to_location
(
  thing_id integer,
  location_id integer,
  CONSTRAINT fk_location_1 FOREIGN KEY (location_id)
      REFERENCES location (id) MATCH SIMPLE
      ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT fk_thing_1 FOREIGN KEY (thing_id)
      REFERENCES thing (id) MATCH SIMPLE
      ON UPDATE NO ACTION ON DELETE NO ACTION
)
WITH (
  OIDS=FALSE
);

CREATE INDEX IF NOT EXISTS fki_location_1
  ON thing_to_location
  USING btree
  (location_id);

CREATE INDEX IF NOT EXISTS fki_thing_1
  ON thing_to_location
  USING btree
  (thing_id);

CREATE TABLE IF NOT EXISTS historicallocation
(
  id serial NOT NULL,
  thing_id integer,
  location_id integer,
  "time" timestamp with time zone,
  CONSTRAINT historicallocation_pkey PRIMARY KEY (id),
  CONSTRAINT fk_location FOREIGN KEY (location_id)
      REFERENCES location (id) MATCH SIMPLE
      ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT fk_thing_hl FOREIGN KEY (thing_id)
      REFERENCES thing (id) MATCH SIMPLE
      ON UPDATE NO ACTION ON DELETE NO ACTION
)
WITH (
  OIDS=FALSE
);

CREATE TABLE IF NOT EXISTS sensor
(
  id serial NOT NULL,
  description character varying(255),
  encodingtype integer,
  metadata character varying(255),
  CONSTRAINT sensor_pkey PRIMARY KEY (id)
)
WITH (
  OIDS=FALSE
);

CREATE TABLE IF NOT EXISTS observedproperty
(
  id serial NOT NULL,
  name character varying(120),
  definition character varying(255),
  description character varying(255),
  CONSTRAINT observedproperty_pkey PRIMARY KEY (id)
)
WITH (
  OIDS=FALSE
);


CREATE TABLE IF NOT EXISTS datastream
(
  id serial NOT NULL,
  description character varying(255),
  unitofmeasurement jsonb,
  observationtype integer,
  observedarea jsonb,
  phenomenontime tstzrange,
  resulttime tstzrange,
  thing_id integer,
  sensor_id integer,
  observerproperty_id integer,
  CONSTRAINT datastream_pkey PRIMARY KEY (id),
  CONSTRAINT fk_observedproperty FOREIGN KEY (observerproperty_id)
      REFERENCES observedproperty (id) MATCH SIMPLE
      ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT fk_sensor FOREIGN KEY (sensor_id)
      REFERENCES sensor (id) MATCH SIMPLE
      ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT fk_thing FOREIGN KEY (thing_id)
      REFERENCES thing (id) MATCH SIMPLE
      ON UPDATE NO ACTION ON DELETE NO ACTION
)
WITH (
  OIDS=FALSE
);

CREATE INDEX IF NOT EXISTS fki_observedproperty
  ON datastream
  USING btree
  (observerproperty_id);

CREATE INDEX IF NOT EXISTS fki_sensor
  ON datastream
  USING btree
  (sensor_id);

CREATE INDEX IF NOT EXISTS fki_thing
  ON datastream
  USING btree
  (thing_id);

CREATE SEQUENCE IF NOT EXISTS observations_id_seq
	START WITH 1
	INCREMENT BY 1
	NO MINVALUE
	NO MAXVALUE
	CACHE 1;

CREATE TABLE IF NOT EXISTS observation
(
  id integer NOT NULL DEFAULT nextval('observations_id_seq'::regclass),
  phenomenontime tstzrange,
  result jsonb,
  resulttime timestamp with time zone,
  resultquality character varying(25),
  validtime tstzrange,
  parameters jsonb,
  stream_id integer,
  featureofinterest_id integer,
  CONSTRAINT fk_datastream FOREIGN KEY (stream_id)
      REFERENCES datastream (id) MATCH SIMPLE
      ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT fk_featureofinterest FOREIGN KEY (featureofinterest_id)
      REFERENCES featureofinterest (id) MATCH SIMPLE
      ON UPDATE NO ACTION ON DELETE NO ACTION
)
WITH (
  OIDS=FALSE
);

CREATE INDEX IF NOT EXISTS fki_datastream
  ON observation
  USING btree
  (stream_id);

CREATE INDEX IF NOT EXISTS fki_featureofinterest
  ON observation
  USING btree
  (featureofinterest_id);

CREATE INDEX IF NOT EXISTS fki_location
  ON historicallocation
  USING btree
  (location_id);

CREATE INDEX IF NOT EXISTS fki_thing_hl
  ON historicallocation
  USING btree
  (thing_id);
