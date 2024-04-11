Feature: Execute manifests
	Execute a collection of recipes from a manifest file

	Scenario: Execute a manifest
		Given a project directory
		And a recipes directory
		And a recipe "foo"
		And recipe "foo" generates file "README.md" with content "initial"
		And a recipe "bar"
		And recipe "bar" generates file "Taskfile.yml" with content "initial"
		And a manifest file that includes recipes
		| foo | bar |
		When I execute the manifest file
		Then no errors were printed
		And CLI produced an output "^Executing manifest with 2 recipes"
		And CLI produced an output "Recipe name: foo"
		And CLI produced an output "Recipe name: bar"
		And the project directory should contain file "README.md"
		And the project directory should contain file "Taskfile.yml"
