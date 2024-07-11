Feature: Execute recipes
	Executing Jalapeno recipes to render templates

	Scenario: Execute single recipe
		Given a recipe "foo"
		And recipe "foo" generates file "README.md" with content "initial"
		When I execute recipe "foo"
		Then execution of the recipe has succeeded
		And the project directory should contain file "README.md"
		And the sauce in index 0 which should have property "Recipe.Name" with value "^foo$"
		And the sauce in index 0 which has a valid ID

	Scenario: Execute single recipe from remote registry
		Given a recipe "foo"
		And recipe "foo" generates file "README.md" with content "initial"
		And a local OCI registry
		And the recipe "foo" is pushed to the local OCI repository "foo:v0.0.1"
		When I execute the recipe from the local OCI repository "foo:v0.0.1"
		Then execution of the recipe has succeeded
		And no errors were printed
		And the project directory should contain file "README.md"
		And the sauce in index 0 which should have property "CheckFrom" with value "^oci://.+/foo$"

	Scenario: Execute multiple recipes
		Given a recipe "foo"
		And recipe "foo" generates file "README.md" with content "initial"
		And a recipe "bar"
		And recipe "bar" generates file "Taskfile.yml" with content "initial"
		When I execute recipe "foo"
		Then execution of the recipe has succeeded
		When I execute recipe "bar"
		Then execution of the recipe has succeeded
		And no errors were printed
		And the project directory should contain file "README.md"
		And the project directory should contain file "Taskfile.yml"
		And the sauce in index 0 which should have property "Recipe.Name" with value "^foo$"
		And the sauce in index 1 which should have property "Recipe.Name" with value "^bar$"

	Scenario: New recipe conflicts with the previous recipe
		Given a recipe "foo"
		And recipe "foo" generates file "README.md" with content "initial"
		And a recipe "bar"
		And recipe "bar" generates file "Taskfile.yml" with content "initial"
		And a recipe "quux"
		And recipe "quux" generates file "Taskfile.yml" with content "initial"
		When I execute recipe "foo"
		And no errors were printed
		Then execution of the recipe has succeeded
		When I execute recipe "bar"
		And no errors were printed
		Then execution of the recipe has succeeded
		When I execute recipe "quux"
		Then CLI produced an error "file 'Taskfile.yml' was already created by other recipe 'bar'"

	Scenario: Execute single recipe to a subpath
		Given a recipe "foo"
		And recipe "foo" generates file "README" with content "initial"
		And recipes will be executed to the subpath "docs"
		When I execute recipe "foo"
		And no errors were printed
		Then execution of the recipe has succeeded
		And CLI produced an output "docs[\S\s]+└── README"
		And the project directory should contain file "docs/README"
		And the sauce in index 0 which should have property "Files.README"
		And the sauce in index 0 which should have property "SubPath" with value "^docs$"

	Scenario: Execute multiple recipes to different subpaths
		Given a recipe "foo"
		And recipe "foo" generates file "README" with content "initial"
		And recipes will be executed to the subpath "foo"
		When I execute recipe "foo"
		Then no errors were printed
		And execution of the recipe has succeeded
		When recipes will be executed to the subpath "bar"
		And I execute recipe "foo"
		Then no errors were printed
		And execution of the recipe has succeeded
		And the project directory should contain file "foo/README"
		And the project directory should contain file "bar/README"
		And the sauce in index 0 which should have property "Files.README"
		And the sauce in index 0 which should have property "SubPath" with value "^foo$"
		And the sauce in index 1 which should have property "Files.README"
		And the sauce in index 1 which should have property "SubPath" with value "^bar$"

	Scenario: Try to execute recipe which escapes the project root
		Given a recipe "foo"
		And recipe "foo" generates file "README.md" with content "initial"
		And recipes will be executed to the subpath "foo/../.."
		When I execute recipe "foo"
		Then CLI produced an error "must point to a directory inside the project root"

	Scenario: Try to execute recipe which uses absolute sub path
		Given a recipe "foo"
		And recipe "foo" generates file "README.md" with content "initial"
		And recipes will be executed to the subpath "/root/foo"
		When I execute recipe "foo"
		Then CLI produced an error "must be a relative path"
