# AggregatorProject
An API project I made in Go for aggregrating and scraping RSS feeds. Users can register an account and add links to RSS feeds, the server scrapes feeds periodically with a worker and users can fetch posts from feeds they follow through the /v1/posts endpoint. User records and feed data is stored in a postgres server.

## Endpoints

### POST /v1/users
Accepts a json object:

```json
{
    "name": "Shwads"
}
```

The server creates a record for the posted user and returns the full record in the response body:

```json
{
    "ID": "e939ef41-4070-4ce0-bfb0-0f84f3aaff90",
    "CreatedAt": "2024-04-30T16:54:01.926718Z",
    "UpdatedAt": "2024-04-30T16:54:01.926718Z",
    "Name": "Shwads",
    "ApiKey": "93c48b9f25f488729bc70ebaf6ee14fb50c1dbb15d14f283c62bc9753a5be77d"
}
``` 

### GET /v1/users
Authorised endpoint, accepts an authorisation header in the form:

```
    Authorization: ApiKey 93c48b9f25f488729bc70ebaf6ee14fb50c1dbb15d14f283c62bc9753a5be77d
```

Server checks the provided key with user's apikey and responds with the full user record if the authorization is successful:

```json
{
    "ID": "e939ef41-4070-4ce0-bfb0-0f84f3aaff90",
    "CreatedAt": "2024-04-30T16:54:01.926718Z", 
    "UpdatedAt": "2024-04-30T16:54:01.926718Z",
    "Name": "Shwads",
    "ApiKey": "3c48b9f25f488729bc70ebaf6ee14fb50c1dbb15d14f283c62bc9753a5be77d"
}
```

### POST /v1/feeds
Authorised endpoint, accepts an authorisation header, as the above endpoint, and a json object in the request body of the form:

```json
{
    "name": "The Boot.dev Blog",
    "url": "https://blog.boot.dev/index.xml"
}
```

Server stores the feed in the database, records that the user who posted the feed is following it and responds with the feed record as it appears in the database:

```json
{
    "ID": "dea35a86-5307-4ea4-a51d-34e092f9520d",
    "CreatedAt": "2024-05-02T15:42:03.14445Z",
    "UpdatedAt": "2024-05-02T15:42:03.14445Z",
    "Name": "The Boot.dev Blog",
    "Url": "https://blog.boot.dev/index.xml",
    "UserID": "e939ef41-4070-4ce0-bfb0-0f84f3aaff90"
}
```

### GET /v1/feeds
Simple GET endpoint. Returns all of the feed records from the database:

```json
[
  {
    "ID": "9c60b4be-5e82-4bc1-8544-281590697141",
    "CreatedAt": "2024-05-06T16:38:17.484683Z",
    "UpdatedAt": "2024-05-16T16:25:50.228634Z",
    "Name": "Ars Technica RSS Feed",
    "Url": "http://feeds.arstechnica.com/arstechnica/index/",
    "UserID": "f784b537-ceca-499c-9310-625ee3cf1298",
    "LastFetchedAt": "2024-05-16T16:25:50.228634Z"
  },
  {
    "ID": "dea35a86-5307-4ea4-a51d-34e092f9520d",
    "CreatedAt": "2024-05-02T15:42:03.14445Z",
    "UpdatedAt": "2024-05-16T16:25:50.362493Z",
    "Name": "The Boot.dev Blog",
    "Url": "https://blog.boot.dev/index.xml",
    "UserID": "e939ef41-4070-4ce0-bfb0-0f84f3aaff90",
    "LastFetchedAt": "2024-05-16T16:25:50.362493Z"
  },
  {
    "ID": "a50cc616-1e17-42e5-a782-44aca36ca4b5",
    "CreatedAt": "2024-05-09T20:46:47.361849Z",
    "UpdatedAt": "2024-05-16T16:25:50.470215Z",
    "Name": "Lane's Blog",
    "Url": "https://wagslane.dev/index.xml",
    "UserID": "e939ef41-4070-4ce0-bfb0-0f84f3aaff90",
    "LastFetchedAt": "2024-05-16T16:25:50.470215Z"
  }
]
```

### GET /v1/feed_follows
Authorised endpoint, accepts an authorisation header with a user's apikey. The server responds with a list of all the 'feed_follows' associated with that user (all of the records indicating which feeds a user is following) in the form:

