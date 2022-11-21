Feature: Jalapenoignore
	Ignore files from upgrades either by the recipe author in the recipe metadata, or by the user in a jalapenoignore file

	Scenario: Ignore in recipe metadata
		Given a project directory
		And a recipes directory
		And a recipe "foo" that generates file "README.md"
		And recipe "foo" ignores pattern "README.md"
		When I execute recipe "foo"
		Then the project directory should contain file "README.md" with "foo"
		When I change recipe "foo" to version "v0.0.2"
		And I change project file "README.md" to contain "bar"
		And I upgrade recipe "foo"
		Then the project directory should contain file "README.md" with "bar"
		And no conflicts were reported