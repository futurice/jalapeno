Feature: Jalapenoignore
	Ignore files from upgrades either by the recipe author in the recipe metadata, or by the user in a jalapenoignore file

	Scenario: Ignore in recipe metadata
		Given a project directory
		And a recipes directory
		And a recipe "foo"
		And recipe "foo" generates file "README.md" with content "initial"
		And recipe "foo" ignores pattern "README.md"
		When I execute recipe "foo"
		Then no errors were printed
		When I change project file "README.md" to contain "modified"
		And I change recipe "foo" to version "v0.0.2"
		And I upgrade recipe "foo"
		Then no conflicts were reported
		And no errors were printed
		And the project directory should contain file "README.md" with "modified"

	Scenario: Ignore with jalapenoignore file
		Given a project directory
		And a recipes directory
		And a recipe "foo"
		And recipe "foo" generates file "README.md" with content "initial"
		When I execute recipe "foo"
		And I change project file "README.md" to contain "modified"
		And I change recipe "foo" to version "v0.0.2"
		And I change project file ".jalapenoignore" to contain "*.md"
		And I upgrade recipe "foo"
		Then no conflicts were reported
		And no errors were printed
		And the project directory should contain file "README.md" with "modified"

	Scenario: Ignored file will not be removed even if new recipe version deprecates it
		Given a project directory
		And a recipes directory
		And a recipe "foo"
		And recipe "foo" generates file "README.md" with content "initial"
		And recipe "foo" generates file "will-be-removed-in-next-version" with content "initial"
		And recipe "foo" ignores pattern "will-be-removed-in-next-version"
		When I execute recipe "foo"
		Then no errors were printed
		When I change project file "will-be-removed-in-next-version" to contain "modified"
		And I change recipe "foo" to version "v0.0.2"
		And I remove file "will-be-removed-in-next-version" from the recipe "foo"
		And I upgrade recipe "foo"
		Then no conflicts were reported
		And no errors were printed
		And the project directory should contain file "will-be-removed-in-next-version" with "modified"
		And the sauce file contains a sauce in index 0 which should not have property "files.will-be-removed-in-next-version"