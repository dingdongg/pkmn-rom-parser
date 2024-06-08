# Known Bugs


| Related feature                  | Erroneous behavior                                    | Desired behavior                                                                 | Steps to reproduce                                        |
| -------------------------------- | ----------------------------------------------------- | -------------------------------------------------------------------------------- | --------------------------------------------------------- |
| Updating a party pokemon's level | doesn't update EXP points accordingly                 | EXP points should be set to 0 (relative to the pokemon's level bar)              | Modify a party pokemon's level.                           |
| Displaying party pokemon         | If a party has < 6 pokemon, may not render correctly? | Pokemon party should render correctly, regardless of the number of party pokemon | read a savefile where the player has < 6 pokemon in party |