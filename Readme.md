## DocDog

Creates RAML 1.0 files from existing source by Swen Kalski

[![Go](https://github.com/skalski/docdog/actions/workflows/go.yml/badge.svg)](https://github.com/skalski/docdog/actions/workflows/go.yml)

## How to use
Simple download the binaries or build from source.

There is one required Flag that should point to the root of source files:
`-path=/path/to/the/root/`

If no `-out` is provided, the RAML File will generated as `out.raml` in the folder, where the binary was called.

## Flags
* `-help` show flags
* `-path` root of source files to scan
* `-out` location and filename of output
* `-verbose` verbose-mode (default:false)
* `-tabs` specify the length of space that represent a tab (default:4)
* `-lang` specify the programming-language, filetype of source (default is `.java`)

## other languages
Actually DogDoc only provides JAVA and SpringBoot(java). Golang and RUST are up on the wishlist.
For SpringBoot just add `-lang=spring`. In this case there are no comments on Endpoint needed. They will ignored and
typical SpringBoot Commands will be used to generate the RAML Document.

## usage
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
*/
```

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
May I take fee later for companies!