```json
[
  {
    "ID": "5cd119a0-215d-4479-a85f-9ca41a9a726d",
    "FeedID": "9c60b4be-5e82-4bc1-8544-281590697141",
    "UserID": "e939ef41-4070-4ce0-bfb0-0f84f3aaff90",
    "CreatedAt": "2024-05-06T16:40:17.31492Z",
    "UpdatedAt": "2024-05-06T16:40:17.31492Z"
  },
  {
    "ID": "c2806c60-da7d-4637-bb8e-4d0886357eb9",
    "FeedID": "a50cc616-1e17-42e5-a782-44aca36ca4b5",
    "UserID": "e939ef41-4070-4ce0-bfb0-0f84f3aaff90",
    "CreatedAt": "2024-05-09T20:46:47.370677Z",
    "UpdatedAt": "2024-05-09T20:46:47.370677Z"
  }
]
```

### POST /v1/feed_follows
Authorised endpoint, accepts an authorisation header and a json object in the request body in the form:

```json
{
    "feed_id": "dea35a86-5307-4ea4-a51d-34e092f9520d"
}
```

The server then records that the posting user is following the feed identified by the request body.

### DELETE /v1/feed_follows/{feedFollowID}
Authorised endpoint (apikey). Server accepts a URL path parameter which identifies a particular feed_follow record in the database. The server then deletes that follow record and responds with a json object:

```json
{
    "status": "OK. Delete successful."
}
```

If the delete was successful.

### GET /v1/posts
Authorised endpoint (apikey). Accepts an optional object in the request body:

```json
{
    "limit": 5
}
```

If no such object is provided then the limit defaults to 10. The server then fetches the top 5 most recent blog posts from the feeds followed by the user and gives them in a response in the form:

```json
[
    {
        "ID": "ea793de3-c4ba-4898-af8e-fca1303be322",
        "CreatedAt": "2024-05-16T21:13:22.138203Z",
        "UpdatedAt": "2024-05-16T21:13:22.138203Z",
        "Title": "It could soon be illegal to publicly wear a mask for health reasons in NC",
        "Url": "https://arstechnica.com/?p=2025011",
        "Description": {
            "String": "Senators skeptical of legal trouble for harmless masking after moving to make it illegal.",
            "Valid": true
        },
        "PublishedAt": "2024-05-16T19:25:21Z",
        "FeedID": "9c60b4be-5e82-4bc1-8544-281590697141"
    },
    {
        "ID": "67529678-89f4-4929-b0df-6d10e6959951",
        "CreatedAt": "2024-05-16T21:13:22.164108Z",
        "UpdatedAt": "2024-05-16T21:13:22.164108Z",
        "Title": "Google Search adds a “web” filter, because it is no longer focused on web results",
        "Url": "https://arstechnica.com/?p=2024844",
        "Description": {
            "String": "Google Search now has an option to search the \"web,\" which is not the default anymore.",
            "Valid": true
        },
        "PublishedAt": "2024-05-16T19:18:49Z",
        "FeedID": "9c60b4be-5e82-4bc1-8544-281590697141"
    },
    {
        "ID": "4466b45a-a66c-40bc-b933-ba058179c883",
        "CreatedAt": "2024-05-16T19:46:09.646032Z",
        "UpdatedAt": "2024-05-16T19:46:09.646032Z",
        "Title": "Pedego Moto review: Fast and furious fun for $4,000",
        "Url": "https://arstechnica.com/?p=2024959",
        "Description": {
            "String": "Pedego's newest e-bike is quality even if a little bit impractical.",
            "Valid": true
        },
        "PublishedAt": "2024-05-16T18:39:25Z",
        "FeedID": "9c60b4be-5e82-4bc1-8544-281590697141"
    },
    {
        "ID": "929e63d3-f9d9-4e55-a319-af521ff56421",
        "CreatedAt": "2024-05-16T19:46:09.712387Z",
        "UpdatedAt": "2024-05-16T19:46:09.712387Z",
        "Title": "Tesla must face fraud suit for claiming its cars could fully drive themselves",
        "Url": "https://arstechnica.com/?p=2024945",
        "Description": {
            "String": "Lawsuit targets 2016 claim that all Tesla cars \"have full self-driving hardware.\"",
            "Valid": true
        },
        "PublishedAt": "2024-05-16T17:56:24Z",
        "FeedID": "9c60b4be-5e82-4bc1-8544-281590697141"
    },
    {
        "ID": "e05b910e-3975-4cc5-a1ce-929093ab07e0",
        "CreatedAt": "2024-05-16T19:46:09.713922Z",
        "UpdatedAt": "2024-05-16T19:46:09.713922Z",
        "Title": "Archie, the Internet’s first search engine, is rescued and running",
        "Url": "https://arstechnica.com/?p=2024876",
        "Description": {
            "String": "A journey through busted tapes, the Internet Old Farts Club, and SPARCstations.",
            "Valid": true
        },
        "PublishedAt": "2024-05-16T17:44:02Z",
        "FeedID": "9c60b4be-5e82-4bc1-8544-281590697141"
    }
]
```


