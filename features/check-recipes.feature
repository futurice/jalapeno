Feature: Check for new recipe versions
  Scenario: Find newer version for a recipe
    Given a project directory
    And a recipes directory
    And a recipe "foo" that generates file "README.md"
    And a local OCI registry
    When I execute recipe "foo"
    And the source of the sauce with recipe "foo" is in the local OCI registry
    Then execution of the recipe has succeeded
    When I push the recipe "foo" to the local OCI repository
    Then no errors were printed
    When I change recipe "foo" to version "v0.0.2"
    And I push the recipe "foo" to the local OCI repository
    And I check new versions for recipe "foo"
    Then CLI produced an output "new versions found: v0.0.2"

  Scenario: Find multiple newer version for a recipe
    Given a project directory
    And a recipes directory
    And a recipe "foo" that generates file "README.md"
    And a local OCI registry
    When I execute recipe "foo"
    And the source of the sauce with recipe "foo" is in the local OCI registry
    Then execution of the recipe has succeeded
    When I push the recipe "foo" to the local OCI repository
    Then no errors were printed
    When I change recipe "foo" to version "v0.0.2"
    And I push the recipe "foo" to the local OCI repository
    And I change recipe "foo" to version "v0.0.3"
    And I push the recipe "foo" to the local OCI repository
    Then I check new versions for recipe "foo"
    Then CLI produced an output "new versions found: v0.0.2, v0.0.3"

  Scenario: Find newer version for multiple recipes
    Given a project directory
    And a recipes directory
    And a recipe "foo" that generates file "foo.md"
    And a recipe "bar" that generates file "bar.md"
    And a local OCI registry
    When I execute recipe "foo"
    And the source of the sauce with recipe "foo" is in the local OCI registry
    Then execution of the recipe has succeeded
    When I execute recipe "bar"
    And the source of the sauce with recipe "bar" is in the local OCI registry
    Then execution of the recipe has succeeded
    When I push the recipe "foo" to the local OCI repository
    When I push the recipe "bar" to the local OCI repository
    Then no errors were printed
    When I change recipe "foo" to version "v0.0.2"
    And I push the recipe "foo" to the local OCI repository
    Then no errors were printed
    When I change recipe "bar" to version "v0.0.2"
    And I push the recipe "bar" to the local OCI repository
    Then no errors were printed
    Then I check new versions
    Then CLI produced an output "foo: new versions found: v0.0.2"
    And CLI produced an output "bar: new versions found: v0.0.2"

  Scenario: Unable to find newer recipe versions
    Given a project directory
    And a recipes directory
    And a recipe "foo" that generates file "README.md"
    And a local OCI registry
    When I push the recipe "foo" to the local OCI repository
    Then no errors were printed
    When I change recipe "foo" to version "v0.0.2"
    And I push the recipe "foo" to the local OCI repository
    And I execute recipe "foo"
    And the source of the sauce with recipe "foo" is in the local OCI registry
    Then execution of the recipe has succeeded
    And I check new versions for recipe "foo"
    Then CLI produced an output "no new versions found"
  
  Scenario: Unable to find newer recipe versions for all recipes
    Given a project directory
    And a recipes directory
    And a recipe "foo" that generates file "foo.md"
    And a recipe "bar" that generates file "bar.md"
    And a local OCI registry
    When I execute recipe "foo"
    And the source of the sauce with recipe "foo" is in the local OCI registry
    Then execution of the recipe has succeeded
    When I execute recipe "bar"
    And the source of the sauce with recipe "bar" is in the local OCI registry
    Then execution of the recipe has succeeded
    When I push the recipe "foo" to the local OCI repository
    When I push the recipe "bar" to the local OCI repository
    Then no errors were printed
    When I change recipe "foo" to version "v0.0.2"
    And I push the recipe "foo" to the local OCI repository
    Then no errors were printed
    Then I check new versions
    Then CLI produced an output "foo: new versions found: v0.0.2"
    And CLI produced an output "bar: no new versions found"

  Scenario: Executing remote recipe automatically adds the repo as source for the sauce
    Given a project directory
    And a recipes directory
    And a recipe "foo" that generates file "README.md"
    And a local OCI registry
    And the recipe "foo" is pushed to the local OCI repository "foo:v0.0.1"
    And I change recipe "foo" to version "v0.0.2"
    And the recipe "foo" is pushed to the local OCI repository "foo:v0.0.2"
    When I execute the recipe from the local OCI repository "foo:v0.0.1"
    Then execution of the recipe has succeeded
    # Note the lack of explicitly setting the source for the sauce
    And I check new versions for recipe "foo"
    Then CLI produced an output "new versions found: v0.0.2"