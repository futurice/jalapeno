Feature: Create new recipes
	Scenario: Using CLI to create a new recipe
		Given a recipes directory
		When I create a recipe with name "foo"
		Then no errors were printed
		And the recipes directory should contain recipe "foo"

	Scenario: Using CLI to create a new recipe test
		Given a recipes directory
		And a recipe "foo"
		And recipe "foo" generates file "README.md" with content "initial"
		When I create a test for recipe "foo"
		Then CLI produced an output "Test 'example' created successfully!"
		And no errors were printed
		And the file "tests/example/test.yml" exist in the recipe "foo"

	Scenario: Using CLI to create a manifest
		Given a manifest directory
		When I create a manifest with the CLI
		Then CLI produced an output "Manifest created successfully!"
		And no errors were printed
		And the manifest directory should contain manifest named "manifest.yml"