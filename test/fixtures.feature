@fixture
Feature: Kactus should provide a fixture system

  Scenario:
  I should be able to initialize test data from files provided using the Load step definition

    Given I load fixture test/fixtures/test.sql into pg.database
    * I load fixture test/fixtures/test2.sql into pg.authent
    * I load fixture test/fixtures/test_data.yml
    # * I load fixture test/fixtures/pubsub.yml into gcp.client

    Then I expect data to be in pg.database:
      | table       | id | val  |
      | test_kactus | 0  | test |
      | test_kactus | 1  | deux |

    * I expect data to be in pg.authent:
      | table       | id | val   |
      | test_kactus | 2  | test2 |
      | test_kactus | 3  | trois |

    * I want to assert picked variables matches:
      | key    | matcher | value                                    |
      | test   |         | true((bool))                             |
      | devil  |         | 666((int))                               |
      | sos    | =       | 1001((int))                              |
      | script |         | To be or not to be. That's the question. |
      | me     |         | George                                   |
      | song   |         | Le gorille                               |
      | sing   |         | false((bool))                            |

    #* I expect message to be received by gcp.client within 5 seconds:
    #  | field | matcher  | value        |
    #  | test  |          | true((bool)) |
    #  | msg   | contains | test message |
#
    #* I expect message to be received by gcp.client within 1 seconds:
    #  | field | matcher | value          |
    #  | test  |         | false((bool))  |
    #  | msg   |         | This is a soup |
#
    #* I expect message to be received by gcp.client within 1 seconds:
    #  | field | matcher | value         |
    #  | test  |         | false((bool)) |
    #  | msg   |         | Alohomora     |
#
    #* I expect message to be received by gcp.client within 1 seconds having metadata:
    #  | field   | matcher | value |
    #  | version |         | test  |
    #  | type    |         | sort  |


  Scenario: From manifest
  I should be able to initialize test data from full manifest

    Given I load manifest test/fixtures/test_manifest.yml

    Then I expect data to be in pg.database:
      | table                | id | val  |
      | test_kactus_manifest | 0  | test |
      | test_kactus_manifest | 1  | deux |

    * I expect data to be in pg.authent:
      | table                | id | val  |
      | test_kactus_manifest | 2  | aaaa |
      | test_kactus_manifest | 3  | troy |

    * I want to assert picked variables matches:
      | key      | matcher | value                    |
      | manifest |         | true((bool))             |
      | sos      |         | _.._                     |
      | value    | =       | 100.56((number))         |
      | test     |         | [troll,de,troy]((array)) |

    #* I expect message to be received by gcp.client within 5 seconds:
    #  | field | matcher  | value         |
    #  | test  |          | false((bool)) |
    #  | msg   | contains | funny message |
#
    #* I expect message to be received by gcp.client within 1 seconds:
    #  | field | matcher | value            |
    #  | test  |         | true((bool))     |
    #  | msg   |         | This is a cookie |
#
    #* I expect message to be received by gcp.client within 1 seconds:
    #  | field | matcher | value        |
    #  | test  |         | true((bool)) |
    #  | msg   |         | AvadaKadavra |

  #@manifest:test/fixtures/test_manifest_tag.yml
  #Scenario: From tag
  #  Given I assume that fixtures was loaded from tag
#
  #  Then I expect data to be in pg.database:
  #    | table           | id | val  |
  #    | test_kactus_tag | 0  | test |
  #    | test_kactus_tag | 1  | deux |
#
  #  * I want to assert picked variables matches:
  #    | key   | matcher | value        |
  #    | tag   |         | true((bool)) |
  #    | suift |         | clean        |
#
  #  * I expect message to be received by gcp.client within 5 seconds:
  #    | field | matcher  | value         |
  #    | test  |          | false((bool)) |
  #    | msg   | contains | tag message   |
  #    | tag   |          | true((bool))  |
#
  #  * I expect message to be received by gcp.client within 1 seconds:
  #    | field | matcher | value             |
  #    | test  |         | false((bool))     |
  #    | msg   |         | This is a lasagna |
  #    | tag   |         | true((bool))      |
#
  #  * I expect message to be received by gcp.client within 1 seconds:
  #    | field | matcher | value         |
  #    | test  |         | false((bool)) |
  #    | msg   |         | stupefix      |
  #    | tag   |         | true((bool))  |
