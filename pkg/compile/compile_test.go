package compile_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/kalo-build/go-util/assertfile"
	rcfg "github.com/kalo-build/morphe-go/pkg/registry/cfg"
	"github.com/kalo-build/plugin-morphe-psql-types/internal/testutils"
	"github.com/kalo-build/plugin-morphe-psql-types/pkg/compile"
	"github.com/kalo-build/plugin-morphe-psql-types/pkg/compile/cfg"
)

type CompileTestSuite struct {
	assertfile.FileSuite

	TestDirPath            string
	TestGroundTruthDirPath string

	ModelsDirPath     string
	EnumsDirPath      string
	StructuresDirPath string
	EntitiesDirPath   string
}

func TestCompileTestSuite(t *testing.T) {
	suite.Run(t, new(CompileTestSuite))
}

func (suite *CompileTestSuite) SetupTest() {
	suite.TestDirPath = testutils.GetTestDirPath()
	suite.TestGroundTruthDirPath = filepath.Join(suite.TestDirPath, "ground-truth", "compile-minimal")

	suite.ModelsDirPath = filepath.Join(suite.TestDirPath, "registry", "minimal", "models")
	suite.EnumsDirPath = filepath.Join(suite.TestDirPath, "registry", "minimal", "enums")
	suite.StructuresDirPath = filepath.Join(suite.TestDirPath, "registry", "minimal", "structures")
	suite.EntitiesDirPath = filepath.Join(suite.TestDirPath, "registry", "minimal", "entities")
}

func (suite *CompileTestSuite) TearDownTest() {
	suite.TestDirPath = ""
}

func (suite *CompileTestSuite) TestMorpheToPSQL() {
	workingDirPath := suite.TestDirPath + "/working"
	suite.Nil(os.Mkdir(workingDirPath, 0644))
	defer os.RemoveAll(workingDirPath)

	config := compile.MorpheCompileConfig{
		MorpheLoadRegistryConfig: rcfg.MorpheLoadRegistryConfig{
			RegistryEnumsDirPath:      suite.EnumsDirPath,
			RegistryStructuresDirPath: suite.StructuresDirPath,
			RegistryModelsDirPath:     suite.ModelsDirPath,
			RegistryEntitiesDirPath:   suite.EntitiesDirPath,
		},
		MorpheConfig: cfg.MorpheConfig{
			MorpheModelsConfig: cfg.MorpheModelsConfig{
				Schema:       "public",
				UseBigSerial: false,
			},
			MorpheStructuresConfig: cfg.MorpheStructuresConfig{
				Schema:            "public",
				UseBigSerial:      false,
				EnablePersistence: true,
			},
			MorpheEnumsConfig: cfg.MorpheEnumsConfig{
				Schema:       "public",
				UseBigSerial: false,
			},
			// MorpheEntitiesConfig: cfg.MorpheEntitiesConfig{
			// 	// 	Schema: "public",
			// },
		},

		ModelWriter: &compile.MorpheTableFileWriter{
			Type:          compile.MorpheTableTypeModels,
			TargetDirPath: workingDirPath + "/models",
		},

		StructureWriter: &compile.MorpheTableFileWriter{
			Type:          compile.MorpheTableTypeStructures,
			TargetDirPath: workingDirPath + "/structures",
		},

		EnumWriter: &compile.MorpheTableFileWriter{
			Type:          compile.MorpheTableTypeEnums,
			TargetDirPath: workingDirPath + "/enums",
		},

		// EntityWriter: &compile.MorpheTableFileWriter{
		// 	Type:          compile.MorpheTableTypeEntities,
		// 	TargetDirPath: workingDirPath + "/entities",
		// },
	}

	compileErr := compile.MorpheToPSQL(config)

	suite.NoError(compileErr)

	modelsDirPath := workingDirPath + "/models"
	gtModelsDirPath := suite.TestGroundTruthDirPath + "/models"
	suite.DirExists(modelsDirPath)

	modelPath0 := modelsDirPath + "/contact_infos.sql"
	gtModelPath0 := gtModelsDirPath + "/contact_infos.sql"
	suite.FileExists(modelPath0)
	suite.FileEquals(modelPath0, gtModelPath0)

	modelPath1 := modelsDirPath + "/companies.sql"
	gtModelPath1 := gtModelsDirPath + "/companies.sql"
	suite.FileExists(modelPath1)
	suite.FileEquals(modelPath1, gtModelPath1)

	modelPath2 := modelsDirPath + "/people.sql"
	gtModelPath2 := gtModelsDirPath + "/people.sql"
	suite.FileExists(modelPath2)
	suite.FileEquals(modelPath2, gtModelPath2)

	enumsDirPath := workingDirPath + "/enums"
	gtEnumsDirPath := suite.TestGroundTruthDirPath + "/enums"
	suite.DirExists(enumsDirPath)

	enumPath0 := enumsDirPath + "/nationalities.sql"
	gtEnumPath0 := gtEnumsDirPath + "/nationalities.sql"
	suite.FileExists(enumPath0)
	suite.FileEquals(enumPath0, gtEnumPath0)

	enumPath1 := enumsDirPath + "/universal_numbers.sql"
	gtEnumPath1 := gtEnumsDirPath + "/universal_numbers.sql"
	suite.FileExists(enumPath1)
	suite.FileEquals(enumPath1, gtEnumPath1)

	// Test structure persistence
	structuresDirPath := workingDirPath + "/structures"
	gtStructuresDirPath := suite.TestGroundTruthDirPath + "/structures"
	suite.DirExists(structuresDirPath)

	structurePath0 := structuresDirPath + "/morphe_structures.sql"
	gtStructurePath0 := gtStructuresDirPath + "/morphe_structures.sql"
	suite.FileExists(structurePath0)
	suite.FileEquals(structurePath0, gtStructurePath0)

	// entitiesDirPath := workingDirPath + "/entities"
	// gtEntitiesDirPath := suite.TestGroundTruthDirPath + "/entities"
	// suite.DirExists(entitiesDirPath)

	// entityPath0 := entitiesDirPath + "/company.go"
	// gtEntityPath0 := gtEntitiesDirPath + "/company.go"
	// suite.FileExists(entityPath0)
	// suite.FileEquals(entityPath0, gtEntityPath0)

	// entityPath1 := entitiesDirPath + "/person.go"
	// gtEntityPath1 := gtEntitiesDirPath + "/person.go"
	// suite.FileExists(entityPath1)
	// suite.FileEquals(entityPath1, gtEntityPath1)
}
