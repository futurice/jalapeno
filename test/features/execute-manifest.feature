Feature: Execute manifests
	Execute a collection of recipes from a manifest file

	Scenario: Execute a manifest
		Given a project directory
		And a recipes directory
		And a recipe "foo"
		And recipe "foo" generates file "README.md" with content "initial"
		And a recipe "foo"
		And recipe "foo" generates file "Taskfile.yml" with content "initial"
		And a manifest file that includes recipes "foo" and "bar"
		When I execute the manifest file
		Then the project directory should contain file "README.md"
		And the project directory should contain file "Taskfile.yml"
