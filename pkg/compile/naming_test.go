package compile_test

import (
	"strings"
	"testing"

	"github.com/kaloseia/plugin-morphe-psql-types/pkg/compile"
	"github.com/stretchr/testify/suite"
)

type NamingTestSuite struct {
	suite.Suite
}

func TestNamingTestSuite(t *testing.T) {
	suite.Run(t, new(NamingTestSuite))
}

func (suite *NamingTestSuite) TestGetTableNameFromModel() {
	// Basic cases
	suite.Equal("basics", compile.GetTableNameFromModel("Basic"))
	suite.Equal("users", compile.GetTableNameFromModel("User"))
	suite.Equal("products", compile.GetTableNameFromModel("Product"))

	// Irregular plurals
	suite.Equal("people", compile.GetTableNameFromModel("Person"))
	suite.Equal("men", compile.GetTableNameFromModel("Man"))
	suite.Equal("women", compile.GetTableNameFromModel("Woman"))
	suite.Equal("children", compile.GetTableNameFromModel("Child"))
	suite.Equal("indices", compile.GetTableNameFromModel("Index"))

	// Words ending with special characters
	suite.Equal("categories", compile.GetTableNameFromModel("Category"))
	suite.Equal("boxes", compile.GetTableNameFromModel("Box"))
	suite.Equal("foxes", compile.GetTableNameFromModel("Fox"))
	suite.Equal("countries", compile.GetTableNameFromModel("Country"))

	// Compound words with irregular endings
	suite.Equal("search_indices", compile.GetTableNameFromModel("SearchIndex"))
	suite.Equal("tech_indices", compile.GetTableNameFromModel("TechIndex"))
	suite.Equal("sales_people", compile.GetTableNameFromModel("SalesPerson"))
	suite.Equal("business_men", compile.GetTableNameFromModel("BusinessMan"))
	suite.Equal("police_women", compile.GetTableNameFromModel("PoliceWoman"))
	suite.Equal("foster_children", compile.GetTableNameFromModel("FosterChild"))

	// Edge cases that could trigger incorrect replacements
	suite.Equal("checkboxes", compile.GetTableNameFromModel("Checkbox"))
	suite.Equal("inboxes", compile.GetTableNameFromModel("Inbox"))
	suite.Equal("mailboxes", compile.GetTableNameFromModel("Mailbox"))
	suite.Equal("oxfords", compile.GetTableNameFromModel("Oxford"))
	suite.Equal("managers", compile.GetTableNameFromModel("Manager"))
	suite.Equal("womans_wears", compile.GetTableNameFromModel("WomansWear"))
}

func (suite *NamingTestSuite) TestGetColumnNameFromField() {
	suite.Equal("id", compile.GetColumnNameFromField("ID"))
	suite.Equal("first_name", compile.GetColumnNameFromField("FirstName"))
	suite.Equal("is_active", compile.GetColumnNameFromField("IsActive"))
	suite.Equal("created_at", compile.GetColumnNameFromField("CreatedAt"))
	suite.Equal("customer_id", compile.GetColumnNameFromField("CustomerID"))
}

func (suite *NamingTestSuite) TestGetForeignKeyColumnName() {
	suite.Equal("user_id", compile.GetForeignKeyColumnName("User", "ID"))
	suite.Equal("product_product_id", compile.GetForeignKeyColumnName("Product", "ProductID"))
	suite.Equal("order_item_id", compile.GetForeignKeyColumnName("OrderItem", "ID"))
	suite.Equal("company_company_id", compile.GetForeignKeyColumnName("Company", "CompanyID"))
}

func (suite *NamingTestSuite) TestGetForeignKeyConstraintName() {
	suite.Equal("fk_orders_user_id", compile.GetForeignKeyConstraintName("orders", "user_id"))
	suite.Equal("fk_order_items_product_id", compile.GetForeignKeyConstraintName("order_items", "product_id"))
	suite.Equal("fk_users_company_id", compile.GetForeignKeyConstraintName("users", "company_id"))
	suite.Equal("fk_order_items_order_id", compile.GetForeignKeyConstraintName("order_items", "order_id"))
}

func (suite *NamingTestSuite) TestGetJunctionTableForeignKeyConstraintName() {
	suite.Equal("fk_user_roles_user_id",
		compile.GetJunctionTableForeignKeyConstraintName("user_roles", "User", "ID"))
	suite.Equal("fk_product_categories_product_product_id",
		compile.GetJunctionTableForeignKeyConstraintName("product_categories", "Product", "ProductID"))
	suite.Equal("fk_article_tags_tag_id",
		compile.GetJunctionTableForeignKeyConstraintName("article_tags", "Tag", "ID"))
}

