Feature: Delete a key

  Scenario: Deleting the key
    Given there are keys:
      | KeyID       | Value         | CreatedAt  | ExpiresAt  |
      | testonekey  | value one     | 2018-01-31 | 2050-01-31 |

    When I send a "DELETE" request to "/v1/keys/testonekey"
    And  the response code should be 204
    Then I send a "GET" request to "/v1/keys"
    Then the response code should be 200
    And the response should match json:
    """
    {
      "items": []
    }
    """
    When I send a "GET" request to "/v1/keys/testonekey"
    Then the response code should be 404
