name: Comment
fields:
  ID:
    type: AutoIncrement
    attributes:
      - mandatory
  Content:
    type: String
  CreatedAt:
    type: String
identifiers:
  primary: ID
related:
  Commentable:
    type: ForOnePoly
    for:
      - Person
      - Company 