func (suite *NamingTestSuite) TestGetIndexName() {
	suite.Equal("idx_users_email", compile.GetIndexName("users", "email"))
	suite.Equal("idx_orders_customer_id", compile.GetIndexName("orders", "customer_id"))
	suite.Equal("idx_products_sku", compile.GetIndexName("products", "sku"))
	suite.Equal("idx_order_items_order_id", compile.GetIndexName("order_items", "order_id"))
}

func (suite *NamingTestSuite) TestGetUniqueConstraintName() {
	suite.Equal("uk_users_email", compile.GetUniqueConstraintName("users", "email"))
	suite.Equal("uk_orders_order_number", compile.GetUniqueConstraintName("orders", "order_number"))
	suite.Equal("uk_user_roles_user_id_role_id", compile.GetUniqueConstraintName("user_roles", "user_id", "role_id"))
	suite.Equal("uk_product_categories_product_id_category_id",
		compile.GetUniqueConstraintName("product_categories", "product_id", "category_id"))
}

func (suite *NamingTestSuite) TestGetJunctionTableName() {
	suite.Equal("user_roles", compile.GetJunctionTableName("User", "Role"))
	suite.Equal("product_categories", compile.GetJunctionTableName("Product", "Category"))
	suite.Equal("article_tags", compile.GetJunctionTableName("Article", "Tag"))
	suite.Equal("order_items", compile.GetJunctionTableName("Order", "Item"))
}

func (suite *NamingTestSuite) TestGetJunctionTableUniqueConstraintName() {
	suite.Equal("uk_user_roles_user_id_role_id",
		compile.GetJunctionTableUniqueConstraintName("user_roles", "User", "ID", "Role", "ID"))
	suite.Equal("uk_product_categories_product_product_id_category_category_id",
		compile.GetJunctionTableUniqueConstraintName("product_categories", "Product", "ProductID", "Category", "CategoryID"))
	suite.Equal("uk_article_tags_article_id_tag_id",
		compile.GetJunctionTableUniqueConstraintName("article_tags", "Article", "ID", "Tag", "ID"))
}

func (suite *NamingTestSuite) TestPluralize() {
	// Regular pluralization
	suite.Equal("users", compile.Pluralize("user"))
	suite.Equal("items", compile.Pluralize("item"))
	suite.Equal("products", compile.Pluralize("product"))

	// Words ending with s, x, z, ch, sh, ss
	suite.Equal("boxes", compile.Pluralize("box"))
	suite.Equal("addresses", compile.Pluralize("address"))
	suite.Equal("watches", compile.Pluralize("watch"))
	suite.Equal("dishes", compile.Pluralize("dish"))
	suite.Equal("classes", compile.Pluralize("class"))
	suite.Equal("buzzes", compile.Pluralize("buzz"))

	// Words ending with consonant + y
	suite.Equal("categories", compile.Pluralize("category"))
	suite.Equal("countries", compile.Pluralize("country"))
	suite.Equal("companies", compile.Pluralize("company"))
	suite.Equal("cities", compile.Pluralize("city"))

	// Words ending with vowel + y
	suite.Equal("boys", compile.Pluralize("boy"))
	suite.Equal("days", compile.Pluralize("day"))

	// Irregular nouns
	suite.Equal("people", compile.Pluralize("person"))
	suite.Equal("children", compile.Pluralize("child"))
	suite.Equal("feet", compile.Pluralize("foot"))
	suite.Equal("teeth", compile.Pluralize("tooth"))
	suite.Equal("geese", compile.Pluralize("goose"))
	suite.Equal("oxen", compile.Pluralize("ox"))
	suite.Equal("mice", compile.Pluralize("mouse"))
	suite.Equal("men", compile.Pluralize("man"))
	suite.Equal("women", compile.Pluralize("woman"))
	suite.Equal("indices", compile.Pluralize("index"))

	// Compound words with irregular endings (PascalCase)
	suite.Equal("searchIndices", compile.Pluralize("searchIndex"))
	suite.Equal("techIndices", compile.Pluralize("techIndex"))
	suite.Equal("salesPeople", compile.Pluralize("salesPerson"))
	suite.Equal("businessMen", compile.Pluralize("businessMan"))
	suite.Equal("policeWomen", compile.Pluralize("policeWoman"))
	suite.Equal("fosterChildren", compile.Pluralize("fosterChild"))

	// Compound words with irregular endings (snake_case)
	suite.Equal("search_indices", compile.Pluralize("search_index"))
	suite.Equal("tech_indices", compile.Pluralize("tech_index"))
	suite.Equal("sales_people", compile.Pluralize("sales_person"))
	suite.Equal("business_men", compile.Pluralize("business_man"))

	// Words containing but not ending with irregular patterns
	suite.Equal("boxes", compile.Pluralize("box"))              // Not "boxen"
	suite.Equal("foxes", compile.Pluralize("fox"))              // Not "foxen"
	suite.Equal("oxfords", compile.Pluralize("oxford"))         // Not "oxenford"
	suite.Equal("manuals", compile.Pluralize("manual"))         // Not "menuals"
	suite.Equal("woman_tools", compile.Pluralize("woman_tool")) // Not "women_tool"
}

