Feature: Upgrade recipe
  Upgrade a rendered Jalapeno recipe

  Scenario: Upgrade rendered recipe
    Given a project directory
    And a recipes directory
    And a recipe "foo" that generates file "README.md"
    When I execute recipe "foo"
    And I change recipe "foo" to version "v0.0.2"
    And I upgrade recipe "foo"
    Then the project directory should contain file ".jalapeno/1-foo.yml" with "version: v0.0.2"
