Feature: Running tests for a recipe
  Scenario: Using CLI to create a new recipe test
    Given a recipes directory
    And a recipe "foo" that generates file "README.md"
    When I create a placeholder test for recipe "foo" using the CLI
    Then CLI produced an output "Test created"
    And no errors were printed
    And the file "tests/example/test.yml" exist in the recipe "foo"

  Scenario: Run tests of the default recipe
    Given a recipes directory
    When I create a recipe with name "foo"
    And I run tests for recipe "foo"
    Then CLI produced an output "✅: defaults"
    And no errors were printed