func (suite *NamingTestSuite) TestAbbreviateIdentifier() {
	// Short identifiers should remain unchanged
	suite.Equal("users", compile.AbbreviateIdentifier("users", false))
	suite.Equal("products", compile.AbbreviateIdentifier("products", false))
	suite.Equal("fk_users_company_id", compile.AbbreviateIdentifier("fk_users_company_id", true))

	// Long identifiers should be abbreviated
	longName := "extremely_long_table_name_that_definitely_exceeds_postgresql_identifier_length_limits_by_a_lot"
	suite.Equal("ex_lo_ta_na_th_de_ex_po_id_le_li_by_a_lo", compile.AbbreviateIdentifier(longName, false))

	// Very long identifier with hash
	veryLongName := "this_is_an_absurdly_long_identifier_name_that_will_absolutely_need_a_hash_suffix_to_avoid_any_potential_collisions_with_similar_identifiers"
	veryLongResult := compile.AbbreviateIdentifier(veryLongName, true)
	suite.LessOrEqual(len(veryLongResult), 63)
	suite.True(strings.HasPrefix(veryLongResult, "th_is_an_ab"),
		"Expected prefix 'th_is_an_ab' but got: %s", veryLongResult)
	hashParts := strings.Split(veryLongResult, "_")
	lastPart := hashParts[len(hashParts)-1]
	suite.Equal(8, len(lastPart), "Expected 8-character hash suffix, got: %s", lastPart)

	// Prefixes like fk_, uk_, idx_ should be preserved
	constraintName := "fk_extremely_long_table_name_referencing_another_extremely_long_table_name_with_a_long_column"
	constraintResult := compile.AbbreviateIdentifier(constraintName, true)
	suite.True(strings.HasPrefix(constraintResult, "fk_"),
		"Expected 'fk_' prefix to be preserved, got: %s", constraintResult)
	suite.LessOrEqual(len(constraintResult), 63)
}

func (suite *NamingTestSuite) TestGetTableNameFromModel_LongNames() {
	veryLongModelName := "ExtremelyLongModelNameThatWouldExceedPostgreSQLIdentifierLengthLimits"

	result := compile.GetTableNameFromModel(veryLongModelName)

	suite.Equal(len(result), 35)
	suite.Equal("ex_lo_mo_na_th_wo_ex_po_sq_id_le_li", result)
}

func (suite *NamingTestSuite) TestGetForeignKeyConstraintName_LongNames() {
	longTableName := "extremely_long_table_name_that_would_exceed_postgresql_identifier_length_limits"
	longColumnName := "another_extremely_long_column_name_that_would_also_exceed_limits"

	result := compile.GetForeignKeyConstraintName(longTableName, longColumnName)

	suite.Equal(len(result), 63)
	suite.Regexp("fk_ex_lo_ta_na_th_wo_ex_po_id_le_li_an_ex_lo_co_na_th__[0-9a-f]{8}", result)
	suite.True(strings.HasPrefix(result, "fk_"),
		"Expected 'fk_' prefix but got: %s", result)
	suite.True(strings.HasPrefix(result, "fk_ex_lo"),
		"Expected 'fk_ex_lo' prefix but got: %s", result)
}

func (suite *NamingTestSuite) TestGetJunctionTableUniqueConstraintName_LongNames() {
	junctionTableName := "very_long_junction_table_name"
	model1Name := "ExtremelyLongFirstModelName"
	model1IdName := "ExtremelyLongFirstModelIdField"
	model2Name := "AnotherExtremelyLongSecondModelName"
	model2IdName := "AnotherExtremelyLongSecondModelIdField"

	result := compile.GetJunctionTableUniqueConstraintName(
		junctionTableName,
		model1Name, model1IdName,
		model2Name, model2IdName,
	)

	suite.LessOrEqual(len(result), 63)
	suite.Regexp("uk_ve_lo_ju_ta_na_ex_lo_fi_mo_na_ex_lo_fi_mo_id_fi_an__[0-9a-f]{8}", result)
	suite.True(strings.HasPrefix(result, "uk_"),
		"Expected 'uk_' prefix but got: %s", result)
	suite.True(strings.Contains(result, "ve_lo_ju_ta_na"),
		"Expected to contain 've_lo_ju_ta_na' but got: %s", result)
}
