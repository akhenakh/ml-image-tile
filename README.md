# ml-image-tile
A tool to split images in tiles in preparation for machine learning work

# detect-blurry
Detect blurry images


## Example

```
 ./ml-image-tile -source /Volumes/Extreme SSD/ML/Bark -dest /Volumes/Extreme SSD/ML/BarkResized -height 224 -width 224 -resize 2 -rejectBlurry true -workerCount 20
{"app":"ml-image-tile","caller":"main.go:206","fileCounter":20253,"level":"info","rejectedBlurryCounter":2348,"tileCounter":695662,"ts":"2021-08-29T01:28:25.458398Z"}
 ```

Generate validation data, by creating a randomized tile.
 ```
 ./ml-image-tile -source /Volumes/Extreme\ SSD/ML/Bark   -dest  /Volumes/Extreme\ SSD/ML/BarkValidation -height 224 -width 224 -resize 2 -validationOnly -validationTileCount 2 -rejectBlurry -workerCount 20
  {"app":"ml-image-tile","caller":"main.go:206","fileCounter":20253,"level":"info","rejectedBlurryCounter":2348,"tileCounter":40300,"ts":"2021-08-29T00:32:05.7698Z"}
 ```
