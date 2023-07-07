Todo List:

- [ ] actually generate a seed from ootr itself
    - [ ] Figure out how to just generate spoilers
    - [ ] are there any good spoiler analyzers?
        - wanna throw a spoiler at a thing have it tell me things
- [ ] parse logic files
    - [ ] python expression parser (neat!)
- [ ] token description pool
    - pretty sure this is all a big dictionary in zootr code
    - crafty json.dumps shenanigans might do the trick
    - add in whatever I need for my own purposes
    - then load from that
- [ ] how to generate my own spoiler log

Idea:

- Tokens are accessed in ECS-esque fashion
```go
tokenPool := &zooty.TokenPool{}
// mint several tokens

// find all available (for placement) tokens that are not labeled advancement
// and are labeled priority. in zootr, ice arrows, stone of agony, double def
matchingSet := tokenPool.Match(
    zooty.NotTagged("advancement"),
    zooty.Tagged("priority"),
    zooty.Tagged("available"),
)
```

- Could be extended to locations as well. For example, some spots are cloakable
  replicating a similar system for locations (or unifying the two) might have
  interesting results:

  ```go
  locationPool := &zooty.Pool{}
  // mint...
  matchingSet := locationPool.Match(
    zooty.Tagged("cloakable"),
    zooty.Tagged("filled"),
    //zooty.Tagged("shop")
  )

  _ := locationPool.Match(
    zooty.NotTagged("filled"),
    zooty.Tagged("song"),
    zooty.Tagged("adult"),
  ) // aka warp songs
  ```
