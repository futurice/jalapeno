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

	Scenario: Tests fail if extra files rendered
		When I create a recipe with name "foo"
		And recipe "foo" generates file "new.md" with content "new file"
		And I run tests for recipe "foo"
		Then CLI produced an output "❌: defaults"
		And CLI produced an error "following files were extra: \[new.md\]"

	Scenario: Tests pass if extra files are ignored
		When I create a recipe with name "foo"
		And recipe "foo" generates file "new.md" with content "new file"
		And extra files in the test are ignored for recipe "foo"
		And I run tests for recipe "foo"
		Then CLI produced an output "✅: defaults"
		And no errors were printed

	Scenario: Tests fail if too few files rendered
		When I create a recipe with name "foo"
		And I remove file "README.md" from the recipe "foo"
		And I run tests for recipe "foo"
		Then CLI produced an output "❌: defaults"
		And CLI produced an error "following files were missing: \[README.md\]"

	Scenario: Tests fail if wrong file rendered
		When I create a recipe with name "foo"
		And I remove file "README.md" from the recipe "foo"
		And recipe "foo" generates file "new.md" with content "Hello World!"
		And I run tests for recipe "foo"
		Then CLI produced an output "❌: defaults"
		And CLI produced an error "recipe did not render file which was expected: file 'README\.md'"

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
