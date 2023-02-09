Feature: Recipes as OCI artifacts
	By pushing recipes as artifacts to OCI compatible repositories, we can improve
	recipe availability and discoverability

	Scenario: Push a recipe to OCI repository
		Given a recipes directory
		And a recipe "foo" that generates file "README.md"
		And a local OCI registry
		When I push the recipe "foo" to the local OCI repository "foo:v0.0.1"
		Then push of the recipe was successful

	Scenario: Pull a recipe from OCI repository
		Given a recipes directory
		And a recipe "foo" that generates file "README.md"
		And a local OCI registry
		And the recipe "foo" is pushed to the local OCI repository "foo:v0.0.1"
		When I pull the recipe "foo" to the local OCI repository "foo:v0.0.1"
		Then pull of the recipe was successful
		And the recipes directory should contain recipe "foo"
	
	Scenario: Push a recipe to OCI repository with authentication
		Given a recipes directory
		And a recipe "foo" that generates file "README.md"
	 	And a local OCI registry with authentication
		When I push the recipe "foo" to the local OCI repository "foo:v0.0.1"
		Then push of the recipe was successful

	Scenario: Pull a recipe from OCI repository with authentication
		Given a recipes directory
		And a recipe "foo" that generates file "README.md"
		And a local OCI registry with authentication
		And the recipe "foo" is pushed to the local OCI repository "foo:v0.0.1"
		When I pull the recipe "foo" to the local OCI repository "foo:v0.0.1"
		Then pull of the recipe was successful
		And the recipes directory should contain recipe "foo"
	