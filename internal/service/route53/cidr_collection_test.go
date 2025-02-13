package route53_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/aws/aws-sdk-go/service/route53"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tfroute53 "github.com/hashicorp/terraform-provider-aws/internal/service/route53"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

func TestAccRoute53CIDRCollection_basic(t *testing.T) {
	ctx := acctest.Context(t)
	var v route53.CollectionSummary
	resourceName := "aws_route53_cidr_collection.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, route53.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckCIDRCollectionDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccCIDRCollection_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCIDRCollectionExists(ctx, resourceName, &v),
					acctest.MatchResourceAttrGlobalARNNoAccount(resourceName, "arn", "route53", regexp.MustCompile(`cidrcollection/.+`)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccRoute53CIDRCollection_disappears(t *testing.T) {
	ctx := acctest.Context(t)
	var v route53.CollectionSummary
	resourceName := "aws_route53_cidr_collection.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, route53.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckCIDRCollectionDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccCIDRCollection_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCIDRCollectionExists(ctx, resourceName, &v),
					acctest.CheckFrameworkResourceDisappears(acctest.Provider, tfroute53.ResourceCIDRCollection, resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccCheckCIDRCollectionDestroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acctest.Provider.Meta().(*conns.AWSClient).Route53Conn()

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "aws_route53_cidr_collection" {
				continue
			}

			_, err := tfroute53.FindCIDRCollectionByID(ctx, conn, rs.Primary.ID)

			if tfresource.NotFound(err) {
				continue
			}

			if err != nil {
				return err
			}

			return fmt.Errorf("Route 53 CIDR Collection %s still exists", rs.Primary.ID)
		}

		return nil
	}
}

func testAccCheckCIDRCollectionExists(ctx context.Context, n string, v *route53.CollectionSummary) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Route 53 CIDR Collection ID is set")
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).Route53Conn()

		output, err := tfroute53.FindCIDRCollectionByID(ctx, conn, rs.Primary.ID)

		if err != nil {
			return err
		}

		*v = *output

		return nil
	}
}

func testAccCIDRCollection_basic(rName string) string {
	return fmt.Sprintf(`
resource "aws_route53_cidr_collection" "test" {
  name = %[1]q
}
`, rName)
}
