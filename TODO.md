# TODO

- [x] DP savefile R/W
- [x] HGSS savefile R/W
- [ ] Finish BW + B2W2 savefile R/W
- [ ] ripper file fetching optimization
    - currently, each call to a ripper function will read the entire ROM into memory. 
      utilize caching to reduce file I/Os for successive function calls
- [ ] update README for example usage
- [ ] release as v8