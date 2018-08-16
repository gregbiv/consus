Feature: Create a key

  Scenario: Creating a key
    When I send a "PUT" request to "/v1/keys?expire_in=60" containing the following JSON:
    """
    {
      "id": "testput",
      "value": "random value"
    }
    """
    Then the response code should be 200
    When I send a "HEAD" request to "/v1/keys/testput"
    And  the response code should be 200

  Scenario: Creating a key with invalid payload
    When I send a "PUT" request to "/v1/keys" containing the following JSON:
    """
    {
      "id": "test invalid id",
      "value": "random value"
    }
    """
    Then the response code should be 400
    And the response should match json:
    """
    {
      "error": {
        "code": "InvalidInput",
        "target": "id",
        "message": "invalid id provided, value should contain only letters"
      }
     }
    """
