Feature: Execute recipes
  Executing Jalapeno recipes to template out projects

  Scenario: Execute single recipe
    Given a project directory
    And a recipes directory
    And a recipe "foo" that generates file "README.md"
    When I execute recipe "foo"
    Then execution of the recipe has succeeded
    And the project directory should contain file "README.md"
    And the project directory should contain file ".jalapeno/1-foo.yml"

  Scenario: Execute multiple recipes
    Given a project directory
    And a recipes directory
    And a recipe "foo" that generates file "README.md"
    And a recipe "bar" that generates file "Taskfile.yml"
    When I execute recipe "foo"
    Then execution of the recipe has succeeded
    When I execute recipe "bar"
    Then execution of the recipe has succeeded
    And the project directory should contain file "README.md"
    And the project directory should contain file "Taskfile.yml"
    And the project directory should contain file ".jalapeno/1-foo.yml"
    And the project directory should contain file ".jalapeno/2-bar.yml"

  Scenario: New recipe conflicts with the previous recipe
    Given a project directory
    And a recipes directory
    And a recipe "foo" that generates file "README.md"
    And a recipe "bar" that generates file "README.md"
    When I execute recipe "foo"
    Then execution of the recipe has succeeded
    When I execute recipe "bar"
    Then execution of the recipe has failed with error "README.md was already created by recipe foo"