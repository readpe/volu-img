## Installation
```bash
$ go get github.com/readpe/volu-img
```

## Usage

The volu-img is a command line tool with two flags, these can be viewed using the `--help` flag as detailed below.

```bash
$ volu-img --help

Usage of ./volu-img:
  -i string
        json input file path
  -p string
        path to images (default ".")
```

The required `-i` flag requests the path to the JSON input file, whos format is detailed below. The `-p` flag request the path to the directory containing the source images, this option will default to the current working directory if no path is provided. 


### Input JSON File

The details for each product image are provided in a JSON file for input using the `-i` command flag, this is required. The `*.jpg` files specified should be located withing the provided `-p` file path.

- `sku`: Represents the product code
- `img`: Is the "main" product image
- `img_alts`: List of product alternative images
- `large`, `medium`, `small`, `tiny`, `thumb` image sizes
    - If height or width is zero, the non-zero value will scale, maintaining aspect ratio. Equivalent to 'max height' or 'max width' option
    - If both height and width are non-zero, the maximum image dimension will scale. Equivalent to 'both' option

```json
[
    {
        "sku":"PRODUCT-1",
        "img":"PRODUCT-1.jpg",
        "img_alts":[
            "PRODUCT-1-ALT.jpg"
        ],
        "large": { "width": 0, "height": 700},
        "medium": { "width": 400, "height": 400},
        "small": { "width": 500, "height": 500},
        "tiny": { "width": 100, "height": 100},
        "thumb": { "width": 0, "height": 50}
    },
    {
        "sku":"PRODUCT-2",
        "img":"PRODUCT-2.jpg",
        "img_alts":[
            "PRODUCT-2-ALT.jpg"
        ],
        "large": { "width": 0, "height": 700},
        "medium": { "width": 400, "height": 400},
        "small": { "width": 500, "height": 500},
        "tiny": { "width": 100, "height": 100},
        "thumb": { "width": 0, "height": 50}
    }
]
```

### Output
A new directory 'volu' will be created within the `-p` path provided, this directory will hold the resized images.

Image nameing convention is based on the following post:

https://helpcenter.volusion.com/en/articles/1773795-product-image-file-names


## Acknowledgements
This program uses the github.com/nfnt/resize module for resizing.
