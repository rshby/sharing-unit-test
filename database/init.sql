CREATE DATABASE sharingunittest;

USE sharingunittest;

CREATE TABLE accounts
(
    id INT PRIMARY KEY AUTO_INCREMENT NOT NULL ,
    email VARCHAR(256) NOT NULL UNIQUE ,
    username VARCHAR(256) NOT NULL UNIQUE ,
    password VARCHAR(256) NOT NULL ,
    full_name VARCHAR(256) DEFAULT NULL,
    gender VARCHAR(1) DEFAULT NULL
);