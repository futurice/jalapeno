Feature: Delete
	Delete removes sauces from the project

	Scenario: Delete all sauce from project
		Given a recipe "foo"
		And recipe "foo" generates file "README.md" with content "initial"
		When I execute recipe "foo"
		And I delete all sauces from the project
		Then no errors were printed
		And there should not be a sauce directory in the project directory

	Scenario: Delete specific sauce from project
		Given a recipe "foo"
		And recipe "foo" generates file "README.md" with content "initial"
		When I execute recipe "foo"
		And I delete the sauce from the index 0
		Then no errors were printed
		And there should not be a sauce directory in the project directory

	Scenario: Delete specific sauce from project when there are multiple sauces
		Given a recipe "foo"
		And a recipe "bar"
		And a recipe "baz"
		And recipe "foo" generates file "foo.md" with content "initial"
		And recipe "bar" generates file "bar.md" with content "initial"
		And recipe "baz" generates file "baz.md" with content "initial"
		When I execute recipe "foo"
		And I execute recipe "bar"
		And I execute recipe "baz"
		When I delete the sauce from the index 1
		Then no errors were printed
		And the project directory should contain file "foo.md"
		And the project directory should not contain file "bar.md"
		And the project directory should contain file "baz.md"
		And the sauce in index 0 should have property "Recipe::Name" with value "^foo$"
		And the sauce in index 1 should have property "Recipe::Name" with value "^baz$"
