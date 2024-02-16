Feature: Running tests for a recipe
	Scenario: Using CLI to create a new recipe test
		Given a recipes directory
		And a recipe "foo" that generates file "README.md" with content "initial"
		When I create a test for recipe "foo"
		Then CLI produced an output "Test 'example' created successfully!"
		And no errors were printed
		And the file "tests/example/test.yml" exist in the recipe "foo"

	Scenario: Run tests of the default recipe
		Given a recipes directory
		When I create a recipe with name "foo"
		And I run tests for recipe "foo"
		Then CLI produced an output "✅: defaults"
		And no errors were printed

	Scenario: Tests fail if templates changes
		Given a recipes directory
		When I create a recipe with name "foo"
		And I change recipe "foo" template "README.md" to render "New version"
		And I run tests for recipe "foo"
		Then CLI produced an output "❌: defaults"
		And CLI produced an error "did not match for file 'README.md'"

	Scenario: Update test file snapshots
		Given a recipes directory
		When I create a recipe with name "foo"
		And I change recipe "foo" template "README.md" to render "New version"
		And I update tests snapshosts for recipe "foo"
		Then CLI produced an output "test snapshots updated"
		And CLI produced an output "README\.md \(modified\)"
		Then I run tests for recipe "foo"
		And CLI produced an output "No snapshot updates required"
