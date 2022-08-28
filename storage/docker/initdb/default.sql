CREATE DATABASE services;
\c 'services';

DROP TABLE IF EXISTS service_packages;

CREATE SEQUENCE services_seq;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE service_packages (
  service_id UUID NOT NULL DEFAULT uuid_generate_v4(),
  name varchar(100) NOT NULL,
  created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
  description varchar(100) NOT NULL DEFAULT '',
  PRIMARY KEY (service_id),
  CONSTRAINT name_unique UNIQUE (name, service_id)
);

ALTER SEQUENCE services_seq RESTART WITH 2006;

INSERT INTO service_packages VALUES 
('a81bc81b-dead-4e5d-abff-90865d1e13b1','s1','2020-1-30 13:10:53.163', '2020-12-30 13:10:53.163','s1d1'),
('bc1bc81b-dead-4e5d-abff-90865d1e13b1','s2','2020-2-20 13:10:53.163', '2020-12-30 13:10:53.163','s2d2'),
('de1bc81b-dead-4e5d-abff-90865d1e13b1','s3','2020-3-21 13:10:53.163', '2020-12-30 13:10:53.163','s3d3');
