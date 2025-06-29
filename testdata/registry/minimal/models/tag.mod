name: Tag
fields:
  ID:
    type: AutoIncrement
    attributes:
      - mandatory
  Name:
    type: String
  Color:
    type: String
identifiers:
  primary: ID
  name: Name
related:
  Taggable:
    type: ForManyPoly
    for:
      - Person
      - Company 