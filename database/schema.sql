CREATE SCHEMA teachme;
SET search_path = teachme;
CREATE SEQUENCE user_group_user_group_id_seq;
CREATE TABLE user_groups (
    user_group_id integer primary key not null default nextval('user_group_user_group_id_seq'::regclass),
    user_group_name varchar(20) NOT NULL UNIQUE,
    created_at timestamp NOT NULL default current_timestamp,
    updated_at timestamp NOT NULL default current_timestamp
);
ALTER SEQUENCE user_group_user_group_id_seq OWNED BY user_groups.user_group_id;

INSERT INTO user_groups (user_group_name) VALUES ('ADMIN');
INSERT INTO user_groups (user_group_name) VALUES ('STUDENT');

CREATE SEQUENCE user_user_id_seq;
CREATE TABLE users (
    user_id integer primary key not null default nextval('user_user_id_seq'::regclass),
    first_name varchar(60) NOT NULL,
    last_name varchar(60) NOT NULL,
    email varchar(256) NOT NULL UNIQUE,
    password varchar(120) NOT NULL,
    user_group_id integer NOT NULL REFERENCES user_groups ON DELETE RESTRICT,
    created_at timestamp NOT NULL default current_timestamp,
    updated_at timestamp NOT NULL default current_timestamp
);
ALTER SEQUENCE user_user_id_seq OWNED BY users.user_id;

CREATE SEQUENCE file_id_seq;
CREATE TABLE files (
    file_id integer primary key not null default nextval('file_id_seq'::regclass),
    filename varchar(60) NOT NULL,
    data bytea NOT NULL,
    updated_at timestamp NOT NULL default current_timestamp
);
ALTER SEQUENCE file_id_seq OWNED BY files.file_id;
