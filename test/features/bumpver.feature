Feature: Bump recipe version and write changelog

    Scenario: Directly bump version
		Given a recipe "foo"
        When I bump recipe "foo" version to "v0.0.2" with message "Test"
		Then recipe "foo" has version "v0.0.2"
		And recipe "foo" has changelog message "Test"
    
	Scenario: Command inits changelog
		Given a recipe "foo"
        When I bump recipe "foo" version to "v0.0.2" with message "Test"
		Then recipe "foo" contains changelog with 2 entries
		And first entry in recipe "foo" changelog has message "Init version"

	Scenario: Invalid semantic version
		Given a recipe "foo"
		When I bump recipe "foo" version to "bar" with message "Test"
		Then CLI produced an error "Error: Invalid Semantic Version"
