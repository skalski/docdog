package something.packages

import this;
import that;

@Getter
@Setter
public class TestObject {

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

private List<OtherObject> testList;

private String[] cars = {"Volvo", "BMW", "Ford", "Mazda"};

    private function getSomeElse(int someElse){
        this.someElse = someElse;
    }
}