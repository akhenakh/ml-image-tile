# ml-image-tile
A tool to split images in tiles in preparation for machine learning work

## Debug
```
 dlv debug --headless --listen ":2345" --log --api-version 2 github.com/akhenakh/ml-image-tile  --  -source=./fakedir -dest=/tmp -logLevel=DEBUG
 ```