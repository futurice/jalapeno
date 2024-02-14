Feature: Validating recipes

	Scenario: Validate a valid recipe
		Given a recipes directory
		And a recipe "foo" that generates file "README.md" with content "initial"
		When I validate recipe "foo"
		Then no errors were printed
