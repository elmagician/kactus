Feature: Picker
  Picker definitions should provide basic operations
  to store and replace values in steps.

  Background:
    Given I want to generate uuids test1, test2, test3
    And wanting to generate uuids test4,abcedaire
    And I set variable foo to bar
    And setting variable PI to 3.14

  Scenario Outline:
    Given I want to debug picker
    And I set variable <name> to <value>
    And I want to assert picked variables matches:
      | key       | matcher   | value                 |
      | foo       |           | bar                   |
      | PI        |           | 3.14((number))        |
      | abcedaire | =         | {{abcedaire}}((uuid)) |
      | <name>    | <matcher> | <matchTo>             |

    Examples:
      | name      | value                     | matcher | matchTo       |
      | myVar     | this is some string value | match   | ([a-z]*.)*    |
      | something | 124                       |         | 124((number)) |

  Scenario Outline:
    Given I set variable <name> to <value>
    When I set variable Alpha to [a-z]*
    And setting variable AlphaTest to alphaisthebest
    Then I want to assert picked variables matches:
      | key       | matcher | value     |
      | AlphaTest | match   | {{Alpha}} |

    Examples:
      | name      | value                     |
      | myVar     | this is some string value |
      | something | 124                       |
