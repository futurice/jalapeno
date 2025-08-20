Feature: List
	List sauces in the project

	Scenario: List all sauces
		Given a recipe "foo"
		And recipe "foo" generates file "README.md" with content "initial"
		When I execute recipe "foo"
		And I list all sauces in the project
		Then no errors were printed
		And CLI produced an output "name: foo"

	Scenario: List specific sauce when there are multiple sauces
		Given a recipe "foo"
		And a recipe "bar"
		And recipe "foo" generates file "foo.md" with content "initial"
		And recipe "bar" generates file "bar.md" with content "initial"
		When I execute recipe "foo"
		And I execute recipe "bar"
		And I list all sauces in the project
		Then no errors were printed
		And CLI produced an output "name: foo"
		And CLI produced an output "name: bar"
