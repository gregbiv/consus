Feature: Get a list of keys
  As an API
  I need to return an existing list of not expired keys

  Background:
    Given there are keys:
      | KeyID       | Value         | CreatedAt  | ExpiresAt  |
      | testonekey  | value one     | 2018-01-31 | 2050-01-31 |
      | testtwokey  | value two     | 2018-01-31 | 2050-01-31 |
      | testthree   | value three   | 2018-01-31 | 2050-01-31 |
      | testexpired | value expired | 2018-01-31 | 2000-01-31 |

  Scenario: Getting the keys list
    When I send a "GET" request to "/v1/keys"
    Then the response code should be 200
    And the response should match json:
    """
    {
      "items": [
        {
          "id": "testonekey",
          "value": "value one",
          "created_at": "2018-01-31T00:00:00Z",
          "updated_at": null,
          "expires_at": "2050-01-31T00:00:00Z"
        },
        {
          "id": "testtwokey",
          "value": "value two",
          "created_at": "2018-01-31T00:00:00Z",
          "updated_at": null,
          "expires_at": "2050-01-31T00:00:00Z"
        },
        {
          "id": "testthree",
          "value": "value three",
          "created_at": "2018-01-31T00:00:00Z",
          "updated_at": null,
          "expires_at": "2050-01-31T00:00:00Z"
        }
      ]
    }
    """
    When I send a "GET" request to "/v1/keys?filter=test$key"
    Then the response code should be 200
    And the response should match json:
    """
    {
      "items": [
        {
          "id": "testonekey",
          "value": "value one",
          "created_at": "2018-01-31T00:00:00Z",
          "updated_at": null,
          "expires_at": "2050-01-31T00:00:00Z"
        },
        {
          "id": "testtwokey",
          "value": "value two",
          "created_at": "2018-01-31T00:00:00Z",
          "updated_at": null,
          "expires_at": "2050-01-31T00:00:00Z"
        }
      ]
    }
    """
    When I send a "GET" request to "/v1/keys/testexpired"
    Then the response code should be 404

  Scenario: Getting a key which doesn't exist
    When I send a "GET" request to "/v1/keys/testnotexists"
    Then the response code should be 404
