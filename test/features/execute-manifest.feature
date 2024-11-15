Feature: Execute manifests
	Execute a collection of recipes from a manifest file

	Scenario: Execute a manifest
		Given a recipe "foo"
		And recipe "foo" generates file "foo.md" with content "initial"
		And a recipe "bar"
		And recipe "bar" generates file "bar.md" with content "initial"
		And a manifest file that includes recipes
		| foo |
		| bar |
		When I execute the manifest file
		Then no errors were printed
		And CLI produced an output "^Executing manifest with 2 recipes"
		And CLI produced an output "Recipe name: foo"
		And CLI produced an output "Recipe name: bar"
		And the project directory should contain file "foo.md"
		And the project directory should contain file "bar.md"

	@docker
	Scenario: Execute a manifest with remote recipes
		Given a local OCI registry
		And a recipe "foo"
		And recipe "foo" generates file "foo.md" with content "initial"
		And the recipe "foo" is pushed to the local OCI repository "foo:v0.0.1"
		And a recipe "bar"
		And recipe "bar" generates file "bar.md" with content "initial"
		And the recipe "bar" is pushed to the local OCI repository "bar:v0.0.1"
		And a manifest file that includes remote recipes
		| foo | v0.0.1 |
		| bar | v0.0.1 |
		When I execute the manifest file with remote recipes
		Then no errors were printed
		And CLI produced an output "^Executing manifest with 2 recipes"
		And CLI produced an output "Recipe name: foo"
		And CLI produced an output "Recipe name: bar"
		And the project directory should contain file "foo.md"
		And the project directory should contain file "bar.md"
	
	Scenario: Execute a manifest with no recipes
		Given a manifest file
		When I execute the manifest file
		Then CLI produced an error "^Error: can not load the manifest: manifest must contain at least one recipe"

	Scenario: Conflicting recipes in manifest results in an error
		Given a recipe "foo"
		And recipe "foo" generates file "foo.md" with content "initial"
		And a recipe "conflicts-with-foo"
		And recipe "conflicts-with-foo" generates file "foo.md" with content "conflict"
		And a manifest file that includes recipes
		| foo |
		| conflicts-with-foo |
		When I execute the manifest file
		Then CLI produced an error "^Error: conflict in recipe 'conflicts-with-foo': file 'foo\.md' was already created by other recipe 'foo'"
