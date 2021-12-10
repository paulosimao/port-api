# ports-api
Marking 2 hours exercise

# Beyond 2 hours exercise
To run the following solution, please use `make docker_run`. It will create a network, the containers, and start them.

To test: 
    - make test_up will upload the json file and upsert it into the DB.
    - make test_down will download data from the DB

Caveats:
    - I am on Apple M1, thus compiling w SQLITE for Docker would take a little longer
    - The format of using s json string for nested data was a quick wayout considering time constraints. In an ideal world, this would be different, either flattening the object to fit into a SQL DB, or using NoSql like Mongo
    - The format adopted to upload data does not help too much, considering that its a KV object. Streaming under these conditions is harder to achieve. So, either we adopt a different format (jsonlines maybe), or we do the parser ourselves. Having said that - the present solution will not fit into memory constrained env, if a too large file is provided. We are limiting files to 10 MB at this stage.
    - Design around DB VOs and ProtoBuffer TOs could be better
# Journal
## 3:15
    Added Signal Handling
    Fixed cancel handling
    Fixed docker-compose, due to M1 constraint
