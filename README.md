# E-Commerce

This repository contains a microservice stack written in Go using Gin.

It's a port of Stephen Grider's [Microservices with Node JS and React](https://www.udemy.com/course/microservices-with-node-js-and-react/).
This repository is just a way to learn about microservices and is complete.

## Stack

- [Gin](https://github.com/gin-gonic/gin) as the HTTP server for all services.
- [Ent](https://entgo.io/) as the ORM.
- [Nats](https://nats.io/) for communiation between each service.
- [Stripe](https://stripe.com/) to handle payment.
- PostgreSQL to store the data.

Each directory contains a different service. The services communicate with each other using Nats.
