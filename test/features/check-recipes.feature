Feature: Check for new recipe versions

	Scenario: Find newer version for a recipe
		Given a recipe "foo"
		And recipe "foo" generates file "README.md" with content "initial"
		And a local OCI registry
		When I execute recipe "foo"
		And the source of the sauce with recipe "foo" is in the local OCI registry
		Then execution of the recipe has succeeded
		When I push the recipe "foo" to the local OCI repository
		Then no errors were printed
		When I change recipe "foo" to version "v0.0.2"
		And I push the recipe "foo" to the local OCI repository
		And I check new versions for recipe "foo"
		Then CLI produced an output "new versions found: v0\.0\.2"
		Then CLI produced an output "To upgrade recipes to the latest version run:\n  (.*) upgrade oci://localhost:\d+/foo:v0.0.2\n"

	Scenario: Find multiple newer version for a recipe
		Given a recipe "foo"
		And recipe "foo" generates file "README.md" with content "initial"
		And a local OCI registry
		When I execute recipe "foo"
		And the source of the sauce with recipe "foo" is in the local OCI registry
		Then execution of the recipe has succeeded
		When I push the recipe "foo" to the local OCI repository
		Then no errors were printed
		When I change recipe "foo" to version "v0.0.2"
		And I push the recipe "foo" to the local OCI repository
		And I change recipe "foo" to version "v0.0.3"
		And I push the recipe "foo" to the local OCI repository
		Then I check new versions for recipe "foo"
		Then CLI produced an output "new versions found: v0\.0\.2, v0\.0\.3"
		Then CLI produced an output "To upgrade recipes to the latest version run:\n  (.*) upgrade oci://localhost:\d+/foo:v0\.0\.3\n"

	Scenario: Find newer version for multiple recipes
		Given a recipe "foo"
		And recipe "foo" generates file "foo.md" with content "initial"
		And a recipe "bar"
		And recipe "bar" generates file "bar.md" with content "initial"
		And a local OCI registry
		When I execute recipe "foo"
		And the source of the sauce with recipe "foo" is in the local OCI registry
		Then execution of the recipe has succeeded
		When I execute recipe "bar"
		And the source of the sauce with recipe "bar" is in the local OCI registry
		Then execution of the recipe has succeeded
		When I push the recipe "foo" to the local OCI repository
		When I push the recipe "bar" to the local OCI repository
		Then no errors were printed
		When I change recipe "foo" to version "v0.0.2"
		And I push the recipe "foo" to the local OCI repository
		Then no errors were printed
		When I change recipe "bar" to version "v0.0.2"
		And I push the recipe "bar" to the local OCI repository
		Then no errors were printed
		Then I check new versions for recipes
		Then CLI produced an output "foo: new versions found: v0\.0\.2"
		And CLI produced an output "bar: new versions found: v0\.0\.2"

	Scenario: Unable to find newer recipe versions
		Given a recipe "foo"
		And recipe "foo" generates file "README.md" with content "initial"
		And a local OCI registry
		When I push the recipe "foo" to the local OCI repository
		Then no errors were printed
		When I change recipe "foo" to version "v0.0.2"
		And I push the recipe "foo" to the local OCI repository
		And I execute recipe "foo"
		And the source of the sauce with recipe "foo" is in the local OCI registry
		Then execution of the recipe has succeeded
		And I check new versions for recipe "foo"
		Then CLI produced an output "no new versions found"
	
	Scenario: Unable to find newer recipe versions for all recipes
		Given a recipe "foo"
		And recipe "foo" generates file "foo.md" with content "initial"
		And a recipe "bar"
		And recipe "bar" generates file "bar.md" with content "initial"
		And a local OCI registry
		When I execute recipe "foo"
		And the source of the sauce with recipe "foo" is in the local OCI registry
		Then execution of the recipe has succeeded
		When I execute recipe "bar"
		And the source of the sauce with recipe "bar" is in the local OCI registry
		Then execution of the recipe has succeeded
		When I push the recipe "foo" to the local OCI repository
		When I push the recipe "bar" to the local OCI repository
		Then no errors were printed
		When I change recipe "foo" to version "v0.0.2"
		And I push the recipe "foo" to the local OCI repository
		Then no errors were printed
		Then I check new versions for recipes
		Then CLI produced an output "foo: new versions found: v0\.0\.2"
		And CLI produced an output "bar: no new versions found"

	Scenario: Executing remote recipe automatically adds the repo as source for the sauce
		Given a recipe "foo"
		And recipe "foo" generates file "README.md" with content "initial"
		And a local OCI registry
		And the recipe "foo" is pushed to the local OCI repository "foo:v0.0.1"
		And I change recipe "foo" to version "v0.0.2"
		And the recipe "foo" is pushed to the local OCI repository "foo:v0.0.2"
		When I execute the recipe from the local OCI repository "foo:v0.0.1"
		Then execution of the recipe has succeeded
		# Note the lack of explicitly setting the source for the sauce
		And I check new versions for recipe "foo"
		Then CLI produced an output "new versions found: v0\.0\.2"

	Scenario: Manually override the check from URL for locally executed recipe
		Given a recipe "foo"
		And recipe "foo" generates file "README.md" with content "initial"
		And a local OCI registry
		When I execute recipe "foo"
		Then execution of the recipe has succeeded
		When the recipe "foo" is pushed to the local OCI repository "foo:v0.0.1"
		And I change recipe "foo" to version "v0.0.2"
		And the recipe "foo" is pushed to the local OCI repository "foo:v0.0.2"
		And I check new versions for recipe "foo" from the local OCI repository "foo"
		Then CLI produced an output "new versions found: v0\.0\.2"
		And the sauce in index 0 which should have property "CheckFrom" with value "^oci://localhost:\d+/foo$"

	Scenario: Find and upgrade newer version for recipes
		Given a recipe "foo"
		And recipe "foo" generates file "README.md" with content "initial"
		And a local OCI registry
		When I execute recipe "foo"
		And the source of the sauce with recipe "foo" is in the local OCI registry
		Then execution of the recipe has succeeded
		When I push the recipe "foo" to the local OCI repository
		Then no errors were printed
		When I change recipe "foo" to version "v0.0.2"
		And I push the recipe "foo" to the local OCI repository
		And I check and upgrade new versions for recipes
		Then CLI produced an output "new versions found: v0\.0\.2"
		Then CLI produced an output "Upgrade completed"

	Scenario: Find and upgrade newer version for a specific recipe
		Given a recipe "foo"
		And recipe "foo" generates file "README.md" with content "initial"
		And a recipe "bar"
		And recipe "bar" generates file "bar.md" with content "initial"
		And a local OCI registry
		When I execute recipe "foo"
		And the source of the sauce with recipe "foo" is in the local OCI registry
		Then execution of the recipe has succeeded
		When I execute recipe "bar"
		And the source of the sauce with recipe "bar" is in the local OCI registry
		Then execution of the recipe has succeeded
		When I push the recipe "foo" to the local OCI repository
		And I push the recipe "bar" to the local OCI repository
		Then no errors were printed
		When I change recipe "foo" to version "v0.0.2"
		And I push the recipe "foo" to the local OCI repository
		And I check and upgrade new version for recipe "foo"
		Then CLI produced an output "new versions found: v0\.0\.2"
		Then CLI produced an output "Upgrade completed"

	Scenario: Find and upgrade newer versions for multiple recipes
		Given a recipe "foo"
		And recipe "foo" generates file "foo.md" with content "initial"
		And a recipe "bar"
		And recipe "bar" generates file "bar.md" with content "initial"
		And a recipe "baz"
		And recipe "baz" generates file "baz.md" with content "initial"
		And a local OCI registry
		When I execute recipe "foo"
		And the source of the sauce with recipe "foo" is in the local OCI registry
		Then execution of the recipe has succeeded
		When I execute recipe "bar"
		And the source of the sauce with recipe "bar" is in the local OCI registry
		Then execution of the recipe has succeeded
		When I execute recipe "baz"
		And the source of the sauce with recipe "baz" is in the local OCI registry
		Then execution of the recipe has succeeded
		When I push the recipe "foo" to the local OCI repository
		When I push the recipe "bar" to the local OCI repository
		When I push the recipe "baz" to the local OCI repository
		Then no errors were printed
		When I change recipe "foo" to version "v0.0.2"
		And I push the recipe "foo" to the local OCI repository
		Then no errors were printed
		When I change recipe "bar" to version "v0.0.2"
		And I push the recipe "bar" to the local OCI repository
		Then no errors were printed
		Then I check and upgrade new versions for recipes
		Then CLI produced an output "foo: new versions found: v0\.0\.2"
		And CLI produced an output "bar: new versions found: v0\.0\.2"
		And CLI produced an output "baz: no new versions found"
		And CLI produced an output "All recipes with newer versions upgraded successfully!"