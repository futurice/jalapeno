Feature: Create new recipes
	Scenario: Using CLI to create a new recipe
		When I create a recipe with name "foo"
		Then no errors were printed
		And the recipes directory should contain recipe "foo"

	Scenario: Using CLI to create a new recipe test
		Given a recipe "foo"
		And recipe "foo" generates file "README.md" with content "initial"
		When I create a test for recipe "foo"
		Then CLI produced an output "Test 'example' created successfully!"
		And no errors were printed
		And the file "tests/example/test.yml" exist in the recipe "foo"

	Scenario: Using CLI to create multiple recipe tests
		Given a recipe "foo"
		And recipe "foo" generates file "README.md" with content "initial"
		When I create a test for recipe "foo"
		Then CLI produced an output "Test 'example' created successfully!"
		When I create a test for recipe "foo"
		Then CLI produced an output "Test 'example_1' created successfully!"
		When I create a test for recipe "foo"
		Then CLI produced an output "Test 'example_2' created successfully!"
		And no errors were printed
		And the file "tests/example/test.yml" exist in the recipe "foo"
		And the file "tests/example_1/test.yml" exist in the recipe "foo"
		And the file "tests/example_2/test.yml" exist in the recipe "foo"

	Scenario: Using CLI to create a manifest
		When I create a manifest with the CLI
		Then CLI produced an output "Manifest created successfully!"
		And no errors were printed
		And the manifest directory should contain manifest named "manifest.yml"