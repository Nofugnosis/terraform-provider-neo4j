package neo4j

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	neo4j "github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

func TestResourceGrant(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testResourceGrantConfig_basic(),
				Check: resource.ComposeTestCheckFunc(
					testAccGrantExists("neo4j_role.test"),
					resource.TestCheckResourceAttr("neo4j_role.test", "name", "testRole"),
				),
			},
		},
	})
}

func TestImportResourceGrant(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:        testResourceGrantConfig_import(),
				ResourceName:  "neo4j_grant.reader",
				ImportState:   true,
				ImportStateId: "ACCESS:*:reader",
			},
			{
				Config:        testResourceGrantConfig_import(),
				ResourceName:  "neo4j_grant.reader_match_all_node",
				ImportState:   true,
				ImportStateId: "MATCH:*:reader_*_NODE:*",
			},
			{
				Config:        testResourceGrantConfig_import(),
				ResourceName:  "neo4j_grant.reader_match_all_relationship",
				ImportState:   true,
				ImportStateId: "MATCH:*:reader_*_RELATIONSHIP:*",
			},
			{
				Config:        testResourceGrantConfig_import(),
				ResourceName:  "neo4j_grant.admin_access_all",
				ImportState:   true,
				ImportStateId: "ACCESS:*:admin",
			},
			{
				Config:        testResourceGrantConfig_import(),
				ResourceName:  "neo4j_grant.admin_transaction_management_all",
				ImportState:   true,
				ImportStateId: "TRANSACTION-MANAGEMENT:*:admin_*",
			},
		},
	})
}

func testAccGrantExists(rn string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[rn]
		if !ok {
			return fmt.Errorf("resource not found: %s", rn)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("role id not set")
		}

		c, err := testAccProvider.Meta().(*Neo4jConfiguration).GetDbConn()
		if err != nil {
			return err
		}
		defer c.Close()
		session := c.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
		defer session.Close()
		result, err := neo4j.Single(session.Run("SHOW ROLE $name PRIVILEGES YIELD role RETURN role LIMIT 1", map[string]interface{}{"username": rs.Primary.ID}))

		fmt.Print(result)

		return nil
	}
}

func testResourceGrantConfig_basic() string {
	return fmt.Sprint(`
	provider "neo4j" {
		host      = "neo4j://localhost:7687"
		username = "neo4j"
		password = "password"
	}
	resource "neo4j_role" "test" {
		name = "testRole"
	}
	resource "neo4j_user" "test" {
		name = "testUser"
		password = "test"
		roles = [
			neo4j_role.test.name
		]
	}
	resource "neo4j_grant" "test" {
		role = "${neo4j_role.test.name}"
		privilege = "READ"
		resource = "*"
		name = "*"
		entity_type = "NODE"
		entity = "*"
	}
	`)
}

func testResourceGrantConfig_import() string {
	return fmt.Sprint(`
	provider "neo4j" {
		host      = "neo4j://localhost:7687"
		username = "neo4j"
		password = "password"
	}
	resource "neo4j_role" "reader" {
		name = "reader"
	}
	resource "neo4j_role" "admin" {
		name = "admin"
	}
	resource "neo4j_grant" "reader" {
		role = "${neo4j_role.reader.name}"
		privilege = "ACCESS"
		name = "*"
	}
	resource "neo4j_grant" "reader_match_all_node" {
		role        = "${neo4j_role.reader.name}"
		privilege   = "MATCH"
		resource    = "*"
		name        = "*"
		entity_type = "NODE"
	}
	resource "neo4j_grant" "reader_match_all_relationship" {
		role        = "${neo4j_role.reader.name}"
		privilege   = "MATCH"
		resource    = "*"
		name        = "*"
		entity_type = "RELATONSHIPO"
	}
	resource "neo4j_grant" "admin_access_all" {
		role        = "admin"
		privilege   = "ACCESS"
		name        = "*"
	}
	resource "neo4j_grant" "admin_transaction_management_all" {
		role        = "admin"
		privilege   = "ACCESS"
		resource = "*"
		name        = "*"
	}
	`)
}
