# 2. Gin boot as main framework 

Date: 2023-11-06

## Status

Accepted

## Context

We need to record the architectural decisions made on this project.

## Decision

Let's adopt the most widely used framework in the Golang community at the moment [Gin](https://gin-gonic.com/docs/), with over 79K stars on [GitHub](https://github.com/gin-gonic) and a proven track record of usage. For the same reasons of adoption, we will also use the gorm ORM with a PostgreSQL database, which will provide us with fast SQL queries on a scalable basis.

We will use Architecture Decision Records,
as [Documentation](https://gin-gonic.com/docs/).

## Consequences

See Michael Nygard's article, linked above. For a lightweight ADR toolset, see Nat
Pryce's [adr-tools](https://github.com/npryce/adr-tools).
