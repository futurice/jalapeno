Feature: Upgrade sauce
	Upgrade a Jalapeno sauce

	Scenario: Upgrade sauce
		Given a recipe "foo"
		And recipe "foo" generates file "README.md" with content "initial"
		And I execute recipe "foo"
		And I change recipe "foo" to version "v0.0.2"
		And I change recipe "foo" template "README.md" to render "New version"
		And recipe "foo" generates file "new_file.md" with content "initial"
		When I upgrade recipe "foo"
		Then no errors were printed
		And CLI produced an output "README\.md \(modified\)"
		And CLI produced an output "new_file\.md \(added\)"
		And the project directory should contain file ".jalapeno/sauces.yml" with "version: v0\.0\.2"
		And the project directory should contain file "README.md" with "New version"
		And no conflicts were reported

	Scenario: Upgrade sauce with same version
		Given a recipe "foo"
		And recipe "foo" generates file "README.md" with content "initial"
		And I execute recipe "foo"
		When I upgrade recipe "foo"
		Then no errors were printed
		And CLI produced an output "no changes were made"
		And the project directory should contain file ".jalapeno/sauces.yml" with "version: v0\.0\.1"
		And no conflicts were reported

	Scenario: Upgrading sauce removes old files from the project directory if deprecated
		Given a recipe "foo"
		And recipe "foo" generates file "README.md" with content "initial"
		And recipe "foo" generates file "will-be-removed-in-next-version.md" with content "removed"
		And I execute recipe "foo"
		And I change recipe "foo" to version "v0.0.2"
		And I change recipe "foo" template "README.md" to render "New version"
		And I remove file "will-be-removed-in-next-version.md" from the recipe "foo"
		When I upgrade recipe "foo"
		Then no errors were printed
		And CLI produced an output "README\.md \(modified\)"
		And CLI produced an output "will-be-removed-in-next-version\.md \(deleted\)"
		And the project directory should not contain file "will-be-removed-in-next-version.md"

	Scenario: Upgrading the recipe does remove old files from the project directory if modified by the user
		Given a recipe "foo"
		And recipe "foo" generates file "README.md" with content "initial"
		And recipe "foo" generates file "will-be-removed-in-next-version.md" with content "removed"
		And I execute recipe "foo"
		And I change recipe "foo" to version "v0.0.2"
		And I change recipe "foo" template "README.md" to render "New version"
		And I remove file "will-be-removed-in-next-version.md" from the recipe "foo"
		And I change project file "will-be-removed-in-next-version.md" to contain "Locally modified"
		When I upgrade recipe "foo"
		Then no errors were printed
		And CLI produced an output "README\.md \(modified\)"
		And CLI produced an output "will-be-removed-in-next-version\.md \(deleted\)"
		And the project directory should not contain file "will-be-removed-in-next-version.md"

	Scenario: Try to upgrade sauces with same recipes without providing sauce ID
		Given a recipe "shared"
		And recipe "shared" generates file "README.md" with content "initial"
		And recipes will be executed to the subpath "foo"
		And I execute recipe "shared"
		Then execution of the recipe has succeeded
		And the sauce in index 0 which should have property "SubPath" with value "^foo$"
		And the project directory should contain file "./foo/README.md"
		When recipes will be executed to the subpath "bar"
		And I execute recipe "shared"
		Then execution of the recipe has succeeded
		And the sauce in index 1 which should have property "SubPath" with value "^bar$"
		And the project directory should contain file "./bar/README.md"
		When I change recipe "shared" to version "v0.0.2"
		And I change recipe "shared" template "README.md" to render "New version"
		And I upgrade recipe "shared"
		Then CLI produced an error "contains multiple sauces with recipe 'shared'. Use --sauce-id"
	
	Scenario: Upgrade sauces with same recipe but with different subpaths
		Given a recipe "shared"
		And recipe "shared" generates file "README.md" with content "initial"
		And recipes will be executed to the subpath "foo"
		And I execute recipe "shared"
		Then execution of the recipe has succeeded
		And the sauce in index 0 which should have property "SubPath" with value "^foo$"
		And the project directory should contain file "./foo/README.md"
		When recipes will be executed to the subpath "bar"
		And I execute recipe "shared"
		Then execution of the recipe has succeeded
		And the sauce in index 1 which should have property "SubPath" with value "^bar$"
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
		Given a local OCI registry
		And a recipe "foo"
		And recipe "foo" generates file "README.md" with content "initial"
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

	Scenario: Attempt upgrade when a file was modified by the user
		Given a recipe "foo"
		And recipe "foo" generates file "README.md" with content "initial"
		And I execute recipe "foo"
		And I change recipe "foo" to version "v0.0.2"
		And I change recipe "foo" template "README.md" to render "New version"
		And I change project file "README.md" to contain "Locally modified"
		When I upgrade recipe "foo"
		Then CLI produced an error "file conflicts"
		And the project directory should contain file "README.md" with "Locally modified"

	Scenario: Attempt upgrade when user keeps the locally modified file
		Given a recipe "foo"
		And recipe "foo" generates file "README.md" with content "initial"
		And I execute recipe "foo"
		And I change recipe "foo" to version "v0.0.2"
		And I change recipe "foo" template "README.md" to render "New version"
		And I change project file "README.md" to contain "Locally modified"
		And I buffer key presses "\r"
		When I upgrade recipe "foo"
		Then CLI produced an output "README\.md: keep old"
		Then CLI produced an output "no changes were made to any files"
		And the project directory should contain file "README.md" with "Locally modified"

	Scenario: Attempt upgrade when user overrides the locally modified file
		Given a recipe "foo"
		And recipe "foo" generates file "README.md" with content "initial"
		And I execute recipe "foo"
		And I change recipe "foo" to version "v0.0.2"
		And I change recipe "foo" template "README.md" to render "New version"
		And I change project file "README.md" to contain "Locally modified"
		And I buffer key presses "→\r"
		When I upgrade recipe "foo"
		Then CLI produced an output "README\.md: use new"
		Then CLI produced an output "README\.md \(modified\)"
		And the project directory should contain file "README.md" with "New version"

	Scenario: Attempt upgrade when user overrides the locally modified file while using arrow keys
		Given a recipe "foo"
		And recipe "foo" generates file "README.md" with content "initial"
		And I execute recipe "foo"
		And I change recipe "foo" to version "v0.0.2"
		And I change recipe "foo" template "README.md" to render "New version"
		And I change project file "README.md" to contain "Locally modified"
		And I buffer key presses "→←→\r"
		When I upgrade recipe "foo"
		Then CLI produced an output "README.md: use new"
		Then CLI produced an output "README\.md \(modified\)"
		And the project directory should contain file "README.md" with "New version"

	Scenario: Attempt upgrade when new file conflicts with existing manually created file
		Given a recipe "foo"
		And recipe "foo" generates file "README.md" with content "initial"
		And I execute recipe "foo"
		And I create a file "new.txt" with contents "manual" to the project directory
		And I change recipe "foo" to version "v0.0.2"
		And I change recipe "foo" template "new.txt" to render "new"
		When I upgrade recipe "foo"
		Then CLI produced an error "file conflicts"

	Scenario: Attempt force upgrade
		Given a recipe "foo"
		And recipe "foo" generates file "README.md" with content "initial"
		And I execute recipe "foo"
		And I change recipe "foo" to version "v0.0.2"
		And I change recipe "foo" template "README.md" to render "New version"
		And I change recipe "foo" template "ANOTHER.md" to render "New version"
		When I upgrade recipe "foo" forcefully
		Then CLI produced an output "README\.md \(modified\)"
		Then CLI produced an output "ANOTHER\.md \(added\)"
		And the project directory should contain file "README.md" with "New version"

	Scenario: Attempt force upgrade when there is a locally modified file
		Given a recipe "foo"
		And recipe "foo" generates file "README.md" with content "initial"
		And I execute recipe "foo"
		And I change recipe "foo" to version "v0.0.2"
		And I change recipe "foo" template "README.md" to render "New version"
		And I change project file "README.md" to contain "Locally modified"
		When I upgrade recipe "foo" forcefully
		Then CLI produced an output "README\.md \(modified\)"
		And the project directory should contain file "README.md" with "New version"
