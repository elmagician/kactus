Feature: Debug
  I should have a set of definition to enabled/disable debug
  for global part of Kactus

  Background:
    Given I want to generate uuids test1, test2, test3
    And wanting to generate uuids test4,abcedaire
    And I set variable foo to bar
    And setting variable PI to 3.14

  Scenario Outline:
    Given I want to debug matchers
    And I want to debug types converters
    And I set variable <name> to <value>
    When I want to assert picked variables matches:
      | key       | matcher   | value                 |
      | foo       |           | bar                   |
      | PI        |           | 3.14((number))        |
      | abcedaire | =         | {{abcedaire}}((uuid)) |
      | <name>    | <matcher> | <matchTo>             |
    And I want to stop debugging types converters
    And I want to stop debugging matchers

    Examples:
      | name      | value                     | matcher | matchTo       |
      | myVar     | this is some string value | match   | ([a-z]*.)*    |
      | something | 124                       |         | 124((number)) |
