Feature: Upgrade sauce
  Upgrade a Jalapeno sauce

  Scenario: Upgrade sauce
    Given a project directory
    And a recipes directory
    And a recipe "foo" that generates file "README.md"
    And I execute recipe "foo"
    And I change recipe "foo" to version "v0.0.2"
    And I change recipe "foo" template "README.md" to render "New version"
    When I upgrade recipe "foo"
    Then no errors were printed
    And the project directory should contain file ".jalapeno/sauces.yml" with "version: v0.0.2"
    And the project directory should contain file "README.md" with "New version"
    And no conflicts were reported

  Scenario: Upgrade sauce from remote recipe
    Given a project directory
    And a recipes directory
    And a local OCI registry
    And a recipe "foo" that generates file "README.md"
    And I execute recipe "foo"
    And I change recipe "foo" to version "v0.0.2"
    And I change recipe "foo" template "README.md" to render "New version"
    And the recipe "foo" is pushed to the local OCI repository "foo:v0.0.2"
    When I upgrade recipe from the local OCI repository "foo:v0.0.2"
    Then no errors were printed
    And the project directory should contain file ".jalapeno/sauces.yml" with "version: v0.0.2"
    And the project directory should contain file "README.md" with "New version"
    And no conflicts were reported

  Scenario: Attempt upgrade when previous sauce file was modified
    Given a project directory
    And a recipes directory
    And a recipe "foo" that generates file "README.md"
    And I execute recipe "foo"
    And I change recipe "foo" to version "v0.0.2"
    And I change recipe "foo" template "README.md" to render "New version"
    And I change project file "README.md" to contain "Locally modified"
    When I upgrade recipe "foo"
    Then CLI produced an error "file conflicts"
    And the project directory should contain file "README.md" with "Locally modified"
  
  Scenario: Attempt upgrade when new file conflicts with existing manually created file
    Given a project directory
    And a recipes directory
    And a recipe "foo" that generates file "README.md"
    And I execute recipe "foo"
    And I create a file "new.txt" with contents "manual" to the project directory
    And I change recipe "foo" to version "v0.0.2"
    And I change recipe "foo" template "new.txt" to render "new"
    When I upgrade recipe "foo"
    Then CLI produced an error "file conflicts"
