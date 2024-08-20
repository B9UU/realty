# Realestate data 
## TODOS
- [x] make API that returns hello world
- [x] add server configuration
- [x] take configuration as flags
- [x] create and connect to database and tables using migrate
- [ ] users, token and realty model
- [x] endpoint to add listings
- [x] return a simple list of 20 listings
- [x] add custom errors
- [x] create a make file
- [x] add json logging (should read more on let's go further)
- [ ] add rate limiter

- [x] testing
- [x] add middleware to log all requests
- [x] add rentals.ca schema to database
- [x] add new realty models
<!-- - [ ] create cities table -->
- [x] update getRealties to use new schema
- [x] endpoint for autocomplete cities
- [x] add new tests for all endpoints
- [x] start taking queries (city) and only return results there
- [x] make getRealties onlly return limited data
- [x] add filters
- [x] returning metadata to getRealties handler
- [-] metadata for autoComplete handler
- [x] create an endpoint to return all the data available for the given realty id
- [-] maybe figure out how to take lng-int and how to search database for that
- [-] figure out how to save images localy with custom url
- [x] add multiple filters
- [x] make getAll return limited data but uses filters
- [x] make get(id) to return all available data for a realty
- [ ] validation

- [x] create users table
- [x] create registerUserActivated handler
- [x] create tokens table
- [x] create login endpoint for activated users to generate tokens
<!--
    since users are activated with registerUserActivated
    they can generate tokens with /login endpoint
-->
- [-] generate and return new auth token on successful user registration
- [x] create registerUser endpoint with activated = false
    - [x] create mailer
    - [x] send the activation token in an email
    - [ ] if account not activated and trying to loging re send
            the activation email and respond with formative message
- [x] create activation endpoint
- [x] rate limiter
    - [-] limit unAuth users to send 10 requests per hour with 1 req/sec
    - [-] limit Auth users to send 100 requests per hour with 5 req/sec
- [ ] add new listings (create a scraper)
- [ ] docs
- [ ] maybe a frontend



