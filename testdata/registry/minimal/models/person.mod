name: Person
fields:
  ID:
    type: AutoIncrement
    attributes:
      - mandatory
  FirstName:
    type: String
  LastName:
    type: String
  Nationality:
    type: Nationality
identifiers:
  primary: ID
  name:
    - FirstName
    - LastName
related:
  ContactInfo:
    type: HasOne
  Company:
    type: ForOne
  Comment:
    type: HasOnePoly
    through: Commentable
  Tag:
    type: HasManyPoly
    through: Taggable