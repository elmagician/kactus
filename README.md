<img src="kactus.png" alt="logo" height="100">

[![Go Reference](https://pkg.go.dev/badge/github.com/elmagician/kactus.svg)](https://pkg.go.dev/github.com/elmagician/kactus)

**Kactus** is a godog library to help write cleaner tests.

<!-- TOC -->

- [Get it](#get-it)

<!-- /TOC -->

## Motivation

Kactus aims to simplify godog test writing for systematic functional tests. It provides helpers to ease API testing and
database integration through the test as well as features around variable managements and value checking.

Use it to as tech addict to create cleaner automated tests on various similar cases.

## Get it

- `go get github.com/elmagician/kactus`

## Simple example

```gherkin
@document
Feature: Document management

  Kactus propose an API to manage documents. As I user, I should be able to create, read and delete documents.

  # CREATION
  @creation
  Scenario Outline: I should be able to create document with name and description
    Given I set request json body:
    """
    {
      "name": "<name>",
      "description": "<description>"
    }
    """
    When I POST http://localhost:8080/documents
    Then response status code should be {{status.Created}}
    And json response should contain:
      | field       | matcher | value         |
      | id          | defined |               |
      | name        | =       | <name>        |
      | description | =       | <description> |

    Examples:
      | name    | description                  |
      | test    | no description               |
      | twice   | func                         |
      | twisted | fate is a funny lol champion |

  @creation
  Scenario Outline: I should be able to create document with name only
    Given I set request json body:
    """
    {
      "name": "<name>"
    }
    """
    When I POST http://localhost:8080/documents
    Then response status code should be {{status.Created}}
    And json response should contain:
      | field | matcher | value  |
      | id    | defined |        |
      | name  | =       | <name> |

    Examples:
      | name    |
      | another |
      | thing   |

  @creation
  Scenario Outline: I should recover errors
    Given I set request json body:
    """
    {
      "name": "<name>",
      "description": "<description>"
    }
    """
    When I POST http://localhost:8080/documents
    Then response status code should be <error_code>
    And json response should contain:
      | field | matcher | value   |
      | error | =       | <error> |

    Examples:
      | name | description    | error_code | error      |
      | test | no description | 400        | duplicated |
      |      | no name        | 400        | invalid    |

  # READ
  @read
  Scenario Outline: I should be able to read document from known ids
    When I GET http://localhost:8080/documents/<id>
    Then response status code should be {{status.Ok}}
    And json response should contain:
      | field       | matcher | value          |
      | id          | =       | <id>((number)) |
      | name        | =       | <name>         |
      | description | =       | <description>  |

    Examples:
      | name    | description                  | id |
      | test    | no description               | 0  |
      | twice   | func                         | 1  |
      | twisted | fate is a funny lol champion | 2  |

  @read
  Scenario: I should be able to know if an error occurs
    When I GET http://localhost:8080/documents/100
    Then response status code should be {{status.NotFound}}
    And json response should contain:
      | field | matcher | value     |
      | error | =       | not found |


  # DELETE
  @delete
  Scenario Outline: I should be able to delete document from known ids
    When I DELETE http://localhost:8080/documents/<id>
    Then response status code should be {{status.Ok}}

    Examples:
      | id |
      | 0  |
      | 1  |
      | 2  |

  @delete
  Scenario: I should be able to know if an error occurs
    When I DELETE http://localhost:8080/documents/1
    Then response status code should be {{status.BadRequest}}
    And json response should contain:
      | field | matcher | value     |
      | error | =       | not found |

```

## Features

### API testing

Kactus provides a set of steps for API management. You can install them through the `InstallAPI` method.
You can also directly use exposed methods in you own steps.

#### Headers management

| Step                                                        | Method                        | Usage                                                                                                   | Example                                                      |
|-------------------------------------------------------------|-------------------------------|---------------------------------------------------------------------------------------------------------|--------------------------------------------------------------|
| `(?:I )?set(?:ting)? request headers:`                      | `api.Client.SetRequestHeader` | Add headers values to API client using a Gherkin table. :warning: This setps overwrite existing headers | `Given I set request headers:`                               |
| `(?:I )?assign(?:ing)? request headers:`                    | `api.Client.AddHeaders`       | Adds or update headers values to API client using a Gherkin table                                       | `Given I assign request headers:`                            |
| `(?:I )?set(?:ing)? ([a-zA-Z0-9-]+) request header to (.+)` | `api.Client.SetHeader`        | Add a single header value to API client                                                                 | `Given I set Authorisation request headers to Bearer XXXXX:` |

###### Credit

Logo: Image
par <a href="https://pixabay.com/fr/users/Katillustrationlondon-10871763/?utm_source=link-attribution&amp;utm_medium=referral&amp;utm_campaign=image&amp;utm_content=4294916">
Katherine Ab</a>
de <a href="https://pixabay.com/fr/?utm_source=link-attribution&amp;utm_medium=referral&amp;utm_campaign=image&amp;utm_content=4294916">
Pixabay</a>
