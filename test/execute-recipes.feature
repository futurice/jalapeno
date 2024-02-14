Feature: Execute recipes
	Executing Jalapeno recipes to template out projects

	Scenario: Execute single recipe
		Given a project directory
		And a recipes directory
		And a recipe "foo" that generates file "README.md" with content "initial"
		When I execute recipe "foo"
		Then execution of the recipe has succeeded
		And the project directory should contain file "README.md"
		And the sauce file contains a sauce in index 0 which should have property "recipe.name" with value "^foo$"
		And the sauce file contains a sauce in index 0 which has a valid ID

	Scenario: Execute single recipe from remote registry
		Given a project directory
		And a recipes directory
		And a recipe "foo" that generates file "README.md" with content "initial"
		And a local OCI registry
		And the recipe "foo" is pushed to the local OCI repository "foo:v0.0.1"
		When I execute the recipe from the local OCI repository "foo:v0.0.1"
		Then execution of the recipe has succeeded
		And no errors were printed
		And the project directory should contain file "README.md"

	Scenario: Execute multiple recipes
		Given a project directory
		And a recipes directory
		And a recipe "foo" that generates file "README.md" with content "initial"
		And a recipe "bar" that generates file "Taskfile.yml" with content "initial"
		When I execute recipe "foo"
		Then execution of the recipe has succeeded
		When I execute recipe "bar"
		Then execution of the recipe has succeeded
		And no errors were printed
		And the project directory should contain file "README.md"
		And the project directory should contain file "Taskfile.yml"
		And the sauce file contains a sauce in index 0 which should have property "recipe.name" with value "^foo$"
		And the sauce file contains a sauce in index 1 which should have property "recipe.name" with value "^bar$"

	Scenario: New recipe conflicts with the previous recipe
		Given a project directory
		And a recipes directory
		And a recipe "foo" that generates file "README.md" with content "initial"
		And a recipe "bar" that generates file "Taskfile.yml" with content "initial"
		And a recipe "quux" that generates file "Taskfile.yml" with content "initial"
		When I execute recipe "foo"
		And no errors were printed
		Then execution of the recipe has succeeded
		When I execute recipe "bar"
		And no errors were printed
		Then execution of the recipe has succeeded
		When I execute recipe "quux"
		Then CLI produced an error "file 'Taskfile.yml' was already created by recipe 'bar'"

	Scenario: Same recipe is executed twice
		Given a project directory
		And a recipes directory
		And a recipe "foo" that generates file "README.md" with content "initial"
		When I execute recipe "foo"
		And no errors were printed
		Then execution of the recipe has succeeded
		When I execute recipe "foo"
		Then CLI produced an error "recipe 'foo' with version 'v0.0.1' has been already executed"

	Scenario: Execute single recipe to a subpath
		Given a project directory
		And a recipes directory
		And a recipe "foo" that generates file "README" with content "initial"
		And recipes will be executed to the subpath "docs"
		When I execute recipe "foo"
		And no errors were printed
		Then execution of the recipe has succeeded
		And CLI produced an output "docs[\S\s]+└── README"
		And the project directory should contain file "docs/README"
		And the sauce file contains a sauce in index 0 which should have property "files.README"
		And the sauce file contains a sauce in index 0 which should have property "subPath" with value "^docs$"

	Scenario: Execute multiple recipes to different subpaths
		Given a project directory
		And a recipes directory
		And a recipe "foo" that generates file "README" with content "initial"
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
		And the sauce file contains a sauce in index 0 which should have property "files.README"
		And the sauce file contains a sauce in index 0 which should have property "subPath" with value "^foo$"
		And the sauce file contains a sauce in index 1 which should have property "files.README"
		And the sauce file contains a sauce in index 1 which should have property "subPath" with value "^bar$"

	Scenario: Try to execute recipe which escapes the project root
		Given a project directory
		And a recipes directory
		And a recipe "foo" that generates file "README.md" with content "initial"
		And recipes will be executed to the subpath "foo/../.."
		When I execute recipe "foo"
		Then CLI produced an error "must point to a directory inside the project root"

	Scenario: Try to execute recipe which uses absolute sub path
		Given a project directory
		And a recipes directory
		And a recipe "foo" that generates file "README.md" with content "initial"
		And recipes will be executed to the subpath "/root/foo"
		When I execute recipe "foo"
		Then CLI produced an error "must be a relative path"
