package payroll.test.controller;

import java.util.List;
import payment.test;

import org.springframework.web.bind.annotation.DeleteMapping;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.PutMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RestController;

@RestController
class EmployeeController {

  private final EmployeeRepository repository;

  EmployeeController(EmployeeRepository repository) {
    this.repository = repository;
  }

  @GetMapping("/employees/{id}")
  function Employee one(@PathVariable Long id) {

    return repository.findById(id)
      .orElseThrow(() -> new EmployeeNotFoundException(id));
  }

  @PutMapping("/employees/add/{id}")
  function Employee replaceEmployee(@RequestBody Employee newEmployee, @PathVariable Long id) {

    return repository.findById(id)
      .map(employee -> {
        employee.setName(newEmployee.getName());
        employee.setRole(newEmployee.getRole());
        return repository.save(employee);
      })
      .orElseGet(() -> {
        newEmployee.setId(id);
        return repository.save(newEmployee);
      });
  }

  @DeleteMapping("/employees/delete/{id}")
  void deleteEmployee(@PathVariable Long id) {
    repository.deleteById(id);
  }


    @RequestMapping(value = "/{deploymentId}/updateDeploymentConfig", method = RequestMethod.POST, consumes = MediaType.APPLICATION_JSON)
    @ResponseBody
    public ConfigurationModel deploymentConfigRest(
                                                                   @PathVariable final long deploymentId,
                                                                   @Valid @RequestBody final RestDeploymentConfigurationModel restDeploymentConfigurationModel,
                                                                   final BindingResult result
                                                                   ) throws ValidationException {
        if (result.hasErrors()) {
            throw new ValidationException(result);
        }
        return updateDeploymentConfig();
    }
}