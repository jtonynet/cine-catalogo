<a id="cine-catalog"></a>
###  :cinema: Go Catalog API for CineTicket pet-project :ticket:

<!-- 
    Logo image generated by Bing IA: https://www.bing.com/images/create/
-->
<img src="./misc/images/src/cinecatalog_title_logo_2.png"/>

<!-- 
    New AutoGo header test
    <img src="./misc/images/src/autogo_title_logo.png"/>
-->

<!-- 
    icons by: https://simpleicons.org
-->
[<img src="./misc/images/icons/go.svg" width="25px" height="25px" alt="go" title="Go">](https://go.dev/) [<img src="./misc/images/icons/docker.svg" width="25px" height="25px" alt="Docker" title="Docker">](https://www.docker.com/) [<img src="./misc/images/icons/dotenv.svg" width="25px" height="25px" alt="DotEnv" title="DotEnv">]() [<img src="./misc/images/icons/github.svg" width="25px" height="25px" alt="GitHub" title="GitHub">](https://github.com/jtonynet) [<img src="./misc/images/icons/miro.svg" width="25px" height="25px" alt="Miro" title="Miro">](https://miro.com/)  [<img src="./misc/images/icons/visualstudiocode.svg" width="25px" height="25px" alt="vscode" title="vscode">](https://code.visualstudio.com/) [<img src="./misc/images/icons/postman.svg" width="25px" height="25px" alt="Postman" title="Postman">](https://blog.postman.com/introducing-the-postman-vs-code-extension/)  [<img src="./misc/images/icons/postgresql.svg" width="25px" height="25px" alt="Postgres" title="Postgres">](https://www.postgresql.org/) [<img src="./misc/images/icons/redis.svg" width="25px" height="25px" alt="Redis" title="Redis">](https://redis.io/) [<img src="./misc/images/icons/swagger.svg" width="25px" height="25px" alt="Swagger" title="Swagger">](https://swagger.io/) [<img src="./misc/images/icons/gatling.svg" width="25px" height="25px" alt="Gatling" title="Gatling">](https://gatling.io/)  [<img src="./misc/images/icons/githubactions.svg" width="25px" height="25px" alt="GithubActions" title="GithubActions">](https://gatling.io/) 




![Badge Status](https://img.shields.io/badge/STATUS-IN_DEVELOPMENT-green)

 __This is an initial readme, here you can find the project's goals, and some features are not yet fully available. *__ 

---


<a id="index"></a>
## :arrow_heading_up: index

- [CineCatalog Microsservice](#cine-catalog)<br/>
  :arrow_heading_up: [index](#arrow_heading_up-index)<br/>
  :green_book: [About](#about)<br/>
  :umbrella: [Event Storming](#event-storming)<br/>
  :computer: [Run the project](#run)<br/>
  :bar_chart: [Diagrams](#diagrams)<br/>
  :newspaper: [API Documentation](#api-docs)<br/>
  :toolbox: [Tools](#tools)<br/>
  :clap: [Best Practices](#best-practices)<br/>
  :brain: [ADR - Architecture Decision Records](#adr)<br/>
  :1234: [Versions](#versions)<br/>

<br/>

[:arrow_heading_up: back to top](#index)

---

<a id="about"></a>
## :green_book: About:

This project aims to address the needs of cataloging cinema halls, movies, and sessions on a cinema ticket e-commerce website. It is part of a broader study of the mentioned e-commerce. However, its responsibility as microservices is to register, maintain, and provide session and seat data.

This is a Golang version of the mentioned service. Swagger Docs, Flow Diagrams, Entity-Relationship Diagrams (DER), and Event Storming provide more context to the service's scenario.

The objective of this system is to maintain a [high level of maturity](https://martinfowler.com/articles/richardsonMaturityModel.html) with a consistent RESTful API, along with the possibility of caching and a robust logging system.

<br/>

[:arrow_heading_up: back to top](#index)

---

<a id="event-storming"></a>
## :umbrella: Event Storming Diagram:

In November 2023, I received assistance from other developers in modeling the events of this project and other parts of 'cine-ticket.' We conducted an extensive remote Event Storming session with the aim of mapping events, commands, aggregates, and their relationships.

The diagram below is a product of this study and is being used as a guide for the development of this API and others that will be part of 'CineTicket.'

At the moment, we are abstracting the authentication flow and the ticket purchase flow.

<img src="./misc/images/src/event_storm_catalog_and_ticket_contexts.png"/>

<br/>

[:arrow_heading_up: back to top](#index)

---
<a id="run"></a>
## :computer: Run the project

Create a copy of the 'sample.env' file with the name '.env' and run the 'docker compose up' command (according to your 'docker compose' version) in the project's root directory:
```bash
$ docker compose up
```

> :writing_hand: **Note**:
>
> :window: Troubleshooting with [Windows](https://stackoverflow.com/questions/53165471/building-docker-images-on-windows-entrypoint-script-no-such-file-or-directory)
> Git attribute settings that might affect the line ending character are not working as expected. To run the project on Windows, you will need to make changes to the './tests/gatling/entrypoint.sh' file. Convert the file from 'LF' to 'CRLF' in your preferred text editor.

<br/>

[:arrow_heading_up: back to top](#index)

---

<a id="diagrams"></a>
## :bar_chart: System Diagrams:

**Flow Diagram:**

```mermaid
graph LR
  subgraph USER FLOW
    F([Logged Admin])
    F --> MAD(Maintain Addresses)
    F --> MFI(Maintain Films)
    F --> MSE(Maintain Sessions)
    F --> MRO(Maintain Rooms)
  end

  subgraph Catalog Microsservice
    TIC[Ticket]
    MFI --> FIL[Films]
    MRO --> ROO[Rooms]
    MSE --> SEA[Seats]
    MAD --> ADR[Adresses]
    MSE --> SES[Sessions]
  end

  %% Create associations between steps
    TIC -->|Associate| SES
    SEA -->|Associate| TIC
    SEA -->|Associate| ROO
    ROO -->|Associate| ADR
    FIL -->|Associate| SES
    SES -->|Associate| ROO
```
<br/><br/>

**DER:**

```mermaid
erDiagram 
    film {
        int id
        UUID uuid
        string description
        int age_rating
        boolean subtitled
        string poster
    }
    session {
        int id
        UUID uuid
        int film_id
        int room_id
        string description
        date date
        timestamptz start_time
        timestamptz end_time
        string time
    }
    room {
        int id
        UUID uuid
        string name
        int capacity
    }
    seat {
        int id
        UUID uuid
        int room_id
        string code
    }
    ticket {
        int id
        UUID uuid
        int session_id
        int seat_id
    }
    address {
        int id
        UUID uuid
        string country
        string state
        string zip_code
        string telephone
        string description
        string postal_code
        string name
    }
    room_address {
        int id
        UUID uuid
        int room_id
        int address_id
    }

    film ||--o{ session : has
    session ||--|| room : occurs
    room ||--|{ seat : has
    session ||--|{ ticket : has
    seat ||--|{ ticket : has
    room ||--|| room_address : located
    room_address ||--|| address : located
```

<br/>

[:arrow_heading_up: back to top](#index)

---
<a id="api-docs"></a>
## :newspaper: API Documentation

####  <img src="./misc/images/icons/swagger.svg" width="20px" height="20px" alt="Swagger" title="Swagger"> Generate Swagger docs:

With the 'cine-catalog' image running, type:

```bash
$ docker exec -ti cine-catalogue swag init --parseDependency --parseInternal  --generalInfo cmd/api/main.go
```

<br/>

[:arrow_heading_up: back to top](#index)

---
<a id="tools"></a>
## :toolbox: Tools

- Language:
  - [Go v1.21.1](https://go.dev/)
  - [GVM v1.0.22](https://github.com/moovweb/gvm)

- Framework & Libs:
  - [Gin](https://gin-gonic.com/)
  - [GORM](https://gorm.io/index.html)
  - [Viper](https://github.com/spf13/viper)
  - [Gin-Swagger](https://github.com/swaggo/gin-swagger)
  - [gjson](https://github.com/tidwall/gjson)
  - [uuid](github.com/google/uuid)

- Infra & Technologies
  - [Docker v24.0.6](https://www.docker.com/)
  - [Docker compose v2.21.0](https://www.docker.com/)
  - [Postgres v16.0](https://www.postgresql.org/)
  - [Redis 6.2](https://redis.io/)
  - [Gatling v3.9.5](https://gatling.io/)


- GUIs:
  - [VsCode](https://code.visualstudio.com/)
  - [Postman](https://blog.postman.com/introducing-the-postman-vs-code-extension/)
  - [DBeaver](https://dbeaver.io/)
  - [another-redis-desktop-manager](https://github.com/qishibo/AnotherRedisDesktopManager)


<br/>

[:arrow_heading_up: voltar](#indice)

---

<a id="best-practices"></a>
## :clap: Best Practices

- [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/)
- [keep a changelog](https://keepachangelog.com/en/1.0.0/)
- [Event Storming](https://en.wikipedia.org/wiki/Event_storming)
- [Miro Diagrams](https://miro.com/)
- [Mermaid Diagrams](https://mermaid.js.org)
- [Swagger](https://swagger.io/)
- [High Rest Maturity Model](https://martinfowler.com/articles/richardsonMaturityModel.html)

<!-- 
- [Load testing](https://en.wikipedia.org/wiki/Load_testing)
- [Go pprof](https://go.dev/blog/pprof)
-->

<br/>

[:arrow_heading_up: back to top](#index)

---

<a id="adr"></a>
## :brain: ADR - Architecture Decision Records:

- [0001: Record architecture decisions](./misc/architecture/decisions/0001-record-architecture-decisions.md)

<br/>

[:arrow_heading_up: back to top](#index)

---

<a id="adr"></a>
## :1234: Versions:

Version tags are being created manually as studies progress with notable improvements in the project. Each feature is developed on a separate branch, and when completed, a tag is generated and merged into the master branch.

For more information, please refer to the [Version History](./CHANGELOG.md)

<br/>

[:arrow_heading_up: back to top](#index)