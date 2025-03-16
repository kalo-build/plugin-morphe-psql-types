package compile

type MorpheTableType string

const (
	MorpheTableTypeEnums      MorpheTableType = "enum"
	MorpheTableTypeModels     MorpheTableType = "model"
	MorpheTableTypeStructures MorpheTableType = "structure"
	MorpheTableTypeEntities   MorpheTableType = "entity"
)
