# 2. Gin, Gorm and Postgres in three tier architecture 

Date: 2023-11-07

## Status

Accepted

## Context

We need to determine the framework, ORM, and database from the market options to start the project development. Additionally, we should define an architectural style that accommodates orderly growth, should the need arise in the future.

## Decision

Let's adopt the most widely used framework in the Golang community at the moment [Gin v1.9.1](https://github.com/gin-gonic/gin), with over __72K__ stars on [GitHub](https://github.com/gin-gonic) and a proven track record of usage. For the same reasons of adoption, we will also use the [Gorm v1.25.4](https://gorm.io/index.html) ORM with a [PostgreSQL v1.5.4](https://www.postgresql.org/) database, which will provide us with fast SQL queries on a scalable basis.

We decided to use [Docker](https://www.docker.com/) and [Docker Compose](https://docs.docker.com/compose/migrate/) from the very beginning to ease development. Additionally, we chose a [Two-tier architecture](https://en.wikipedia.org/wiki/Multitier_architecture#Three-tier_architecture) to maintain proper separation of concerns and adhered to the [Standard Go Project Layout](https://github.com/golang-standards/project-layout/blob/master/README.md) with over __42k__ stars on GitHub to ensure development uniformity.

As we develop, we intend to maintain the [RESTful](https://restfulapi.net/) standard as close to its original specifications as possible, using [HATEOAS](https://restfulapi.net/hateoas/) and achieving a [High Level of Maturity](https://martinfowler.com/articles/richardsonMaturityModel.html).

[Fiber](https://github.com/gofiber/fiber) and [Beego](https://github.com/beego/beego) were considered as frameworks but have lower adoption compared to our choice. The other decisions are based on the past experience of our developers and have proven to be successful in previous projects.


## Consequences

Ease and speed in development by using tools and patterns widely adopted by the community and with abundant reference material.
