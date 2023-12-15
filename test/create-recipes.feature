Feature: Create new recipes
  Scenario: Using CLI to create a new recipe
    Given a recipes directory
    When I create a recipe with name "foo"
    Then no errors were printed
    And the recipes directory should contain recipe "foo"
