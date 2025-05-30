# Soda
Soda is a Self-hOsted Database As a service.

## A what?!?
A self-hosted database as a service.

But, isn't that a pointless contradiction? Why would you want to self-host
a database as a service?

Well, maybe you want to allow your team to provision databases in a
controlled way, or maybe you want to play around with terraform in
the lab like I do.

## So what's the plan?
The plan is to implement a simple to deploy, simple to use web app, with 
an API that can be used from a terraform provider and the like.

## This is a learning project
Soda is primarily a vehicle to allow me to learn Go, so priorities might seem
wrong, such as passwords not being encrypted in the database. It's not that I 
don't know that's bad, I'm just focussed on learning basic syntax and structure 
for now and I'll go down that rabit hole when I get to it. Before then I want
to get a good general knowledge of Go fundamentals.

> [!CAUTION]
> Soda is not production ready and is not meant to be
> deployed anywhere of any importance.

## Development
Below are some of the things I might get round to, in rough, but not definate order.

To run the project locally for development do

```
$ make start-soda-db
$ make soda-dev
```

### To Do
- [x] Create database on target server
- [x] Form validation
    - [x] New Server
    - [x] New Database
    - [x] Unit tests
- [x] API
- [x] Build docker container
    - [x] Probably need to fix path issues with templates and static files, so can be run outside of Air.
    - [x] Migrate system database on startup
- [ ] Git versioning
  - [x] Docker container
  - [ ] Application binary
- [x] Upgrade logging
- [x] Sanitize queries
- [ ] Transactions around create database operations
- [ ] User documentation
- [x] Error handling (better than `log.Fatal`)
- [ ] Server password encryption
- [ ] Dependency injection
- [ ] Migrate to one of the web frameworks
- [ ] Support commandline args for logging to file
- [ ] Authn/Authz
    - [ ] Users
    - [ ] Roles
    - [ ] Permissions
    - [ ] API Tokens
- [ ] Server status
  - [ ] Server list
  - [ ] Server details
  - [ ] Database list
  - [ ] Database details
  - [ ] API??
- [ ] Server stats
- [ ] Database stats
- [ ] Support multiple databases / provider model
- [ ] Row versioning (for change history / revert)

### Database Migrations
Migrations arevia https://github.com/golang-migrate/migrate. 
There are a bunch of targets in the makefile to make working with 
the database simpler.

- `make migration {name}` - adds a migration.
- `make migrate-up [count=$count]` - migrates up `$count` migrations, or all 
the way up if `count` is omitted.
- `make migrate-up-one` - migrates up one version.
- `make migrate-down [count=$count]` - migrates down `$count` migrations, or all
 the way down if `count` is omitted.
- `make migrate-down-one` - migrates down one version.
- `make migrate-test-last` - test the most recent migration by applying it, 
rolling it back, and applying it again.