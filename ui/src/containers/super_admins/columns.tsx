import { V1Beta1ServiceUser, V1Beta1User } from "@raystack/frontier";
import { ColumnDef } from "@tanstack/react-table";
import { Link } from "react-router-dom";

export const getColumns: () => ColumnDef<
  V1Beta1User | V1Beta1ServiceUser,
  any
>[] = () => {
  return [
    {
      header: "Title",
      accessorKey: "title",
      filterVariant: "text",
      cell: (info) => info.getValue() || "-",
    },
    {
      header: "Email",
      accessorKey: "email",
      filterVariant: "text",
      cell: (info) => info.getValue() || "-",
    },
    {
      header: "Status",
      accessorKey: "state",
      meta: {
        data: [
          { label: "Enabled", value: "enabled" },
          { label: "Disabled", value: "disabled" },
        ],
      },
      cell: (info) => info.getValue(),
      footer: (props) => props.column.id,
      filterFn: (row, id, value) => {
        return value.includes(row.getValue(id));
      },
    },
    {
      header: "Organization",
      accessorKey: "org_id",
      cell: (info) => {
        const org_id = info.getValue();
        return org_id ? (
          <Link to={`/organisations/${org_id}`}>{org_id}</Link>
        ) : (
          "-"
        );
      },
    },
  ];
};
