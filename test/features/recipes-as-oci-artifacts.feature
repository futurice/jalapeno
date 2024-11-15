Feature: Recipes as OCI artifacts
	By pushing and pulling recipes as artifacts to OCI compatible repositories, we can improve
	recipe availability and discoverability

	@registry
	Scenario: Push a recipe to OCI repository
		Given a recipe "foo"
		And recipe "foo" generates file "README.md" with content "initial"
		And a local OCI registry
		When I push the recipe "foo" to the local OCI repository
		Then no errors were printed
	@registry
	Scenario: Pull a recipe from OCI repository
		Given a recipe "foo"
		And recipe "foo" generates file "README.md" with content "initial"
		And a local OCI registry
		And the recipe "foo" is pushed to the local OCI repository "foo:v0.0.1"
		When I pull recipe from the local OCI repository "foo:v0.0.1"
		Then no errors were printed
		And the project directory should contain file "foo/recipe.yml"

	@registry
	Scenario: Push a recipe to OCI repository using the 'latest' tag
		Given a recipe "foo"
		And recipe "foo" generates file "README.md" with content "initial"
		And a local OCI registry
		When I push the recipe "foo" to the local OCI repository with '--latest' flag
		Then no errors were printed
		When I pull recipe from the local OCI repository "foo:latest"
		Then no errors were printed
		And the project directory should contain file "foo/recipe.yml" with "version: v0.0.1"

	@registry
	Scenario: Pushing a recipe to OCI repository using the 'latest' tag pushes the version tag also
		Given a recipe "foo"
		And recipe "foo" generates file "README.md" with content "initial"
		And a local OCI registry
		When I push the recipe "foo" to the local OCI repository with '--latest' flag
		Then no errors were printed
		When I pull recipe from the local OCI repository "foo:v0.0.1"
		Then no errors were printed
		And the project directory should contain file "foo/recipe.yml" with "version: v0.0.1"

	@registry
	Scenario: Pushing a recipe to OCI repository using the 'latest' tag replaces the previous tag
		Given a recipe "foo"
		And recipe "foo" generates file "README.md" with content "initial"
		And a local OCI registry
		When I push the recipe "foo" to the local OCI repository with '--latest' flag
		Then no errors were printed
		When I change recipe "foo" to version "v0.0.2"
		And I change recipe "foo" template "README.md" to render "New version"
		And I push the recipe "foo" to the local OCI repository with '--latest' flag
		And I pull recipe from the local OCI repository "foo:latest"
		Then no errors were printed
		And the project directory should contain file "foo/recipe.yml" with "version: v0.0.2"
	
	@registry
	Scenario: Push a recipe to OCI repository with authentication
		Given a recipe "foo"
		And recipe "foo" generates file "README.md" with content "initial"
		And a local OCI registry with authentication
		When I push the recipe "foo" to the local OCI repository
		Then no errors were printed

	@registry
	Scenario: Pull a recipe from OCI repository with authentication
		Given a recipe "foo"
		And recipe "foo" generates file "README.md" with content "initial"
		And a local OCI registry with authentication
		And the recipe "foo" is pushed to the local OCI repository "foo:v0.0.1"
		When I pull recipe from the local OCI repository "foo:v0.0.1"
		Then no errors were printed
		And the project directory should contain file "foo/recipe.yml"

	@registry
	Scenario: Try to push a recipe to OCI repository without authentication
		Given a recipe "foo"
		And recipe "foo" generates file "README.md" with content "initial"
		And a local OCI registry with authentication
		And registry credentials are not provided by the command
		When I push the recipe "foo" to the local OCI repository
		Then CLI produced an error "basic credential not found"

	@registry
	Scenario: Try to pull a recipe from OCI repository which not exist
		Given a local OCI registry with authentication
		When I pull recipe from the local OCI repository "foo:v0.0.1"
		Then CLI produced an error "recipe not found"

	@registry
	Scenario: Push a recipe from OCI repository using credentials from config file
		Given a recipe "foo"
		And recipe "foo" generates file "README.md" with content "initial"
		And a local OCI registry with authentication
		And registry credentials are provided by config file
		And registry credentials are not provided by the command
		When I push the recipe "foo" to the local OCI repository
		Then no errors were printed
