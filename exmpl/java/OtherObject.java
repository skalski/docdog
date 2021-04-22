package something.packages

import that;

@Getter
@Setter
public class OtherObject {

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