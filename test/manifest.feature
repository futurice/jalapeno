Feature: Execute manifests
  Execute a collection of recipes from a manifest file
  
  Scenario: Execute a manifest
    Given a project directory
    And a recipes directory
    And a recipe "foo" that generates file "README.md"
    And a recipe "bar" that generates file "Taskfile.yml"
    And a manifest file that includes recipes "foo" and "bar"
    When I execute the manifest file
    Then the project directory should contain file "README.md"
    And the project directory should contain file "Taskfile.yml"
