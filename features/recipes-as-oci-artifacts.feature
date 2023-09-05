Feature: Recipes as OCI artifacts
	By pushing and pulling recipes as artifacts to OCI compatible repositories, we can improve
	recipe availability and discoverability

	Scenario: Push a recipe to OCI repository
		Given a recipes directory
		And a recipe "foo" that generates file "README.md"
		And a local OCI registry
		When I push the recipe "foo" to the local OCI repository
		Then push of the recipe was successful

	Scenario: Pull a recipe from OCI repository
		Given a recipes directory
		And a recipe "foo" that generates file "README.md"
		And a local OCI registry
		And the recipe "foo" is pushed to the local OCI repository "foo:v0.0.1"
		When I pull the recipe "foo" from the local OCI repository "foo:v0.0.1"
		Then pull of the recipe was successful
		And the recipes directory should contain recipe "foo"
	
	Scenario: Push a recipe to OCI repository with authentication
		Given a recipes directory
		And a recipe "foo" that generates file "README.md"
	 	And a local OCI registry with authentication
		When I push the recipe "foo" to the local OCI repository
		Then push of the recipe was successful

	Scenario: Pull a recipe from OCI repository with authentication
		Given a recipes directory
		And a recipe "foo" that generates file "README.md"
		And a local OCI registry with authentication
		And the recipe "foo" is pushed to the local OCI repository "foo:v0.0.1"
		When I pull the recipe "foo" from the local OCI repository "foo:v0.0.1"
		Then pull of the recipe was successful
		And the recipes directory should contain recipe "foo"
	
	Scenario: Try to push a recipe to OCI repository without authentication
		Given a recipes directory
		And a recipe "foo" that generates file "README.md"
	 	And a local OCI registry with authentication
	 	And registry credentials are not provided by the command
		When I push the recipe "foo" to the local OCI repository
		Then CLI produced an error "Error: unauthorized"

	Scenario: Try to pull a recipe from OCI repository which not exist
		Given a recipes directory
	 	And a local OCI registry with authentication
		When I pull the recipe "foo" from the local OCI repository "foo:v0.0.1"
		Then CLI produced an error "Error: recipe not found"

	Scenario: Pull a recipe from OCI repository using credentials from config file
		Given a recipes directory
		And a recipe "foo" that generates file "README.md"
	 	And a local OCI registry with authentication
		And registry credentials are provided by config file
		And registry credentials are not provided by the command
		When I push the recipe "foo" to the local OCI repository
		Then push of the recipe was successful

	Scenario: Pull a recipe from OCI repository using credentials from default config file
		Given a recipes directory
		And a recipe "foo" that generates file "README.md"
	 	And a local OCI registry with authentication
		And registry credentials are provided by default config file
		And registry credentials are not provided by the command
		When I push the recipe "foo" to the local OCI repository
		Then push of the recipe was successful
