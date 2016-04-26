package gostdb

import "fmt"

func GetCreateDatabaseQuery(schema string) string {
	return fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %s;" +
	"SET search_path = %s;" +
	"" +
	"CREATE TABLE IF NOT EXISTS datastream (" +
	"    id integer NOT NULL," +
	"    description character varying(255)," +
	"    unitofmeasurement json," +
	"    observationtype integer," +
	"    observedarea json," +
	"    phenomenontime tstzrange," +
	"    resulttime tstzrange," +
	"    thing_id integer," +
	"    sensor_id integer," +
	"    observerproperty_id integer" +
	");" +
	"" +
	"CREATE SEQUENCE IF NOT EXISTS datastream_id_seq" +
	"    START WITH 1" +
	"    INCREMENT BY 1" +
	"    NO MINVALUE" +
	"    NO MAXVALUE" +
	"    CACHE 1;" +
	"" +
	"CREATE TABLE IF NOT EXISTS featureofinterest (" +
	"    id integer NOT NULL," +
	"    description character varying(255)," +
	"    encodingtype integer," +
	"    feature json" +
	");" +
	"" +
	"CREATE SEQUENCE IF NOT EXISTS featureofinterest_id_seq" +
	"    START WITH 1" +
	"    INCREMENT BY 1" +
	"    NO MINVALUE" +
	"    NO MAXVALUE" +
	"    CACHE 1;" +
	"" +
	"CREATE TABLE IF NOT EXISTS historicallocation (" +
	"    id integer NOT NULL," +
	"    thing_id integer," +
	"    location_id integer," +
	"    \"time\" timestamp with time zone" +
	");" +
	"" +
	"CREATE SEQUENCE IF NOT EXISTS historicallocation_id_seq" +
	"    START WITH 1" +
	"    INCREMENT BY 1" +
	"    NO MINVALUE" +
	"    NO MAXVALUE" +
	"    CACHE 1;" +
	"" +
	"CREATE TABLE IF NOT EXISTS location (" +
	"    id integer NOT NULL," +
	"    description character varying(255)," +
	"    encodingtype integer," +
	"    location character varying(255)" +
	");" +
	"" +
	"CREATE SEQUENCE IF NOT EXISTS location_id_seq" +
	"    START WITH 1" +
	"    INCREMENT BY 1" +
	"    NO MINVALUE" +
	"    NO MAXVALUE" +
	"    CACHE 1;" +
	"" +
	"CREATE TABLE IF NOT EXISTS observation (" +
	"    id integer NOT NULL," +
	"    phenomenontime tstzrange," +
	"    result json," +
	"    resulttime timestamp with time zone," +
	"    resultquality character varying(25)," +
	"    validtime tstzrange," +
	"    parameters json," +
	"    stream_id integer," +
	"    featureofinterest_id integer" +
	");" +
	"" +
	"CREATE SEQUENCE IF NOT EXISTS observations_id_seq" +
	"    START WITH 1" +
	"    INCREMENT BY 1" +
	"    NO MINVALUE" +
	"    NO MAXVALUE" +
	"    CACHE 1;" +
	"" +
	"CREATE TABLE IF NOT EXISTS observedproperty (" +
	"    id integer NOT NULL," +
	"    name character varying(120)," +
	"    definition character varying(255)," +
	"    description character varying(255)" +
	");" +
	"" +
	"CREATE SEQUENCE IF NOT EXISTS observedproperty_id_seq" +
	"    START WITH 1" +
	"    INCREMENT BY 1" +
	"    NO MINVALUE" +
	"    NO MAXVALUE" +
	"    CACHE 1;" +
	"" +
	"CREATE TABLE IF NOT EXISTS sensor (" +
	"    id integer NOT NULL," +
	"    description character varying(255)," +
	"    encodingtype integer," +
	"    metadata json" +
	");" +
	"" +
	"CREATE SEQUENCE IF NOT EXISTS sensor_id_seq" +
	"    START WITH 1" +
	"    INCREMENT BY 1" +
	"    NO MINVALUE" +
	"    NO MAXVALUE" +
	"    CACHE 1;" +
	"" +
	"CREATE TABLE IF NOT EXISTS thing (" +
	"    id integer NOT NULL," +
	"    description character varying(255)," +
	"    properties json" +
	");" +
	"" +
	"CREATE SEQUENCE IF NOT EXISTS thing_id_seq" +
	"    START WITH 1" +
	"    INCREMENT BY 1" +
	"    NO MINVALUE" +
	"    NO MAXVALUE" +
	"    CACHE 1;" +
	"" +
	"CREATE TABLE IF NOT EXISTS thing_to_location (" +
	"    thing_id integer," +
	"    location_id integer" +
	");" +
	"" +
	"" +
	"ALTER TABLE ONLY datastream ALTER COLUMN id SET DEFAULT nextval('datastream_id_seq'::regclass);" +
	"ALTER TABLE ONLY featureofinterest ALTER COLUMN id SET DEFAULT nextval('featureofinterest_id_seq'::regclass);" +
	"ALTER TABLE ONLY historicallocation ALTER COLUMN id SET DEFAULT nextval('historicallocation_id_seq'::regclass);" +
	"ALTER TABLE ONLY location ALTER COLUMN id SET DEFAULT nextval('location_id_seq'::regclass);" +
	"ALTER TABLE ONLY observation ALTER COLUMN id SET DEFAULT nextval('observations_id_seq'::regclass);" +
	"ALTER TABLE ONLY observedproperty ALTER COLUMN id SET DEFAULT nextval('observedproperty_id_seq'::regclass);" +
	"ALTER TABLE ONLY sensor ALTER COLUMN id SET DEFAULT nextval('sensor_id_seq'::regclass);" +
	"ALTER TABLE ONLY thing ALTER COLUMN id SET DEFAULT nextval('thing_id_seq'::regclass);" +
	"ALTER TABLE ONLY datastream ADD CONSTRAINT datastream_pkey PRIMARY KEY (id);" +
	"ALTER TABLE ONLY featureofinterest ADD CONSTRAINT featureofinterest_pkey PRIMARY KEY (id);" +
	"ALTER TABLE ONLY historicallocation ADD CONSTRAINT historicallocation_pkey PRIMARY KEY (id);" +
	"ALTER TABLE ONLY location ADD CONSTRAINT location_pkey PRIMARY KEY (id);" +
	"ALTER TABLE ONLY observedproperty ADD CONSTRAINT observedproperty_pkey PRIMARY KEY (id);" +
	"ALTER TABLE ONLY sensor ADD CONSTRAINT sensor_pkey PRIMARY KEY (id);" +
	"ALTER TABLE ONLY thing ADD CONSTRAINT thing_pkey PRIMARY KEY (id);" +
	"" +
	"CREATE INDEX IF NOT EXISTS fki_datastream ON observation USING btree (stream_id);" +
	"CREATE INDEX IF NOT EXISTS fki_featureofinterest ON observation USING btree (featureofinterest_id);" +
	"CREATE INDEX IF NOT EXISTS fki_location ON historicallocation USING btree (location_id);" +
	"CREATE INDEX IF NOT EXISTS fki_location_1 ON thing_to_location USING btree (location_id);" +
	"CREATE INDEX IF NOT EXISTS fki_observedproperty ON datastream USING btree (observerproperty_id);" +
	"CREATE INDEX IF NOT EXISTS fki_sensor ON datastream USING btree (sensor_id);" +
	"CREATE INDEX IF NOT EXISTS fki_thing ON datastream USING btree (thing_id);" +
	"CREATE INDEX IF NOT EXISTS fki_thing_1 ON thing_to_location USING btree (thing_id);" +
	"CREATE INDEX IF NOT EXISTS fki_thing_hl ON historicallocation USING btree (thing_id);" +
	"" +
	"CREATE OR REPLACE FUNCTION create_constraint_if_not_exists (t_name text, c_name text, constraint_sql text)" +
	"  RETURNS void" +
	"AS" +
	"$func$" +
	"    CASE WHEN not exists (select constraint_name from information_schema.constraint_column_usage where table_name = t_name and constraint_name = c_name) THEN" +
	"        execute 'ALTER TABLE ' || t_name || ' ADD CONSTRAINT ' || c_name || ' ' || constraint_sql;" +
	"    END" +
	"$func$ LANGUAGE sql;" +
	"" +
	"SELECT create_constraint_if_not_exists('observation', 'fk_datastream', 'ALTER TABLE ONLY observation ADD CONSTRAINT fk_datastream FOREIGN KEY (stream_id) REFERENCES datastream(id);');" +
	"SELECT create_constraint_if_not_exists('observation', 'fk_featureofinterest', 'ALTER TABLE ONLY observation ADD CONSTRAINT fk_featureofinterest FOREIGN KEY (featureofinterest_id) REFERENCES featureofinterest(id);');" +
	"SELECT create_constraint_if_not_exists('historicallocation', 'fk_location', 'ALTER TABLE ONLY historicallocation ADD CONSTRAINT fk_location FOREIGN KEY (location_id) REFERENCES location(id);');" +
	"SELECT create_constraint_if_not_exists('thing_to_location', 'fk_location_1', 'ALTER TABLE ONLY thing_to_location ADD CONSTRAINT fk_location_1 FOREIGN KEY (location_id) REFERENCES location(id);');" +
	"SELECT create_constraint_if_not_exists('datastream', 'fk_observedproperty', 'ALTER TABLE ONLY datastream ADD CONSTRAINT fk_observedproperty FOREIGN KEY (observerproperty_id) REFERENCES observedproperty(id);');" +
	"SELECT create_constraint_if_not_exists('datastream', 'fk_sensor', 'ALTER TABLE ONLY datastream ADD CONSTRAINT fk_sensor FOREIGN KEY (sensor_id) REFERENCES sensor(id);');" +
	"SELECT create_constraint_if_not_exists('datastream', 'fk_thing', 'ALTER TABLE ONLY datastream ADD CONSTRAINT fk_thing FOREIGN KEY (thing_id) REFERENCES thing(id);');" +
	"SELECT create_constraint_if_not_exists('thing_to_location', 'fk_thing_1', 'ALTER TABLE ONLY thing_to_location ADD CONSTRAINT fk_thing_1 FOREIGN KEY (thing_id) REFERENCES thing(id);');" +
	"SELECT create_constraint_if_not_exists('historicallocation', 'fk_thing_hl', 'ALTER TABLE ONLY historicallocation ADD CONSTRAINT fk_thing_hl FOREIGN KEY (thing_id) REFERENCES thing(id);');", schema, schema)
}
