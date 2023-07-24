import { PlusIcon } from "@radix-ui/react-icons";

import { Button, DataTable, Flex, Text, useTable } from "@raystack/apsara";
import { useNavigate } from "react-router-dom";

export const ProjectsHeader = () => {
  const navigate = useNavigate();
  const { filteredColumns, table } = useTable();
  const isFiltered = filteredColumns.length > 0;

  return (
    <>
      <Flex align="center" justify="between">
        <Text style={{ fontSize: "14px", fontWeight: "500" }}>projects</Text>
        <Flex gap="small">
          {isFiltered ? <DataTable.ClearFilter /> : <DataTable.FilterOptions />}
          <DataTable.ViewOptions />
          <DataTable.GloabalSearch placeholder="Search projects..." />
          <Button
            variant="secondary"
            onClick={() => navigate("/console/projects/create")}
            style={{ width: "100%" }}
          >
            <Flex
              direction="column"
              align="center"
              style={{ paddingRight: "var(--pd-4)" }}
            >
              <PlusIcon />
            </Flex>
            new project
          </Button>
        </Flex>
      </Flex>
    </>
  );
};
