Feature: Recipes as OCI artifacts
	By pushing recipes as artifacts to OCI compatible repositories, we can improve
	recipe availability and discoverability

	Scenario: Push a recipe to OCI repository
		Given a recipes directory
		And a recipe "foo" that generates file "README.md"
		And a local OCI registry
		When I push the recipe "foo" to the local OCI repository "foo:v0.0.1"
		Then the recipe "foo" should exist in the local OCI repository "foo:v0.0.1"

	Scenario: Pull a recipe from OCI repository
		Given a recipes directory
		And a recipe "foo" that generates file "README.md"
		And a local OCI registry
		And the recipe "foo" is pushed to the local OCI repository "foo:v0.0.1"
		When I pull the recipe "foo" to the local OCI repository "foo:v0.0.1"
		Then the recipes directory should contain recipe "foo"