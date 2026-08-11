name: Person
fields:
  ID:
    type: AutoIncrement
  Name:
    type: String
  Email:
    type: String
identifiers:
  primary: ID
related:
  WorkContact:
    type: ForOne
    aliased: Contact
  PersonalContact:
    type: ForOne
    aliased: Contact