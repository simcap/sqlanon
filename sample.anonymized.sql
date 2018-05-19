CREATE TABLE names (
    id serial NOT NULL,
    "language" character varying(100) NOT NULL,
    entity_type character varying(100) NOT NULL,
    entity_id integer NOT NULL,
    full_name character varying(100),
    first_name character varying(100),
    middle_name character varying(100),
    last_name character varying(100),
    alias character varying(100),
    abbreviation character varying(100),
    short_name character varying(100),
    prefix character varying(20),
    suffix character varying(20)
);

INSERT INTO names (id, "language", entity_type, entity_id, full_name, first_name, middle_name, last_name, alias, abbreviation, short_name, prefix, suffix) VALUES (1, 'en-US', 'affiliations', 1, 'Major League Baseball', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL);
INSERT INTO names (id, "language", entity_type, entity_id, full_name, first_name, middle_name, last_name, alias, abbreviation, short_name, prefix, suffix) VALUES (2, 'en-US', 'affiliations', 2, 'Baseball', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL);
INSERT INTO names (id, "language", entity_type, entity_id, full_name, first_name, middle_name, last_name, alias, abbreviation, short_name, prefix, suffix) VALUES (34, 'en-US', 'affiliations', 4, 'American Football', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL);
INSERT INTO names (id, "language", entity_type, entity_id, full_name, first_name, middle_name, last_name, alias, abbreviation, short_name, prefix, suffix) VALUES (33, 'en-US', 'affiliations', 3, 'National Football League', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL);
INSERT INTO names (id, "language", entity_type, entity_id, full_name, first_name, middle_name, last_name, alias, abbreviation, short_name, prefix, suffix) VALUES (67, 'en-US', 'affiliations', 5, 'American', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL);
INSERT INTO names (id, "language", entity_type, entity_id, full_name, first_name, middle_name, last_name, alias, abbreviation, short_name, prefix, suffix) VALUES (68, 'en-US', 'affiliations', 6, 'AFC East', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL);
INSERT INTO names (id, "language", entity_type, entity_id, full_name, first_name, middle_name, last_name, alias, abbreviation, short_name, prefix, suffix) VALUES (69, 'en-US', 'affiliations', 7, 'AFC North', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL);
INSERT INTO names (id, "language", entity_type, entity_id, full_name, first_name, middle_name, last_name, alias, abbreviation, short_name, prefix, suffix) VALUES (70, 'en-US', 'affiliations', 8, 'AFC South', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL);
INSERT INTO names (id, "language", entity_type, entity_id, full_name, first_name, middle_name, last_name, alias, abbreviation, short_name, prefix, suffix) VALUES (71, 'en-US', 'affiliations', 9, 'AFC West', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL);
INSERT INTO names (id, "language", entity_type, entity_id, full_name, first_name, middle_name, last_name, alias, abbreviation, short_name, prefix, suffix) VALUES (72, 'en-US', 'affiliations', 10, 'National', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL);
INSERT INTO names (id, "language", entity_type, entity_id, full_name, first_name, middle_name, last_name, alias, abbreviation, short_name, prefix, suffix) VALUES (73, 'en-US', 'affiliations', 11, 'NFC East', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL);
INSERT INTO names (id, "language", entity_type, entity_id, full_name, first_name, middle_name, last_name, alias, abbreviation, short_name, prefix, suffix) VALUES (74, 'en-US', 'affiliations', 12, 'NFC North', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL);
INSERT INTO names (id, "language", entity_type, entity_id, full_name, first_name, middle_name, last_name, alias, abbreviation, short_name, prefix, suffix) VALUES (75, 'en-US', 'affiliations', 13, 'NFC South', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL);
