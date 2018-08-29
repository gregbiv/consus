Feature: Checking system status
  In order to know that the service is functioning as expected
  As a system administrator
  I need to be able to receive correct status responses.

  Scenario: Display metrics of the service
    When I request the metrics url of the system
    Then I should receive a response with all metrics successful

  Scenario: Request an invalid route
    When I request a url that does not exist in the system
    Then I should see that the request page cannot be found

  Scenario: API docs are served
    When I request the api documentation
    Then I should see that the api documentation

