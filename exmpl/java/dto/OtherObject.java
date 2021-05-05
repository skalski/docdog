package com.test.dto;

@Getter
@Setter
class OtherObject implements OtherAbsObject {

public String somePublic = "test";

private String someVar;

/*
    @DD:DESCRIPTION 'some var we use'
    @DD:NOTNULL
*/
private int someElse;

/*
    @DD:IGNORE
*/
private bool ignoreThis;
}