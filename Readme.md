
## DocDog - the code to RAML converter
![logo](./artwork/docdog-small.png) 

Creates RAML 1.0 files from existing Java and Spring Boot sources by Swen Kalski

Twitter: https://twitter.com/KalskiSwen

[![Go](https://github.com/skalski/docdog/actions/workflows/go.yml/badge.svg)](https://github.com/skalski/docdog/actions/workflows/go.yml)

## Features

* make RAML 1.0 files including the scheme from plain JAVA
* make RAML 1.0 files including the scheme without using tags from Spring Boot 
* understand deeper structures of java including Interfaces and implementations
* generated files can used in API Console without hassle
* can pipe to CLI for other cmd tools
* It's BETA so everything that happen is a feature :-)


## Download
[See latest release](https://github.com/skalski/docdog/releases) 

Binaries available for Mac OS (ARM/intel), Linux (amd/intel 64bit) and Windows 64bit

## Build
To build the latest source for your own, you can use `./build_all.sh` on Linux, Mac OS or even
Windows with Linux Subsystem.
You have to use Go 1.13 at least.

## How to use
Simply download the binaries or build from source.

There is one required Flag that should point to the root of source files:
`-path=/path/to/the/root/`

If no `-out` is provided, the RAML File will be generated as `out.raml` in the folder, where the binary was called.

## Flags
* `-help` show help
* `-path` root of source files to scan
* `-out` location and filename of output
* `-verbose` verbose-mode for debug purpose (default:false)
* `-tabs` specify the length of spaces that represent a tab (default:4)
* `-lang` specify the programming-language, filetype of source (default is `.java`)
* `-print` print the output directly to the cli for piping (no file will be created)

If `-print=true` the `-out` flag will be ignored. No output file will be created.

## other languages
Actually DogDoc only provides JAVA and SpringBoot(java). Golang and RUST are up on the wishlist.
For SpringBoot just add `-lang=spring`. In this case there are no comments on Endpoint needed. They will ignored and
typical SpringBoot Commands will be used to generate the RAML Document.

## How to use Comment-Annotations
You can see, how the Comments for DocDog are used in the `\exmpl\java` folder.

### API Endpoints
Every endpoint of your API must be Marked for DocDog to find it:

Example:
```
/*
@DD:ENDPOINT 'api/testpoint'
@DD:DESCRIPTION 'important testendpoint'
@DD:PARAM int id 'id of user'
@DD:PARAM string token 'security-token' @DD:NOTNULL
@DD:PAYLOAD testObject 'json object'
@DD:TYPE post
@DD:RESPONSE 200 json ResponseObject
@DD:RESPONSE 500 text
*/
```

The types of annotations that could be used for Endpoints.
* `@DD:ENDPOINT '<string>'` Declare an Endpoint and use <string> as part of URL
* `@DD:DESCRIPTION '<string>'` Set a description for this endpoint
* `@DD:PARAM <dataTyp> <varName> '<string>' <notNull:optional>` Add a param for the endpoint (@DD:NOTNULL add a required-tag)
* `@DD:PAYLOAD <dataTyp> '<string>'` add a body and a description for this payload
* `@DD:TYPE post` http request type (post, get, delete ..)
* `@DD:RESPONSE <int> <json/text> <dataType:when json>` add a type of response and his type. If json is used you must provide a datatype.


### Objects 
Objects will be found by DocDog himself.
You can markup variables in Objects too.

Example to add a Description and mark it as required:
```java
/*
    @DD:DESCRIPTION 'some var we use'
    @DD:NOTNULL
*/
private int someElse;
```

Also, you can say DocDog to ignore a variable in an Object:
```java
/*
    @DD:IGNORE
*/
private bool ignoreThis;
```

## License

Actually I got no idea.
For private this will stay free.
This means - you can use it, fork it or whatever. If the source is used elsewhere the original developer (me) must be mentiond.
If it is used in a commercial software you must pay me... that's simple :)

## Contributions 

The "DocDog Dog" was created by Jasmin Kalski
