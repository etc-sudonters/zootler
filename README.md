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
p := zootler.NewHashPool()
for i = 0; i <= 10_000; i++ {
    entity := p.Creat()
    if i % 13 == 0 {
        entity.Add(LuckyComponent{})
    }
}

for _, lucky := p.Query(LuckyComponent{}) {
    ...
}
```

- Could be extended to locations as well. For example, some spots are cloakable
  replicating a similar system for locations (or unifying the two) might have
  interesting results:

  ```go
  p := zootler.NewHashPool()
  for _, unfilledCloakableSpots := p.Query(
    CheckComponent,
    zootler.Include[CloakableComponent]{},
    zootler.Exclude[FilledComponent]{},
    zootler.Include[Reachable]{},
    ) {
        ...
    }
    
  for _, unfilledAdultSongSpots := p.Query(
    CheckComponent,
    zootler.Exclude[FilledComponent]{},
    zootler.Include[SongComponent]{},
    zootler.Include[AdultExclusiveComponent]{},
    ) { ... }
  ```
