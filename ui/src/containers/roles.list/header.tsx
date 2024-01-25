import { DataTable, useTable } from "@raystack/apsara";
import PageHeader from "~/components/page-header";

const pageHeader = {
  title: "Roles",
  breadcrumb: [],
};

export const RolesHeader = () => {
  const { filteredColumns, table } = useTable();
  const isFiltered = filteredColumns.length > 0;

  return (
    <PageHeader title={pageHeader.title} breadcrumb={pageHeader.breadcrumb}>
      {isFiltered ? <DataTable.ClearFilter /> : <DataTable.FilterOptions />}
      <DataTable.ViewOptions />
      <DataTable.GloabalSearch placeholder="Search roles..." />
    </PageHeader>
  );
};
