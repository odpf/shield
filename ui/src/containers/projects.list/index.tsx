import { EmptyState, Flex, Table } from "@odpf/apsara";
import { Outlet } from "react-router-dom";
import useSWR from "swr";
import { tableStyle } from "~/styles";
import { fetcher } from "~/utils/helper";
import { columns } from "./columns";
import { ProjectsHeader } from "./header";

export default function ProjectList() {
  const { data, error } = useSWR("/admin/v1beta1/projects", fetcher);
  const { projects = [] } = data || { projects: [] };
  return (
    <Flex direction="row" css={{ height: "100%", width: "100%" }}>
      <Table
        css={tableStyle}
        columns={columns}
        data={projects ?? []}
        noDataChildren={noDataChildren}
      >
        <Table.TopContainer>
          <ProjectsHeader />
        </Table.TopContainer>
      </Table>
      <Outlet />
    </Flex>
  );
}

export const noDataChildren = (
  <EmptyState>
    <div className="svg-container"></div>
    <h3>0 project created</h3>
    <div className="pera">Try creating a new project.</div>
  </EmptyState>
);
