CREATE DATABASE IF NOT exists rate_limiter_db;

USE rate_limiter_db;

CREATE TABLE user_details (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id VARCHAR(255) NOT NULL,
    api_name VARCHAR(255) NOT NULL,
    last_request_time TIMESTAMP NOT NULL,
    request_count INT NOT NULL
);

