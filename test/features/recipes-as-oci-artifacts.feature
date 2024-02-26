Feature: Recipes as OCI artifacts
	By pushing and pulling recipes as artifacts to OCI compatible repositories, we can improve
	recipe availability and discoverability

	Scenario: Push a recipe to OCI repository
		Given a recipes directory
		And a recipe "foo"
		And recipe "foo" generates file "README.md" with content "initial"
		And a local OCI registry
		When I push the recipe "foo" to the local OCI repository
		Then no errors were printed

	Scenario: Pull a recipe from OCI repository
		Given a recipes directory
		And a project directory
		And a recipe "foo"
		And recipe "foo" generates file "README.md" with content "initial"
		And a local OCI registry
		And the recipe "foo" is pushed to the local OCI repository "foo:v0.0.1"
		When I pull the recipe "foo" from the local OCI repository "foo:v0.0.1"
		Then no errors were printed
		And the project directory should contain file "foo/recipe.yml"
	
	Scenario: Push a recipe to OCI repository with authentication
		Given a recipes directory
		And a recipe "foo"
		And recipe "foo" generates file "README.md" with content "initial"
		And a local OCI registry with authentication
		When I push the recipe "foo" to the local OCI repository
		Then no errors were printed

	Scenario: Pull a recipe from OCI repository with authentication
		Given a recipes directory
		And a project directory
		And a recipe "foo"
		And recipe "foo" generates file "README.md" with content "initial"
		And a local OCI registry with authentication
		And the recipe "foo" is pushed to the local OCI repository "foo:v0.0.1"
		When I pull the recipe "foo" from the local OCI repository "foo:v0.0.1"
		Then no errors were printed
		And the project directory should contain file "foo/recipe.yml"
	
	Scenario: Try to push a recipe to OCI repository without authentication
		Given a recipes directory
		And a recipe "foo"
		And recipe "foo" generates file "README.md" with content "initial"
		And a local OCI registry with authentication
		And registry credentials are not provided by the command
		When I push the recipe "foo" to the local OCI repository
		Then CLI produced an error "unauthorized"

	Scenario: Try to pull a recipe from OCI repository which not exist
		Given a recipes directory
		And a project directory
		And a local OCI registry with authentication
		When I pull the recipe "foo" from the local OCI repository "foo:v0.0.1"
		Then CLI produced an error "recipe not found"

	Scenario: Push a recipe from OCI repository using credentials from config file
		Given a recipes directory
		And a recipe "foo"
		And recipe "foo" generates file "README.md" with content "initial"
		And a local OCI registry with authentication
		And registry credentials are provided by config file
		And registry credentials are not provided by the command
		When I push the recipe "foo" to the local OCI repository
		Then no errors were printed
