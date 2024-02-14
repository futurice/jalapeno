Feature: Eject
	Ejecting removes all traces of Jalapeno from a project

	Scenario: Eject project
		Given a project directory
		And a recipes directory
		And a recipe "foo" that generates file "README.md" with content "initial"
		When I execute recipe "foo"
		And I eject Jalapeno from the project
		Then no errors were printed
		And there should not be a sauce directory in the project directory
