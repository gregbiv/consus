Feature: Check if key exists

  Scenario: Checking the key
    Given there are keys:
      | KeyID       | Value         | CreatedAt  | ExpiresAt  |
      | testonekey  | value one     | 2018-01-31 | 2050-01-31 |

    When I send a "HEAD" request to "/v1/keys/testonekey"
    And  the response code should be 200
    Then I send a "GET" request to "/v1/keys/random"
    Then the response code should be 404
