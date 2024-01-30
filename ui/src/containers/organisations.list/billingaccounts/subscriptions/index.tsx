import { DataTable, EmptyState, Flex } from "@raystack/apsara";
import { V1Beta1Organization } from "@raystack/frontier";
import { useFrontier } from "@raystack/frontier/react";
import { useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import { OrganizationsHeader } from "../../header";
import { getColumns } from "./columns";

export default function OrganisationBASubscriptions() {
  const { client } = useFrontier();
  let { organisationId, billingaccountId } = useParams();
  const [organisation, setOrganisation] = useState<V1Beta1Organization>();
  const [subscriptions, setSubscriptions] = useState([]);

  const pageHeader = {
    title: "Organizations",
    breadcrumb: [
      {
        href: `/organisations`,
        name: `Organizations list`,
      },
      {
        href: `/organisations/${organisationId}`,
        name: `${organisation?.name}`,
      },
      {
        href: `/organisations/${organisationId}/billingaccounts/${billingaccountId}`,
        name: `${billingaccountId}`,
      },
      {
        href: "",
        name: `Organizations Billing Account's subsctriptions`,
      },
    ],
  };

  useEffect(() => {
    async function getOrganization() {
      const {
        // @ts-ignore
        data: { organization },
      } = await client?.frontierServiceGetOrganization(organisationId ?? "");
      setOrganisation(organization);
    }
    getOrganization();
  }, [organisationId]);

  useEffect(() => {
    async function getOrganizationSubscriptions() {
      const {
        // @ts-ignore
        data: { subscriptions },
      } = await client?.frontierServiceListSubscriptions(
        organisationId ?? "",
        billingaccountId ?? ""
      );
      setSubscriptions(subscriptions);
    }
    getOrganizationSubscriptions();
  }, [organisationId ?? ""]);

  let { userId } = useParams();
  const tableStyle = subscriptions?.length
    ? { width: "100%" }
    : { width: "100%", height: "100%" };

  return (
    <Flex direction="row" style={{ height: "100%", width: "100%" }}>
      <DataTable
        data={subscriptions ?? []}
        // @ts-ignore
        columns={getColumns(subscriptions)}
        emptyState={noDataChildren}
        parentStyle={{ height: "calc(100vh - 60px)" }}
        style={tableStyle}
      >
        <DataTable.Toolbar>
          <OrganizationsHeader header={pageHeader} />
          <DataTable.FilterChips style={{ paddingTop: "16px" }} />
        </DataTable.Toolbar>
      </DataTable>
    </Flex>
  );
}

export const noDataChildren = (
  <EmptyState>
    <div className="svg-container"></div>
    <h3>0 subsctription created</h3>
  </EmptyState>
);

export const TableDetailContainer = ({ children }: any) => (
  <div>{children}</div>
);
