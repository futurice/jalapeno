Feature: Running tests for a recipe

	Scenario: Run tests of the default recipe
		When I create a recipe with name "foo"
		And I run tests for recipe "foo"
		Then CLI produced an output "✅: defaults"
		And no errors were printed

	Scenario: Tests fail if templates changes
		When I create a recipe with name "foo"
		And I change recipe "foo" template "README.md" to render "New version"
		And I run tests for recipe "foo"
		Then CLI produced an output "❌: defaults"
		And CLI produced an error "did not match for file 'README.md'"

	Scenario: Update test file snapshots
		When I create a recipe with name "foo"
		And I change recipe "foo" template "README.md" to render "New version"
		And I update tests snapshosts for recipe "foo"
		Then CLI produced an output "test snapshots updated"
		And CLI produced an output "README\.md \(modified\)"
		Then I run tests for recipe "foo"
		And CLI produced an output "✅: defaults"

	Scenario: Expect specific initHelp
		When I create a recipe with name "foo"
		And I expect recipe "foo" initHelp to match "The recipe user will see this message after the recipe execution. Templating is supported, for example: The recipe name is foo."
		And I run tests for recipe "foo"
		Then no errors were printed
		And CLI produced an output "✅: defaults"

	Scenario: Expected initHelp did not match
		When I create a recipe with name "foo"
		And I expect recipe "foo" initHelp to match "Not found"
		And I run tests for recipe "foo"
		And CLI produced an output "❌: defaults"
		Then CLI produced an error "expected init help did not match"
