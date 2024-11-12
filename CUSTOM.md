How to customize
================

- clone repo (can be private)
- inside repo copy cmd/swamp into your new custom app - i.e. cmd/lake
- add/edit stuff to your cmd/lake
  - default swamp_repos.yml
  - static/ - files for web
  - templates/ - you would need to copy templates from root of repo. Note you may use swamp-ui-dev which generates fake data but uses templates from current working directory
  - remove annoying intro from main page (index.html)
- commit/tag to your repo
- workflow would create new relase for you
