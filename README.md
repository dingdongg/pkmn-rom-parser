# Pokemon Platinum ROM Parser
Application used to parse in-game information in generation IV/V Pokemon games.

## Currently supported games
- Pokemon Platinum only (support for other games will be implemented in future releases)

## Features
- Read/update party pokemon stats, including:
    - level
    - name
    - EVs/IVs
    - held item
    - nature
    - battle stats
- checksum validations, safe from memory corruptions!

## TODO
- extend support for other gen. 4/5 games
- Read/update PC system pokemon too

### Credits
---
The information in `char_encoder/char_encoder.go` was extracted from [this Bulbapedia article](https://bulbapedia.bulbagarden.net/wiki/Character_encoding_(Generation_IV)) using a custom script.

The `tableValues` information in `shuffler/shuffler.go` was extracted from [Project Pokemon](https://projectpokemon.org/home/docs/gen-4/pkm-structure-r65/) using a custom HTML parsing script.