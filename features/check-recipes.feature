Feature: Check for new recipe versions
  Scenario: Find newer version for a recipe
    Given a project directory
    And a recipes directory
    And a recipe "foo" that generates file "README.md"
    And a local OCI registry
    And the source of the recipe "foo" is in the local OCI registry
    When I execute recipe "foo"
    Then execution of the recipe has succeeded
    When I push the recipe "foo" to the local OCI repository
    Then push of the recipe was successful
    When I change recipe "foo" to version "v0.0.2"
    And I push the recipe "foo" to the local OCI repository
    And I check new versions for recipe "foo"
    Then CLI produced an output "New versions found"

  Scenario: Unable to find newer recipe versions
    Given a project directory
    And a recipes directory
    And a recipe "foo" that generates file "README.md"
    And a local OCI registry
    And the source of the recipe "foo" is in the local OCI registry
    When I push the recipe "foo" to the local OCI repository
    Then push of the recipe was successful
    When I change recipe "foo" to version "v0.0.2"
    And I push the recipe "foo" to the local OCI repository
    And I execute recipe "foo"
    Then execution of the recipe has succeeded
    And I check new versions for recipe "foo"
    Then CLI produced an output "No new versions found"
