# go-hackernews

Quiet Hackernews is a clone of Hackernews, but with Caching and Concurrency.\
It implements a Active Expiry caching mechanism that expires the cache after 10 seconds, and re-populates the cache 5 seconds before it expires to provide latest data with 0 delay.

## Tech Stack

**Go** is used as the backend, making use of its concurrency techniques to make page reloads blazing fast. The frontend is powered by **HTML** with Go templating.
