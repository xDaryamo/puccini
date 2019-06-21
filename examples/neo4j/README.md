Neo4j Example
-------------

There are several ways to get data into [Neo4j](https://neo4j.com/).

We have here chosen to use the [HTTP API](https://neo4j.com/docs/http-api/3.5/) to send a series
of [Cypher CREATE](https://neo4j.com/docs/cypher-manual/3.5/clauses/create/) commands packed
into a JSON POST, which would create the data in a single transaction.

(Note that Neo4j's [RESTful API](https://neo4j.com/docs/rest-docs/3.5/) is being deprecated and
should not be used.)

Another approach could be to use the `neo4j-admin` tool to
[import data](https://neo4j.com/docs/operations-manual/3.5/tools/import/) in CSV format.

You can visualize your results using Bolt, Neo4j's web GUI, at 
[`http://localhost:7474/browser/`](http://localhost:7474/browser/).
