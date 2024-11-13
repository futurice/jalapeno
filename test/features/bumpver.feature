Feature: Bump recipe version and write changelog

    Scenario: Directly bump version
		Given a recipe "foo"
		When I bump recipe "foo" version to "v0.0.2" with message "Test"
		Then no errors were printed
		And recipe "foo" has version "v0.0.2"
		And CLI produced an output "Recipe version bumped: v0.0.1 => v0.0.2"
		And recipe "foo" has changelog message "Test"
    
	Scenario: Command inits changelog
		Given a recipe "foo"
		When I bump recipe "foo" version to "v0.0.2" with message "Test"
		Then no errors were printed
		And recipe "foo" contains changelog with 2 entries
		And CLI produced an output "Recipe version bumped: v0.0.1 => v0.0.2"
		And first entry in recipe "foo" changelog has message "Init version"

	Scenario: Invalid semantic version
		Given a recipe "foo"
		When I bump recipe "foo" version to "not-valid-semver" with message "Test"
		Then CLI produced an error "Error: provided version is not valid semver"
