DO $$
BEGIN
     IF EXISTS(
        SELECT schema_name
          FROM information_schema.schemata
          WHERE schema_name = '%s'
      )
    THEN
      EXECUTE 'drop schema %s cascade';
    END IF;
END
$$;

CREATE SCHEMA %s;
SET search_path = %s;

CREATE TABLE featureofinterest
(
  id bigserial NOT NULL,
  name character varying(255),
  description character varying(500),
  encodingtype integer,
  feature public.geometry(geometry,4326),
  original_location_id bigint,
  CONSTRAINT featureofinterest_pkey PRIMARY KEY (id)
)
WITH (
  OIDS=FALSE
);

CREATE TABLE thing
(
  id bigserial NOT NULL,
  name character varying(255),
  description character varying(500),
  properties jsonb,
  CONSTRAINT thing_pkey PRIMARY KEY (id)
)
WITH (
  OIDS=FALSE
);

CREATE TABLE location
(
  id bigserial NOT NULL,
  name character varying(255),
  description character varying(500),
  encodingtype integer,
  location public.geometry(geometry,4326),
  CONSTRAINT location_pkey PRIMARY KEY (id)
)
WITH (
  OIDS=FALSE
);

CREATE TABLE thing_to_location
(
  thing_id bigint,
  location_id bigint,
  CONSTRAINT fk_location_1 FOREIGN KEY (location_id)
      REFERENCES location (id) MATCH SIMPLE
      ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT fk_thing_1 FOREIGN KEY (thing_id)
      REFERENCES thing (id) MATCH SIMPLE
      ON UPDATE NO ACTION ON DELETE CASCADE
)
WITH (
  OIDS=FALSE
);

CREATE INDEX fki_location_1
  ON thing_to_location
  USING btree
  (location_id);

CREATE INDEX fki_thing_1
  ON thing_to_location
  USING btree
  (thing_id);

CREATE TABLE historicallocation
(
  id bigserial NOT NULL,
  thing_id bigint,
  location_id bigint,
  "time" timestamp with time zone,
  CONSTRAINT historicallocation_pkey PRIMARY KEY (id),
  CONSTRAINT fk_location FOREIGN KEY (location_id)
      REFERENCES location (id) MATCH SIMPLE
      ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT fk_thing_hl FOREIGN KEY (thing_id)
      REFERENCES thing (id) MATCH SIMPLE
      ON UPDATE NO ACTION ON DELETE CASCADE
)
WITH (
  OIDS=FALSE
);

CREATE TABLE sensor
(
  id bigserial NOT NULL,
  name character varying(255),
  description character varying(500),
  encodingtype integer,
  metadata character varying(255),
  CONSTRAINT sensor_pkey PRIMARY KEY (id)
)
WITH (
  OIDS=FALSE
);

CREATE TABLE observedproperty
(
  id bigserial NOT NULL,
  name character varying(120),
  definition character varying(255),
  description character varying(500),
  CONSTRAINT observedproperty_pkey PRIMARY KEY (id)
)
WITH (
  OIDS=FALSE
);


CREATE TABLE datastream
(
  id bigserial NOT NULL,
  name character varying(255),
  description character varying(500),
  unitofmeasurement jsonb,
  observationtype integer,
  observedarea public.geometry(geometry,4326),
  phenomenontime tstzrange,
  resulttime tstzrange,
  thing_id bigint,
  sensor_id bigint,
  observedproperty_id bigint,
  CONSTRAINT datastream_pkey PRIMARY KEY (id),
  CONSTRAINT fk_observedproperty FOREIGN KEY (observedproperty_id)
      REFERENCES observedproperty (id) MATCH SIMPLE
      ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT fk_sensor FOREIGN KEY (sensor_id)
      REFERENCES sensor (id) MATCH SIMPLE
      ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT fk_thing FOREIGN KEY (thing_id)
      REFERENCES thing (id) MATCH SIMPLE
      ON UPDATE NO ACTION ON DELETE CASCADE
)
WITH (
  OIDS=FALSE
);

CREATE INDEX fki_observedproperty
  ON datastream
  USING btree
  (observedproperty_id);

CREATE INDEX fki_sensor
  ON datastream
  USING btree
  (sensor_id);

CREATE INDEX fki_thing
  ON datastream
  USING btree
  (thing_id);

CREATE SEQUENCE observations_id_seq
	START WITH 1
	INCREMENT BY 1
	NO MINVALUE
	NO MAXVALUE
	CACHE 1;

CREATE UNLOGGED TABLE observation
(
  id bigint NOT NULL DEFAULT nextval('observations_id_seq'::regclass),
  data jsonb,
  stream_id bigint,
  featureofinterest_id bigint,
  CONSTRAINT fk_datastream FOREIGN KEY (stream_id)
      REFERENCES datastream (id) MATCH SIMPLE
      ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT fk_featureofinterest FOREIGN KEY (featureofinterest_id)
      REFERENCES featureofinterest (id) MATCH SIMPLE
      ON UPDATE NO ACTION ON DELETE CASCADE
)
WITH (
  OIDS=FALSE
);

CREATE INDEX fki_datastream
  ON observation
  USING btree
  (stream_id);

CREATE INDEX fki_featureofinterest
  ON observation
  USING btree
  (featureofinterest_id);

CREATE INDEX fki_location
  ON historicallocation
  USING btree
  (location_id);

CREATE INDEX fki_thing_hl
  ON historicallocation
  USING btree
  (thing_id);
