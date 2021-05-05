package payment.test

@Getter
@Setter
public class RestDeploymentConfigurationModel {

public String somePublic = "test";

private String someVar;

@NotNull
private int someElse;

@JsonIgnore
private bool ignoreThis;

private List<otherObject> testList;

private String[] cars = {"Volvo", "BMW", "Ford", "Mazda"};

    private function getSomeElse(int someElse){
        this.someElse = someElse;
    }
}