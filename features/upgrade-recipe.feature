Feature: Upgrade sauce
  Upgrade a Jalapeno sauce

  Scenario: Upgrade sauce
    Given a project directory
    And a recipes directory
    And a recipe "foo" that generates file "README.md"
    And I execute recipe "foo"
    And I change recipe "foo" to version "v0.0.2"
    And I change recipe "foo" template "README.md" to render "New version"
    When I upgrade sauce "foo"
    Then the project directory should contain file ".jalapeno/sauce.yml" with "version: v0.0.2"
    And no errors were printed
    And the project directory should contain file "README.md" with "New version"
    And no conflicts were reported

  Scenario: Attempt upgrade when target file modified
    Given a project directory
    And a recipes directory
    And a recipe "foo" that generates file "README.md"
    And I execute recipe "foo"
    And I change recipe "foo" to version "v0.0.2"
    And I change recipe "foo" template "README.md" to render "New version"
    And I change project file "README.md" to contain "Locally modified"
    When I upgrade sauce "foo"
    Then conflicts are reported
    And the project directory should contain file "README.md" with "Locally modified"
