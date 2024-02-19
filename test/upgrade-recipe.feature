Feature: Upgrade sauce
	Upgrade a Jalapeno sauce

	Scenario: Upgrade sauce
		Given a project directory
		And a recipes directory
		And a recipe "foo" that generates file "README.md" with content "initial"
		And I execute recipe "foo"
		And I change recipe "foo" to version "v0.0.2"
		And I change recipe "foo" template "README.md" to render "New version"
		When I upgrade recipe "foo"
		Then no errors were printed
		And CLI produced an output "README\.md \(modified\)"
		And the project directory should contain file ".jalapeno/sauces.yml" with "version: v0\.0\.2"
		And the project directory should contain file "README.md" with "New version"
		And no conflicts were reported

	Scenario: Upgrade sauce with same version
		Given a project directory
		And a recipes directory
		And a recipe "foo" that generates file "README.md" with content "initial"
		And I execute recipe "foo"
		When I upgrade recipe "foo"
		Then no errors were printed
		And CLI produced an output "no changes were made"
		And the project directory should contain file ".jalapeno/sauces.yml" with "version: v0\.0\.1"
		And no conflicts were reported

	Scenario: Upgrading sauce removes old files from the project directory if deprecated
		Given a project directory
		And a recipes directory
		And a recipe "foo" that generates file "README.md" with content "initial"
		And a recipe "foo" that generates file "will-be-removed-in-next-version.md" with content "removed"
		And I execute recipe "foo"
		And I change recipe "foo" to version "v0.0.2"
		And I change recipe "foo" template "README.md" to render "New version"
		And I remove file "will-be-removed-in-next-version.md" from the recipe "foo"
		When I upgrade recipe "foo"
		Then no errors were printed
		And CLI produced an output "README\.md \(modified\)"
		And the project directory should not contain file "will-be-removed-in-next-version.md"

	Scenario: Upgrading the recipe does remove old files from the project directory if modified by the user
		Given a project directory
		And a recipes directory
		And a recipe "foo" that generates file "README.md" with content "initial"
		And a recipe "foo" that generates file "will-be-removed-in-next-version.md" with content "removed"
		And I execute recipe "foo"
		And I change recipe "foo" to version "v0.0.2"
		And I change recipe "foo" template "README.md" to render "New version"
		And I remove file "will-be-removed-in-next-version.md" from the recipe "foo"
		And I change project file "will-be-removed-in-next-version.md" to contain "Locally modified"
		When I upgrade recipe "foo"
		Then no errors were printed
		And the project directory should not contain file "will-be-removed-in-next-version.md"

	Scenario: Try to upgrade sauces with same recipes without providing sauce ID
		Given a project directory
		And a recipes directory
		And a recipe "shared" that generates file "README.md" with content "initial"
		And recipes will be executed to the subpath "foo"
		And I execute recipe "shared"
		Then execution of the recipe has succeeded
		And the sauce file contains a sauce in index 0 which should have property "subPath" with value "^foo$"
		And the project directory should contain file "./foo/README.md"
		When recipes will be executed to the subpath "bar"
		And I execute recipe "shared"
		Then execution of the recipe has succeeded
		And the sauce file contains a sauce in index 1 which should have property "subPath" with value "^bar$"
		And the project directory should contain file "./bar/README.md"
		When I change recipe "shared" to version "v0.0.2"
		And I change recipe "shared" template "README.md" to render "New version"
		And I upgrade recipe "shared"
		Then CLI produced an error "contains multiple sauces with recipe 'shared'. Use --sauce-id"
	
	Scenario: Upgrade sauces with same recipe but with different subpaths
		Given a project directory
		And a recipes directory
		And a recipe "shared" that generates file "README.md" with content "initial"
		And recipes will be executed to the subpath "foo"
		And I execute recipe "shared"
		Then execution of the recipe has succeeded
		And the sauce file contains a sauce in index 0 which should have property "subPath" with value "^foo$"
		And the project directory should contain file "./foo/README.md"
		When recipes will be executed to the subpath "bar"
		And I execute recipe "shared"
		Then execution of the recipe has succeeded
		And the sauce file contains a sauce in index 1 which should have property "subPath" with value "^bar$"
		And the project directory should contain file "./bar/README.md"
		When I change recipe "shared" to version "v0.0.2"
		And I change recipe "shared" template "README.md" to render "New version"
		And I select sauce in index 0 for the upgrade
		And I upgrade recipe "shared"
		Then no errors were printed
		Then CLI produced an output "README\.md \(modified\)"
		And the project directory should contain file "./foo/README.md" with "New version"
		And the project directory should contain file "./bar/README.md" with "initial"

	Scenario: Upgrade sauce from remote recipe
		Given a project directory
		And a recipes directory
		And a local OCI registry
		And a recipe "foo" that generates file "README.md" with content "initial"
		And I execute recipe "foo"
		And I change recipe "foo" to version "v0.0.2"
		And I change recipe "foo" template "README.md" to render "New version"
		And the recipe "foo" is pushed to the local OCI repository "foo:v0.0.2"
		When I upgrade recipe from the local OCI repository "foo:v0.0.2"
		Then no errors were printed
		Then CLI produced an output "README\.md \(modified\)"
		And the project directory should contain file ".jalapeno/sauces.yml" with "version: v0\.0\.2"
		And the project directory should contain file ".jalapeno/sauces.yml" with "from: oci://localhost:\d+/foo"
		And the project directory should contain file "README.md" with "New version"
		And no conflicts were reported

	Scenario: Attempt upgrade when previous sauce file was modified
		Given a project directory
		And a recipes directory
		And a recipe "foo" that generates file "README.md" with content "initial"
		And I execute recipe "foo"
		And I change recipe "foo" to version "v0.0.2"
		And I change recipe "foo" template "README.md" to render "New version"
		And I change project file "README.md" to contain "Locally modified"
		When I upgrade recipe "foo"
		Then CLI produced an error "file conflicts"
		And the project directory should contain file "README.md" with "Locally modified"

	Scenario: Attempt upgrade when user keeps the locally modified file
		Given a project directory
		And a recipes directory
		And a recipe "foo" that generates file "README.md" with content "initial"
		And I execute recipe "foo"
		And I change recipe "foo" to version "v0.0.2"
		And I change recipe "foo" template "README.md" to render "New version"
		And I change project file "README.md" to contain "Locally modified"
		And I buffer key presses "n\r"
		When I upgrade recipe "foo"
		Then CLI produced an output "README\.md: keep"
		Then CLI produced an output "no changes were made to any files"
		And the project directory should contain file "README.md" with "Locally modified"

	Scenario: Attempt upgrade when user overrides the locally modified file
		Given a project directory
		And a recipes directory
		And a recipe "foo" that generates file "README.md" with content "initial"
		And I execute recipe "foo"
		And I change recipe "foo" to version "v0.0.2"
		And I change recipe "foo" template "README.md" to render "New version"
		And I change project file "README.md" to contain "Locally modified"
		And I buffer key presses "y\r"
		When I upgrade recipe "foo"
		Then CLI produced an output "README\.md: override"
		Then CLI produced an output "README\.md \(modified\)"
		And the project directory should contain file "README.md" with "New version"

	Scenario: Attempt upgrade when user overrides the locally modified file while using arrow keys
		Given a project directory
		And a recipes directory
		And a recipe "foo" that generates file "README.md" with content "initial"
		And I execute recipe "foo"
		And I change recipe "foo" to version "v0.0.2"
		And I change recipe "foo" template "README.md" to render "New version"
		And I change project file "README.md" to contain "Locally modified"
		And I buffer key presses "→←→\r"
		When I upgrade recipe "foo"
		Then CLI produced an output "README.md: override"
		Then CLI produced an output "README\.md \(modified\)"
		And the project directory should contain file "README.md" with "New version"

	Scenario: Attempt upgrade when new file conflicts with existing manually created file
		Given a project directory
		And a recipes directory
		And a recipe "foo" that generates file "README.md" with content "initial"
		And I execute recipe "foo"
		And I create a file "new.txt" with contents "manual" to the project directory
		And I change recipe "foo" to version "v0.0.2"
		And I change recipe "foo" template "new.txt" to render "new"
		When I upgrade recipe "foo"
		Then CLI produced an error "file conflicts"
