# Pokemon Platinum ROM Parser
Application used to parse in-game information in the Pokemon Platinum ROM for the NDS.

## TODO
- extend support for other gen. 4/5 games
- read from the more recent savefile chunk, instead of whatever is stored in the first small block in memory
- ROM writer package?

### Credits
---
The information in `char_encoder/char_encoder.go` was extracted from [this Bulbapedia article](https://bulbapedia.bulbagarden.net/wiki/Character_encoding_(Generation_IV)) using a custom script.

The `tableValues` information in `shuffler/shuffler.go` was extracted from [Project Pokemon](https://projectpokemon.org/home/docs/gen-4/pkm-structure-r65/) using a custom HTML parsing script.