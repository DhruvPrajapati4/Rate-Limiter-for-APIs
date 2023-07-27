# Rate Limiter Spring Boot Application

This project is a Rate Limiter implementation using Spring Boot. It provides a simple API rate limiting functionality to control the number of incoming requests for specific users and APIs.

## Table of Contents
- [Introduction](#introduction)
- [Technologies Used](#technologies-used)
- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Installation](#installation)
- [Usage](#usage)
- [API Endpoints](#api-endpoints)
- [Rate Limiting Logic](#rate-limiting-logic)
- [Configuration](#configuration)

## Introduction

Rate limiting is a technique used to control the rate of incoming requests to prevent server overload. This project provides a simple rate limiter for four sample APIs (GET, POST, PUT, DELETE) that can be used to limit the number of requests made by a specific user within a certain time frame.

The rate limiter implementation includes three types of checks:
1. API-specific rate limiter: Limit the number of requests a user can make for a specific API within a minute.
2. Global user rate limiter: Limit the total number of requests a user can make across all APIs within a minute.
3. Global API rate limiter: Limit the total number of requests any user can make for a specific API within a second.

## Technologies Used

The project is built using the following technologies:
- Java (version 11 or higher)
- Spring Boot (version 2.5.4 or higher)
- MySQL Database

## Getting Started

### Prerequisites

To run this application, you will need the following installed on your machine:
- Java Development Kit (JDK) 11 or higher
- MySQL Database
- Spring-boot v2.7.8 or Above

### Installation

1. Clone the Git repository to your local machine: git clone https://github.com/DhruvPrajapati4/Rate-limiter-for-APIs.git

2. Setup the MySQL database using the DDL.sql file

3. Navigate to the project directory: cd src

4. Create/Modify the `application.properties` file in the `src/main/resources` directory and configure the MySQL database connection settings according to your local environment.

5. Build the application using Maven: mvn clean install

6. Run the application: mvn spring-boot:run

## Usage

Once the application is running, you can use any API testing tool (e.g., Postman) or your web browser to test the rate limiter functionality.

## API Endpoints

The application provides the following API endpoints:

1. GET `/api-alpha/{userId}`: Sample API to be rate-limited for a user.
2. POST `/api-beta/{userId}`: Sample API to be rate-limited for a user.
3. PUT `/api-gamma/{userId}`: Sample API to be rate-limited for a user with a specific `userId`.
4. DELETE `/api-delta/{userId}`: Sample API to be rate-limited for a user.

All API endpoints accept a query parameter `userId` of type `long` to identify the user.

## Rate Limiting Logic

The rate limiter logic is implemented in the `RateLimiterService` class. It performs API-specific rate limiting, global user rate limiting, and global API rate limiting checks for each incoming request.

API-specific rate limiting allows a user to call a specific API only a certain number of times within a minute. Global user rate limiting restricts the total number of requests a user can make across all APIs within a minute. Global API rate limiting restricts the total number of requests any user can make for a specific API within a second.

## Configuration

You can customize the rate limit values for API-specific, global user, and global API rate limiting by modifying the constants in the `RateLimiterService` class.
