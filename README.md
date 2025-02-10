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


## Development

### To Do

- [x] Create database on target server
- [x] Form validation
    - [x] New Server
    - [x] New Database
    - [x] Unit tests
- [ ] Sanitize queries
- [ ] Transactions around create database operations
- [ ] Git versioning
- [ ] Build docker container
    - [ ] Probably need to fix path issues with templates and static files, so can be run outside of Air.
- [ ] User documentation
- [ ] Error handling (better than `log.Fatal`)
- [ ] API
- [ ] Server password encryption
- [ ] Authn/Authz
    - [ ] Users
    - [ ] Roles
    - [ ] Permissions
    - [ ] API Tokens
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